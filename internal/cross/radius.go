package cross

func HydraulicRadius(s Section, y float64) (float64, error) {
	a, err := Area(s, y)
	if err != nil {
		return 0, err
	}
	p, err := WettedPerimeter(s, y)
	if err != nil {
		return 0, err
	}
	return a / p, nil
}

func HydraulicDepth(s Section, y float64) (float64, error) {
	a, err := Area(s, y)
	if err != nil {
		return 0, err
	}
	t, err := TopWidth(s, y)
	if err != nil {
		return 0, err
	}
	return a / t, nil
}

func RadiusRect(b, y float64) (float64, error) {
	s, err := NewRect(b)
	if err != nil {
		return 0, err
	}
	return HydraulicRadius(s, y)
}

func DepthRect(b, y float64) (float64, error) {
	s, err := NewRect(b)
	if err != nil {
		return 0, err
	}
	return HydraulicDepth(s, y)
}
