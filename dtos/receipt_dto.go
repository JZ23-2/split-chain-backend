package dtos

type ReceiptItem struct {
	Name     string  `json:"name" example:"Front and rear brake cables"`
	Quantity float32 `json:"quantity" example:"1"`
	Price    float32 `json:"price" example:"100"`
}

type ReceiptResponse struct {
	StoreName string        `json:"storeName" example:"Nigger Store"`
	BillDate  string        `json:"billDate" example:"2025-10-02"`
	Tax       float32       `json:"tax" example:"9.10"`
	Service   float32       `json:"service"`
	Items     []ReceiptItem `json:"items"`
}
