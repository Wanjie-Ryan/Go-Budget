package services

import (
	"errors"
	"strings"
	"time"

	request "github.com/Wanjie-Ryan/Go-Budget/cmd/api/requests"
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
