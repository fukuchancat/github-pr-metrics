package utils

import (
	"fmt"
)

// APIError represents an error from the GitHub API
type APIError struct {
	StatusCode int
	Message    string
}

// Error returns the error message
func (e *APIError) Error() string {
	return fmt.Sprintf("GitHub API error (status code %d): %s", e.StatusCode, e.Message)
}

// RateLimitError represents a rate limit error from the GitHub API
type RateLimitError struct {
	ResetTime string
}

// Error returns the error message
func (e *RateLimitError) Error() string {
	return fmt.Sprintf("GitHub API rate limit exceeded. Reset at %s", e.ResetTime)
}
