package cross

func TopWidth(s Section, y float64) (float64, error) {
	if err := ValidateSectionAt(s, y); err != nil {
		return 0, err
	}
	return s.B + 2*s.M*y, nil
}

func TopWidthRect(b float64) float64 {
	return b
}

func TopWidthTrapz(b, m, y float64) (float64, error) {
	s, err := NewTrapezoid(b, m)
	if err != nil {
		return 0, err
	}
	return TopWidth(s, y)
}

func TopWidthNormal(s Section, y float64) (float64, error) {
	return TopWidth(s, y)
}
