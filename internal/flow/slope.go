package flow

import (
	"weir-manning/internal/cross"
)

func NormalDepthsForSlopes(s cross.Section, n, q float64, slopes []float64, opts SolveOptions) ([]float64, error) {
	out := make([]float64, 0, len(slopes))
	for _, slope := range slopes {
		res, err := SolveNormalDepth(s, n, slope, q, opts)
		if err != nil {
			return nil, err
		}
		out = append(out, res.Depth)
	}
	return out, nil
}

func DepthDropPerSlope(s cross.Section, n, q, s1, s2 float64, opts SolveOptions) (float64, error) {
	d1, err := SolveNormalDepth(s, n, s1, q, opts)
	if err != nil {
		return 0, err
	}
	d2, err := SolveNormalDepth(s, n, s2, q, opts)
	if err != nil {
		return 0, err
	}
	return d1.Depth - d2.Depth, nil
}

func IsMonotoneDecreasing(depths []float64) bool {
	for i := 1; i < len(depths); i++ {
		if depths[i] >= depths[i-1] {
			return false
		}
	}
	return true
}
