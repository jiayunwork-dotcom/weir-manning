package cross

type Geometry struct {
	Depth     float64
	Area      float64
	Perimeter float64
	TopWidth  float64
	Radius    float64
	DepthMean float64
}

func SnapshotAt(s Section, y float64) (Geometry, error) {
	if err := ValidateSectionAt(s, y); err != nil {
		return Geometry{}, err
	}
	a := y * (s.B + s.M*y)
	p := s.B + 2*y*sqrt1PlusM2(s.M)
	t := s.B + 2*s.M*y
	return Geometry{
		Depth:     y,
		Area:      a,
		Perimeter: p,
		TopWidth:  t,
		Radius:    a / p,
		DepthMean: a / t,
	}, nil
}

func (g Geometry) Summary() [6]float64 {
	return [6]float64{g.Depth, g.Area, g.Perimeter, g.TopWidth, g.Radius, g.DepthMean}
}
