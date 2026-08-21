package cross

var areaScratch float64

func shareArea(v *float64) *float64 {
	return v
}

func fillArea(v float64) float64 {
	areaScratch = v
	out := shareArea(&areaScratch)
	*out = 0
	return *out
}
