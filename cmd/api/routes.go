package main

import handler "github.com/Wanjie-Ryan/Go-Budget/cmd/api/handlers"

func (app *Application) routes(handler handler.Handler) {
	app.server.GET("/health", handler.HealthCheck)

}
