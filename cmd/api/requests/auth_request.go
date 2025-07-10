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
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=4"`
	NewPassword     string `json:"newPassword" validate:"required,min=4"`
	ConfirmPassword string `json:"confirmPassword" validate:"eqfield=NewPassword"`
	// the confirm password field should be equal to the new password
}

type ForgotPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	FrontendUrl string `json:"frontendurl" validate:"required,url"`
}
