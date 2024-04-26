package tax

import (
	"errors"
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

func TestValidateTaxRequest(t *testing.T) {

	t.Run("Given: TotalIncome is -1, should return ErrTotalIncome / ErrInvalidWHTMoreThanTotalIncome", func(t *testing.T) {
		tax := TaxRequest{
			TotalIncome: -1.0,
			WHT:         0.0,
		}
		err := ValidateTaxRequest(tax)
		assert.NotNil(t, err)
		assert.Equal(t, true, errors.Is(err, ErrTotalIncomeLessThanZero))
		assert.Equal(t, true, errors.Is(err, ErrWHTMoreThanTotalIncome))
	})

	t.Run("Given: WHT is -1, should return ErrInvalidWHTMoreThanTotalIncome", func(t *testing.T) {
		tax := TaxRequest{
			TotalIncome: 1.0,
			WHT:         -1.0,
		}
		err := ValidateTaxRequest(tax)
		assert.NotNil(t, err)
		assert.Equal(t, true, errors.Is(err, ErrWHTLessThanZero))
	})

	t.Run("Given: Allowance is SSF, should return ErrNotSupportAllowanceType", func(t *testing.T) {
		tax := TaxRequest{
			TotalIncome: 1.0,
			WHT:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "SSF",
					Amount:        0.0,
				},
			},
		}
		err := ValidateTaxRequest(tax)
		assert.NotNil(t, err)
		assert.Equal(t, true, errors.Is(err, ErrNotSupportAllowanceType))
	})

	t.Run("Given: Allowance amount type donation is -1, should return ErrAllowanceAmountLessThanZero", func(t *testing.T) {
		tax := TaxRequest{
			TotalIncome: 1.0,
			WHT:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        -1.0,
				},
			},
		}
		err := ValidateTaxRequest(tax)
		assert.NotNil(t, err)
		assert.Equal(t, true, errors.Is(err, ErrAllowanceAmountLessThanZero))
	})

	t.Run("Given: Allowance amount type donations is -1 / k-recipt is -1 , should return ErrAllowanceAmountLessThanZero", func(t *testing.T) {
		tax := TaxRequest{
			TotalIncome: 1.0,
			WHT:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donations",
					Amount:        -1.0,
				},
				{
					AllowanceType: "k-recipts",
					Amount:        -1.0,
				},
			},
		}
		err := ValidateTaxRequest(tax)
		assert.NotNil(t, err)
		assert.Equal(t, true, errors.Is(err, ErrNotSupportAllowanceType))
		assert.Equal(t, true, errors.Is(err, ErrAllowanceAmountLessThanZero))
	})

}

func TestImportTaxCSV(t *testing.T) {
	t.Run("Story: EXP04-0: As user, I want to import tax csv file with good format, should success", func(t *testing.T) {
		testFile := "../sampleCSV/taxes.csv"

		taxRequestsActual, err := ImportTaxCSV(testFile)

		taxRequestsExp := []TaxRequest{
			{
				TotalIncome: 500000.0,
				WHT:         0.0,
				Allowances: []Allowance{
					{
						AllowanceType: "donation",
						Amount:        0.0,
					},
				},
			},
			{
				TotalIncome: 600000.0,
				WHT:         40000.0,

				Allowances: []Allowance{
					{
						AllowanceType: "donation",
						Amount:        20000.0,
					},
				},
			},
			{
				TotalIncome: 750000.0,
				WHT:         50000.0,
				Allowances: []Allowance{
					{
						AllowanceType: "donation",
						Amount:        15000.0,
					},
				},
			},
		}

		assert.Nil(t, err)
		assert.Equal(t, 3, len(taxRequestsActual))

		assert.Equal(t, taxRequestsExp, taxRequestsActual)
	})
}
