package deductions

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	repo "github.com/pingkunga/assessment-tax/postgres"
	"github.com/stretchr/testify/assert"
)

func TestSetPersonalDeductionHandler(t *testing.T) {
	// <base_url>/deductions/personal
	t.Run("Given SetPersonalDeduction with match Schema, Should Successfully", func(t *testing.T) {
		// Your test code here
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", io.NopCloser(strings.NewReader(`{"amount": 70000}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a new StubRepository
		repo := &StubRepository{
			GetDeductionStub: 100000,
			DeductionConfigStub: repo.DeductionConfig{
				DeductionType: "personal",
				DeductionMin:  10000,
				DeductionMax:  100000,
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)
		handler := NewHandler(service)

		err := handler.SetPersonalDeductionHandler(c)

		assert.NoError(t, err)
		assert.EqualValues(t, http.StatusOK, rec.Code)
		assert.EqualValues(t, "{\"personalDeduction\":70000}\n", rec.Body.String())
	})

	//Validate Error
	t.Run("Given SetPersonalDeduction with input not in range, Should Return Error", func(t *testing.T) {
		// Your test code here
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", io.NopCloser(strings.NewReader(`{"amount": 4000}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a new StubRepository
		repo := &StubRepository{
			GetDeductionStub: 100000,
			DeductionConfigStub: repo.DeductionConfig{
				DeductionType: "personal",
				DeductionMin:  10000,
				DeductionMax:  100000,
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)
		handler := NewHandler(service)

		err := handler.SetPersonalDeductionHandler(c)

		assert.NoError(t, err)
		assert.EqualValues(t, http.StatusBadRequest, rec.Code)
		assert.EqualValues(t, "{\"message\":\"personal deduction must be greater than 10,000.00\"}\n", rec.Body.String())
	})
}

func TestSetKReceiptDeductionHandler(t *testing.T) {

	t.Run("Given Set k-receipt with match Schema, Should Successfully", func(t *testing.T) {
		// Your test code here
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", io.NopCloser(strings.NewReader(`{"amount": 60000}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a new StubRepository
		repo := &StubRepository{
			GetDeductionStub: 100000,
			DeductionConfigStub: repo.DeductionConfig{
				DeductionType: "k-receipt",
				DeductionMin:  0,
				DeductionMax:  100000,
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)
		handler := NewHandler(service)

		err := handler.SetKReceiptDeductionHandler(c)

		assert.NoError(t, err)
		assert.EqualValues(t, http.StatusOK, rec.Code)
		assert.EqualValues(t, "{\"kReceipt\":60000}\n", rec.Body.String())
	})

	//Validate Error
	t.Run("Given SetPersonalDeduction with input not in range, Should Return Error", func(t *testing.T) {
		// Your test code here
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", io.NopCloser(strings.NewReader(`{"amount": 100001}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a new StubRepository
		repo := &StubRepository{
			GetDeductionStub: 100000,
			DeductionConfigStub: repo.DeductionConfig{
				DeductionType: "k-receipt",
				DeductionMin:  0,
				DeductionMax:  100000,
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)
		handler := NewHandler(service)

		err := handler.SetKReceiptDeductionHandler(c)

		assert.NoError(t, err)
		assert.EqualValues(t, http.StatusBadRequest, rec.Code)
		assert.EqualValues(t, "{\"message\":\"k-receipt deduction must be less than 100,000.00\"}\n", rec.Body.String())
	})
}

func TestDeductionConfigsHandler(t *testing.T) {
	t.Run("Given DeductionConfigsHandler in general case, Should Successfully with mock data 2 record", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/deductions", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Create a new StubRepository
		repo := &StubRepository{
			GetDeductionStub: 100000,
			DeductionConfigsStub: []repo.DeductionConfig{
				{
					DeductionType: "personal",
					DeductionMin:  10000,
					DeductionMax:  100000,
				},
				{
					DeductionType: "k-receipt",
					DeductionMin:  0,
					DeductionMax:  100000,
				},
			},
		}

		// Create a new DeductionsService
		service := NewService(repo)
		handler := NewHandler(service)

		err := handler.DeductionConfigsHandler(c)

		assert.NoError(t, err)
		assert.EqualValues(t, http.StatusOK, rec.Code)
		assert.EqualValues(t, "[{\"id\":0,\"deductionType\":\"personal\",\"deductionAmount\":0,\"deductionMin\":10000,\"adeductionMax\":100000},{\"id\":0,\"deductionType\":\"k-receipt\",\"deductionAmount\":0,\"deductionMin\":0,\"adeductionMax\":100000}]\n", rec.Body.String())
	})
}
