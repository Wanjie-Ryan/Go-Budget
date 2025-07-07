package handler

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Wanjie-Ryan/Go-Budget/common"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// how the error will be structured
// type ValidationError struct {
// 	Error     string `json:"error"`     //human readable error, eg. email is required
// 	Key       string `json:"key"`       // json key that failed
// 	Condition string `json:"condition"` // which rule failed, required or email
// }

// the function below will be returning an array of validationError
// the function will accept the echo context and the payload as arguments
func (h *Handler) ValidateBodyRequest(c echo.Context, payload interface{}) []*common.ValidationError {

	// create a new validator instance
	var validate *validator.Validate
	// enables the required tag on structs themselves, not just on fields
	validate = validator.New(validator.WithRequiredStructEnabled())
	var errors []*common.ValidationError
	// run the validator, will check the validate tags
	err := validate.Struct(payload)

	// attempt to cast the errors to validationErrors
	// if the error really came from validation, you get a slice of field level errors, otherwise ok is false
	validationErrors, ok := err.(validator.ValidationErrors)
	//loop through the validator errors
	if ok {
		// to get the full json payload that is failing, enables you to read email instead of Email
		reflected := reflect.ValueOf(payload)
		// the _ is the index
		for _, validationErr := range validationErrors {
			fmt.Println(reflected.Type().FieldByName(validationErr.StructField()))

			// find the struct field by the Go field name
			field, _ := reflected.Type().FieldByName(validationErr.StructField())

			// extract the json tag
			key := field.Tag.Get("json")
			if key == "" {
				key = strings.ToLower(validationErr.StructField())
			}
			// which validation rule failed
			condition := validationErr.Tag()

			// extracting the param from the error
			param := validationErr.Param()
			// building the error message
			errMessage := key + " field is " + condition

			// switch the condition, if condition is required, show this message, if email show this message, etc etc

			switch condition {
			case "required":
				errMessage = key + " is required"
			case "email":
				errMessage = key + " must be a valid email address"
			case "min":
				errMessage = key + " must be at least " + param + " characters long"
			}

			// fmt.Println(validationErr.Tag()) // what failed
			// fmt.Println(validationErr.ActualTag()) // json field that failed
			// fmt.Println("failed field",validationErr.Field())
			// //build own error object
			currentValidationError := &common.ValidationError{
				Error:     errMessage, // email is required
				Key:       key,        // email
				Condition: condition,  //required
			}
			errors = append(errors, currentValidationError)
		}

	}
	return errors

}
