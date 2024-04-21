package tax

func CalculateTax(tax TaxRequest) float64 {

	netTaxAmount := tax.TotalIncome - PersonalDeduction()

	var taxAmount float64
	for _, step := range TaxStepList() {
		taxAmount = taxAmount + (netTaxAmount * step.TaxRate)
		netTaxAmount = netTaxAmount - step.MaxIncome

		if netTaxAmount < 0 {
			break
		}
	}

	taxAmount = taxAmount - tax.WHT

	return taxAmount
}

func PersonalDeduction() float64 {
	return 60000
}

/*
การคำนวนภาษีตามขั้นบันใด
รายได้ 0 - 150,000 ได้รับการยกเว้น
150,001 - 500,000 อัตราภาษี 10%
500,001 - 1,000,000 อัตราภาษี 15%
1,000,001 - 2,000,000 อัตราภาษี 20%
มากกว่า 2,000,000 อัตราภาษี 35%
*/
func TaxStepList() []TaxStep {
	return []TaxStep{
		{seq: 1, MinIncome: 0, MaxIncome: 150000, TaxRate: 0},
		{seq: 2, MinIncome: 150001, MaxIncome: 500000, TaxRate: 0.1},
		{seq: 3, MinIncome: 500001, MaxIncome: 1000000, TaxRate: 0.15},
		{seq: 4, MinIncome: 1000001, MaxIncome: 2000000, TaxRate: 0.2},
		{seq: 5, MinIncome: 2000001, MaxIncome: 999999999, TaxRate: 0.35},
	}
}
