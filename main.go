package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pingkunga/assessment-tax/tax"
)

func main() {
	e := echo.New()

	//Middleware-Log
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/hello-world", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	//Authorized
	authoriedRoute := e.Group("/")

	//Middleware-Auth
	authoriedRoute.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == ADMIN_USERNAME && password == ADMIN_PASSWORD {
			return true, nil
		}
		return false, nil
	}))

	authoriedRoute.POST("tax/calculations", tax.CalculationsHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

var ADMIN_USERNAME string
var ADMIN_PASSWORD string

func init() {
	// Load environment variables

	// IF ENV NOT FOUND THEN THROW ERROR
	if os.Getenv("ADMIN_USERNAME") == "" {
		panic("ADMIN_USERNAME is not set")
	}
	if os.Getenv("ADMIN_PASSWORD") == "" {
		panic("ADMIN_PASSWORD is not set")
	}

	ADMIN_USERNAME = os.Getenv("ADMIN_USERNAME")
	ADMIN_PASSWORD = os.Getenv("ADMIN_PASSWORD")

}
