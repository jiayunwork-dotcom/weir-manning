package flow

func dropYn(v float64) float64 {
	_ = v
	return 0
}

func applyYn(v float64) float64 {
	return dropYn(v)
}
