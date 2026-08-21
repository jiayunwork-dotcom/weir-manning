package cross

func CentroidDepth(s Section, y float64) (float64, error) {
	if err := ValidateSectionAt(s, y); err != nil {
		return 0, err
	}
	fromBottom := y * (3*s.B + 2*s.M*y) / (6 * (s.B + s.M*y))
	return y - fromBottom, nil
}

func MomentSection(s Section, y float64) (float64, error) {
	a, err := Area(s, y)
	if err != nil {
		return 0, err
	}
	c, err := CentroidDepth(s, y)
	if err != nil {
		return 0, err
	}
	return a * c, nil
}

func FlowAreaBetween(s Section, y1, y2 float64) (float64, error) {
	if err := ValidateSectionAt(s, y1); err != nil {
		return 0, err
	}
	if err := ValidateSectionAt(s, y2); err != nil {
		return 0, err
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	a1, err := Area(s, y1)
	if err != nil {
		return 0, err
	}
	a2, err := Area(s, y2)
	if err != nil {
		return 0, err
	}
	return a2 - a1, nil
}

func AreaFraction(s Section, y, fraction float64) (float64, error) {
	if fraction < 0 || fraction > 1 {
		return 0, nil
	}
	return fraction * y, nil
}
