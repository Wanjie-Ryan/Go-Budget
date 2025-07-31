package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	request "github.com/Wanjie-Ryan/Go-Budget/cmd/api/requests"
	"github.com/Wanjie-Ryan/Go-Budget/common"
	"github.com/Wanjie-Ryan/Go-Budget/internal/models"
	"gorm.io/gorm"
)

type BudgetService struct {
	DB *gorm.DB
}

func NewBudgetService(db *gorm.DB) *BudgetService {
	return &BudgetService{DB: db}
}

func (b *BudgetService) CreateBudget(payload *request.CreateBudgetRequest, userId uint) (*models.BudgetModel, error) {

	slug := strings.ToLower(payload.Title)
	slug = strings.Replace(slug, " ", "_", -1)

	budgetModel := &models.BudgetModel{

		Amount:      payload.Amount,
		UserID:      userId,
		Title:       payload.Title,
		Slug:        slug,
		Description: payload.Description,
	}
	if payload.Date == "" {
		currentDate := time.Now()
		budgetModel.Date = currentDate
	}

	// next is to extract the month and the year
	// the purpose for wrapping the month in uint is so that to convert it to the same data type as the model, of which the month was declared as uint
	budgetMonth := uint(budgetModel.Date.Month())
	budgetYear := uint16(budgetModel.Date.Year())

	budgetModel.Month = budgetMonth
	budgetModel.Year = budgetYear

	// next is to check if the budget being created actually exists, with the combination of the slug, month, year and user_id, the combination should be unique

	budgetExist, err := b.budgetExistByID_Slug_Month_Year(budgetModel.UserID, budgetModel.Month, budgetModel.Year, budgetModel.Slug)

	if err != nil {
		// if the record has not been found, then that means that we should create the budget, cause it will not be a duplicate
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result := b.DB.Create(budgetModel)
			if result.Error != nil {
				return nil, result.Error
			}

			return budgetModel, nil

		}
		return nil, err
	}

	// if the mother error is nil, then it means that the budget with the combination of slug, month, year and user_id already exists
	return budgetExist, errors.New("budget with the combination of slug, month, year and user_id already exists")

}

func (b *BudgetService) budgetExistByID_Slug_Month_Year(UserID uint, month uint, year uint16, slug string) (*models.BudgetModel, error) {

	singleBudget := &models.BudgetModel{}

	result := b.DB.Where("user_id = ? AND month = ? AND year = ? AND slug = ?", UserID, month, year, slug).First(&singleBudget)

	// if there is an error, it means the user has somehow created a duplicate budget that has the same combination of slug, month, year and user_id
	if result.Error != nil {
		return nil, result.Error
	}
	return singleBudget, nil
}

// getting all the budgets in the database
func (b *BudgetService) GetAllBudgets(paginator *common.Pagination, budget []*models.BudgetModel) (common.Pagination, error) {
	b.DB.Scopes(paginator.Paginate()).Find(&budget)
	paginator.Items = budget
	return *paginator, nil
}

// function to update a budget

func (b *BudgetService) UpdateBudget(budget *models.BudgetModel, payload *request.UpdateBudgetRequest, budgetId uint) (*models.BudgetModel, error) {

	// var budget models.BudgetModel

	// getBudgetByUser := b.DB.Where("user_id = ?", userId).First(&budget)

	// if getBudgetByUser.Error != nil{
	// 	if errors.Is(getBudgetByUser.Error, gorm.ErrRecordNotFound){
	// 		return nil, errors.New("budget was not found")
	// 	}
	// 	return nil, getBudgetByUser.Error
	// }

	// if the date is not empty
	if payload.Date != "" {

		// time.Parse is the std way to turn a date/time string into time.Time object
		// time.DateOnly is just a handy constnat in the std library defined as the layout string "2006-01-01"
		timeParsed, err := time.Parse(time.DateOnly, payload.Date)
		if err != nil {
			return nil, err
		}
		budget.Date = timeParsed
	}

	// next check if actually amount was passed in the payload
	if payload.Amount > 0 {
		budget.Amount = payload.Amount
	}

	if payload.Description != nil {
		budget.Description = payload.Description
	}

	if payload.Title != "" {
		budget.Title = payload.Title
		slug := strings.ToLower(payload.Title)
		slug = strings.Replace(slug, " ", "_", -1)
		budget.Slug = slug
	}

	// if the count is greater than 0, then it means that there is another budget with the same combination of slug, month, year and user_id
	count := b.CountForYearAndMonthAndSlugAndUserIDExcludeBudgetID(budget.UserID, budget.Month, budget.Year, budget.Slug, budgetId)

	if count > 0 {
		return nil, errors.New("budget with the combination of slug, month, year and user_id already exists")
	}

	b.DB.Model(budget).Updates(budget)
	return budget, nil

}

// function to check if a budget truly exists by its id
func (b *BudgetService) GetBudgetById(budgetId uint) (*models.BudgetModel, error) {
	var budget models.BudgetModel

	// OPTION A
	// result := b.DB.First(&budget, budgetId)

	// OPTION B
	result := b.DB.Where("id = ?", budgetId).First(&budget)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("budget was not found")
		}
		return nil, result.Error
	}

	return &budget, nil
}

func (b *BudgetService) CountForYearAndMonthAndSlugAndUserIDExcludeBudgetID(userId uint, month uint, year uint16, slug string, budgetId uint) int64 {

	fmt.Println("userId is", userId, "month is", month, "year is", year, "slug is", slug, "budgetId is", budgetId)

	var count int64

	// the <> sign means NOT EQUAL TO, one can alternatively use !=
	// the query will count the rows where all these conditions are met
	// useful when you want to exclude a specific record identified by budgetID while counting other records that match other conditions
	// in our case scenario, the title, slug, month, year should be unique, therefore check each field where the id of the budget does not match the budgetId, for those mentioned fields, and ensure that no field matches the same value
	b.DB.Model(models.BudgetModel{}).Where("user_id = ? AND month = ? AND year = ? AND slug = ? AND id <> ?", userId, month, year, slug, budgetId).Count(&count)

	fmt.Println("count is", count)
	return count

}
