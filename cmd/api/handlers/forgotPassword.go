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

	encodedEmail := base64.StdEncoding.EncodeToString([]byte(retrievedUser.Email))
	frontendUrl, err := url.Parse(payload.FrontendUrl)

	if err != nil {
		return common.SendBadRequestResponse(c, "Invalid Frontend URL")
	}

	// if the url passed is correct, build a query from it

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
