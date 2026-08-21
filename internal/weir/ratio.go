package weir

import "math"

func SeriesForHead(heads []float64, b, cd float64) ([]float64, error) {
	out := make([]float64, 0, len(heads))
	for _, h := range heads {
		q, err := Discharge(b, cd, h)
		if err != nil {
			return nil, err
		}
		out = append(out, q)
	}
	return out, nil
}

func ScaleExponent(head1, head2, q1, q2 float64) float64 {
	return math.Log(q2/q1) / math.Log(head2/head1)
}

func ExponentFalls(heads, flows []float64, exponent, tolerance float64) bool {
	if len(heads) < 2 {
		return false
	}
	for i := 1; i < len(heads); i++ {
		e := ScaleExponent(heads[i-1], heads[i], flows[i-1], flows[i])
		if math.Abs(e-exponent) > tolerance {
			return false
		}
	}
	return true
}
