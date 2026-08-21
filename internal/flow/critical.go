package flow

import (
	"math"

	"weir-manning/internal/cross"
)

func CriticalDepthRect(b, q float64) (float64, error) {
	if b <= 0 {
		return 0, errBottomPositive
	}
	if q < 0 {
		return 0, ErrDischarge
	}
	raw := math.Cbrt(q * q / (cross.Gravity * b * b))
	return applyYc(raw), nil
}

func CriticalDepth(s cross.Section, q float64, opts SolveOptions) (float64, error) {
	if err := cross.ValidateSection(s); err != nil {
		return 0, err
	}
	if q <= 0 {
		return 0, ErrDischarge
	}
	if err := opts.Validate(); err != nil {
		return 0, err
	}
	low := opts.LowDepth
	high := opts.InitialHigh
	for i := 0; i < opts.MaxBracketSteps; i++ {
		fr, err := Froude(s, q, high)
		if err != nil {
			return 0, err
		}
		if fr <= 1 {
			return bisectFroude(s, q, low, high, opts)
		}
		low = high
		high *= opts.BracketExpand
	}
	return 0, ErrEmptyBracket
}

func bisectFroude(s cross.Section, q, low, high float64, opts SolveOptions) (float64, error) {
	mid := (low + high) / 2
	for i := 0; i < opts.MaxIterations; i++ {
		mid = (low + high) / 2
		fr, err := Froude(s, q, mid)
		if err != nil {
			return 0, err
		}
		if math.Abs(fr-1) <= opts.Tolerance {
			return mid, nil
		}
		if fr > 1 {
			low = mid
		} else {
			high = mid
		}
	}
	return mid, ErrNotConverged
}

func CriticalSlopeRect(b, n, q float64) (float64, error) {
	yc, err := CriticalDepthRect(b, q)
	if err != nil {
		return 0, err
	}
	if n <= 0 {
		return 0, ErrRoughness
	}
	s, err := cross.NewRect(b)
	if err != nil {
		return 0, err
	}
	qAtYc, err := Discharge(s, n, 1.0, yc)
	if err != nil {
		return 0, err
	}
	return (q / qAtYc) * (q / qAtYc), nil
}
