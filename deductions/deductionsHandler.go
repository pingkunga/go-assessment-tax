package deductions

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pingkunga/assessment-tax/common"
	repo "github.com/pingkunga/assessment-tax/postgres"
)

type DeductionHandler struct {
	service IDeductionService
}

type IDeductionService interface {
	SetPersonalDeduction(request DebuctionRequest) (PersonalDeductionResponse, error)
	SetKPlustDeduction(request DebuctionRequest) (KReceiptResponse, error)
	GetPersonalDeduction() (PersonalDeductionResponse, error)
	DeductionConfigs() ([]repo.DeductionConfig, error)
}

func NewHandler(service IDeductionService) *DeductionHandler {
	return &DeductionHandler{service: service}
}

func (h *DeductionHandler) SetPersonalDeductionHandler(c echo.Context) error {

	var debuctionRequest DebuctionRequest
	if err := c.Bind(&debuctionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	deductionsResponse, error := h.service.SetPersonalDeduction(debuctionRequest)
	if error != nil {
		return c.JSON(http.StatusBadRequest, common.Err{Message: error.Error()})
	}

	return c.JSON(http.StatusOK, deductionsResponse)
}

func (h *DeductionHandler) DeductionConfigsHandler(c echo.Context) error {

	deduction, error := h.service.DeductionConfigs()
	if error != nil {
		return c.JSON(http.StatusBadRequest, common.Err{Message: error.Error()})
	}

	return c.JSON(http.StatusOK, deduction)
}
