package main

import (
	"fmt"

	"github.com/Wanjie-Ryan/Go-Budget/common"
	"github.com/Wanjie-Ryan/Go-Budget/internal/models"
)

func main() {

	db, err := common.NewMySql()
	if err != nil {
		// panic stops normal execution and unwinds the stack, prints the error and stack trace to the console
		panic(err)
	}
	err = db.AutoMigrate(&models.UserModel{}, &models.AppTokenModel{}, &models.CategoryModel{}, &models.BudgetModel{})
	// err = db.Migrator().AlterColumn(&models.UserModel{}, "Firstname")
	// err = db.Migrator().AlterColumn(&models.UserModel{}, "Lastname")
	// err = db.Migrator().AlterColumn(&models.UserModel{}, "Email")
	// err = db.Migrator().AlterColumn(&models.UserModel{}, "Gender")
	// err = db.Migrator().AlterColumn(&models.UserModel{}, "password")
	if err != nil {
		panic(err)
	}
	fmt.Println("Database Migrated")
}
