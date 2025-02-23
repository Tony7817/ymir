package logic

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ymir.com/app/order/model"
	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"
	"ymir.com/pkg/id"
	"ymir.com/pkg/vars"
	"ymir.com/pkg/xerr"

	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/threading"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	if len(in.Items) > 10 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.ErrorRequestOrderMaximunReach), "Request too many orders, order num: %d", len(in.Items))
	}

	barrier, err := dtmgrpc.BarrierFromGrpc(l.ctx)
	if err != nil {
		l.Logger.Errorf("[CreateOrder] BarrierFromGrpc error: %+v", err)
		return nil, status.Error(codes.Aborted, dtmcli.ResultFailure)
	}

	err = barrier.CallWithDB(l.svcCtx.DB, func(tx *sql.Tx) error {
		err := mr.Finish(func() error {
			var totalPrice int64
			for i := 0; i < len(in.Items); i++ {
				totalPrice += in.Items[i].Price
			}
			_, err := l.svcCtx.OrderModel.SFInsertTx(l.ctx, tx, &model.Order{
				Id:         in.OrderId,
				UserId:     in.UserId,
				Status:     model.OrderStatusPending,
				TotalPrice: totalPrice,
			})
			if err != nil {
				l.Logger.Errorf("[CreateOrder] OrderModel.SFInsertTx error: %+v", err)
				return status.Error(codes.Aborted, dtmcli.ResultFailure)
			}
			return nil
		}, func() error {
			err := l.createOrderItems(tx, in.OrderId, in.Items)
			if err != nil {
				return status.Error(codes.Aborted, dtmcli.ResultFailure)
			}
			return nil
		})
		if err != nil {
			return err
		}
		_, err = l.createPaypalOrder(in.Items, in.OrderId, in.RequestId)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{
		OrderId: in.OrderId,
	}, nil
}

func (l *CreateOrderLogic) createOrderItems(tx *sql.Tx, oId int64, os []*order.OrderItem) error {
	return mr.MapReduceVoid(func(source chan<- *order.OrderItem) {
		for i := 0; i < len(os); i++ {
			source <- os[i]
		}
	},
		func(o *order.OrderItem, writer mr.Writer[any], cancel func(error)) {
			_, err := l.svcCtx.OrderItemModel.SFInsertTx(l.ctx, tx, &model.OrderItem{
				Id:             o.OrderItemId,
				OrderId:        oId,
				ProductId:      o.ProductId,
				ProductColorId: o.ColorId,
				Size:           o.Size,
				Quantity:       o.Qunantity,
				Price:          o.Price,
				Subtotal:       o.Qunantity * o.Price,
			})
			if err != nil {
				l.Logger.Errorf("[CreateOrder] OrderItemModel.SFInsertTx error: %+v", err)
				cancel(status.Error(codes.Aborted, dtmcli.ResultFailure))
			}
		},
		func(pipe <-chan any, cancel func(error)) {})
}

func (l *CreateOrderLogic) createPaypalOrder(orderItems []*order.OrderItem, orderId int64, requestId string) (string, error) {
	tokenStr, err := l.svcCtx.Redis.GetCtx(l.ctx, vars.PaypalAccessTokenKey)
	var token *vars.PaypalTokenResponse
	if err != nil {
		return "", err
	}
	if tokenStr == "" {
		token, err = getPaypalToken()
		if err != nil {
			return "", err
		}
		_, err = l.svcCtx.Redis.SetnxExCtx(l.ctx, vars.PaypalAccessTokenKey, token.AccessToken, token.ExpiresIn)
		if err != nil {
			return "", err
		}
	} else {
		err = json.Unmarshal([]byte(tokenStr), token)
		if err != nil {
			return "", err
		}
	}

	var req = vars.PaypalCreateOrderRequest{
		Intent:        "CAPTURE",
		PurchaseUnits: make([]vars.PaypalPurchaseUnit, 1),
	}

	var totalPrice float64

	for i := 0; i < len(orderItems); i++ {
		price := float64(orderItems[i].Price) / 100
		totalPrice += price
	}
	req.PurchaseUnits[0] = vars.PaypalPurchaseUnit{
		Amount: vars.PaypalAmount{
			CurrencyCode: "USD",
			Value:        strconv.FormatFloat(totalPrice, 'f', 2, 64),
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", "https://api.sandbox.paypal.com/v2/checkout/orders", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer request.Body.Close()

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("PayPal-Request-Id", requestId)
	request.Header.Set("access_token", token.AccessToken)

	var client = &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var createOrderResp vars.PaypalCreateOrderResponse
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		err = json.NewDecoder(resp.Body).Decode(&createOrderResp)
		if err != nil {
			return "", err
		}
		_, err = l.svcCtx.PaypalModel.Insert(l.ctx, &model.Paypal{
			Id:            id.NewSnowFlake().GenerateID(),
			OrderId:       orderId,
			RequestId:     requestId,
			PaypalOrderId: createOrderResp.OrderId,
			ReqBody:       sql.NullString{String: string(body), Valid: true},
		})
		if err != nil {
			return "", err
		}

		// set the order in redis
		threading.GoSafe(func() {
			_, _ = l.svcCtx.PaypalModel.FindOneByRequestId(l.ctx, requestId)
		})
	} else {
		return "", errors.New("create paypal order failed")
	}

	return createOrderResp.OrderId, nil
}

func getPaypalToken() (*vars.PaypalTokenResponse, error) {
	var clientId = os.Getenv("PAYPAL_CLIENT_ID")
	var clientSecret = os.Getenv("PAYPAL_SECRET")
	if clientId == "" || clientSecret == "" {
		return nil, errors.New("paypal client id or secret is empty")
	}

	var auth = base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
	var data = url.Values{}
	data.Set("grant_type", "client_credentials")
	var body = strings.NewReader(data.Encode())

	req, err := http.NewRequest("POST", vars.PaypalGetTokenUrl, body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResp vars.PaypalTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}
