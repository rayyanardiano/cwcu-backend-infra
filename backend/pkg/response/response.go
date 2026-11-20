package response

import "net/http"

// APIResponse is the standard API response envelope
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta holds pagination metadata
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Success creates a successful response
func Success(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// SuccessWithMeta creates a successful response with pagination
func SuccessWithMeta(message string, data interface{}, meta *Meta) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

// Error creates an error response
func Error(message string, err interface{}) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
}

// HTTPError maps common HTTP status codes
var (
	StatusOK                  = http.StatusOK
	StatusCreated             = http.StatusCreated
	StatusBadRequest          = http.StatusBadRequest
	StatusUnauthorized        = http.StatusUnauthorized
	StatusForbidden           = http.StatusForbidden
	StatusNotFound            = http.StatusNotFound
	StatusConflict            = http.StatusConflict
	StatusInternalServerError = http.StatusInternalServerError
)
