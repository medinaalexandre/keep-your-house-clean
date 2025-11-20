package auth

import (
	"errors"
	"net/http"
)

type HTTPError struct {
	Message    string
	StatusCode int
}

func (e *HTTPError) Error() string {
	return e.Message
}

func (e *HTTPError) HTTPCode() int {
	return e.StatusCode
}

var (
	ErrInvalidCredentials = &HTTPError{
		Message:    "invalid credentials",
		StatusCode: http.StatusUnauthorized,
	}
	ErrUserInactive = &HTTPError{
		Message:    "user account is inactive",
		StatusCode: http.StatusForbidden,
	}
	ErrDomainExists = &HTTPError{
		Message:    "domain already exists",
		StatusCode: http.StatusConflict,
	}
	ErrEmailExists = &HTTPError{
		Message:    "email already exists",
		StatusCode: http.StatusConflict,
	}
	ErrPasswordHashFailed = &HTTPError{
		Message:    "failed to hash password",
		StatusCode: http.StatusInternalServerError,
	}
)

func IsHTTPError(err error) (*HTTPError, bool) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr, true
	}
	return nil, false
}

