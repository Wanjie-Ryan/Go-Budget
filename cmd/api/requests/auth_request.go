package request

// this acts like some DTO
// defines the shape of the JSON payload one expects from client to send when they register
type RegisterUserRequest struct {
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=4"`
}

// minimum length of password is 4
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required,min=4"`
}
