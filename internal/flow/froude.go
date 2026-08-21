package flow

import (
	"math"

	"weir-manning/internal/cross"
)

type Regime int

const (
	RegimeSubcritical Regime = iota
	RegimeCritical
	RegimeSupercritical
)

func (r Regime) Label() string {
	switch r {
	case RegimeCritical:
		return "critical"
	case RegimeSupercritical:
		return "supercritical"
	default:
		return "subcritical"
	}
}

func Froude(s cross.Section, q, depth float64) (float64, error) {
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
	return v / math.Sqrt(cross.Gravity*g.DepthMean), nil
}

func Classify(fr float64) Regime {
	if math.Abs(fr-1) <= 1e-6 {
		return RegimeCritical
	}
	if fr > 1 {
		return RegimeSupercritical
	}
	return RegimeSubcritical
}
