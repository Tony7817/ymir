package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"ymir.com/app/order/model"
	"ymir.com/app/order/rpc/internal/svc"
	"ymir.com/app/order/rpc/order"
	"ymir.com/pkg/id"
	"ymir.com/pkg/paypal"
	"ymir.com/pkg/util"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CaptureOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCaptureOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptureOrderLogic {
	return &CaptureOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CaptureOrderLogic) CaptureOrder(in *order.CapturePaypalOrderRequest) (*order.CapturePaypalOrderResposne, error) {
	token, err := paypal.GetToken(l.ctx, l.svcCtx.Redis)
	if err != nil {
		return nil, err
	}

	var url = paypal.PaypalCaptureOrderUrl(in.PaypalOrderId)
	var body = []byte(`{}`)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	var client = &http.Client{}
	respb, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer respb.Body.Close()

	bodyStr, err := io.ReadAll(respb.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body error")
	}

	var res paypal.CaptureOrderCreatedResponse
	var errMsg paypal.ErrorResponse
	if respb.StatusCode == http.StatusCreated {
		err := json.Unmarshal(bodyStr, &res)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal response body error")
		}
		_, err = l.svcCtx.OrderAddressModel.Insert(l.ctx, &model.OrderAddress{
			Id:         id.SF.GenerateID(),
			OrderId:    in.OrderId,
			UserId:     in.UserId,
			FullName:   res.PurchaseUnits[0].Shipping.Name.FullName,
			Email:      res.Payer.EmailAddress,
			Address1:   util.IsNullString(res.PurchaseUnits[0].Shipping.Address.AddressLine1),
			Address2:   util.IsNullString(res.PurchaseUnits[0].Shipping.Address.AddressLine2),
			AdminArea1: util.IsNullString(res.PurchaseUnits[0].Shipping.Address.AdminArea1),
			AdminArea2: util.IsNullString(res.PurchaseUnits[0].Shipping.Address.AdminArea2),
			PostalCode: res.PurchaseUnits[0].Shipping.Address.PostalCode,
		})
		if err != nil {
			return nil, err
		}
	} else if respb.StatusCode == http.StatusOK {
		return nil, xerr.NewErrCode(xerr.ErrorOrderPaied)
	} else {
		err = json.NewDecoder(respb.Body).Decode(&errMsg)
		if err != nil {
			return nil, errors.Wrap(err, "decode error response body error")
		}
	}

	return &order.CapturePaypalOrderResposne{
		OrderId:    in.OrderId,
		Intent:     res.Intent,
		Status:     res.Status,
		CreateTime: res.CreateTime,
		PayerInfo: &order.PaypalPayer{
			Email:         res.Payer.EmailAddress,
			PayerFullName: res.PurchaseUnits[0].Shipping.Name.FullName,
			PayerAddress: &order.PayerAddress{
				AddressLine1: res.PurchaseUnits[0].Shipping.Address.AddressLine1,
				AddressLine2: res.PurchaseUnits[0].Shipping.Address.AddressLine2,
				AdminArea1:   res.PurchaseUnits[0].Shipping.Address.AdminArea1,
				AdminArea2:   res.PurchaseUnits[0].Shipping.Address.AdminArea2,
				PostalCode:   res.PurchaseUnits[0].Shipping.Address.PostalCode,
				CountryCode:  res.PurchaseUnits[0].Shipping.Address.CountryCode,
			},
		},
	}, nil
}
