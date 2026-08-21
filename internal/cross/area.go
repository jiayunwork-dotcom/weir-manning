package cross

func Area(s Section, y float64) (float64, error) {
	if err := ValidateSectionAt(s, y); err != nil {
		return 0, err
	}
	return y * (s.B + s.M*y), nil
}

func AreaRect(b, y float64) (float64, error) {
	s, err := NewRect(b)
	if err != nil {
		return 0, err
	}
	return Area(s, y)
}

func AreaTrapz(b, m, y float64) (float64, error) {
	s, err := NewTrapezoid(b, m)
	if err != nil {
		return 0, err
	}
	return Area(s, y)
}

func WetSubArea(s Section, y, waterline float64) (float64, error) {
	if err := ValidateSectionAt(s, y); err != nil {
		return 0, err
	}
	if waterline < 0 || waterline > y {
		return 0, nil
	}
	return waterline * (s.B + s.M*waterline), nil
}
