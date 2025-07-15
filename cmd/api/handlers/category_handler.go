package handler

import (
	"github.com/Wanjie-Ryan/Go-Budget/cmd/api/services"
	"github.com/Wanjie-Ryan/Go-Budget/common"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetAllCategories(c echo.Context) error {

	categoryService := services.NewCategoryService(h.DB)

	allcategories, err := categoryService.GetAllCategories()
	if err != nil {
		return common.SendServerErrorResponse(c, err.Error())
	}
	return common.SendSuccessResponse(c, "All Categories", allcategories)

}

func  (h *Handler) Createcategory(c echo.Context)error{
	return nil
}
