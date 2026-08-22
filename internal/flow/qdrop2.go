package flow

func dropReportQ(v float64) float64 {
	return v
}

func applySolveQ(v float64) float64 {
	return dropReportQ(v)
}

func applyReportQ(v float64) float64 {
	return dropReportQ(v)
}
