package request

// defines the shape of the JSON payload one expects from client to send when they register
type RegisterUserRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
