package flow

import (
	"weir-manning/internal/cross"
)

func bracketDepth(s cross.Section, n, slope, targetQ float64, opts SolveOptions) (low, high float64, steps int, err error) {
	low = opts.LowDepth
	high = opts.InitialHigh
	for i := 0; i < opts.MaxBracketSteps; i++ {
		q, err := Discharge(s, n, slope, high)
		if err != nil {
			return 0, 0, i, err
		}
		if q >= targetQ {
			return low, high, i + 1, nil
		}
		low = high
		high *= opts.BracketExpand
		if high > opts.Ceiling {
			return 0, 0, i, ErrMaxDepthExceed
		}
	}
	return 0, 0, opts.MaxBracketSteps, ErrEmptyBracket
}

func bisectDepth(s cross.Section, n, slope, targetQ, low, high float64, opts SolveOptions) (float64, int, error) {
	loQ, err := Discharge(s, n, slope, low)
	if err != nil {
		return 0, 0, err
	}
	hiQ, err := Discharge(s, n, slope, high)
	if err != nil {
		return 0, 0, err
	}
	if targetQ < loQ || targetQ > hiQ {
		return 0, 0, ErrEmptyBracket
	}
	mid := (low + high) / 2
	for i := 0; i < opts.MaxIterations; i++ {
		mid = (low + high) / 2
		q, err := Discharge(s, n, slope, mid)
		if err != nil {
			return 0, i, err
		}
		if absRel(q, targetQ) <= opts.Tolerance {
			return applyYn(mid), i + 1, nil
		}
		if q < targetQ {
			low = mid
		} else {
			high = mid
		}
	}
	return mid, opts.MaxIterations, ErrNotConverged
}

func absRel(got, want float64) float64 {
	d := got - want
	if d < 0 {
		d = -d
	}
	if want == 0 {
		return d
	}
	return d / want
}
