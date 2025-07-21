package models

import "time"

type BudgetModel struct {
	BaseModel
	// title will be created as a default index by gorm
	Title       string  `gorm:"index;type:varchar(200);not null" json:"title"`
	Slug        string  `gorm:"index;uniqueIndex:unique_user_id_slug_year_month;type:varchar(200);not null" json:"slug"`
	Description *string `gorm:"type:text" json:"description"`
	UserID      uint    `gorm:"not null;uniqueIndex:unique_user_id_slug_year_month;column:user_id" json:"user_id"`
	// the amount field, the decimal means it will hold values with values upto 10  digits but with only 2 DPs max
	Amount float64 `gorm:"type:decimal(10,2);not null" json:"amount"`
	// will create a join table called budget_category
	// You can still Preload("Categories") to get categories when querying budgets.
	Categories []CategoryModel `gorm:"constraint:OnDelete:CASCADE;many2many:budget_category" json:"categories"`
	Date       time.Time       `gorm:"type:datetime;not null" json:"date"`
	Month      uint            `gorm:"type:TINYINT UNSIGNED;not null;index:idx_month_year;uniqueIndex:unique_user_id_slug_year_month" json:"month"`
	Year       uint16          `gorm:"type:TINYINT UNSIGNED;not null;index:idx_month_year;uniqueIndex:unique_user_id_slug_year_month" json:"year"`
}

// we need to ensure that the combination of slug, year, month and user_id is unique,
//  hence create a unique undex with a combination of the 3 fields, there will be no field of the 3 mentioned that has a duplicate value, all values in them will be unique

func (BudgetModel) TableName() string {
	return "budget"
}
