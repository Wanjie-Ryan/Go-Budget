package common

import "gorm.io/gorm"

// a scope function in GORM is a reusable query snippet - a function that returns a closure (another function) which modifies the GORM query builder
func WhereUserIDScope(UserId uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", UserId)
	}
}

// the function returns a function that, when applied to a GORM query, adds a WHERE user_id = ? caluse with the give userId
// instead of doing this everywhere
// db.Where("user_id = ?", userID).Find(&budgets)
// you can write
// db.Scopes(common.WhereUserIDScope(userID)).Find(&budgets)
// You give the value (userId) — GORM gives the connection (*gorm.DB) — the scope glues them together into a filtered query.
