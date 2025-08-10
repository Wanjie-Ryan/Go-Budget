package services

import (
	"errors"
	"fmt"
	"strings"

	request "github.com/Wanjie-Ryan/Go-Budget/cmd/api/requests"
	"github.com/Wanjie-Ryan/Go-Budget/common"
	"github.com/Wanjie-Ryan/Go-Budget/internal/models"
	"gorm.io/gorm"
)

type CategoryService struct {
	DB *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {

	return &CategoryService{DB: db}
}

// function to get all categories from the DB
func (cs *CategoryService) GetAllCategories(paginator *common.Pagination, categories []*models.CategoryModel) (common.Pagination, error) {

	// create a variable to hold the retrieved categoryModel data
	// var categories []*models.CategoryModel

	// result := cs.DB.Find(&categories)

	// this is simply calling the method to fill the categories slice and move on
	cs.DB.Scopes(paginator.Paginate()).Find(&categories)

	// if result.Error != nil {
	// 	// return nil, errors.New("failed to fetch categories")
	// 	fmt.Println(result.Error.Error())
	// 	return nil, errors.New(result.Error.Error())
	// }

	// in this case the catgeories was declared and in the handler and passed down
	paginator.Items = categories

	// return categories, nil
	return *paginator, nil

}

// the function below will return an instance of the categoryModel and an error, and will accept the payload as an argument
func (cs *CategoryService) Createcategory(categoryPayload *request.Categoryrequest) (*models.CategoryModel, error) {

	//create a slug from the category name being passed by user
	// 1. transform the name to lowercase using the go string methods
	slug := strings.ToLower(categoryPayload.Name)
	// 2. replace all spaces with underscore till the end of the string
	slug = strings.Replace(slug, " ", "_", -1)

	categoryModelCreated := models.CategoryModel{
		Name:     categoryPayload.Name,
		Slug:     slug,
		IsCustom: categoryPayload.IsCustom,
	}

	// result := cs.DB.Create(&categoryModelCreated)

	// because the slug and the name of the categories are supposed to be unique, then we should first check if they exist in the DB, if they DO NOT exist, create them, if they exist, it loads into object, therefore we will not use the .create method directly, we will find the slug in the DB that matches the slug that was created
	//.firstorcreate, either gets the first slug that matches the passed slug or creates it
	// the query below will be checking where the slug and the name are equal to the slug and the name that was created

	result := cs.DB.Where(models.CategoryModel{Slug: slug, Name: categoryModelCreated.Name}).FirstOrCreate(&categoryModelCreated)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		// the duplicate error key lets you get the specific error message in case there is a duplication
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, errors.New("category already exists")
		}
		// return nil, errors.New(result.Error.Error())
		return nil, errors.New("failed to create category")
	}
	return &categoryModelCreated, nil

}

// function to delete a category
func (cs *CategoryService) DeleteCategory(id uint) error {

	// result := cs.DB.Where(models.CategoryModel{Name:category.Name}).Delete(&models.CategoryModel{})
	// if result.Error !=nil{
	// 	return result.Error
	// }

	// DELETING FROM ID

	var category models.CategoryModel

	singleCategoryResult := cs.DB.First(&category, id)

	if singleCategoryResult.Error != nil {
		if errors.Is(singleCategoryResult.Error, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return singleCategoryResult.Error
	}

	// without the unscoped, it will only do a soft delete, but with the unscoped, it will permanently delete from the DB
	result := cs.DB.Unscoped().Delete(&category)
	fmt.Println("result after deleting", *result)
	return result.Error

}

// function to get a single category
func (cs *CategoryService) GetSingleCategory(d uint) (*models.CategoryModel, error) {
	var category models.CategoryModel

	result := cs.DB.First(&category, d)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, result.Error
	}

	return &category, nil
}

// function to update a category
// func (cs *CategoryService) UpdateCategory(categoryPayload *request.Categoryrequest, id uint)(*models.CategoryModel, error){

// 	_, err :=cs.GetSingleCategory(id)

// 	if err != nil{
// 		return nil, err
// 	}

// 	slug := strings.ToLower(categoryPayload.Name)
// 	slug = strings.Replace(slug, " ", "_", -1)

// 	// retrievedCategory.Name = categoryPayload.Name
// 	// retrievedCategory.Slug = slug
// 	// retrievedCategory.IsCustom = categoryPayload.IsCustom

// 	// result := cs.DB.Save(retrievedCategory)

// 	modelToUpdate := models.CategoryModel{
// 		Name:     categoryPayload.Name,
// 		Slug:     slug,
// 		IsCustom: categoryPayload.IsCustom,
// 	}

// 	result := cs.DB.Model(&models.CategoryModel{}).Where("id = ?", id).Updates(modelToUpdate)

// 	if result.Error != nil {
// 		fmt.Println(result.Error.Error())
// 		return nil, result.Error
// 	}
// 	return &modelToUpdate, nil
// }

// New function to update category

func (cs *CategoryService) UpdateCategory(categoryPayload *request.Categoryrequest, id uint) (*models.CategoryModel, error) {
	// Fetch the existing category by ID
	var category models.CategoryModel
	result := cs.DB.First(&category, id)
	if result.Error != nil {
		return nil, result.Error // Return the error if not found
	}

	// Now update the fields from the payload
	category.Name = categoryPayload.Name
	category.Slug = strings.ToLower(categoryPayload.Name)
	category.Slug = strings.Replace(category.Slug, " ", "_", -1)
	category.IsCustom = categoryPayload.IsCustom // Update the 'IsCustom' field

	// Save the updated category back to the database
	updateResult := cs.DB.Save(&category)
	if updateResult.Error != nil {
		return nil, updateResult.Error // Return error if save fails
	}

	return &category, nil
}

func (cs *CategoryService) GetMultipleCategories(loadedCategories *request.CreateBudgetRequest) ([]*models.CategoryModel, error) {

	var categories []*models.CategoryModel

	// the where method of GORM filters the categories by the IDs provided in the loadedcategories.Categories. The IN query checks if the the id is one of the specified IDs in the array
	result := cs.DB.Where("id IN ? ", loadedCategories.Categories).Find(&categories)

	if result.Error != nil {
		return nil, result.Error
	}
	if len(categories) == 0 {
		return nil, errors.New("no categories found")
	}
	return categories, nil

}
