package weir

import (
	"math"

	"weir-manning/internal/cross"
)

func CoefficientForDischarge(b, head, q float64) (float64, error) {
	if b <= 0 {
		return 0, ErrWidth
	}
	if head <= 0 {
		return 0, ErrHead
	}
	if q <= 0 {
		return 0, ErrDischarge
	}
	return q / (b * math.Sqrt(2*cross.Gravity) * math.Pow(head, 1.5)), nil
}

func WidthForDischarge(cd, head, q float64) (float64, error) {
	if cd <= 0 {
		return 0, ErrCoefficient
	}
	if head <= 0 {
		return 0, ErrHead
	}
	if q <= 0 {
		return 0, ErrDischarge
	}
	return q / (cd * math.Sqrt(2*cross.Gravity) * math.Pow(head, 1.5)), nil
}

func SeriesResults(heads []float64, b, cd float64) ([]DischargeResult, error) {
	out := make([]DischargeResult, 0, len(heads))
	for _, h := range heads {
		r, err := Result(b, cd, h)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func CurveSummary(heads []float64, b, cd float64) (minQ, maxQ float64, err error) {
	rs, err := SeriesResults(heads, b, cd)
	if err != nil {
		return 0, 0, err
	}
	minQ, maxQ = rs[0].Q, rs[0].Q
	for _, r := range rs {
		if r.Q < minQ {
			minQ = r.Q
		}
		if r.Q > maxQ {
			maxQ = r.Q
		}
	}
	return minQ, maxQ, nil
}
