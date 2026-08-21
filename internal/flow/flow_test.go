package flow

import (
	"errors"
	"math"
	"testing"

	"weir-manning/internal/cross"
)

func rectSection(t *testing.T) cross.Section {
	t.Helper()
	s, err := cross.NewRect(2.0)
	if err != nil {
		t.Fatalf("NewRect failed: %v", err)
	}
	return s
}

func trapzSection(t *testing.T) cross.Section {
	t.Helper()
	s, err := cross.NewTrapezoid(3.0, 1.0)
	if err != nil {
		t.Fatalf("NewTrapezoid failed: %v", err)
	}
	return s
}

func TestNormalDepthBackSubstitutionRecoversDischarge(t *testing.T) {
	s := rectSection(t)
	res, err := SolveNormalDepth(s, 0.015, 0.001, 2.0, DefaultOptions())
	if err != nil {
		t.Fatalf("SolveNormalDepth failed: %v", err)
	}
	if res.Discharge <= 0 {
		t.Fatalf("solved discharge %v must be positive", res.Discharge)
	}
	rel, err := BackSubstitutionError(s, 0.015, 0.001, 2.0, res.Depth)
	if err != nil {
		t.Fatalf("BackSubstitutionError failed: %v", err)
	}
	if rel > 1e-6 {
		t.Errorf("back-substituted Q relative error = %v, want below 1e-6", rel)
	}
}

func TestRectCriticalDepthClosedForm(t *testing.T) {
	b, q := 2.0, 2.0
	yc, err := CriticalDepthRect(b, q)
	if err != nil {
		t.Fatalf("CriticalDepthRect failed: %v", err)
	}
	want := math.Cbrt(q * q / (cross.Gravity * b * b))
	if math.Abs(yc-want) > 1e-12 {
		t.Errorf("yc = %v, want closed form %v", yc, want)
	}
	s := rectSection(t)
	fr, err := Froude(s, q, yc)
	if err != nil {
		t.Fatalf("Froude failed: %v", err)
	}
	if math.Abs(fr-1) > 1e-6 {
		t.Errorf("Froude at critical depth = %v, want 1", fr)
	}
	if Classify(fr) != RegimeCritical {
		t.Errorf("regime at yc = %v, want critical", Classify(fr))
	}
}

func TestSteeperSlopeLowersNormalDepth(t *testing.T) {
	s := rectSection(t)
	depths, err := NormalDepthsForSlopes(s, 0.015, 2.0, []float64{0.0005, 0.001, 0.002, 0.004}, DefaultOptions())
	if err != nil {
		t.Fatalf("NormalDepthsForSlopes failed: %v", err)
	}
	if !IsMonotoneDecreasing(depths) {
		t.Errorf("normal depths must fall as slope rises, got %v", depths)
	}
	if depths[0] <= depths[len(depths)-1] {
		t.Errorf("flattest slope must carry the deepest flow, got %v", depths)
	}
}

func TestInvalidManningInputsError(t *testing.T) {
	s := rectSection(t)
	cases := []struct {
		name string
		n    float64
		sl   float64
		q    float64
	}{
		{"zero roughness", 0, 0.001, 2.0},
		{"negative roughness", -0.015, 0.001, 2.0},
		{"zero slope", 0.015, 0, 2.0},
		{"negative slope", 0.015, -0.001, 2.0},
		{"zero discharge", 0.015, 0.001, 0},
		{"negative discharge", 0.015, 0.001, -2.0},
	}
	for _, c := range cases {
		if _, err := SolveNormalDepth(s, c.n, c.sl, c.q, DefaultOptions()); err == nil {
			t.Errorf("%s: expected an error, got none", c.name)
		}
	}
}

func TestFroudeRegimeClassification(t *testing.T) {
	s := rectSection(t)
	q := 2.0
	frDeep, err := Froude(s, q, 1.0)
	if err != nil {
		t.Fatalf("Froude failed: %v", err)
	}
	if Classify(frDeep) != RegimeSubcritical {
		t.Errorf("deep flow Fr=%v must be subcritical", frDeep)
	}
	frShallow, err := Froude(s, q, 0.3)
	if err != nil {
		t.Fatalf("Froude failed: %v", err)
	}
	if Classify(frShallow) != RegimeSupercritical {
		t.Errorf("shallow flow Fr=%v must be supercritical", frShallow)
	}
	yc, err := CriticalDepthRect(2.0, q)
	if err != nil {
		t.Fatalf("CriticalDepthRect failed: %v", err)
	}
	if frShallow <= frDeep {
		t.Errorf("Froude must fall as depth rises, got %v then %v", frShallow, frDeep)
	}
	_ = yc
}

func TestExampleRectCanalHandMagnitude(t *testing.T) {
	s := rectSection(t)
	res, err := SolveNormalDepth(s, 0.015, 0.001, 2.0, DefaultOptions())
	if err != nil {
		t.Fatalf("SolveNormalDepth failed: %v", err)
	}
	if res.Depth < 0.78 || res.Depth > 0.84 {
		t.Errorf("yn = %v, want about 0.81 m for the rect-canal example", res.Depth)
	}
	if res.Velocity < 1.0 || res.Velocity > 1.5 {
		t.Errorf("v = %v, want about 1.23 m/s", res.Velocity)
	}
	if res.Froude >= 1 {
		t.Errorf("Fr = %v, example must be subcritical", res.Froude)
	}
	if res.Radius <= 0 || res.Radius >= res.Depth {
		t.Errorf("R = %v must lie in (0, yn=%v)", res.Radius, res.Depth)
	}
}

func TestNonConvergenceReported(t *testing.T) {
	s := rectSection(t)
	opts := DefaultOptions()
	opts.MaxIterations = 1
	opts.Tolerance = 1e-12
	_, err := SolveNormalDepth(s, 0.015, 0.001, 2.0, opts)
	if err == nil {
		t.Fatal("a single bisection step with tiny tolerance must not converge")
	}
	if !errors.Is(err, ErrNotConverged) {
		t.Errorf("must report ErrNotConverged, got %v", err)
	}

	opts = DefaultOptions()
	opts.InitialHigh = 1.0
	opts.BracketExpand = 2.0
	opts.Ceiling = 4.0
	_, err = SolveNormalDepth(s, 0.015, 0.001, 1e6, opts)
	if err == nil {
		t.Fatal("a discharge that outruns the ceiling must fail")
	}
	if !errors.Is(err, ErrMaxDepthExceed) && !errors.Is(err, ErrEmptyBracket) {
		t.Errorf("must report a bracket failure, got %v", err)
	}
}

func TestNormalDepthPipelineSharesSectionGeometry(t *testing.T) {
	for _, tc := range []struct {
		name string
		s    cross.Section
		n    float64
		sl   float64
		q    float64
	}{
		{"rect", rectSection(t), 0.015, 0.001, 2.0},
		{"trapz", trapzSection(t), 0.025, 0.002, 5.0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, err := SolveNormalDepth(tc.s, tc.n, tc.sl, tc.q, DefaultOptions())
			if err != nil {
				t.Fatalf("SolveNormalDepth failed: %v", err)
			}
			yn := res.Depth
			aMan := yn * (tc.s.B + tc.s.M*yn)
			pMan := tc.s.B + 2*yn*math.Sqrt(1+tc.s.M*tc.s.M)
			rMan := aMan / pMan
			qMan := (1 / tc.n) * aMan * math.Pow(rMan, 2.0/3.0) * math.Sqrt(tc.sl)
			if rel := math.Abs(qMan-tc.q) / tc.q; rel > 1e-6 {
				t.Errorf("independent perimeter/area recompute gives Q=%v, want %v (rel %v)", qMan, tc.q, rel)
			}
			if math.Abs(res.Radius-rMan) > 1e-9 {
				t.Errorf("result R=%v must equal section radius %v", res.Radius, rMan)
			}
			dMan := aMan / (tc.s.B + 2*tc.s.M*yn)
			if math.Abs(res.DepthMean-dMan) > 1e-9 {
				t.Errorf("result depth mean %v must equal section value %v", res.DepthMean, dMan)
			}
			frMan := res.Velocity / math.Sqrt(cross.Gravity*dMan)
			if math.Abs(res.Froude-frMan) > 1e-9 {
				t.Errorf("result Fr=%v must equal section-based Fr=%v", res.Froude, frMan)
			}
			if rMan >= yn {
				t.Errorf("radius %v must stay below depth %v; both side walls must be in the wetted perimeter", rMan, yn)
			}
		})
	}
}

func TestTrapezoidNormalDepthBackSubstitution(t *testing.T) {
	s := trapzSection(t)
	res, err := SolveNormalDepth(s, 0.025, 0.002, 5.0, DefaultOptions())
	if err != nil {
		t.Fatalf("SolveNormalDepth failed: %v", err)
	}
	rel, err := BackSubstitutionError(s, 0.025, 0.002, 5.0, res.Depth)
	if err != nil {
		t.Fatalf("BackSubstitutionError failed: %v", err)
	}
	if rel > 1e-6 {
		t.Errorf("trapezoid back-substitution relative error = %v, want below 1e-6", rel)
	}
}

func TestMinimumEnergyMatchesCriticalDepth(t *testing.T) {
	for _, tc := range []struct {
		name string
		s    cross.Section
		q    float64
	}{
		{"rect", rectSection(t), 2.0},
		{"trapz", trapzSection(t), 5.0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			yc, err := CriticalDepth(tc.s, tc.q, DefaultOptions())
			if err != nil {
				t.Fatalf("CriticalDepth failed: %v", err)
			}
			ymin, err := MinimumEnergyDepth(tc.s, tc.q, DefaultOptions())
			if err != nil {
				t.Fatalf("MinimumEnergyDepth failed: %v", err)
			}
			if math.Abs(ymin-yc)/yc > 1e-6 {
				t.Errorf("minimum-energy depth %v must match critical depth %v", ymin, yc)
			}
		})
	}
}

func TestAlternateDepthsShareEnergy(t *testing.T) {
	s := rectSection(t)
	q := 2.0
	energy, err := SpecificEnergy(s, q, 1.5)
	if err != nil {
		t.Fatalf("SpecificEnergy failed: %v", err)
	}
	sub, super, err := AlternateDepths(s, q, energy, DefaultOptions())
	if err != nil {
		t.Fatalf("AlternateDepths failed: %v", err)
	}
	yc, err := CriticalDepth(s, q, DefaultOptions())
	if err != nil {
		t.Fatalf("CriticalDepth failed: %v", err)
	}
	if super >= yc || sub <= yc {
		t.Errorf("alternate depths must straddle yc=%v, got super=%v sub=%v", yc, super, sub)
	}
	eSub, err := SpecificEnergy(s, q, sub)
	if err != nil {
		t.Fatalf("SpecificEnergy failed: %v", err)
	}
	eSuper, err := SpecificEnergy(s, q, super)
	if err != nil {
		t.Fatalf("SpecificEnergy failed: %v", err)
	}
	if math.Abs(eSub-eSuper)/energy > 1e-6 {
		t.Errorf("alternate depths must share specific energy, got %v and %v", eSub, eSuper)
	}
}

func TestConjugateDepthsConserveMomentum(t *testing.T) {
	s := rectSection(t)
	q := 2.0
	y1, y2, err := MomentumPair(s, q, 0.3, DefaultOptions())
	if err != nil {
		t.Fatalf("MomentumPair failed: %v", err)
	}
	if y2 <= y1 {
		t.Errorf("conjugate depth %v must exceed upstream depth %v", y2, y1)
	}
	m1, err := MomentumFunction(s, q, y1)
	if err != nil {
		t.Fatalf("MomentumFunction failed: %v", err)
	}
	m2, err := MomentumFunction(s, q, y2)
	if err != nil {
		t.Fatalf("MomentumFunction failed: %v", err)
	}
	if rel := math.Abs(m1-m2) / m1; rel > 1e-6 {
		t.Errorf("momentum not conserved across conjugate depths: %v vs %v", m1, m2)
	}
}

func TestUniformStateIdentitiesHold(t *testing.T) {
	for _, tc := range []struct {
		name  string
		s     cross.Section
		n     float64
		sl    float64
		depth float64
	}{
		{"rect", rectSection(t), 0.015, 0.001, 1.0},
		{"trapz", trapzSection(t), 0.025, 0.002, 1.2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := VerifyUniformState(tc.s, tc.n, tc.sl, tc.depth); err != nil {
				t.Errorf("uniform-state identities violated: %v", err)
			}
		})
	}
}

func TestRatingCurveMonotonic(t *testing.T) {
	s := rectSection(t)
	pts, err := RatingCurve(s, 0.015, 0.001, DepthsRisingTo(0.2, 2.0, 10))
	if err != nil {
		t.Fatalf("RatingCurve failed: %v", err)
	}
	for i := 1; i < len(pts); i++ {
		if pts[i].Discharge <= pts[i-1].Discharge {
			t.Errorf("discharge must rise with depth, got %v then %v", pts[i-1].Discharge, pts[i].Discharge)
		}
	}
}

func TestSlopeInverseRoundTrip(t *testing.T) {
	s := rectSection(t)
	n, q, depth := 0.015, 2.0, 1.0
	slope, err := SlopeForDischarge(s, n, q, depth)
	if err != nil {
		t.Fatalf("SlopeForDischarge failed: %v", err)
	}
	got, err := Discharge(s, n, slope, depth)
	if err != nil {
		t.Fatalf("Discharge failed: %v", err)
	}
	if math.Abs(got-q)/q > 1e-9 {
		t.Errorf("round-trip discharge = %v, want %v", got, q)
	}
}
