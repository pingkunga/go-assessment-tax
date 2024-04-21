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

	taxAmount := CalculateTax(taxRequest)

	TaxResponse := TaxResponse{Tax: taxAmount}
	return c.JSON(http.StatusOK, TaxResponse)
}
