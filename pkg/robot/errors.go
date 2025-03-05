package robot

import "fmt"

type HTTPError struct {
	StatusCode int         `json:"-"`
	Message    interface{} `json:"message"`
	Err        error
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("HTTP error with status code: %d", e.StatusCode)
}

// Unwrap returns the underlying error
func (e *HTTPError) Unwrap() error {
	return e.Err
}
