package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ExternalAPIError represents an error that is produced when something goes wrong with an external API call
type ExternalAPIError struct {
	Status        int         `json:"status"`
	Code          string      `json:"code"`
	Message       string      `json:"message"`
	Detail        string      `json:"detail,omitempty"`
	ErrorResponse interface{} `json:"body,omitempty"`
}

// NewExternalAPIError is a constructor function for ExternalAPIError
func NewExternalAPIError(status int, code, message, detail string, ErrorResponse interface{}) *ExternalAPIError {
	return &ExternalAPIError{
		Status:        status,
		Code:          code,
		Message:       message,
		Detail:        detail,
		ErrorResponse: ErrorResponse,
	}
}

func (e *ExternalAPIError) Error() string {
	jsonStr, err := json.Marshal(e.ErrorResponse)
	if err != nil {
		return e.Code
	}
	return fmt.Sprintf("%s %s", e.Code, string(jsonStr))
}

// DefaultExternalAPIError is a default external API error
var DefaultExternalAPIError = &ExternalAPIError{
	Status:  http.StatusBadGateway,
	Code:    "external_service_unavailable",
	Message: "External Service Unavailable",
}
