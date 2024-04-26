package utils

type ErrResponse struct {
	Error interface{} `json:"error,omitempty"`
}
