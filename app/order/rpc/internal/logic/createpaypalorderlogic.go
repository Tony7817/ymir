package logic

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"ymir.com/app/order/model"
	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"
	"ymir.com/pkg/id"
	"ymir.com/pkg/paypal"
	"ymir.com/pkg/xerr"

	"github.com/go-sql-driver/mysql"
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
	oItems, err := l.svcCtx.OrderItemModel.FindOneByOrderId(l.ctx, in.OrderId)
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
	paypalOrderId, err := l.createPaypalOrder(orderItems, in.OrderId, in.UserId, in.RequestId)
	if err != nil {
		return nil, errors.Wrapf(err, "create paypal order failed")
	}

	return &order.CreatePaypalOrderResponse{
		OrderId:       in.OrderId,
		PaypalOrderId: paypalOrderId,
	}, nil
}

func (l *CreatePaypalOrderLogic) createPaypalOrder(orderItems []*order.OrderItem, orderId, userId int64, requestId string) (string, error) {
	token, err := paypal.GetToken(l.ctx, l.svcCtx.Redis)
	if err != nil {
		return "", err
	}
	var req = paypal.PaypalCreateOrderRequest{
		Intent:        "CAPTURE",
		PurchaseUnits: make([]paypal.PaypalPurchaseUnit, 1),
	}

	var totalPrice float64

	for i := 0; i < len(orderItems); i++ {
		price := float64(orderItems[i].Price) / 100
		totalPrice += price
	}
	req.PurchaseUnits[0] = paypal.PaypalPurchaseUnit{
		Amount: paypal.PaypalAmount{
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
		return "", err
	}
	defer request.Body.Close()

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("PayPal-Request-Id", requestId)
	request.Header.Set("Authorization", "Bearer "+token)

	var client = &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var createOrderResp paypal.PaypalCreateOrderResponse
	var errMsg *paypal.ErrorResponse
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		err = json.NewDecoder(resp.Body).Decode(&createOrderResp)
		if err != nil {
			l.Logger.Errorf("[CreateOrder] decode paypal create order response error: %+v", err)
			return "", err
		}

	} else {
		err = json.NewDecoder(resp.Body).Decode(&errMsg)
		if err != nil {
			l.Logger.Errorf("[CreateOrder] decode paypal create order error response error: %+v", err)
			return "", err
		}
		l.Logger.Errorf("[CreateOrder] request paypal create order failed, resp:%+v", resp)
	}

	_, err = l.svcCtx.PaypalModel.Insert(l.ctx, &model.Paypal{
		Id:            id.NewSnowFlake().GenerateID(),
		OrderId:       orderId,
		UserId:        userId,
		RequestId:     requestId,
		PaypalOrderId: createOrderResp.OrderId,
		ReqBody:       sql.NullString{String: string(body), Valid: true},
		ErrMsg:        paypal.SetErrorInNullString(errMsg),
	})
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return "", xerr.NewErrCode(xerr.ErrorOrderCreated)
			}
		}
		return "", err
	}

	// set the order in redis
	threading.GoSafe(func() {
		_, _ = l.svcCtx.PaypalModel.FindOneByRequestId(l.ctx, requestId)
	})

	if errMsg != nil {
		return "", xerr.NewErrCode(xerr.ErrorCreateOrder)
	}

	return createOrderResp.OrderId, nil
}
