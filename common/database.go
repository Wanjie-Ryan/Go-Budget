package common

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger"
)

func NewMySql() (*gorm.DB, error) {

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	host := os.Getenv("DB_HOST")
	db_name := os.Getenv("DB_NAME")
	db_username := os.Getenv("DB_USERNAME")
	db_password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_username, db_password, host, db_name)
	// fmt.Println(dsn)
	//tries to connect to the database
	// gorm.Open returns (*gorm.DB, error)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Default().Println("Connected to database")
	return db, nil

}
