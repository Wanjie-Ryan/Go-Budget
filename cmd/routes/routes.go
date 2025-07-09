package routes

import (
	"context"

	handler "github.com/Wanjie-Ryan/Go-Budget/cmd/api/handlers"
	middlewares "github.com/Wanjie-Ryan/Go-Budget/cmd/api/middleware"
	"github.com/Wanjie-Ryan/Go-Budget/internal/mailer"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	// middlewares "github.com/Wanjie-Ryan/Go-Budget/cmd/api/middleware"
)

// what will do here now is dividing the routes into the ones that require authentication, and those that dont't.
// also we'll add a prefix to the routes

type Application struct {
	logger        echo.Logger
	server        *echo.Echo
	handler       handler.Handler
	appMiddleware middlewares.AppMiddleware
}

func (app *Application) Initial(db *gorm.DB, c context.Context) {

	app.server = echo.New()
	app.handler = handler.Handler{DB: db, Mailer: mailer.NewMailer()}
	app.appMiddleware = middlewares.AppMiddleware{DB: db}

	go app.Routes(app.handler)

}
func (app *Application) Routes(handler handler.Handler) {

	apiGroup := app.server.Group("/api")
	publicAuthRoutes := apiGroup.Group("/auth")
	// {
	publicAuthRoutes.POST("/register", handler.Registerhandler)
	publicAuthRoutes.POST("/login", handler.Loginhandler)
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
