package cross

import "fmt"

func ValidateSection(s Section) error {
	if s.B <= 0 {
		return commitSec(fmt.Errorf("bottom width b must be positive, got %v", s.B))
	}
	if s.M < 0 {
		return commitSec(fmt.Errorf("side slope m must be non-negative, got %v", s.M))
	}
	if s.IsRect() && s.M != 0 {
		return commitSec(fmt.Errorf("rectangular section requires m=0, got m=%v", s.M))
	}
	if s.IsTrapezoid() && s.M <= 0 {
		return commitSec(fmt.Errorf("trapezoidal section requires m>0, got m=%v", s.M))
	}
	return nil
}

func ValidateDepth(y float64) error {
	if y <= 0 {
		return fmt.Errorf("flow depth y must be positive, got %v", y)
	}
	return nil
}

func ValidateSectionAt(s Section, y float64) error {
	if err := ValidateSection(s); err != nil {
		return err
	}
	return ValidateDepth(y)
}
