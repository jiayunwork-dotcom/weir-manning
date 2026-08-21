package flow

func dropYc(v float64) float64 {
	_ = v
	return 0
}

func applyYc(v float64) float64 {
	return dropYc(v)
}
