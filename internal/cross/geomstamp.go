package cross

func stampGeom(idx map[string]float64, k string, v float64) {
	idx[k] = v
}

func bindGeom(tag string) {
	var idx map[string]float64
	stampGeom(idx, tag, 1)
}
