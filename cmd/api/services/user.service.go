package services

import (
	"fmt"

	request "github.com/Wanjie-Ryan/Go-Budget/cmd/api/requests"
	"github.com/Wanjie-Ryan/Go-Budget/internal/models"
	"gorm.io/gorm"
)

type Userservice struct {
	db *gorm.DB
}

func NewUserservice(db *gorm.DB) *Userservice {
	return &Userservice{db: db}
}

func (u Userservice) RegisterUser(user request.RegisterUserRequest) (*models.UserModel, error) {
	fmt.Println("do me")
	return nil, nil
}

// the function below expects you to return a model of user, and an error
func (u Userservice) GetUserByEmail(email string) (*models.UserModel, error) {
	var user models.UserModel
	// will try and get the user by the provided email
	result := u.db.Where("email = ?", email).First(&user)

	// if result.RowsAffected == 0{

	// }
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
