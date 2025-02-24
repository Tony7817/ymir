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

	"ymir.com/app/order/model"
	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"
	"ymir.com/pkg/id"
	"ymir.com/pkg/vars"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type CreatePaypalOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePaypalOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaypalOrderLogic {
	return &CreatePaypalOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePaypalOrderLogic) CreatePaypalOrder(in *order.CreatePaypalOrderRequest) (*order.CreatePaypalOrderResponse, error) {
	oItems, err := l.svcCtx.OrderItemModel.FindOneByOrderIdNotSoftDelete(l.ctx, in.OrderId)
	if err != nil {
		return nil, err
	}
	var orderItems = make([]*order.OrderItem, len(oItems))
	for i := 0; i < len(oItems); i++ {
		orderItems[i] = &order.OrderItem{
			ProductId:   oItems[i].ProductId,
			ColorId:     oItems[i].ProductColorId,
			Size:        oItems[i].Size,
			Qunantity:   oItems[i].Quantity,
			Price:       oItems[i].Price,
			OrderItemId: oItems[i].Id,
		}
	}
	paypalOrderId, err := l.createPaypalOrder(orderItems, in.OrderId, in.RequestId)
	if err != nil {
		return nil, errors.Wrapf(err, "create paypal order failed")
	}

	return &order.CreatePaypalOrderResponse{
		OrderId:       in.OrderId,
		PaypalOrderId: paypalOrderId,
	}, nil
}

func (l *CreatePaypalOrderLogic) createPaypalOrder(orderItems []*order.OrderItem, orderId int64, requestId string) (string, error) {
	tokenStr, err := l.svcCtx.Redis.GetCtx(l.ctx, vars.PaypalAccessTokenKey)
	var token = &vars.PaypalTokenResponse{}
	if err != nil {
		return "", err
	}
	if tokenStr == "" {
		token, err = getPaypalToken()
		if err != nil {
			return "", err
		}
		_, err = l.svcCtx.Redis.SetnxExCtx(l.ctx, vars.PaypalAccessTokenKey, token.AccessToken, token.ExpiresIn-60)
		if err != nil {
			l.Logger.Errorf("[CreateOrder] set paypal token to redis error: %+v", err)
			return "", err
		}
	} else {
		token.AccessToken = tokenStr
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
		l.Logger.Errorf("[CreateOrder] marshal paypal create order request error: %+v", err)
		return "", err
	}

	request, err := http.NewRequest("POST", "https://api.sandbox.paypal.com/v2/checkout/orders", bytes.NewBuffer(body))
	if err != nil {
		l.Logger.Errorf("[CreateOrder] create paypal order request error: %+v", err)
		return "", err
	}
	defer request.Body.Close()

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("PayPal-Request-Id", requestId)
	request.Header.Set("Authorization", "Bearer "+token.AccessToken)

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
			l.Logger.Errorf("[CreateOrder] decode paypal create order response error: %+v", err)
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
		l.Logger.Errorf("[CreateOrder] request paypal create order failed, resp:%+v", resp)
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
