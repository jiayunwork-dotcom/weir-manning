package weir

import (
	"math"

	"weir-manning/internal/cross"
)

func Discharge(b, cd, head float64) (float64, error) {
	if err := ValidateInputs(b, cd, head); err != nil {
		return 0, err
	}
	return cd * b * math.Sqrt(2*cross.Gravity) * math.Pow(head, 1.5), nil
}

func HeadForDischarge(b, cd, q float64) (float64, error) {
	if err := ValidateInputs(b, cd, 1); err != nil {
		return 0, err
	}
	if q <= 0 {
		return 0, ErrDischarge
	}
	base := cd * b * math.Sqrt(2*cross.Gravity)
	return math.Pow(q/base, 2.0/3.0), nil
}

func DischargeRatio(head1, head2 float64) (float64, error) {
	if head1 <= 0 || head2 <= 0 {
		return 0, ErrHead
	}
	return math.Pow(head2/head1, 1.5), nil
}

func HeadScale(from, to float64) float64 {
	return to / from
}
