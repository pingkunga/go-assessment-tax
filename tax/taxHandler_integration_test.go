//go:build integration

package tax

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/pingkunga/assessment-tax/common"
	"github.com/stretchr/testify/assert"
)

// ================================================================
// Test Method Section
// ================================================================
func TestITTaxCalculations(t *testing.T) {
	//Arrange
	var actualTaxReponse TaxResponse

	//Act
	res := common.ClientRequest(http.MethodPost,
		common.Uri("tax/calculations"),
		io.NopCloser(strings.NewReader(
			`{
			"totalIncome": 500000.0
			,"wht": 0.0
			,"allowances": [
				{"allowanceType": "donation"
				,"amount": 200000.0}]}`)))
	err := res.Decode(&actualTaxReponse)

	//Assert
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	//ดูว่ามีของคืนมาไหม ไม่สนใจว่า Value ถูกไหม

	var expectedTaxReponse = TaxResponse{
		Tax:       19000.0,
		TaxRefund: 0.0,
		TaxLevels: []TaxLevel{
			{
				Level: "0-150,000",
				Tax:   0.0,
			},
			{
				Level: "150,001-500,000",
				Tax:   19000.0,
			},
			{
				Level: "500,001-1,000,000",
				Tax:   0.0,
			},
			{
				Level: "1,000,001-2,000,000",
				Tax:   0.0,
			},
			{
				Level: "2,000,001 ขึ้นไป",
				Tax:   0.0,
			},
		},
	}

	assert.EqualValues(t, expectedTaxReponse, actualTaxReponse)
}
