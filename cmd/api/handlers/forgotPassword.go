package handler

import (
	"errors"
	"fmt"

	request "github.com/Wanjie-Ryan/Go-Budget/cmd/api/requests"
	"github.com/Wanjie-Ryan/Go-Budget/cmd/api/services"
	"github.com/Wanjie-Ryan/Go-Budget/common"
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
	return nil

}
