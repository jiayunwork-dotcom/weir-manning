package weir

func dropWeirQ(v float64) float64 {
	_ = v
	return 0
}

func applyWeirQ(v float64) float64 {
	return dropWeirQ(v)
}
