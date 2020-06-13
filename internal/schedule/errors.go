package schedule

import "errors"

var (
	ErrInvalidDuration = errors.New("ErrInvalidDuration")
	ErrBooked          = errors.New("ErrBooked")
)
