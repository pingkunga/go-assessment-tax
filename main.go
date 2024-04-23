package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pingkunga/assessment-tax/deductions"
	"github.com/pingkunga/assessment-tax/postgres"
	"github.com/pingkunga/assessment-tax/tax"
)

func main() {
	db, err := postgres.New()
	if err != nil {
		panic(err)
	}

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

	repo := postgres.NewRepository(db)
	service := deductions.NewService(repo)
	handler := deductions.NewHandler(service)

	authoriedRoute.POST("admin/deductions/personal", handler.SetPersonalDeductionHandler)

	e.Logger.Fatal(e.Start(":" + APP_PORT))
}

var ADMIN_USERNAME string
var ADMIN_PASSWORD string
var APP_PORT string

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

	if os.Getenv("PORT") == "" {
		APP_PORT = "8080"
	} else {
		APP_PORT = os.Getenv("PORT")
		_, err := strconv.Atoi(APP_PORT)
		if err != nil {
			panic("PORT must be a number")
		}
	}

	if os.Getenv("DATABASE_URL") == "" {
		panic("DATABASE_URL is not set")
	}
}
