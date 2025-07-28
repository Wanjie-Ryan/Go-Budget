package handler

import (
	"fmt"

	request "github.com/Wanjie-Ryan/Go-Budget/cmd/api/requests"
	"github.com/Wanjie-Ryan/Go-Budget/cmd/api/services"
	"github.com/Wanjie-Ryan/Go-Budget/common"
	"github.com/Wanjie-Ryan/Go-Budget/internal/models"
	"github.com/labstack/echo/v4"
)

// function to create a budget
func (h *Handler) CreateBudget(c echo.Context) error {

	// because we are storing the authenticated user in our context, then we have access to the user and all its properties
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, "User Authentication Failed")
	}

	createBudgetPayload := new(request.CreateBudgetRequest)
	fmt.Println("create budget payload", *createBudgetPayload)

	if err := (&echo.DefaultBinder{}).BindBody(c, createBudgetPayload); err != nil {
		fmt.Println("error binding body", err)
		return common.SendBadRequestResponse(c, "Invalid Budget Request Body")
	}

	validationErr := h.ValidateBodyRequest(c, createBudgetPayload)
	if validationErr != nil {
		return common.SendFailedvalidationResponse(c, validationErr)
	}

	budgetService := services.NewBudgetService(h.DB)
	categoryService := services.NewCategoryService(h.DB)

	createdBudget, err := budgetService.CreateBudget(createBudgetPayload, user.ID)

	if err != nil {
		fmt.Println("error creating budget", err)
		return common.SendServerErrorResponse(c, err.Error())
	}

	// associating the categories to the budget

	// categories, err := categoryService.GetMultipleCategories(createBudgetPayload.Categories)
	categories, err := categoryService.GetMultipleCategories(createBudgetPayload)

	if err != nil {
		fmt.Println("error getting categories", err)
		return common.SendServerErrorResponse(c, err.Error())
	}

	err = budgetService.DB.Model(createdBudget).Association("Categories").Replace(categories)

	if err != nil {
		fmt.Println("error associating categories to budget", err)
		return common.SendServerErrorResponse(c, err.Error())
	}

	createdBudget.Categories = categories

	return common.SendSuccessResponse(c, "Budget Created", createdBudget)
}

// function to list budgets

func (h *Handler) GetAllBudgets(c echo.Context) error {

	// session
	user, ok := c.Get("user").(models.UserModel)
	fmt.Println("user", user.ID)
	if !ok {
		return common.SendUnauthorizedResponse(c, "User Authentication Failed")
	}

	// get the model first
	var budgetModel []*models.BudgetModel
	budgetService := services.NewBudgetService(h.DB)
	//the items being retreived from the database will have the categories preloaded
	query := h.DB.Preload("Categories").Scopes(common.WhereUserIDScope(user.ID))
	paginator := common.NewPagination(budgetModel, c.Request(), query)

	paginatedBudget, err := budgetService.GetAllBudgets(paginator, budgetModel)

	if err != nil {
		return common.SendServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "All Budgets", paginatedBudget)

}
