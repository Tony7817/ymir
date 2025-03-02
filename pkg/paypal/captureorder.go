package paypal

type CaptureOrderIdempotentResponse struct {
	CreateTime *string `json:"create_time"`
	Orderid    string  `json:"id"`
	Payer      Payer   `json:"payer"`
	Status     string  `json:"status"`
	Intent     string  `json:"intent"`
}

type PayerName struct {
	GivenName string `json:"given_name"`
	Surname   string `json:"surname"`
}

type TaxInfo struct {
	TaxId     string `json:"tax_id"`
	TaxIdType string `json:"tax_id_type"`
}

type PayerAddress struct {
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	AdminArea1   string `json:"admin_area_1"`
	AdminArea2   string `json:"admin_area_2"`
	PostalCode   string `json:"postal_code"`
	CountryCode  string `json:"country_code"`
}

type CaptureOrderPurchaseUnit struct {
	Shipping struct {
		SType string `json:"type"`
		Name  struct {
			FullName string `json:"full_name"`
		} `json:"name"`
		Address PayerAddress `json:"address"`
	} `json:"shipping"`
	Payment struct {
		Capture struct {
			Id     string `json:"id"`
			Status string `json:"status"`
			Amount struct {
				CurrencyCode string `json:"currency_code"`
				Value        string `json:"value"`
			} `json:"amount"`
		}
	} `json:"payment"`
}

type CaptureOrderCreatedResponse struct {
	CreateTime    string                     `json:"create_time,omitempty"`
	OrderId       string                     `json:"id,omitempty"`
	Intent        string                     `json:"intent"`
	Payer         Payer                      `json:"payer,omitempty"`
	PurchaseUnits []CaptureOrderPurchaseUnit `json:"purchase_units,omitempty"`
	Status        string                     `json:"status,omitempty"`
}

type Payer struct {
	Id           string    `json:"payer_id,omitempty"`
	EmailAddress string    `json:"email_address,omitempty"`
	Name         PayerName `json:"name,omitempty"`
	Phone        struct {
		PhoneType   string `json:"phone_type"`
		PhoneNumber string `json:"phone_number"`
	} `json:"phone,omitempty"`
	TaxInfo TaxInfo      `json:"tax_info,omitempty"`
	Address PayerAddress `json:"address,omitempty"`
}
