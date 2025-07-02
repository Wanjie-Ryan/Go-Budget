package models

// import "gorm.io/gorm"

type UserModel struct {
	// if you see fields may become null when entering data, assign them as pointers  using the *
	// the gorm.Model in it carries the id, updated at, created at and deleted at
	// gorm.Model
	BaseModel
	Firstname *string `gorm:"type:varchar(200)" json:"firstname"`
	Lastname  *string `gorm:"type:varchar(200)" json:"lastname"`
	Email     string  `gorm:"type:varchar(100); not null; unique" json:"email"`
	Gender    *string `gorm:"type:varchar(50)" json:"gender"`
	Password  string  `gorm:"type:varchar(200); not null" json:"-"`
}

// the - in the password means that go won't serialize it and send it to the user along with the other attributes.

// if you want to change the table name, to a custom name, the function below helps, instead of the table being named UserModel, it will be named users
// func (receiver UserModel) TableName() string{
// 	return "users"
// }
