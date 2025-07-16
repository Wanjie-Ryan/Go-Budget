package handler

import (
	"fmt"

	request "github.com/Wanjie-Ryan/Go-Budget/cmd/api/requests"
	"github.com/Wanjie-Ryan/Go-Budget/cmd/api/services"
	"github.com/Wanjie-Ryan/Go-Budget/common"
	"github.com/Wanjie-Ryan/Go-Budget/internal/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetAllCategories(c echo.Context) error {
	_, ok:=c.Get("user").(models.UserModel)
	if !ok{
		return common.SendUnauthorizedResponse(c, "User Authentication Failed")
	}

	categoryService := services.NewCategoryService(h.DB)

	allcategories, err := categoryService.GetAllCategories()
	if err != nil {
		return common.SendServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "All Categories", allcategories)

}

func (h *Handler) Createcategory(c echo.Context) error {

	// the code below helps to get the current authenticated user via the context
	// if the user context is not being utilized put an underscore
	_, ok :=c.Get("user").(models.UserModel)
	if !ok{
		return common.SendUnauthorizedResponse(c, "User Authentication Failed")
	}

	categoryPayload := new(request.Categoryrequest)
	if err := (&echo.DefaultBinder{}).BindBody(c, categoryPayload); err != nil {
		return common.SendBadRequestResponse(c, "Invalid Category Request Body")
	}
	fmt.Println("category payload", categoryPayload)

	validationErr := h.ValidateBodyRequest(c, categoryPayload)

	if validationErr != nil {
		return common.SendFailedvalidationResponse(c, validationErr)
	}

	categoryService := services.NewCategoryService(h.DB)

	result, err := categoryService.Createcategory(categoryPayload)
	if err != nil {
		return common.SendServerErrorResponse(c, err.Error())
	} 

	return common.SendSuccessResponse(c, "Category Created", result)
}
