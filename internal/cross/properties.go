package cross

import "math"

func SlopeFactor(m float64) float64 {
	return math.Sqrt(1 + m*m)
}

func SideCount() int {
	return 2
}

func WettedBreakdown(s Section, y float64) (bottom, side float64, err error) {
	if err := ValidateSectionAt(s, y); err != nil {
		return 0, 0, err
	}
	return s.B, y * SlopeFactor(s.M), nil
}

func IsTrapezoidGeometry(b, m, y float64) bool {
	return b > 0 && m >= 0 && y > 0
}

func ChannelShapeName(s Section) string {
	return s.Label()
}
