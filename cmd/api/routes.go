package main

import (
	handler "github.com/Wanjie-Ryan/Go-Budget/cmd/api/handlers"
	// middlewares "github.com/Wanjie-Ryan/Go-Budget/cmd/api/middleware"
)

// what will do here now is dividing the routes into the ones that require authentication, and those that dont't.
// also we'll add a prefix to the routes
func (app *Application) routes(handler handler.Handler) {

	apiGroup := app.server.Group("/api")
	publicAuthRoutes := apiGroup.Group("/auth")
	// {
	publicAuthRoutes.POST("/register", handler.Registerhandler)
	publicAuthRoutes.POST("/login", handler.Loginhandler)
	publicAuthRoutes.POST("/reset-token", handler.ForgotPassword)
	publicAuthRoutes.POST("/reset-password", handler.ResetPasswordHandler)
	// }

	profileAuthRoutes := apiGroup.Group("/profile", app.appMiddleware.AuthMiddleware)
	{
		profileAuthRoutes.GET("/authenticated/user", handler.GetAuthUserHandler)
		profileAuthRoutes.PATCH("/update/password", handler.UpdateUserPassword)

	}

	app.server.GET("/health", handler.HealthCheck)
	// ORIGINAL
	// app.server.POST("/register", handler.Registerhandler)
	// app.server.POST("/login", handler.Loginhandler)

	// app.server.GET("/authenticated/user", handler.GetAuthUserHandler, app.appMiddleware.AuthMiddleware)

}
