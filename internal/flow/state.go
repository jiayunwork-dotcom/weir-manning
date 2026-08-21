package flow

import (
	"weir-manning/internal/cross"
)

func UniformState(s cross.Section, n, slope, depth float64) (UniformResult, error) {
	if err := ValidateInputs(n, slope, 0); err != nil {
		return UniformResult{}, err
	}
	if err := cross.ValidateSectionAt(s, depth); err != nil {
		return UniformResult{}, err
	}
	return BuildResult(s, n, slope, 0, depth)
}

func SlopeForDischarge(s cross.Section, n, targetQ, depth float64) (float64, error) {
	if targetQ <= 0 {
		return 0, ErrDischarge
	}
	if err := cross.ValidateSectionAt(s, depth); err != nil {
		return 0, err
	}
	k, err := Conveyance(s, n, depth)
	if err != nil {
		return 0, err
	}
	return (targetQ / k) * (targetQ / k), nil
}
