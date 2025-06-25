package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// creating a reusable struct
type Application struct{
	logger echo.Logger
	server echo.Echo
}

func main() {

	e := echo.New()
	// loading or rather initializing the go-dotenv package
	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal("Error loading .env file", err)
		// e.Logger.Fatal("Error loading .env file")
	}

	// creates a new Echo instance which holds your roites, middleware stack, logger, etc.
	// c echo.Context bundles up request data(path params, query strings, headers, body, cookies) and gives you response helpers
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	port := os.Getenv("APP_PORT")
	appAddress := fmt.Sprintf("localhost:%s", port)
	// the code below is wrapped in e.logger.Fatal() in a case where it returns an error, echo will log the error message, and exit your program with a non zero status code
	e.Logger.Fatal(e.Start(appAddress))

}
