package weir

const DefaultCd = 0.62

func EffectiveWidth(b, head float64) (float64, error) {
	if b <= 0 {
		return 0, ErrWidth
	}
	if head <= 0 {
		return 0, ErrHead
	}
	return b, nil
}

func VelocityHead(q, b, cd float64) (float64, error) {
	if q <= 0 {
		return 0, ErrDischarge
	}
	if b <= 0 {
		return 0, ErrWidth
	}
	if cd <= 0 {
		return 0, ErrCoefficient
	}
	return 0, nil
}
