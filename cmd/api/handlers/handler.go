package handler

import (
	"github.com/Wanjie-Ryan/Go-Budget/internal/mailer"
	"gorm.io/gorm"
)

type Handler struct {
	// points to gorm.DB, how you talk to a DB
	// *gorm.DB is what GORM uses to talk to your database—opening connections, building and running SQL queries, managing transactions
	DB *gorm.DB
	Mailer mailer.Mailer
}
