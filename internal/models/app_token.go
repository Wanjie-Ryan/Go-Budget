package models

import "time"

type AppTokenModel struct {
	BaseModel
	Type      string    `json:"-" gorm:"index; not null;type:varchar(255)"`
	Used      bool      `json:"-" gorm:"index;not null;type:bool"`
	ExpiresAt time.Time `json:"-" gorm:"index;not null;"`
}
