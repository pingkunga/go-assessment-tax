package tax

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CalculationsHandler(c echo.Context) error {
	var taxRequest TaxRequest
	if err := c.Bind(&taxRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := ValidateTaxRequest(taxRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	TaxResponse := CalculateTax(taxRequest)
	return c.JSON(http.StatusOK, TaxResponse)
}
