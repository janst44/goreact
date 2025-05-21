package main

// ErrorResponse represents a generic error response.
//
// swagger:model
type ErrorResponse struct {
	Code    ErrorCode `json:"code" example:"Non-http status code machine-readable error"`
	Message string    `json:"message" example:"human readable error message"`
	Details string    `json:"details,omitempty" example:"optional details about the error"`
}

type ValidationResponse struct {
	Code    ErrorCode         `json:"code" example:"VALIDATION_ERROR"`
	Message string            `json:"message" example:"Validation failed"`
	Errors  []ValidationError `json:"errors"`
}

type ValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"must be a valid email address"`
}
