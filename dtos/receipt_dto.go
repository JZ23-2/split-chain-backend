package dtos

type ReceiptItem struct {
	Name          string  `json:"name" example:"Front and rear brake cables"`
	Quantity      float32 `json:"quantity" example:"1"`
	UnitPrice     float32 `json:"unitPrice" example:"100"`
	TotalPrice    float32 `json:"totalPrice" example:"100"`
	PriceAfterTax float32 `json:"priceAfterTax" example:"105.88083"`
	PriceInHBAR   float32 `json:"priceInHBAR" example:"0"`
}

type ReceiptResponse struct {
	StoreName   string        `json:"storeName" example:"Nigger Store"`
	Date        string        `json:"date" example:"2025-10-02"`
	Tax         float32       `json:"tax" example:"9.10"`
	TotalAmount float32       `json:"totalAmount" example:"15.40"`
	Items       []ReceiptItem `json:"items"`
}
