package tax

type TaxRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxResponse struct {
	Tax       float64 `json:"tax"`
	TaxRefund float64 `json:"taxRefund"`
}

type TaxStep struct {
	seq       int
	MinIncome float64
	MaxIncome float64
	TaxRate   float64
}
