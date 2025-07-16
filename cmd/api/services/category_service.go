package services

import (
	"errors"
	"fmt"
	"strings"

	request "github.com/Wanjie-Ryan/Go-Budget/cmd/api/requests"
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

// the function below will return an instance of the categoryModel and an error, and will accept the payload as an argument
func (cs *CategoryService) Createcategory(categoryPayload *request.Categoryrequest) (*models.CategoryModel, error) {

	//create a slug from the category name being passed by user
	// 1. transform the name to lowercase using the go string methods
	slug := strings.ToLower(categoryPayload.Name)
	// 2. replace all spaces with underscore till the end of the string
	slug = strings.Replace(slug, " ", "_", -1)

	categoryModelCreated := models.CategoryModel{
		Name: categoryPayload.Name,
		Slug: slug,
		IsCustom: categoryPayload.IsCustom,
	}

	// result := cs.db.Create(&categoryModelCreated)

	// because the slug and the name of the categories are supposed to be unique, then we should first check if they exist in the DB, if they DO NOT exist, create them, if they exist, it loads into object, therefore we will not use the .create method directly, we will find the slug in the DB that matches the slug that was created
	//.firstorcreate, either gets the first slug that matches the passed slug or creates it

	result := cs.db.Where(models.CategoryModel{Slug: slug}).FirstOrCreate(&categoryModelCreated)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		// return nil, errors.New(result.Error.Error())
		return nil, errors.New("failed to create category")
	}
	return &categoryModelCreated, nil

 }
  