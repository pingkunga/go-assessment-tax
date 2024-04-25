package deductions

type DebuctionRequest struct {
	Amount float64 `json:"amount"`
}

type PersonalDeductionResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type KReceiptResponse struct {
	KReceipt float64 `json:"kReceipt"`
}
