package cross

import (
	"math"
	"testing"
)

func TestRectGeometryValues(t *testing.T) {
	s, err := NewRect(2.0)
	if err != nil {
		t.Fatalf("NewRect failed: %v", err)
	}
	a, err := Area(s, 0.8)
	if err != nil {
		t.Fatalf("Area failed: %v", err)
	}
	if math.Abs(a-1.6) > 1e-12 {
		t.Errorf("area = %v, want 1.6", a)
	}
	p, err := WettedPerimeter(s, 0.8)
	if err != nil {
		t.Fatalf("WettedPerimeter failed: %v", err)
	}
	if math.Abs(p-3.6) > 1e-12 {
		t.Errorf("perimeter = %v, want 3.6", p)
	}
	tw, err := TopWidth(s, 0.8)
	if err != nil {
		t.Fatalf("TopWidth failed: %v", err)
	}
	if math.Abs(tw-2.0) > 1e-12 {
		t.Errorf("top width = %v, want 2.0", tw)
	}
	r, err := HydraulicRadius(s, 0.8)
	if err != nil {
		t.Fatalf("HydraulicRadius failed: %v", err)
	}
	if math.Abs(r-1.6/3.6) > 1e-12 {
		t.Errorf("radius = %v, want %v", r, 1.6/3.6)
	}
	d, err := HydraulicDepth(s, 0.8)
	if err != nil {
		t.Fatalf("HydraulicDepth failed: %v", err)
	}
	if math.Abs(d-0.8) > 1e-12 {
		t.Errorf("depth mean = %v, want 0.8", d)
	}
}

func TestTrapezoidGeometryMatchesHandCalc(t *testing.T) {
	s, err := NewTrapezoid(3.0, 1.0)
	if err != nil {
		t.Fatalf("NewTrapezoid failed: %v", err)
	}
	a, err := Area(s, 2.0)
	if err != nil {
		t.Fatalf("Area failed: %v", err)
	}
	if math.Abs(a-10.0) > 1e-12 {
		t.Errorf("area = %v, want 10", a)
	}
	p, err := WettedPerimeter(s, 2.0)
	if err != nil {
		t.Fatalf("WettedPerimeter failed: %v", err)
	}
	wantP := 3.0 + 2*2.0*math.Sqrt(2)
	if math.Abs(p-wantP) > 1e-12 {
		t.Errorf("perimeter = %v, want %v", p, wantP)
	}
	tw, err := TopWidth(s, 2.0)
	if err != nil {
		t.Fatalf("TopWidth failed: %v", err)
	}
	if math.Abs(tw-7.0) > 1e-12 {
		t.Errorf("top width = %v, want 7", tw)
	}
	r, err := HydraulicRadius(s, 2.0)
	if err != nil {
		t.Fatalf("HydraulicRadius failed: %v", err)
	}
	if math.Abs(r-10.0/wantP) > 1e-12 {
		t.Errorf("radius = %v, want %v", r, 10.0/wantP)
	}
}

func TestRectangleReducesFromTrapezoidFormula(t *testing.T) {
	s, err := NewRect(2.5)
	if err != nil {
		t.Fatalf("NewRect failed: %v", err)
	}
	for _, y := range []float64{0.3, 0.8, 1.7} {
		a, err := Area(s, y)
		if err != nil {
			t.Fatalf("Area failed: %v", err)
		}
		if math.Abs(a-2.5*y) > 1e-12 {
			t.Errorf("rect area at y=%v = %v, want %v", y, a, 2.5*y)
		}
		p, err := WettedPerimeter(s, y)
		if err != nil {
			t.Fatalf("WettedPerimeter failed: %v", err)
		}
		if math.Abs(p-(2.5+2*y)) > 1e-12 {
			t.Errorf("rect perimeter at y=%v = %v, want %v", y, p, 2.5+2*y)
		}
		tw, err := TopWidth(s, y)
		if err != nil {
			t.Fatalf("TopWidth failed: %v", err)
		}
		if math.Abs(tw-2.5) > 1e-12 {
			t.Errorf("rect top width at y=%v = %v, want %v", y, tw, 2.5)
		}
	}
}

func TestInvalidSectionInputsRejected(t *testing.T) {
	cases := []struct {
		name  string
		shape Shape
		b     float64
		m     float64
	}{
		{"zero width", ShapeRect, 0, 0},
		{"negative width", ShapeRect, -1, 0},
		{"negative slope", ShapeRect, 2, -0.5},
		{"rect with slope", ShapeRect, 2, 0.5},
		{"trapz zero slope", ShapeTrapezoid, 2, 0},
		{"trapz negative slope", ShapeTrapezoid, 2, -1},
	}
	for _, c := range cases {
		if _, err := NewSection(c.shape, c.b, c.m); err == nil {
			t.Errorf("%s: expected an error, got none", c.name)
		}
	}
}

func TestGeometryMonotonicInDepth(t *testing.T) {
	rect, err := NewRect(2.0)
	if err != nil {
		t.Fatalf("NewRect failed: %v", err)
	}
	if err := MonotonicRise(rect, 0.01, 4.0, 40); err != nil {
		t.Errorf("rect geometry not monotonic: %v", err)
	}
	trapz, err := NewTrapezoid(3.0, 1.5)
	if err != nil {
		t.Fatalf("NewTrapezoid failed: %v", err)
	}
	if err := MonotonicRise(trapz, 0.01, 4.0, 40); err != nil {
		t.Errorf("trapz geometry not monotonic: %v", err)
	}
}

func TestSnapshotConsistent(t *testing.T) {
	s, err := NewTrapezoid(3.0, 1.0)
	if err != nil {
		t.Fatalf("NewTrapezoid failed: %v", err)
	}
	g, err := SnapshotAt(s, 2.0)
	if err != nil {
		t.Fatalf("SnapshotAt failed: %v", err)
	}
	if math.Abs(g.Area-g.Perimeter*g.Radius) > 1e-12 {
		t.Errorf("area %v != perimeter*radius %v", g.Area, g.Perimeter*g.Radius)
	}
	if math.Abs(g.DepthMean*g.TopWidth-g.Area) > 1e-12 {
		t.Errorf("depth mean * top width %v != area %v", g.DepthMean*g.TopWidth, g.Area)
	}
	sum := g.Summary()
	want := [6]float64{2.0, 10.0, 3.0 + 4*math.Sqrt(2), 7.0, 10.0 / (3.0 + 4*math.Sqrt(2)), 10.0 / 7.0}
	for i := range want {
		if math.Abs(sum[i]-want[i]) > 1e-12 {
			t.Errorf("summary[%d] = %v, want %v", i, sum[i], want[i])
		}
	}
}
