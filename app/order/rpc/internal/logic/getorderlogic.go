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
	"ymir.com/pkg/paypal"
	"ymir.com/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderLogic) GetOrder(in *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	var orderRes *model.Order
	var oiRes []*order.OrderItem
	err := mr.Finish(func() error {
		var err error
		orderRes, err = l.svcCtx.OrderModel.FindOne(l.ctx, in.OrderId)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.Wrapf(err, "find order by id: %d failed, user id: %d", in.OrderId, in.UserId)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Wrapf(xerr.NewErrCode(xerr.ErrorOrderNotExist), "cannot find order by orderId: %d, userId: %d", in.OrderId, in.UserId)
		}
		return nil
	}, func() error {
		var err error
		res, err := l.svcCtx.OrderItemModel.FindOrderItemsByOrderIdUserId(l.ctx, in.OrderId, in.UserId)
		if err != nil {
			return errors.Wrapf(err, "find order items by orderId: %d failed", in.OrderId)
		}
		for i := range res {
			oiRes = append(oiRes, &order.OrderItem{
				ProductId:            res[i].ProductId,
				ColorId:              res[i].ProductColorId,
				Size:                 res[i].Size,
				Qunantity:            res[i].Quantity,
				Price:                res[i].Price,
				OrderItemId:          res[i].Id,
				ProductDescription:   res[i].ProductDescription,
				Color:                res[i].ProductColor,
				ProductColorCoverUrl: res[i].ProductColorCoverUrl,
			})
		}
		return nil
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrapf(err, "get order by orderId: %d failed", in.OrderId)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return &order.GetOrderResponse{
			Order:      nil,
			OrderItems: nil,
		}, nil
	}

	if orderRes.Status == model.OrderStatusPending {
		orderDetail, err := l.getPaypalOrderDetail(in.OrderId)
		if err != nil {
			return nil, errors.Wrap(err, "get paypal order detail failed")
		}
		if orderDetail.Status == model.OrderStatusCompleted {
			orderRes.Status = model.OrderStatusCompleted
			if err := l.svcCtx.OrderModel.UpdateOrderPartial(l.ctx, orderRes.Id, orderRes.UserId, &model.OrderPartial{
				Staus: &model.OrderStatusCompleted,
			}); err != nil {
				return nil, errors.Wrap(err, "update order status failed")
			}
		}
	}

	return &order.GetOrderResponse{
		Order: &order.OrderContent{
			RequestId:     orderRes.RequestId,
			OrderId:       orderRes.Id,
			UserId:        orderRes.UserId,
			TotalPrice:    orderRes.TotalPrice,
			Status:        orderRes.Status,
			Unit:          orderRes.Unit,
			PaypalOrderId: orderRes.PaypalOrderId.String,
			CreatedAt:     orderRes.CreatedAt.Unix(),
		},
		OrderItems: oiRes,
	}, nil
}

func (l *GetOrderLogic) getPaypalOrderDetail(orderId int64) (*paypal.CaptureOrderCreatedResponse, error) {
	token, err := paypal.GetToken(l.ctx, l.svcCtx.Redis)
	if err != nil {
		return nil, err
	}

	var url = paypal.PaypalShowOrderDetailUrl(strconv.Itoa(int(orderId)))
	var body = []byte(`{}`)
	request, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
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

	var res paypal.CaptureOrderCreatedResponse
	if err := json.NewDecoder(respb.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
