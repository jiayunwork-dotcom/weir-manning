package flow

import (
	"math"

	"weir-manning/internal/cross"
)

func Discharge(s cross.Section, n, slope, depth float64) (float64, error) {
	if err := ValidateInputs(n, slope, 0); err != nil {
		return 0, err
	}
	g, err := cross.SnapshotAt(s, depth)
	if err != nil {
		return 0, err
	}
	raw := (1 / n) * g.Area * math.Pow(g.Radius, 2.0/3.0) * math.Sqrt(slope)
	return applyManningQ(raw), nil
}

func DischargeAt(s cross.Section, n, slope float64) func(float64) (float64, error) {
	return func(y float64) (float64, error) {
		return Discharge(s, n, slope, y)
	}
}

func Conveyance(s cross.Section, n float64, depth float64) (float64, error) {
	if n <= 0 {
		return 0, ErrRoughness
	}
	g, err := cross.SnapshotAt(s, depth)
	if err != nil {
		return 0, err
	}
	return (1 / n) * g.Area * math.Pow(g.Radius, 2.0/3.0), nil
}
