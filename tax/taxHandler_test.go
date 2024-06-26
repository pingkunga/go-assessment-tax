package tax

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pingkunga/assessment-tax/common"
	"github.com/stretchr/testify/assert"
)

func TestTaxHandler_CalculationsHandler(t *testing.T) {
	t.Run("Given CalculationsHandler with match Schema, Should Successfully", func(t *testing.T) {
		e := echo.New()

		/*
			{
			"totalIncome": 500000.0,
			"wht": 0.0,
			"allowances": [
				{
				"allowanceType": "donation",
				"amount": 200000.0
				}
			]
			}
		*/
		req := httptest.NewRequest(http.MethodPost,
			"/tax/calculations",
			io.NopCloser(strings.NewReader(
				`{
					"totalIncome": 500000.0
					,"wht": 0.0
					,"allowances": [
						{"allowanceType": "donation"
						,"amount": 200000.0}]}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		service := NewService(MockDeductuinService())
		handler := NewHandler(service)

		err := handler.CalculationsHandler(c)

		assert.NoError(t, err)
		assert.EqualValues(t, http.StatusOK, rec.Code)

		var actualTaxReponse TaxResponse
		json.NewDecoder(rec.Body).Decode(&actualTaxReponse)

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

		//2 อันนี้เหมือนกัน
		if !reflect.DeepEqual(expectedTaxReponse, actualTaxReponse) {
			t.Errorf("expected list of wallets %v, got %v", expectedTaxReponse, actualTaxReponse)
		}
		assert.EqualValues(t, expectedTaxReponse, actualTaxReponse)
	})

	t.Run("Given CalculationsHandler with not match Schema, Should Return Error", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost,
			"/tax/calculations",
			io.NopCloser(strings.NewReader(
				`{
					"totalIncome": -1
					,"wht": 0.0
					,"allowances": [
						{"allowanceType": "donation"
						,"amount": 200000.0}]}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		service := NewService(MockDeductuinService())
		handler := NewHandler(service)

		err := handler.CalculationsHandler(c)

		var actualTaxReponse common.Err
		json.NewDecoder(rec.Body).Decode(&actualTaxReponse)

		var expectedTaxReponse = common.Err{
			Message: "Validate Err:TotalIncome is less than 0\nWHT is more than TotalIncome",
		}
		assert.NoError(t, err)
		assert.EqualValues(t, http.StatusBadRequest, rec.Code)
		assert.EqualValues(t, expectedTaxReponse, actualTaxReponse)
	})
}

func TestTaxHandler_BatchCalculationsHandler(t *testing.T) {
	t.Run("Given BatchCalculationsHandler with match Schema, Should Successfully", func(t *testing.T) {
		e := echo.New()

		//https://stackoverflow.com/questions/43904974/testing-go-http-request-formfile
		//Create Multi-Part
		path := "../sampleCSV/taxes.csv"

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("taxFile", path)
		assert.NoError(t, err)
		sample, err := os.Open(path)
		assert.NoError(t, err)

		_, err = io.Copy(part, sample)
		assert.NoError(t, err)
		assert.NoError(t, writer.Close())

		// Source
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		service := NewService(MockDeductuinService())
		handler := NewHandler(service)

		err = handler.BatchCalculationsHandler(c)

		assert.NoError(t, err)
		assert.EqualValues(t, http.StatusOK, rec.Code)

		var actualTaxReponse TaxBatchsResponse
		json.NewDecoder(rec.Body).Decode(&actualTaxReponse)

		/*
			"taxes": [
				{
				"totalIncome": 500000,
				"tax": 29000,
				"taxRefund": 0
				},
				{
				"totalIncome": 600000,
				"tax": 0,
				"taxRefund": 3000
				},
				{
				"totalIncome": 750000,
				"tax": 3750,
				"taxRefund": 0
				}
			]
		*/
		var expectedTaxReponse = TaxBatchsResponse{
			Taxes: []TaxBatchResponse{
				{
					TotalIncome: 500000.0,
					Tax:         29000.0,
					TaxRefund:   0.0,
				},
				{
					TotalIncome: 600000.0,
					Tax:         0.0,
					TaxRefund:   3000.0,
				},
				{
					TotalIncome: 750000.0,
					Tax:         3750.0,
					TaxRefund:   0.0,
				},
			},
		}

		assert.Equal(t, expectedTaxReponse, actualTaxReponse)

	})

}
