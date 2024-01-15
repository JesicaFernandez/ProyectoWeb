package web

type ErrorResponse struct {
	// Status is the status of the error
	Status int `json:"status"`
	// Message is the message of the error
	Message string `json:"message"`
}

// Success is a function that returns a success response
func Success(status int, msg string) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Message: msg,
	}
}

// Error is a function that returns an error response
func Error(status int, msg string) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Message: msg,
	}
}
