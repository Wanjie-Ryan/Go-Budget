package handler

import (
	"fmt"
	"net/http"

	request "github.com/Wanjie-Ryan/Go-Budget/cmd/api/requests"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Registerhandler(c echo.Context) error {

	// 1. bind request body
	payload := new(request.RegisterUserRequest)
	// echo.DefaultBinder is Echo's built-in helper that knows how to::: i. Read the requests Content-type ii. parse the JSON body iii. match JSON keys (firstname, email) to your struct fields via the json tags.
	// bind body unmarshals the raw JSON into your payload object
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		fmt.Println("register error", err)
		return c.String(http.StatusBadRequest, "Invalid request body")
	}
	fmt.Println("register payload", payload)

	// 2. validate request body
	// 3. create user and other stuff
	return c.String(http.StatusOK, "Success")

}
