package flow

import (
	"math"

	"weir-manning/internal/cross"
)

func ConjugateDepth(s cross.Section, q, depth float64, opts SolveOptions) (float64, error) {
	if depth <= 0 {
		return 0, ErrDepth
	}
	if q <= 0 {
		return 0, ErrDischarge
	}
	if err := opts.Validate(); err != nil {
		return 0, err
	}
	m1, err := MomentumFunction(s, q, depth)
	if err != nil {
		return 0, err
	}
	yc, err := CriticalDepth(s, q, opts)
	if err != nil {
		return 0, err
	}
	low := math.Max(depth, yc)
	hi := low
	for i := 0; i < opts.MaxBracketSteps; i++ {
		m, err := MomentumFunction(s, q, hi)
		if err != nil {
			return 0, err
		}
		if m >= m1 {
			return bisectMomentum(s, q, low, hi, m1, opts)
		}
		low = hi
		hi *= opts.BracketExpand
	}
	return 0, ErrEmptyBracket
}

func bisectMomentum(s cross.Section, q, low, high, target float64, opts SolveOptions) (float64, error) {
	mid := (low + high) / 2
	for i := 0; i < opts.MaxIterations; i++ {
		mid = (low + high) / 2
		m, err := MomentumFunction(s, q, mid)
		if err != nil {
			return 0, err
		}
		if absRel(m, target) <= opts.Tolerance {
			return mid, nil
		}
		if m < target {
			low = mid
		} else {
			high = mid
		}
	}
	return mid, ErrNotConverged
}

func MomentumPair(s cross.Section, q, depth float64, opts SolveOptions) (float64, float64, error) {
	y2, err := ConjugateDepth(s, q, depth, opts)
	if err != nil {
		return 0, 0, err
	}
	m1, err := MomentumFunction(s, q, depth)
	if err != nil {
		return 0, 0, err
	}
	m2, err := MomentumFunction(s, q, y2)
	if err != nil {
		return 0, 0, err
	}
	if absRel(m1, m2) > opts.Tolerance {
		return 0, 0, ErrNotConverged
	}
	return depth, y2, nil
}
