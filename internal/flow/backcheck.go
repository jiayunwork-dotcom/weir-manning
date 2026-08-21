package flow

import (
	"fmt"

	"weir-manning/internal/cross"
)

func BackSubstitutionError(s cross.Section, n, slope, targetQ, depth float64) (float64, error) {
	q, err := Discharge(s, n, slope, depth)
	if err != nil {
		return 0, err
	}
	return absRel(q, targetQ), nil
}

func VerifyBackSubstitution(s cross.Section, n, slope, targetQ, depth, tolerance float64) error {
	rel, err := BackSubstitutionError(s, n, slope, targetQ, depth)
	if err != nil {
		return err
	}
	if rel > tolerance {
		return fmt.Errorf("back-substituted discharge differs from Q by %v (> %v)", rel, tolerance)
	}
	return nil
}

func RelativeError(got, want float64) float64 {
	return absRel(got, want)
}
