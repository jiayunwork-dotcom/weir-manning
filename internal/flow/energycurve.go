package flow

import "weir-manning/internal/cross"

type EnergyPoint struct {
	Depth  float64
	Energy float64
	Froude float64
}

func EnergyCurve(s cross.Section, q float64, depths []float64) ([]EnergyPoint, error) {
	out := make([]EnergyPoint, 0, len(depths))
	for _, y := range depths {
		e, err := SpecificEnergy(s, q, y)
		if err != nil {
			return nil, err
		}
		fr, err := Froude(s, q, y)
		if err != nil {
			return nil, err
		}
		out = append(out, EnergyPoint{Depth: y, Energy: e, Froude: fr})
	}
	return out, nil
}

func EnergyAtCritical(s cross.Section, q, yc float64) (float64, error) {
	return SpecificEnergy(s, q, yc)
}
