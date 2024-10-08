package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	NotFoundErrorType            = "NotFoundError"
	InternalServerErrorType      = "InternalServerError"
	BadRequestErrorType          = "BadRequestError"
	UnauthorizedErrorType        = "UnauthorizedError"
	ForbiddenErrorType           = "ForbiddenError"
	ConflictErrorType            = "ConflictError"
	MethodNotAllowedErrorType    = "MethodNotAllowedError"
	RequestTimeoutErrorType      = "RequestTimeoutError"
	UnprocessableEntityErrorType = "UnprocessableEntityError"
	TooManyRequestsErrorType     = "TooManyRequestsError"
)

type ErrorOption func(*ApiError)

type ApiErrors interface {
	Error() string // This is the method required by the error interface in Go
	Type() string
	Code() int
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

// ApiError represents a structured error for the API.
type ApiError struct {
	ErrorType  string `json:"error_type"`
	Message    string `json:"message"`
	ErrorCode  int    `json:"error_code"`
	InnerError error  `json:"-"`
}

// make sure ApiError implements ApiErrors interface in compile time
var _ ApiErrors = (*ApiError)(nil)

// Type return ApiError Type.
func (e *ApiError) Type() string {
	return e.ErrorType
}

// Code return ApiError code.
func (e *ApiError) Code() int {
	return e.ErrorCode
}

// Error implements the error interface for ApiError
func (e *ApiError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.ErrorCode, e.Message)
}

// HTTPError generates an HTTP error response
func (e *ApiError) HTTPError() (int, string) {
	return e.ErrorCode, e.Message
}

// InternalError return ApiError message.
func (e *ApiError) InternalError() error {
	return e.InnerError
}

// MarshalJSON customizes the JSON serialization for ApiError.
func (e *ApiError) MarshalJSON() ([]byte, error) {
	type Alias ApiError // Create an alias to avoid recursion
	var internalErrorStr string
	if e.InnerError != nil {
		internalErrorStr = e.InnerError.Error()
	}

	return json.Marshal(&struct {
		InternalError string `json:"internal_error,omitempty"`
		*Alias
	}{
		InternalError: internalErrorStr,
		Alias:         (*Alias)(e),
	})
}

// UnmarshalJSON customizes the JSON deserialization for ApiError.
func (e *ApiError) UnmarshalJSON(data []byte) error {
	type Alias ApiError
	aux := &struct {
		InternalError *string `json:"internal_error,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.InternalError != nil {
		e.InnerError = fmt.Errorf(*aux.InternalError)
	}
	return nil
}

// ErrorType represents an error type configuration
type ErrorType struct {
	ErrorCode int
	Message   string
}

// ErrorRegistry is a map of error types and their properties
var ErrorRegistry = map[string]ErrorType{
	NotFoundErrorType:            {http.StatusNotFound, "Resource not found"},
	InternalServerErrorType:      {http.StatusInternalServerError, "Internal server error"},
	BadRequestErrorType:          {http.StatusBadRequest, "Bad request"},
	UnauthorizedErrorType:        {http.StatusUnauthorized, "Unauthorized access"},
	ForbiddenErrorType:           {http.StatusForbidden, "Forbidden"},
	ConflictErrorType:            {http.StatusConflict, "Conflict occurred"},
	MethodNotAllowedErrorType:    {http.StatusMethodNotAllowed, "Method not allowed"},
	RequestTimeoutErrorType:      {http.StatusRequestTimeout, "Request timed out"},
	UnprocessableEntityErrorType: {http.StatusUnprocessableEntity, "Unprocessable entity"},
	TooManyRequestsErrorType:     {http.StatusTooManyRequests, "Too many requests"},
	// You can add more error types as needed...
}

// NewApiError creates a new ApiError based on the error type.
func NewApiError(errorType string, userMessage string, options ...ErrorOption) *ApiError {
	apiError := &ApiError{
		ErrorType: "GenericError",
		Message:   userMessage,
		ErrorCode: http.StatusInternalServerError,
	}
	if errType, exists := ErrorRegistry[errorType]; exists {
		apiError = &ApiError{
			ErrorType: errorType,
			Message:   userMessage,
			ErrorCode: errType.ErrorCode,
		}
	}
	for _, option := range options {
		option(apiError)
	}
	return apiError
}

// RegisterErrorType Optional: Method to add new error types to the registry at runtime
func RegisterErrorType(name string, errorCode int, message string) {
	ErrorRegistry[name] = ErrorType{errorCode, message}
}

// WithInternalError to wrap internal errors
func WithInternalError(err error) ErrorOption {
	return func(ae *ApiError) {
		ae.InnerError = err
	}
}
