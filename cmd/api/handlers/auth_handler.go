package handler

import (
	"fmt"
	"net/http"

	// handler "github.com/Wanjie-Ryan/Go-Budget/cmd/api/handlers"
	request "github.com/Wanjie-Ryan/Go-Budget/cmd/api/requests"
	// "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// echo context is used in places where the handler is called, it carries everything needed for http request/response
func (h *Handler) Registerhandler(c echo.Context) error {

	// 1. bind request body
	// allocate a new pointer to the RegisterUserRequest
	// the new keyword returns *RegisteruserRequest, which is a pointer
	payload := new(request.RegisterUserRequest)
	// echo.DefaultBinder is Echo's built-in helper that knows how to::: i. Read the requests Content-type ii. parse the JSON body iii. match JSON keys (firstname, email) to your struct fields via the json tags.
	// bind body unmarshals the raw JSON into your payload object
	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		fmt.Println("register error", err)
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}
	fmt.Println("register payload", payload)
	// var validate *validator.Validate
	// validate = validator.New(validator.WithRequiredStructEnabled())
	// validationErrors := validate.Struct(payload)
	// validationErrors := err.(validator.ValidationErrors)
	// fmt.Println("validation errors",validationErrors)
	// for _, err := range validationErrors.(validator.ValidationErrors){
	// 	fmt.Println(err.Namespace())
	// 	fmt.Println(err.Field())
	// 	fmt.Println(err.StructNamespace())
	// 	fmt.Println(err.StructField())
	// 	fmt.Println(err.Tag())
	// 	fmt.Println(err.ActualTag())
	// 	fmt.Println(err.Kind())
	// 	fmt.Println(err.Type())
	// 	fmt.Println(err.Value())
	// 	fmt.Println(err.Param())
	// 	fmt.Println()
	// }
	// fmt.Println("validation errors",validationErrors)

	// validationErrors := handler.ValidateBodyRequest(c, *payload)
	validationErrors := h.ValidateBodyRequest(c, *payload)
	fmt.Println("validation errors", validationErrors)
	if validationErrors != nil {

		return c.JSON(http.StatusBadRequest, validationErrors)
	}
	// 2. validate request body
	// 3. create user and other stuff
	// return c.String(http.StatusOK, "Success")

	return c.JSON(http.StatusCreated, "Registration Successful")
}
