package middlewares

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func CustomMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	//the next keyword in  a middleware is put, so that it can execute the next middleware after it
	fmt.Println("Middleware is running")
	// the middleware below injects a response header on every request
	return func(c echo.Context) error {
		// sets the server HTTP header to "Echo/3.0 on every outgoing reponse"
		c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
		return next(c)
		// the next calls the next handler/middleware in the chain, the request continues its normal flow
	}
}
