package main

import (
	handler "github.com/Wanjie-Ryan/Go-Budget/cmd/api/handlers"
	middlewares "github.com/Wanjie-Ryan/Go-Budget/cmd/api/middleware"
)

func (app *Application) routes(handler handler.Handler) {
	app.server.GET("/health", handler.HealthCheck)
	app.server.POST("/register", handler.Registerhandler)
	app.server.POST("/login", handler.Loginhandler)
	app.server.GET("/authenticated/user", handler.GetAuthUserHandler, middlewares.AuthMiddleware)

}
