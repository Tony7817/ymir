// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3
// Source: order.proto

package orderclient

import (
	"context"

	"ymir.com/app/order/rpc/order"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CapturePaypalOrderRequest      = order.CapturePaypalOrderRequest
	CapturePaypalOrderResposne     = order.CapturePaypalOrderResposne
	CheckOrderBelongTouserRequest  = order.CheckOrderBelongTouserRequest
	CheckOrderBelongTouserResponse = order.CheckOrderBelongTouserResponse
	CheckOrderIdempotentRequest    = order.CheckOrderIdempotentRequest
	CheckOrderIdempotentResponse   = order.CheckOrderIdempotentResponse
	CreateOrderRequest             = order.CreateOrderRequest
	CreateOrderResponse            = order.CreateOrderResponse
	CreatePaypalOrderRequest       = order.CreatePaypalOrderRequest
	CreatePaypalOrderResponse      = order.CreatePaypalOrderResponse
	DeleteOrderRequest             = order.DeleteOrderRequest
	DeleteOrderResponse            = order.DeleteOrderResponse
	GetOrderAddressRequest         = order.GetOrderAddressRequest
	GetOrderAddressResponse        = order.GetOrderAddressResponse
	GetOrderItemRequest            = order.GetOrderItemRequest
	GetOrderItemResponse           = order.GetOrderItemResponse
	GetOrderListRequest            = order.GetOrderListRequest
	GetOrderListResponse           = order.GetOrderListResponse
	GetOrderRequest                = order.GetOrderRequest
	GetOrderResponse               = order.GetOrderResponse
	OrderAddress                   = order.OrderAddress
	OrderContent                   = order.OrderContent
	OrderItem                      = order.OrderItem
	PayOrderRequest                = order.PayOrderRequest
	PayOrderResponse               = order.PayOrderResponse
	PayerAddress                   = order.PayerAddress
	PayerName                      = order.PayerName
	PayerPhone                     = order.PayerPhone
	PayerTaxInfo                   = order.PayerTaxInfo
	Paypal                         = order.Paypal
	PaypalOrderResponse            = order.PaypalOrderResponse
	PaypalOrderReuqest             = order.PaypalOrderReuqest
	PaypalPayer                    = order.PaypalPayer
	UpdateOrderRequest             = order.UpdateOrderRequest
	UpdateOrderResponse            = order.UpdateOrderResponse

	Order interface {
		CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
		CreateOrderRollback(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
		DeleteOrder(ctx context.Context, in *DeleteOrderRequest, opts ...grpc.CallOption) (*DeleteOrderResponse, error)
		DeleteOrderRollback(ctx context.Context, in *DeleteOrderRequest, opts ...grpc.CallOption) (*DeleteOrderResponse, error)
		GetOrder(ctx context.Context, in *GetOrderRequest, opts ...grpc.CallOption) (*GetOrderResponse, error)
		UpdateOrder(ctx context.Context, in *UpdateOrderRequest, opts ...grpc.CallOption) (*UpdateOrderResponse, error)
		PayOrder(ctx context.Context, in *PayOrderRequest, opts ...grpc.CallOption) (*PayOrderResponse, error)
		OrderList(ctx context.Context, in *GetOrderListRequest, opts ...grpc.CallOption) (*GetOrderListResponse, error)
		PaypalOrder(ctx context.Context, in *PaypalOrderReuqest, opts ...grpc.CallOption) (*PaypalOrderResponse, error)
		CreatePaypalOrder(ctx context.Context, in *CreatePaypalOrderRequest, opts ...grpc.CallOption) (*CreatePaypalOrderResponse, error)
		CaptureOrder(ctx context.Context, in *CapturePaypalOrderRequest, opts ...grpc.CallOption) (*CapturePaypalOrderResposne, error)
		OrderAddress(ctx context.Context, in *GetOrderAddressRequest, opts ...grpc.CallOption) (*GetOrderAddressResponse, error)
		GetOrderItem(ctx context.Context, in *GetOrderItemRequest, opts ...grpc.CallOption) (*GetOrderItemResponse, error)
		CheckOrderIdempotent(ctx context.Context, in *CheckOrderIdempotentRequest, opts ...grpc.CallOption) (*CheckOrderIdempotentResponse, error)
		CheckOrderBelongToUser(ctx context.Context, in *CheckOrderBelongTouserRequest, opts ...grpc.CallOption) (*CheckOrderBelongTouserResponse, error)
	}

	defaultOrder struct {
		cli zrpc.Client
	}
)

func NewOrder(cli zrpc.Client) Order {
	return &defaultOrder{
		cli: cli,
	}
}

func (m *defaultOrder) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.CreateOrder(ctx, in, opts...)
}

func (m *defaultOrder) CreateOrderRollback(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.CreateOrderRollback(ctx, in, opts...)
}

func (m *defaultOrder) DeleteOrder(ctx context.Context, in *DeleteOrderRequest, opts ...grpc.CallOption) (*DeleteOrderResponse, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.DeleteOrder(ctx, in, opts...)
}

func (m *defaultOrder) DeleteOrderRollback(ctx context.Context, in *DeleteOrderRequest, opts ...grpc.CallOption) (*DeleteOrderResponse, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.DeleteOrderRollback(ctx, in, opts...)
}

func (m *defaultOrder) GetOrder(ctx context.Context, in *GetOrderRequest, opts ...grpc.CallOption) (*GetOrderResponse, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.GetOrder(ctx, in, opts...)
}

func (m *defaultOrder) UpdateOrder(ctx context.Context, in *UpdateOrderRequest, opts ...grpc.CallOption) (*UpdateOrderResponse, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.UpdateOrder(ctx, in, opts...)
}

func (m *defaultOrder) PayOrder(ctx context.Context, in *PayOrderRequest, opts ...grpc.CallOption) (*PayOrderResponse, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.PayOrder(ctx, in, opts...)
}

func (m *defaultOrder) OrderList(ctx context.Context, in *GetOrderListRequest, opts ...grpc.CallOption) (*GetOrderListResponse, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.OrderList(ctx, in, opts...)
}

func (m *defaultOrder) PaypalOrder(ctx context.Context, in *PaypalOrderReuqest, opts ...grpc.CallOption) (*PaypalOrderResponse, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.PaypalOrder(ctx, in, opts...)
}

func (m *defaultOrder) CreatePaypalOrder(ctx context.Context, in *CreatePaypalOrderRequest, opts ...grpc.CallOption) (*CreatePaypalOrderResponse, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.CreatePaypalOrder(ctx, in, opts...)
}

func (m *defaultOrder) CaptureOrder(ctx context.Context, in *CapturePaypalOrderRequest, opts ...grpc.CallOption) (*CapturePaypalOrderResposne, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.CaptureOrder(ctx, in, opts...)
}

func (m *defaultOrder) OrderAddress(ctx context.Context, in *GetOrderAddressRequest, opts ...grpc.CallOption) (*GetOrderAddressResponse, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.OrderAddress(ctx, in, opts...)
}

func (m *defaultOrder) GetOrderItem(ctx context.Context, in *GetOrderItemRequest, opts ...grpc.CallOption) (*GetOrderItemResponse, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.GetOrderItem(ctx, in, opts...)
}

func (m *defaultOrder) CheckOrderIdempotent(ctx context.Context, in *CheckOrderIdempotentRequest, opts ...grpc.CallOption) (*CheckOrderIdempotentResponse, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.CheckOrderIdempotent(ctx, in, opts...)
}

func (m *defaultOrder) CheckOrderBelongToUser(ctx context.Context, in *CheckOrderBelongTouserRequest, opts ...grpc.CallOption) (*CheckOrderBelongTouserResponse, error) {
	client := order.NewOrderClient(m.cli.Conn())
	return client.CheckOrderBelongToUser(ctx, in, opts...)
}
