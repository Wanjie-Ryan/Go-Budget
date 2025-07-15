package services

import (
	"errors"
	"fmt"

	"github.com/Wanjie-Ryan/Go-Budget/internal/models"
	"gorm.io/gorm"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {

	return &CategoryService{db: db}
}

// function to get all categories from the DB
func (cs *CategoryService) GetAllCategories() ([]*models.CategoryModel, error) {

	// create a variable to hold the retrieved categoryModel data
	var categories []*models.CategoryModel

	result := cs.db.Find(&categories)

	if result.Error != nil {
		// return nil, errors.New("failed to fetch categories")
		fmt.Println(result.Error.Error())
		return nil, errors.New(result.Error.Error())
	}
	return categories, nil

}

func (cs *CategoryService) Createcategory() {

}
