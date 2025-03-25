package paypal

type PaypalCreateOrderRequest struct {
	Intent        string               `json:"intent"`
	PurchaseUnits []PaypalPurchaseUnit `json:"purchase_units"`
}

type PaypalPurchaseUnit struct {
	Amount PaypalAmount `json:"amount"`
}

type PaypalAmount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type PaypalCreateOrderResponse struct {
	PaypalOrderId string `json:"id"`
}
