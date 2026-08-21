package flow

import (
	"fmt"
	"math"

	"weir-manning/internal/cross"
)

type IdentityCheck struct {
	Name string
	OK   bool
	Got  float64
	Want float64
}

func CheckUniformState(s cross.Section, n, slope, depth float64) ([]IdentityCheck, error) {
	res, err := UniformState(s, n, slope, depth)
	if err != nil {
		return nil, err
	}
	g, err := cross.SnapshotAt(s, depth)
	if err != nil {
		return nil, err
	}
	checks := []IdentityCheck{
		{
			Name: "Q=K*sqrt(S)",
			Got:  res.Discharge,
			Want: (1 / n) * g.Area * math.Pow(g.Radius, 2.0/3.0) * math.Sqrt(slope),
		},
		{Name: "A=P*R", Got: g.Area, Want: g.Perimeter * g.Radius},
		{Name: "A=T*D", Got: g.Area, Want: g.TopWidth * g.DepthMean},
		{
			Name: "Fr=v/sqrt(g*D)",
			Got:  res.Froude,
			Want: res.Velocity / math.Sqrt(cross.Gravity*g.DepthMean),
		},
	}
	for i := range checks {
		checks[i].OK = absRel(checks[i].Got, checks[i].Want) <= 1e-9
	}
	return checks, nil
}

func VerifyUniformState(s cross.Section, n, slope, depth float64) error {
	checks, err := CheckUniformState(s, n, slope, depth)
	if err != nil {
		return err
	}
	for _, c := range checks {
		if !c.OK {
			return fmt.Errorf("identity %s violated: got %v want %v", c.Name, c.Got, c.Want)
		}
	}
	return nil
}

func AllIdentitiesHold(checks []IdentityCheck) bool {
	for _, c := range checks {
		if !c.OK {
			return false
		}
	}
	return true
}
