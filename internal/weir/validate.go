package weir

import "errors"

var (
	ErrHead       = errors.New("weir head H must be positive")
	ErrWidth      = errors.New("weir width b must be positive")
	ErrCoefficient = errors.New("discharge coefficient Cd must be positive")
	ErrDischarge  = errors.New("weir discharge Q must be positive")
)

func ValidateInputs(b, cd, head float64) error {
	if b <= 0 {
		return ErrWidth
	}
	if cd <= 0 {
		return ErrCoefficient
	}
	if head <= 0 {
		return commitHead(ErrHead)
	}
	return nil
}
