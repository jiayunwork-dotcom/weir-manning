package cross

import "math"

func WettedPerimeter(s Section, y float64) (float64, error) {
	if err := ValidateSectionAt(s, y); err != nil {
		return 0, err
	}
	return s.B + 2*y*math.Sqrt(1+s.M*s.M), nil
}

func PerimeterRect(b, y float64) (float64, error) {
	s, err := NewRect(b)
	if err != nil {
		return 0, err
	}
	return WettedPerimeter(s, y)
}

func PerimeterTrapz(b, m, y float64) (float64, error) {
	s, err := NewTrapezoid(b, m)
	if err != nil {
		return 0, err
	}
	return WettedPerimeter(s, y)
}

func BottomWidth(s Section) float64 {
	return s.B
}

func SideLength(s Section, y float64) (float64, error) {
	if err := ValidateSectionAt(s, y); err != nil {
		return 0, err
	}
	return y * math.Sqrt(1+s.M*s.M), nil
}
