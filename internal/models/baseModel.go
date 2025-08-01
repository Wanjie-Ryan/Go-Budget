package models

import (
	// "database/sql"
	"time"

	"gorm.io/gorm"
)

// type DeletedAt sql.NullTime

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
