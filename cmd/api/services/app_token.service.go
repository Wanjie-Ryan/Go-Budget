package services

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/Wanjie-Ryan/Go-Budget/internal/models"
	"gorm.io/gorm"
)

type AppTokenService struct {
	db *gorm.DB
}

// this is for the user WHO IS NOT AUTHENTICATED

func NewAppTokenService(db *gorm.DB) *AppTokenService {
	return &AppTokenService{db: db}
}

// function to generate a random token, and the functino will return an integer
// it is a method on the AppTokenService struct that holds the db
func (ats *AppTokenService) GenerateToken() int {

	// the method below initializes the random number generator with a unique value based on the current time in nanoseconds
	// without the above line, rand.Int would always generate the same sequence of numbers each time your program runs. Seeding below ensures different results on every run.
	rand.Seed(time.Now().UnixNano())
	// min:=10000
	// max:=99999

	// generates a random integer in the range 0 up to but not including n.
	// generates random 5 digit integer btn 10k and 99,999(inclusive)
	return rand.Intn(99999-10000+1) + 10000
}

func (ats *AppTokenService) GenerateresetPasswordToken(user models.UserModel) (*models.AppTokenModel, error) {
	// for the token we will pass whatever token was generated in the function below
	tokenCreated := models.AppTokenModel{
		TargetId:  user.ID,
		Type:      "reset-password",
		Token:     strconv.Itoa(ats.GenerateToken()), // the method strconv helps us convert an integer to a string
		Used:      false,
		ExpiresAt: time.Now().Add(time.Hour * 1), // the token will expire after 24 hours

	}
	// ats.db.Create returns only one value which is the gorm.db, and in it it has the error, result, everything

	result := ats.db.Create(&tokenCreated)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tokenCreated, nil

}

// validating the token that has been created

func (ats *AppTokenService) ValidateToken(user models.UserModel, token string) (*models.AppTokenModel, error) {

	var retrievedToken models.AppTokenModel
	// we will search the db for where the token is equal to the passed params, and return the first value that is retrieved, and store the value in the retrievedToken variable
	// the result variable holds a *gorm.Db object, the variable holds metadata about the query execution::: result.Error, result.RowsAffected, result.statement
	// the result variable DOES NOT directly hold the data.
	result := ats.db.Where(&models.AppTokenModel{TargetId: user.ID, Type: "reset-password", Token: token}).First(&retrievedToken)

	// the outer if is  a general statement to catch any error, the inner if statement is to catch the specific gorm.recordnotfound error
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("Invalid password reset token")
		}
		return nil, result.Error

	}

	// if used is true, return an error
	if retrievedToken.Used {
		return nil, errors.New("Invalid password reset token")
	}

	if retrievedToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("Password reset token has expired")
	}

	return &retrievedToken, nil

}

// the db.Model explicitly sets the table/model you want to work on
// tells gorm sth like All ops from this point in the chain should use this models table and schema
// the statement tells gorm, target the AppTokenModel table, Apply this where condition, and run an update SQL query directly on that table.
func (ats *AppTokenService) InvalidateToken(user_id int, appToken models.AppTokenModel) {
	ats.db.Model(&models.AppTokenModel{}).Where("target_id=? AND token=?", user_id, appToken.Token).Update("used", true)

}
