package flow

var velScratch float64

func shareVel(v *float64) *float64 {
	return v
}

func fillVel(v float64) float64 {
	velScratch = v
	out := shareVel(&velScratch)
	return *out
}
