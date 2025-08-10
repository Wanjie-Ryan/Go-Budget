package handler

import (
	"fmt"

	request "github.com/Wanjie-Ryan/Go-Budget/cmd/api/requests"
	"github.com/Wanjie-Ryan/Go-Budget/cmd/api/services"
	"github.com/Wanjie-Ryan/Go-Budget/common"
	"github.com/Wanjie-Ryan/Go-Budget/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// function to create a budget
func (h *Handler) CreateBudget(c echo.Context) error {

	// because we are storing the authenticated user in our context, then we have access to the user and all its properties
	user, ok := c.Get("user").(models.UserModel)
	if !ok {
		return common.SendUnauthorizedResponse(c, "User Authentication Failed")
	}

	createBudgetPayload := new(request.CreateBudgetRequest)
	// fmt.Println("create budget payload", *createBudgetPayload)

	if err := (&echo.DefaultBinder{}).BindBody(c, createBudgetPayload); err != nil {
		// fmt.Println("error binding body", err)
		return common.SendBadRequestResponse(c, "Invalid Budget Request Body")
	}

	validationErr := h.ValidateBodyRequest(c, createBudgetPayload)
	if validationErr != nil {
		return common.SendFailedvalidationResponse(c, validationErr)
	}

	budgetService := services.NewBudgetService(h.DB)
	categoryService := services.NewCategoryService(h.DB)
	categories, err := categoryService.GetMultipleCategories(createBudgetPayload)

	if err != nil {
		fmt.Println("error getting categories", err)
		return common.SendServerErrorResponse(c, err.Error())
	}

	createdBudget := &models.BudgetModel{}

	// START OF PERFORMING TRANSACTIONS
	err = h.DB.Transaction(func(tx *gorm.DB) error {

		// to do transactions in here, you have to utilize the tx, not the db you instantiated, therefore re-assign the DB to the tx

		budgetService.DB = tx
		categoryService.DB = tx

		createdBudget, err = budgetService.CreateBudget(createBudgetPayload, user.ID)

		if err != nil {
			fmt.Println("error creating budget", err)
			return common.SendServerErrorResponse(c, err.Error())
		}

		err = tx.Model(createdBudget).Association("Categories").Replace(categories)

		if err != nil {
			fmt.Println("error associating categories to budget", err)
			return common.SendServerErrorResponse(c, err.Error())
		}

		// fmt.Println("transaction", tx)
		return nil
	})

	// END OF PERFORMING TRANSACTIONS

	// ERROR FOR CAPTURING THE ERROR THAT HAPPENS DURING THE TRANSACTION
	if err != nil {
		return common.SendServerErrorResponse(c, err.Error())
	}
	// ERROR FOR CAPTURING THE ERROR THAT HAPPENS DURING THE TRANSACTION

	// associating the categories to the budget

	// categories, err := categoryService.GetMultipleCategories(createBudgetPayload.Categories)

	createdBudget.Categories = categories

	return common.SendSuccessResponse(c, "Budget Created", createdBudget)
}

// function to list budgets

func (h *Handler) GetAllBudgets(c echo.Context) error {

	// session
	user, ok := c.Get("user").(models.UserModel)
	// fmt.Println("user", user.ID)
	if !ok {
		return common.SendUnauthorizedResponse(c, "User Authentication Failed")
	}

	// get the model first
	var budgetModel []*models.BudgetModel
	budgetService := services.NewBudgetService(h.DB)
	//the items being retreived from the database will have the categories preloaded
	// query := h.DB.Preload("Categories").Scopes(common.WhereUserIDScope(user.ID))
	query := h.DB.Scopes(common.WhereUserIDScope(user.ID))
	paginator := common.NewPagination(budgetModel, c.Request(), query)

	paginatedBudget, err := budgetService.GetAllBudgets(paginator, budgetModel)

	if err != nil {
		return common.SendServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "All Budgets", paginatedBudget)

}

// function to update budget
func (h *Handler) UpdateBudget(c echo.Context) error {

	User, ok := c.Get("user").(models.UserModel)

	if !ok {
		return common.SendUnauthorizedResponse(c, "User Authentication Failed")
	}
	// bind the budgetID from the parameter first

	budgetId := new(request.IDParamRequest)

	if err := (&echo.DefaultBinder{}).BindPathParams(c, budgetId); err != nil {
		return common.SendBadRequestResponse(c, "Invalid ID Parameter")
	}

	// before doing the validation, we need to confirm, check if the budget exists or not, using the id coming from the params

	budgetService := services.NewBudgetService(h.DB)
	// categoryService := services.NewCategoryService(h.DB)

	budgetByid, err := budgetService.GetBudgetById(budgetId.ID)

	if err != nil {
		if err.Error() == "budget was not found" {
			return common.SendNotFoundResponse(c, "Budget Not Found")
		}
		return common.SendBadRequestResponse(c, err.Error())
	}

	if User.ID != budgetByid.UserID {
		// return common.SendUnauthorizedResponse(c, "Cannot perform this action")
		return common.SendBadRequestResponse(c, "Cannot perform this action")
	}

	updateBudgetPayload := new(request.UpdateBudgetRequest)

	if err := (&echo.DefaultBinder{}).BindBody(c, updateBudgetPayload); err != nil {
		return common.SendBadRequestResponse(c, "Invalid Budget Request Body")
	}

	validationErr := h.ValidateBodyRequest(c, updateBudgetPayload)
	if validationErr != nil {
		return common.SendFailedvalidationResponse(c, validationErr)
	}

	// by passing the budgetById, we actually pass the budget model retrieved to utilize the preloaded budgets

	budgetUpdate, err := budgetService.UpdateBudget(budgetByid, updateBudgetPayload, budgetId.ID)

	if err != nil {
		if err.Error() == "budget with the combination of slug, month, year and user_id already exists" {
			return common.SendBadRequestResponse(c, "Budget with the combination of slug, month, year and user_id already exists")
		}

		return common.SendServerErrorResponse(c, err.Error())

	}

	// categories, err := categoryService.GetMultipleCategories(updateBudgetPayload)

	return common.SendSuccessResponse(c, "Budget Updated", budgetUpdate)

}

// function to delete a budget, in the model, the category had a ondelete cascade, meaning that when the budget is deleted, the associated categories too are deleted
func (h *Handler) DeleteBudget(c echo.Context) error {

	User, ok := c.Get("user").(models.UserModel)

	if !ok {
		return common.SendUnauthorizedResponse(c, "User Authentication Failed")
	}

	budgetId := new(request.IDParamRequest)

	if err := (&echo.DefaultBinder{}).BindPathParams(c, budgetId); err != nil {
		return common.SendBadRequestResponse(c, "Invalid ID Parameter")
	}

	budgetService := services.NewBudgetService(h.DB)

	singleBudget, err := budgetService.GetBudgetById(budgetId.ID)
	if err != nil {
		if err.Error() == "budget was not found" {
			return common.SendNotFoundResponse(c, err.Error())
		}
		return common.SendServerErrorResponse(c, err.Error())
	}

	if User.ID != singleBudget.UserID {
		return common.SendBadRequestResponse(c, "Cannot perform this action")
	}

	err = budgetService.DeleteBudget(singleBudget)

	if err != nil {
		return common.SendServerErrorResponse(c, err.Error())
	}

	return common.SendSuccessResponse(c, "Budget Deleted", nil)

}
