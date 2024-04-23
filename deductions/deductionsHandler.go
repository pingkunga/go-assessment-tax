package deductions

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// type DeductionHandler struct {
// 	store DeductionStorer
// }

// func New(db DeductionStorer) *DeductionHandler {
// 	return &DeductionHandler{store: db}
// }

// for implement interface
type DeductionStorer interface {
	getAllowanceConfig(pKey string) (float64, error)
	setAllowanceConfig(pKey string, newValue float64) error
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SetPersonalDeductionHandler(c echo.Context) error {

	var debuctionRequest DebuctionRequest
	if err := c.Bind(&debuctionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	deductionsResponse, error := h.service.SetPersonalDeduction(debuctionRequest)
	if error != nil {
		return c.JSON(http.StatusBadRequest, error)
	}

	return c.JSON(http.StatusOK, deductionsResponse)
}

func (h *Handler) AllowanceConfigsHandler(c echo.Context) error {

	allowances, error := h.service.AllowanceConfigs()
	if error != nil {
		return c.JSON(http.StatusBadRequest, error)
	}

	return c.JSON(http.StatusOK, allowances)
}
