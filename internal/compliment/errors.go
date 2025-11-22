package compliment

import "errors"

var (
	ErrUserNotAuthenticated = errors.New("user not authenticated")
	ErrComplimentNotFound   = errors.New("compliment not found")
	ErrInvalidPoints        = errors.New("points must be between 0 and 5")
	ErrInvalidUser          = errors.New("cannot compliment yourself")
	ErrUserNotFound         = errors.New("user not found")
)

