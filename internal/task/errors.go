package task

import "errors"

var (
	ErrUserNotAuthenticated        = errors.New("user not authenticated")
	ErrTaskNotFound                = errors.New("task not found")
	ErrTaskAlreadyCompleted        = errors.New("task already completed")
	ErrTaskNotCompleted            = errors.New("task is not completed")
	ErrFrequencyNotDefined         = errors.New("frequency unit or frequency value not defined")
)
