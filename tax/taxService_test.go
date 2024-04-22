package tax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTax(t *testing.T) {
	t.Run("Story: EXP01-0: As user, I want to calculate my tax", func(t *testing.T) {
		tax := TaxRequest{
			TotalIncome: 500000.0,
			WHT:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        0.0,
				},
			},
		}
		taxResponse := CalculateTax(tax)

		assert.EqualValues(t, 29000.0, taxResponse.Tax)
		assert.EqualValues(t, 0.0, taxResponse.TaxRefund)

	})

	t.Run("Story: EXP01-1: As user, I want to calculate my tax with Income 2,000,001", func(t *testing.T) {
		tax := TaxRequest{
			TotalIncome: 2000001.0,
			WHT:         0.0,
		}
		taxResponse := CalculateTax(tax)

		assert.EqualValues(t, 258000.2, taxResponse.Tax)
		assert.EqualValues(t, 0.0, taxResponse.TaxRefund)
	})

	t.Run("Story: EXP02-0: As user, I want to calculate my tax with WHT", func(t *testing.T) {
		tax := TaxRequest{
			TotalIncome: 500000.0,
			WHT:         25000.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        0.0,
				},
			},
		}
		taxResponse := CalculateTax(tax)

		assert.EqualValues(t, 4000.0, taxResponse.Tax)
		assert.EqualValues(t, 0.0, taxResponse.TaxRefund)

	})

	t.Run("Story: EXP03-0: As user, I want to calculate my tax with donation", func(t *testing.T) {
		tax := TaxRequest{
			TotalIncome: 500000.0,
			WHT:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        200000.0,
				},
			},
		}

		taxResponse := CalculateTax(tax)

		assert.EqualValues(t, 19000.0, taxResponse.Tax)
		assert.EqualValues(t, 0.0, taxResponse.TaxRefund)

	})

	t.Run("Story: EXP03-1: As user, I want to calculate my tax with donation and WHT, should return taxRefund", func(t *testing.T) {
		tax := TaxRequest{
			TotalIncome: 500000.0,
			WHT:         20000.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        200000.0,
				},
			},
		}

		taxResponse := CalculateTax(tax)

		assert.EqualValues(t, 0.0, taxResponse.Tax)
		assert.EqualValues(t, 1000.0, taxResponse.TaxRefund)
	})

}
