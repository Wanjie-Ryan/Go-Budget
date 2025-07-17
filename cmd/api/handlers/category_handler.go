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
	_, ok := c.Get("user").(models.UserModel)
	if !ok {
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
	_, ok := c.Get("user").(models.UserModel)
	if !ok {
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

// we need the id of the specific category, and for this, we create the param_request
func (h *Handler) DeleteCategory(c echo.Context) error {
	_, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, "User Authentication Failed")
	}

	// path like /category/:id
	// binding the incoming id param request to the handler
	// var categoryId request.IDParamRequest
	categoryId := new(request.IDParamRequest)
	paramErr := (&echo.DefaultBinder{}).BindPathParams(c, categoryId)

	if paramErr != nil {
		fmt.Println("parameter error", paramErr)
		return common.SendBadRequestResponse(c, "Invalid ID Parameter")
	}
	fmt.Println("category id", categoryId)

	// categoryPayload := new(request.Categoryrequest)
	// if err := (&echo.DefaultBinder{}).BindBody(c, categoryPayload); err !=nil{
	// 	return common.SendBadRequestResponse(c, "Invalid Category Request Body")
	// }

	// validationErr := h.ValidateBodyRequest(c, categoryPayload)
	// if validationErr != nil {
	// 	return common.SendFailedvalidationResponse(c, validationErr)
	// }

	categoryService := services.NewCategoryService(h.DB)

	error := categoryService.DeleteCategory(categoryId.ID)
	if error != nil {
		if error.Error() == "category not found" {
			return common.SendNotFoundResponse(c, "Category Not Found")
		}
		return common.SendServerErrorResponse(c, error.Error())
	}
	return common.SendSuccessResponse(c, "Category Deleted", nil)
}

//GET a single category handler

func (h *Handler) GetSingleCategory(c echo.Context) error {
	_, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, "User Authentication Failed")
	}

	categoryId := new(request.IDParamRequest)

	if err := (&echo.DefaultBinder{}).BindPathParams(c, categoryId); err != nil {
		return common.SendBadRequestResponse(c, "Invalid ID Parameter")
	}

	categoryService := services.NewCategoryService(h.DB)

	category, err := categoryService.GetSingleCategory(categoryId.ID)

	if err != nil {
		if err.Error() == "category not found" {
			return common.SendNotFoundResponse(c, "Category Not Found")
		}
		return common.SendServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Category Found", category)
}

// function to update a category

func (h *Handler) UpdateCategory(c echo.Context) error {

	_, ok := c.Get("user").(models.UserModel)

	if !ok {
		return common.SendUnauthorizedResponse(c, "User Authentication Failed")
	}

	// getting the id from the params and binding it to the param request
	categoryId := new(request.IDParamRequest)

	if err := (&echo.DefaultBinder{}).BindPathParams(c, categoryId); err != nil {
		return common.SendBadRequestResponse(c, "Invalid ID Parameter")
	}

	// getting the category payload to update from the body, and bind to the bind body

	updatePayload := new(request.Categoryrequest)

	if err := (&echo.DefaultBinder{}).BindBody(c, updatePayload); err != nil {
		return common.SendBadRequestResponse(c, "Invalid Category Request Body")
	}

	validationErr := h.ValidateBodyRequest(c, updatePayload)
	if validationErr != nil {
		return common.SendFailedvalidationResponse(c, validationErr)
	}

	categoryService := services.NewCategoryService(h.DB)

	updatedcategory, err := categoryService.UpdateCategory(updatePayload, categoryId.ID)

	if err != nil {
		if err.Error() == "category not found" {
			return common.SendNotFoundResponse(c, "Category Not Found")
		}
		return common.SendServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Category Updated", updatedcategory)

}
