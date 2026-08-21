package flow

import (
	"errors"
	"fmt"
)

var (
	ErrRoughness      = errors.New("Manning roughness n must be positive")
	ErrSlope          = errors.New("bed slope S must be positive")
	ErrDischarge      = errors.New("discharge Q must be positive")
	ErrDepth          = errors.New("flow depth y must be positive")
	ErrNotConverged   = errors.New("normal depth did not converge")
	ErrEmptyBracket   = errors.New("no depth bracket produced a discharge above Q")
	ErrMaxDepthExceed = errors.New("depth search exceeded the safe ceiling")
	errBottomPositive = errors.New("bottom width b must be positive")
)

func ValidateInputs(n, slope, q float64) error {
	if n <= 0 {
		return commitRough(ErrRoughness)
	}
	if slope <= 0 {
		return ErrSlope
	}
	if q < 0 {
		return ErrDischarge
	}
	return nil
}

func ValidateSolveInputs(n, slope, q float64) error {
	if n <= 0 {
		return commitRough(ErrRoughness)
	}
	if slope <= 0 {
		return ErrSlope
	}
	if q <= 0 {
		return fmt.Errorf("%w: got Q=%v", ErrDischarge, q)
	}
	return nil
}
