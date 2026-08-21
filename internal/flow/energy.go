package flow

import (
	"math"

	"weir-manning/internal/cross"
)

func SpecificEnergy(s cross.Section, q, depth float64) (float64, error) {
	if q < 0 {
		return 0, ErrDischarge
	}
	g, err := cross.SnapshotAt(s, depth)
	if err != nil {
		return 0, err
	}
	v := q / g.Area
	return depth + v*v/(2*cross.Gravity), nil
}

func MomentumFunction(s cross.Section, q, depth float64) (float64, error) {
	if q < 0 {
		return 0, ErrDischarge
	}
	g, err := cross.SnapshotAt(s, depth)
	if err != nil {
		return 0, err
	}
	m, err := cross.MomentSection(s, depth)
	if err != nil {
		return 0, err
	}
	return q*q/(cross.Gravity*g.Area) + m, nil
}

func EnergyDerivative(s cross.Section, q, depth float64) (float64, error) {
	if q < 0 {
		return 0, ErrDischarge
	}
	g, err := cross.SnapshotAt(s, depth)
	if err != nil {
		return 0, err
	}
	if g.DepthMean <= 0 {
		return 0, ErrDepth
	}
	v := q / g.Area
	fr := v / math.Sqrt(cross.Gravity*g.DepthMean)
	return 1 - fr*fr, nil
}

func MinimumEnergyDepth(s cross.Section, q float64, opts SolveOptions) (float64, error) {
	if q <= 0 {
		return 0, ErrDischarge
	}
	if err := cross.ValidateSection(s); err != nil {
		return 0, err
	}
	if err := opts.Validate(); err != nil {
		return 0, err
	}
	low := opts.LowDepth
	high := opts.InitialHigh
	for i := 0; i < opts.MaxBracketSteps; i++ {
		d, err := EnergyDerivative(s, q, high)
		if err != nil {
			return 0, err
		}
		if d > 0 {
			return bisectEnergy(s, q, low, high, opts)
		}
		low = high
		high *= opts.BracketExpand
	}
	return 0, ErrEmptyBracket
}

func bisectEnergy(s cross.Section, q, low, high float64, opts SolveOptions) (float64, error) {
	mid := (low + high) / 2
	for i := 0; i < opts.MaxIterations; i++ {
		mid = (low + high) / 2
		d, err := EnergyDerivative(s, q, mid)
		if err != nil {
			return 0, err
		}
		if math.Abs(d) <= opts.Tolerance {
			return mid, nil
		}
		if d > 0 {
			high = mid
		} else {
			low = mid
		}
	}
	return mid, ErrNotConverged
}

func AlternateDepths(s cross.Section, q, energy float64, opts SolveOptions) (sub, super float64, err error) {
	yc, err := CriticalDepth(s, q, opts)
	if err != nil {
		return 0, 0, err
	}
	lo := yc
	hi := opts.InitialHigh
	for i := 0; i < opts.MaxBracketSteps; i++ {
		e, err := SpecificEnergy(s, q, hi)
		if err != nil {
			return 0, 0, err
		}
		if e >= energy {
			sub, err = bisectEnergyLevel(s, q, lo, hi, energy, opts, true)
			if err != nil {
				return 0, 0, err
			}
			break
		}
		lo = hi
		hi *= opts.BracketExpand
	}
	hi = yc
	lo = opts.LowDepth
	for i := 0; i < opts.MaxBracketSteps; i++ {
		e, err := SpecificEnergy(s, q, lo)
		if err != nil {
			return 0, 0, err
		}
		if e >= energy {
			super, err = bisectEnergyLevel(s, q, lo, hi, energy, opts, false)
			if err != nil {
				return 0, 0, err
			}
			break
		}
		hi = lo
		lo *= 0.5
	}
	if sub <= 0 || super <= 0 {
		return 0, 0, ErrNotConverged
	}
	return sub, super, nil
}

func bisectEnergyLevel(s cross.Section, q, low, high, energy float64, opts SolveOptions, rising bool) (float64, error) {
	mid := (low + high) / 2
	for i := 0; i < opts.MaxIterations; i++ {
		mid = (low + high) / 2
		e, err := SpecificEnergy(s, q, mid)
		if err != nil {
			return 0, err
		}
		if absRel(e, energy) <= opts.Tolerance {
			return mid, nil
		}
		if rising {
			if e < energy {
				low = mid
			} else {
				high = mid
			}
		} else {
			if e < energy {
				high = mid
			} else {
				low = mid
			}
		}
	}
	return mid, ErrNotConverged
}
