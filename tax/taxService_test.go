package tax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTax(t *testing.T) {
	t.Run("As user, I want to calculate my tax", func(t *testing.T) {
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
		taxAmount := CalculateTax(tax)

		assert.EqualValues(t, 29000.0, taxAmount)
	})

}
