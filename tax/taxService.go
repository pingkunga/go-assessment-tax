package tax

import (
	"errors"
	"fmt"

	"github.com/dustin/go-humanize"
)

func CalculateTax(tax TaxRequest) TaxResponse {

	netIncome := calculateNetIncome(tax)

	var taxAmount float64
	var taxLevels []TaxLevel
	for _, step := range TaxStepList() {
		stepTaxAmount := 0.0
		if netIncome <= 0 {
			stepTaxAmount = 0
		} else if netIncome > step.MaxIncome {
			stepTaxAmount = (step.MaxIncome * step.TaxRate)
			netIncome = netIncome - step.MaxIncome
		} else {
			stepTaxAmount = (netIncome * step.TaxRate)
			netIncome = netIncome - step.MaxIncome
		}

		taxLevels = append(taxLevels, TaxLevel{Level: step.FriendlyMessage, Tax: stepTaxAmount})
		taxAmount = taxAmount + stepTaxAmount
	}

	taxAmount = taxAmount - tax.WHT

	if taxAmount < 0 {
		return TaxResponse{Tax: 0, TaxRefund: taxAmount * -1, TaxLevels: taxLevels}
	}
	return TaxResponse{Tax: taxAmount, TaxRefund: 0, TaxLevels: taxLevels}
}

func calculateNetIncome(tax TaxRequest) float64 {
	netIncome := tax.TotalIncome - PersonalDeduction()

	for _, allowance := range tax.Allowances {
		netIncome = netIncome - allowanceAmountGuard(allowance)
	}

	return netIncome
}

func allowanceAmountGuard(allowance Allowance) float64 {

	if allowance.AllowanceType == DEDUCTION_DONATION && allowance.Amount > DEDUCTION_DONATION_MAX {
		return DEDUCTION_DONATION_MAX
	}
	return allowance.Amount
}

const DEDUCTION_DONATION = "donation"
const DEDUCTION_DONATION_MAX = 100000

const DEDUCTION_K_RECEIPT = "k-receipt"

var ErrTotalIncome = errors.New("TotalIncome is less than 0")
var ErrInvalidWHT = errors.New("WHT is less than 0")
var ErrInvalidWHTMoreThanTotalIncome = errors.New("WHT is more than TotalIncome")
var ErrInvalidAllowanceAmount = errors.New("Allowance amount is less than 0")
var ErrNotSupportAllowanceType = errors.New("Allowance type support: donation, k-receipt")

func ValidateTaxRequest(tax TaxRequest) (err error) {
	if tax.TotalIncome < 0 {
		err = errors.Join(err, ErrTotalIncome)
	}

	if tax.WHT < 0 {
		err = errors.Join(err, ErrTotalIncome)
	}

	if tax.TotalIncome < tax.WHT {
		err = errors.Join(err, ErrInvalidWHTMoreThanTotalIncome)
	}

	for _, allowance := range tax.Allowances {

		if allowance.AllowanceType != DEDUCTION_DONATION && allowance.AllowanceType != DEDUCTION_K_RECEIPT {
			err = errors.Join(err, ErrNotSupportAllowanceType)
		}

		if allowance.Amount < 0 {
			err = errors.Join(err, ErrInvalidAllowanceAmount)
		}

	}

	return nil
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
	taxStep := []TaxStep{
		{seq: 1, MinIncome: 0, MaxIncome: 150000, TaxRate: 0},
		{seq: 2, MinIncome: 150001, MaxIncome: 500000, TaxRate: 0.1},
		{seq: 3, MinIncome: 500001, MaxIncome: 1000000, TaxRate: 0.15},
		{seq: 4, MinIncome: 1000001, MaxIncome: 2000000, TaxRate: 0.2},
		{seq: 5, MinIncome: 2000001, MaxIncome: 999999999, TaxRate: 0.35},
	}

	//Update FriendlyMessage
	for i, step := range taxStep {
		if taxStep[i].MaxIncome == 999999999 {
			taxStep[i].FriendlyMessage = fmt.Sprintf("%s ขึ้นไป", humanize.FormatFloat("#,###.##", step.MinIncome))
		} else {
			taxStep[i].FriendlyMessage = fmt.Sprintf("%s - %s", humanize.FormatFloat("#,###.##", step.MinIncome), humanize.FormatFloat("#,###.##", step.MaxIncome))
		}
	}

	return taxStep
}
