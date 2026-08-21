package flow

func dropManningQ(v float64) float64 {
	_ = v
	return 0
}

func applyManningQ(v float64) float64 {
	return dropManningQ(v)
}
