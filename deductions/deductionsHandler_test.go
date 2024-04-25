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

	//Validate Error
}
