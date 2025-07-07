package handler

import (
	"errors"
	"fmt"

	// "net/http"

	// "net/http"

	// handler "github.com/Wanjie-Ryan/Go-Budget/cmd/api/handlers"
	request "github.com/Wanjie-Ryan/Go-Budget/cmd/api/requests"
	"github.com/Wanjie-Ryan/Go-Budget/cmd/api/services"
	"github.com/Wanjie-Ryan/Go-Budget/common"
	"github.com/Wanjie-Ryan/Go-Budget/internal/mailer"
	"gorm.io/gorm"

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
		// return c.JSON(http.StatusBadRequest, "Invalid request body")
		return common.SendBadRequestResponse(c, "Invalid Request Body")
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

		// return c.JSON(http.StatusBadRequest, validationErrors)
		return common.SendFailedvalidationResponse(c, validationErrors)
	}

	userService := services.NewUserservice(h.DB)
	userExist, err := userService.GetUserByEmail(payload.Email)
	fmt.Println("does user exist by email", userExist, err)
	// 2. validate request body
	// 3. create user and other stuff
	// return c.String(http.StatusOK, "Success")

	// how to check the type of an error, specifically, in our case, email not found

	// the check below says that (ORIGINALLY)  if the email is NOT found, then negate that to false, and throw a common error saying email already exists
	// if no match is found, gorm sets result.Error to gorm.ErrRecordNotFound

	if errors.Is(err, gorm.ErrRecordNotFound) == false {

		return common.SendBadRequestResponse(c, "Email already exists")
	}

	registeredUser, err := userService.RegisterUser(*payload)
	if err != nil {
		// return common.SendErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return common.SendServerErrorResponse(c, err.Error())
	}

	mailData := mailer.EmailData{
		Subject: "Welcome to Budget Tracker",
		Meta: struct {
			FirstName string
			LoginLink string
		}{
			FirstName: *registeredUser.Firstname,
			LoginLink: "#",
		},
	}

	err = h.Mailer.Send(payload.Email, "welcome.html", mailData)
	if err != nil {
		fmt.Println("mail error", err)
		// return err
	}

	// return c.JSON(http.StatusCreated, "Registration Successful")
	return common.SendSuccessResponse(c, "Registration Successful", registeredUser)
}

// CREATING THE LOGIN  HANDLER
func (h *Handler) Loginhandler(c echo.Context) error{
	// bind data
	loginPayload := new(request.LoginUserRequest)

	if err :=(&echo.DefaultBinder{}).BindBody(c, loginPayload);
	err !=nil{
		fmt.Println("login error", err)
		return common.SendBadRequestResponse(c, "Invalid Request Body")
	}
	fmt.Println("login payload", loginPayload)
	// validate data
	loginValidationErrors := h.ValidateBodyRequest(c, *loginPayload)
	fmt.Println("login validation errors", loginValidationErrors)
	if loginValidationErrors !=nil{
		return common.SendFailedvalidationResponse(c, loginValidationErrors)
	}

	userService := services.NewUserservice(h.DB)


	// if the user with supplied mail exist

	user, err := userService.GetUserByEmail(loginPayload.Email)

	// errors.Is is a way to check if a certain error matches a specific known error type, even if its wrapped inside other errors.
	// when you use GORM and you try to query for sth in the DB, if no record is found, GORM returns an error called gorm.ErrRecordNotFound
	// when the condition below is set to TRUE, it means that the email WAS NOT FOUND, otherwise when FALSE, it means user with that email ALREADY EXISTS (Gorm did not return 'record not found')
	if errors.Is(err, gorm.ErrRecordNotFound) {

		return common.SendBadRequestResponse(c, "Invalid Email or Password")
	}

	fmt.Println("retrieved user", *user)

	// compare the passwords
	// if the comparison of the passwords do not match (false) return an error
	if common.CheckPasswordHash(loginPayload.Password, user.Password) == false{
		return common.SendBadRequestResponse(c, "Invalid Email or Password")
	}

	// generate access token



	// return response with user token
	return common.SendSuccessResponse(c, "Login Successful", user)
}