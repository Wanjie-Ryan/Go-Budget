package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	middlewares "github.com/Wanjie-Ryan/Go-Budget/cmd/api/middleware"
	"github.com/Wanjie-Ryan/Go-Budget/cmd/routes"
	"github.com/Wanjie-Ryan/Go-Budget/common"

	// "github.com/Wanjie-Ryan/Go-Budget/internal/mailer"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// // creating a reusable struct
// type Application struct {
// 	logger        echo.Logger
// 	server        *echo.Echo
// 	handler       handler.Handler
// 	appMiddleware middlewares.AppMiddleware
// }

func main() {

	e := echo.New()
	c := context.Background()
	// loading or rather initializing the go-dotenv package
	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal("Error loading .env file", err)
		// e.Logger.Fatal("Error loading .env file")
	}
	//db is *gorm.DB
	db, err := common.NewMySql()
	if err != nil {
		e.Logger.Fatal(err.Error())
	}
	// common.NewMySql()

	// creates a new Echo instance which holds your roites, middleware stack, logger, etc.
	// c echo.Context bundles up request data(path params, query strings, headers, body, cookies) and gives you response helpers
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	// appMailer := mailer.NewMailer()

	// h := handler.Handler{DB: db, Mailer: appMailer}
	// appMiddleware := middlewares.AppMiddleware{DB: db}
	// app := Application{
	// 	logger:        e.Logger,
	// 	server:        e,
	// 	handler:       h,
	// 	appMiddleware: appMiddleware,
	// }
	// fmt.Println(app)
	e.Use(middleware.Logger())
	e.Use(middlewares.CustomMiddleware)
	var a routes.Application
	a.Initial(db, c)
	port := os.Getenv("APP_PORT")
	appAddress := fmt.Sprintf("localhost:%s", port)
	// the code below is wrapped in e.logger.Fatal() in a case where it returns an error, echo will log the error message, and exit your program with a non zero status code
	e.Logger.Fatal(e.Start(appAddress))

}
