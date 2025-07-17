package request

// this path param that holds an id will be bound to the categoryHandler
type IDParamRequest struct {
	ID uint `param:"id" binding:"required"`
}
