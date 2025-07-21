package handler

import (
	request "github.com/Wanjie-Ryan/Go-Budget/cmd/api/requests"
	"github.com/Wanjie-Ryan/Go-Budget/common"
	"github.com/Wanjie-Ryan/Go-Budget/internal/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateBudget(c echo.Context) error {

	_, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, "User Authentication Failed")
	}

	createBudgetPayload := new(request.CreateBudgetRequest)

	if err := (&echo.DefaultBinder{}).BindBody(c, createBudgetPayload); err != nil {
		return common.SendBadRequestResponse(c, "Invalid Budget Request Body")
	}

	validationErr := h.ValidateBodyRequest(c, createBudgetPayload)
	if validationErr != nil {
		return common.SendFailedvalidationResponse(c, validationErr)
	}

	return nil
}
