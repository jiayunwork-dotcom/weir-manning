package weir

import (
	"math"
	"testing"
)

func TestWeirDischargeFormulaValue(t *testing.T) {
	q, err := Discharge(1.0, 0.62, 0.3)
	if err != nil {
		t.Fatalf("Discharge failed: %v", err)
	}
	want := 0.62 * 1.0 * math.Sqrt(2*9.81) * math.Pow(0.3, 1.5)
	if math.Abs(q-want) > 1e-12 {
		t.Errorf("Q = %v, want %v", q, want)
	}
}

func TestWeirHeadScalingThreeHalves(t *testing.T) {
	heads := []float64{0.2, 0.3, 0.4, 0.5}
	flows, err := SeriesForHead(heads, 1.5, 0.62)
	if err != nil {
		t.Fatalf("SeriesForHead failed: %v", err)
	}
	for i := 1; i < len(heads); i++ {
		want := math.Pow(heads[i]/heads[i-1], 1.5)
		got := flows[i] / flows[i-1]
		if math.Abs(got-want) > 1e-12 {
			t.Errorf("Q ratio head %v->%v = %v, want 3/2 power %v", heads[i-1], heads[i], got, want)
		}
	}
	if e := ScaleExponent(heads[0], heads[2], flows[0], flows[2]); math.Abs(e-1.5) > 1e-12 {
		t.Errorf("scale exponent = %v, want 1.5", e)
	}
	if !ExponentFalls(heads, flows, 1.5, 1e-9) {
		t.Error("all head steps must keep the 3/2 power exponent")
	}
}

func TestWeirInvalidInputsError(t *testing.T) {
	cases := []struct {
		name string
		b    float64
		cd   float64
		h    float64
	}{
		{"zero head", 1.5, 0.62, 0},
		{"negative head", 1.5, 0.62, -0.2},
		{"zero width", 0, 0.62, 0.25},
		{"negative width", -1.5, 0.62, 0.25},
		{"zero coefficient", 1.5, 0, 0.25},
		{"negative coefficient", 1.5, -0.62, 0.25},
	}
	for _, c := range cases {
		if _, err := Discharge(c.b, c.cd, c.h); err == nil {
			t.Errorf("%s: expected an error, got none", c.name)
		}
	}
}

func TestWeirHeadForDischargeInverts(t *testing.T) {
	b, cd, h := 1.5, 0.62, 0.25
	q, err := Discharge(b, cd, h)
	if err != nil {
		t.Fatalf("Discharge failed: %v", err)
	}
	back, err := HeadForDischarge(b, cd, q)
	if err != nil {
		t.Fatalf("HeadForDischarge failed: %v", err)
	}
	if math.Abs(back-h)/h > 1e-9 {
		t.Errorf("inverted head = %v, want %v", back, h)
	}
}

func TestWeirResultExponent(t *testing.T) {
	res, err := Result(1.5, 0.62, 0.25)
	if err != nil {
		t.Fatalf("Result failed: %v", err)
	}
	e, err := res.Exponent(0.5)
	if err != nil {
		t.Fatalf("Exponent failed: %v", err)
	}
	if math.Abs(e-1.5) > 1e-9 {
		t.Errorf("result exponent = %v, want 1.5", e)
	}
	if res.Rate <= 0 {
		t.Errorf("rate %v must be positive", res.Rate)
	}
}

func TestInverseCoefficientRoundTrip(t *testing.T) {
	b, cd, h := 1.5, 0.62, 0.25
	q, err := Discharge(b, cd, h)
	if err != nil {
		t.Fatalf("Discharge failed: %v", err)
	}
	gotCD, err := CoefficientForDischarge(b, h, q)
	if err != nil {
		t.Fatalf("CoefficientForDischarge failed: %v", err)
	}
	if math.Abs(gotCD-cd)/cd > 1e-9 {
		t.Errorf("inverted Cd = %v, want %v", gotCD, cd)
	}
	got, err := Discharge(b, gotCD, h)
	if err != nil {
		t.Fatalf("Discharge failed: %v", err)
	}
	if math.Abs(got-q)/q > 1e-9 {
		t.Errorf("round-trip discharge = %v, want %v", got, q)
	}
}

func TestWidthForDischargeRoundTrip(t *testing.T) {
	b, cd, h := 1.5, 0.62, 0.25
	q, err := Discharge(b, cd, h)
	if err != nil {
		t.Fatalf("Discharge failed: %v", err)
	}
	w, err := WidthForDischarge(cd, h, q)
	if err != nil {
		t.Fatalf("WidthForDischarge failed: %v", err)
	}
	if math.Abs(w-b)/b > 1e-9 {
		t.Errorf("inverted width = %v, want %v", w, b)
	}
}

func TestSeriesCurveSummary(t *testing.T) {
	heads := []float64{0.2, 0.3, 0.4, 0.5}
	rs, err := SeriesResults(heads, 1.5, 0.62)
	if err != nil {
		t.Fatalf("SeriesResults failed: %v", err)
	}
	if len(rs) != len(heads) {
		t.Fatalf("series length = %d, want %d", len(rs), len(heads))
	}
	for i := 1; i < len(rs); i++ {
		if rs[i].Q <= rs[i-1].Q {
			t.Errorf("series discharge must rise with head, got %v then %v", rs[i-1].Q, rs[i].Q)
		}
	}
	minQ, maxQ, err := CurveSummary(heads, 1.5, 0.62)
	if err != nil {
		t.Fatalf("CurveSummary failed: %v", err)
	}
	if minQ != rs[0].Q || maxQ != rs[len(rs)-1].Q {
		t.Errorf("summary min=%v max=%v, want %v and %v", minQ, maxQ, rs[0].Q, rs[len(rs)-1].Q)
	}
}
