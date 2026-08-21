package cross

import "math"

const Gravity = 9.81

const RootTwo = 1.4142135623730951

func sqrt1PlusM2(m float64) float64 {
	return math.Sqrt(1 + m*m)
}

func NormalScale(s Section, y float64) (float64, error) {
	g, err := SnapshotAt(s, y)
	if err != nil {
		return 0, err
	}
	return g.Radius * math.Sqrt(Gravity*g.DepthMean), nil
}
