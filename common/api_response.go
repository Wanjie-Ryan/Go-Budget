package common

import (
	"net/http"

	// handler "github.com/Wanjie-Ryan/Go-Budget/cmd/api/handlers"
	// handler "github.com/Wanjie-Ryan/Go-Budget/cmd/api/handlers"
	"github.com/labstack/echo/v4"
)

// key will be string, but value can be anything
type ApiResponse map[string]any

// create a success struct
type JSONSuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type JSONFailedValidationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	// Errors  []*handler.ValidationError `json:"errors"`
	Errors []*ValidationError `json:"errors"`
}

type JSONErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	// Errors  handler.ValidationError `json:"errors"`
}

// the data being passed here will be of type interface, can be anything, and this function expects to be returned to an error
func SendSuccessResponse(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, JSONSuccessResponse{Success: true, Message: message, Data: data})
}

// for the function below it will accept errors as an argument, and the errors will be of type ValidationError which will enter as an array
func SendFailedvalidationResponse(c echo.Context, errors []*ValidationError) error {
	return c.JSON(http.StatusUnprocessableEntity, JSONFailedValidationResponse{Success: false, Message: "Validation Failed", Errors: errors})

}

func SendErrorResponse(c echo.Context, message string, statusCode int) error {
	return c.JSON(statusCode, JSONErrorResponse{Success: false, Message: message})

}

func SendBadRequestResponse(c echo.Context, message string) error {
	// return c.JSON(http.StatusBadRequest, JSONErrorResponse{Success: false, Message: message})
	return SendErrorResponse(c, message, http.StatusBadRequest)
}

func SendNotFoundResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, message, http.StatusNotFound)
}
