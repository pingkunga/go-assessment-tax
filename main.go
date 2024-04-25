package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pingkunga/assessment-tax/deductions"
	"github.com/pingkunga/assessment-tax/postgres"
	"github.com/pingkunga/assessment-tax/tax"
)

func main() {
	db, err := postgres.New()
	if err != nil {
		panic("Cannot connect to Database, Detail" + err.Error())
	}

	e := echo.New()

	//Middleware-Log
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/hello-world", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	e.POST("/tax/calculations", tax.CalculationsHandler)

	//Authorized
	authoriedRoute := e.Group("/admin")

	//Middleware-Auth
	authoriedRoute.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == ADMIN_USERNAME && password == ADMIN_PASSWORD {
			return true, nil
		}
		return false, nil
	}))

	repo := postgres.NewRepository(db)
	service := deductions.NewService(repo)
	handler := deductions.NewHandler(service)

	authoriedRoute.GET("/deductions", handler.DeductionConfigsHandler)
	authoriedRoute.POST("/deductions/personal", handler.SetPersonalDeductionHandler)

	/*
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		//Start server
		go func() {
			if err := e.Start(":" + APP_PORT); err != nil && err != http.ErrServerClosed { // Start server
				e.Logger.Fatal("shutting down the server")
			}
		}()
		<-shutdown

		// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	*/

	// Start server in a goroutine so that it doesn't block
	go func() {
		err := e.Start(":" + APP_PORT)
		if err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	// 1 = buffer นะ
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
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
