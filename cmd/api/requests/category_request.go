package request

type Categoryrequest struct {
	Name string `json:"name" validate:"required"`
}
