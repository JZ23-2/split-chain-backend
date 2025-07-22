package dtos

type ReceiptItem struct {
	Name          string  `json:"name"`
	Quantity      float32 `json:"quantity"`
	UnitPrice     float32 `json:"unitPrice"`
	TotalPrice    float32 `json:"totalPrice"`
	PriceAfterTax float32 `json:"priceAfterTax"`
	PriceInHBAR   float32 `json:"priceInHBAR"`
}

type ReceiptResponse struct {
	StoreName   string        `json:"storeName"`
	Date        string        `json:"date"`
	Tax         float32       `json:"tax"`
	TotalAmount float32       `json:"totalAmount"`
	Items       []ReceiptItem `json:"items"`
}
