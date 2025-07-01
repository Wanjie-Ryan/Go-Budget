package common

type ValidationError struct {
	Error     string `json:"error"`     //human readable error, eg. email is required
	Key       string `json:"key"`       // json key that failed
	Condition string `json:"condition"` // which rule failed, required or email
}