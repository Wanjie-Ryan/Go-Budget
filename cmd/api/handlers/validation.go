package handler

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// how the error will be structured
type ValidationError struct {
	Error     string `json:"error"`
	Key       string `json:"key"`
	Condition string `json:"condition"`
}

// the function below will be returning an array of validationError
// the function will accept the echo context and the payload as arguments
func (h *Handler) ValidateBodyRequest(c echo.Context, payload interface{}) []*ValidationError {

	// create a new validator instance
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())
	var errors []*ValidationError
	// run the validator, will check the validate tags
	err := validate.Struct(payload)

	// attempt to cast the errors to validationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	//loop through the validator errors
	if ok {
		// to get the full json payload that is failing
		reflected := reflect.ValueOf(payload)
		// the _ is the index
		for _, validationErr := range validationErrors {
			fmt.Println(reflected.Type().FieldByName(validationErr.StructField()))
			field, _ := reflected.Type().FieldByName(validationErr.StructField())

			key := field.Tag.Get("json")
			if key == "" {
				key = strings.ToLower(validationErr.StructField())
			}
			condition := validationErr.Tag()

			// building the error message
			errMessage := key + " field is " + condition

			switch condition {
			case "required":
				errMessage = key + " is required"
			case "email":
				errMessage = key + " must be a valid email address"
			}

			// fmt.Println(validationErr.Tag()) // what failed
			// fmt.Println(validationErr.ActualTag()) // json field that failed
			// fmt.Println("failed field",validationErr.Field())
			// //build own error object
			currentValidationError := &ValidationError{
				Error:     errMessage, // email is required
				Key:       key,        // email
				Condition: condition,  //required
			}
			errors = append(errors, currentValidationError)
		}

	}
	return errors

}
