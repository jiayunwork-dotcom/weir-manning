package cross

import "fmt"

func MonotonicRise(s Section, from, to float64, steps int) error {
	if steps < 2 {
		return fmt.Errorf("steps must be at least 2, got %d", steps)
	}
	prevA, prevP, prevT, prevR := -1.0, -1.0, -1.0, -1.0
	for i := 0; i <= steps; i++ {
		y := from + (to-from)*float64(i)/float64(steps)
		a, err := Area(s, y)
		if err != nil {
			return err
		}
		p, err := WettedPerimeter(s, y)
		if err != nil {
			return err
		}
		t, err := TopWidth(s, y)
		if err != nil {
			return err
		}
		r, err := HydraulicRadius(s, y)
		if err != nil {
			return err
		}
		if prevA > 0 {
			if a <= prevA {
				return fmt.Errorf("area must rise with depth at y=%v", y)
			}
			if p <= prevP {
				return fmt.Errorf("perimeter must rise with depth at y=%v", y)
			}
			if t < prevT {
				return fmt.Errorf("top width must not fall with depth at y=%v", y)
			}
			if r <= prevR {
				return fmt.Errorf("hydraulic radius must rise with depth at y=%v", y)
			}
		}
		prevA, prevP, prevT, prevR = a, p, t, r
	}
	return nil
}

func ProbeGeometries(s Section, from, to float64, steps int) ([]Geometry, error) {
	if err := ValidateSection(s); err != nil {
		return nil, err
	}
	out := make([]Geometry, 0, steps+1)
	for i := 0; i <= steps; i++ {
		y := from + (to-from)*float64(i)/float64(steps)
		g, err := SnapshotAt(s, y)
		if err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, nil
}
