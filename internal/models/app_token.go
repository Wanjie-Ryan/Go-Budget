package models

import "time"

// uint data type means the integer cannot be negative, the difference btn INT and UINT is that INT can hold both +ve and -ve values, while UINT can only hold +ve values
type AppTokenModel struct {
	BaseModel
	TargetId  uint      `json:"target_id" gorm:"index;not null"`
	Type      string    `json:"-" gorm:"index; not null;type:varchar(255)"`
	Token     string    `json:"-"gorm: "type:varchar(255)"`
	Used      bool      `json:"-" gorm:"index;not null;type:bool"`
	ExpiresAt time.Time `json:"-" gorm:"index;not null;"`
}
