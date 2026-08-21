package flow

import "weir-manning/internal/cross"

type CurvePoint struct {
	Depth     float64
	Discharge float64
	Velocity  float64
	Froude    float64
	Regime    Regime
}

func RatingCurve(s cross.Section, n, slope float64, depths []float64) ([]CurvePoint, error) {
	out := make([]CurvePoint, 0, len(depths))
	for _, y := range depths {
		p, err := PointAt(s, n, slope, y)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func PointAt(s cross.Section, n, slope, depth float64) (CurvePoint, error) {
	res, err := UniformState(s, n, slope, depth)
	if err != nil {
		return CurvePoint{}, err
	}
	return CurvePoint{
		Depth:     depth,
		Discharge: res.Discharge,
		Velocity:  res.Velocity,
		Froude:    res.Froude,
		Regime:    res.Regime,
	}, nil
}

func RatingTableAround(s cross.Section, n, slope, targetQ float64, steps int) ([]CurvePoint, error) {
	if steps < 2 {
		steps = 8
	}
	res, err := SolveNormalDepth(s, n, slope, targetQ, DefaultOptions())
	if err != nil {
		return nil, err
	}
	lo := res.Depth * 0.4
	hi := res.Depth * 2.0
	depths := make([]float64, 0, steps+1)
	for i := 0; i <= steps; i++ {
		y := lo + (hi-lo)*float64(i)/float64(steps)
		depths = append(depths, y)
	}
	return RatingCurve(s, n, slope, depths)
}

func SolveSeriesQ(s cross.Section, n, slope float64, discharges []float64, opts SolveOptions) ([]UniformResult, error) {
	out := make([]UniformResult, 0, len(discharges))
	for _, q := range discharges {
		res, err := SolveNormalDepth(s, n, slope, q, opts)
		if err != nil {
			return nil, err
		}
		out = append(out, res)
	}
	return out, nil
}

func DepthsRisingTo(lo, hi float64, steps int) []float64 {
	depths := make([]float64, 0, steps+1)
	for i := 0; i <= steps; i++ {
		depths = append(depths, lo+(hi-lo)*float64(i)/float64(steps))
	}
	return depths
}
