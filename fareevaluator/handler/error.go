package handler

import "errors"

var (
	ErrInvalidDataCount   = errors.New("invalid data count")
	ErrInvalidFormat      = errors.New("invalid format")
	ErrInvalidMileage     = errors.New("invalid mileage")
	ErrInvalidElapsedTime = errors.New("invalid elapsed time")
)
