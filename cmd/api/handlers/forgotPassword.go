package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"

	request "github.com/Wanjie-Ryan/Go-Budget/cmd/api/requests"
	"github.com/Wanjie-Ryan/Go-Budget/cmd/api/services"
	"github.com/Wanjie-Ryan/Go-Budget/common"
	"github.com/Wanjie-Ryan/Go-Budget/internal/mailer"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (h *Handler) ForgotPassword(c echo.Context) error {

	// get the payload
	payload := new(request.ForgotPasswordRequest)

	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		return common.SendBadRequestResponse(c, "Invalid Request Body")
	}

	validationErrors := h.ValidateBodyRequest(c, *payload)
	if validationErrors != nil {
		return common.SendFailedvalidationResponse(c, validationErrors)
	}

	userService := services.NewUserservice(h.DB)
	appTokenservice := services.NewAppTokenService(h.DB)

	// try n get the email of the user from the payload being passed
	retrievedUser, err := userService.GetUserByEmail(payload.Email)

	// if error is not equal to nil, it will look first for the gor record not found error , if its not found, then it will throw the internal server error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return common.SendNotFoundResponse(c, "Email Does not Exist, Please Register")
		}
		return common.SendServerErrorResponse(c, "An Unexpected error occurred")
	}
	fmt.Println("retrieved user", retrievedUser)

	// call the service function to create for creating a token
	token, err := appTokenservice.GenerateresetPasswordToken(*retrievedUser)
	if err != nil {
		return common.SendServerErrorResponse(c, "An Error occurred, please try again later")
	}

	// if valid, send an email to the user with the token.
	//1 . Generate an encoded mail, which will be attached to the url
	// users email is included in the reset url as a query parameter

	// encodedEmail := base64.StdEncoding.EncodeToString([]byte(retrievedUser.Email))
	encodedEmail := base64.RawURLEncoding.EncodeToString([]byte(retrievedUser.Email))
	frontendUrl, err := url.Parse(payload.FrontendUrl)

	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid Frontend URL")
	}

	// if the url passed is correct, build a query from it

	// create a query object that will hold parameters
	// final url will look sth like
	// https://myapp.com/reset-password?email=dGVzdEBleGFtcGxlLmNvbQ==&token=ABCDEF123456

	query := url.Values{}
	query.Set("email", encodedEmail)
	query.Set("token", token.Token)
	frontendUrl.RawQuery = query.Encode()

	mailData := mailer.EmailData{
		Subject: "Reset Password For Budget Tracker",
		Meta: struct {
			Token       string
			FrontendUrl string
		}{
			Token:       token.Token,
			FrontendUrl: frontendUrl.String(),
		},
	}

	err = h.Mailer.Send(payload.Email, "forgot-password.html", mailData)
	if err != nil {
		fmt.Println("mail error", err)
		// return err
	}

	return common.SendSuccessResponse(c, "Forgot Password Email Sent Successfully", nil)

}

func (h *Handler) ResetPasswordHandler(c echo.Context) error {

	resetPayload := new(request.ResetPasswordRequest)
	if err := (&echo.DefaultBinder{}).BindBody(c, resetPayload); err != nil {
		return common.SendBadRequestResponse(c, "Invalid Request Body")
	}

	validationErr := h.ValidateBodyRequest(c, *resetPayload)
	if validationErr != nil {
		return common.SendFailedvalidationResponse(c, validationErr)
	}

	// now, we will get the email of the user from the meta which is being passed as payload too when the user wants to reset password.
	// and from that email got from the meta payload, decode it, and use it to get a specific user.

	// email, err := base64.StdEncoding.DecodeString(resetPayload.Meta)
	email, err := base64.RawURLEncoding.DecodeString(resetPayload.Meta)
	if err != nil {
		fmt.Println("decode error", err)
		return common.SendServerErrorResponse(c, "An Error occurred, please try again later")
	}
	// email being returned is going to be of data type byte, not string
	fmt.Println("decoded email", string(email))

	userService := services.NewUserservice(h.DB)
	// convert the email to string coz, its being returned as byte data type after being decoded
	userByMail, err := userService.GetUserByEmail(string(email))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return common.SendNotFoundResponse(c, "Email Does not Exist, Please Register")
		}

		return common.SendServerErrorResponse(c, "An Error occurred, please try again later")
	}

	appTokenService := services.NewAppTokenService(h.DB)

	// after getting the user by mail, we then move forward and check if the token is valid

	token, err := appTokenService.ValidateToken(*userByMail, resetPayload.Token)

	if err != nil {
		return common.SendNotFoundResponse(c, err.Error())
	}

	// if the token is valid and has been found in the DB, then the next step is changing the user password
	// the changePassword function only returns an error
	err = userService.ChangePassword(resetPayload.Password, *userByMail)
	if err != nil {
		return common.SendServerErrorResponse(c, err.Error())
	}

	// after the password was changed succesfully, we then invalidate the token

	appTokenService.InvalidateToken(int(userByMail.ID), *token)

	// after all that is successful, now return a message to the user informing the user of reseting password being successful
	return common.SendSuccessResponse(c, "Password Reset Successful", nil)

}
