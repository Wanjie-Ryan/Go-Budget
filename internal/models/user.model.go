package models

import "gorm.io/gorm"

type UserModel struct {
	// if you see fields may become null when entering data, assign them as pointers  using the *
	// the gorm.Model in it carries the id, updated at, created at and deleted at
	gorm.Model
	Firstname *string
	Lastname  *string
	Email     string
	Gender    *string
	Password  string
}
