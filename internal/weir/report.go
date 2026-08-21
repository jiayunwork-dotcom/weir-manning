package weir

type DischargeResult struct {
	Width float64
	Cd    float64
	Head  float64
	Q     float64
	Rate  float64
}

func Result(b, cd, head float64) (DischargeResult, error) {
	q, err := Discharge(b, cd, head)
	if err != nil {
		return DischargeResult{}, err
	}
	return DischargeResult{
		Width: b,
		Cd:    cd,
		Head:  head,
		Q:     q,
		Rate:  q / head,
	}, nil
}

func (r DischargeResult) Exponent(head2 float64) (float64, error) {
	if head2 <= 0 {
		return 0, ErrHead
	}
	q2, err := Discharge(r.Width, r.Cd, head2)
	if err != nil {
		return 0, err
	}
	return ScaleExponent(r.Head, head2, r.Q, q2), nil
}
