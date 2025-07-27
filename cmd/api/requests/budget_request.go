package request

// request to store and update a budget
type CreateBudgetRequest struct {

	// the dive attribute goes through the array (each of the item) in categories and ensures that all of them are of the type uint64
	// uint --> unsigned integer, cannot be NEGATIVE
	Categories []uint64 `json:"categories" validate:"required,dive,min=1"`
	Amount     float64  `json:"amount" validate:"required,numeric,min=1"`
	//omitempty --> its ok if the user doesn't provide this
	// the format date is type YYYY-MM-DD
	Date        string  `json:"date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Title       string  `json:"title" validate:"required,min=2,max=200"`
	Description *string `json:"description" validate:"omitempty,min=2,max=1000"`
}

type UpdateBudgetRequest struct {
	Categories []uint64 `json:"categories" validate:"omitempty,dive,min=1"`
	Amount     float64  `json:"amount" validate:"omitempty,numeric,min=1"`
	Date       string   `json:"date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Title      string   `json:"title" validate:"omitempty,min=2,max=200"`
}

// the omitempty being in the json tag means if the field is empty, don't include it in the json output
//in the validate field, the omitempty means that if the field is empty, don't validate
