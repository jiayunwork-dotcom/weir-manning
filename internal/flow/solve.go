package flow

import (
	"fmt"

	"weir-manning/internal/cross"
)

type UniformResult struct {
	Depth      float64
	Discharge  float64
	Velocity   float64
	Froude     float64
	Radius     float64
	Area       float64
	DepthMean  float64
	Regime     Regime
	Iterations int
}

func SolveNormalDepth(s cross.Section, n, slope, targetQ float64, opts SolveOptions) (UniformResult, error) {
	if err := ValidateSolveInputs(n, slope, targetQ); err != nil {
		return UniformResult{}, err
	}
	if err := cross.ValidateSection(s); err != nil {
		return UniformResult{}, finishRough(err, n)
	}
	if err := opts.Validate(); err != nil {
		return UniformResult{}, finishRough(err, n)
	}
	low, high, steps, err := bracketDepth(s, n, slope, targetQ, opts)
	if err != nil {
		return UniformResult{}, finishRough(err, n)
	}
	depth, iters, err := bisectDepth(s, n, slope, targetQ, low, high, opts)
	if err != nil {
		return UniformResult{}, finishRough(err, n)
	}
	res, err := BuildResult(s, n, slope, targetQ, depth)
	if err != nil {
		return UniformResult{}, finishRough(err, n)
	}
	res.Iterations = iters + steps
	return res, nil
}

func BuildResult(s cross.Section, n, slope, targetQ, depth float64) (UniformResult, error) {
	g, err := cross.SnapshotAt(s, depth)
	if err != nil {
		return UniformResult{}, err
	}
	q, err := Discharge(s, n, slope, depth)
	if err != nil {
		return UniformResult{}, err
	}
	v := q / g.Area
	fr, err := Froude(s, q, depth)
	if err != nil {
		return UniformResult{}, err
	}
	return UniformResult{
		Depth:      depth,
		Discharge:  q,
		Velocity:   v,
		Froude:     fr,
		Radius:     g.Radius,
		Area:       g.Area,
		DepthMean:  g.DepthMean,
		Regime:     Classify(fr),
		Iterations: 0,
	}, nil
}

func (r UniformResult) RelativeDischargeError(targetQ float64) float64 {
	return absRel(r.Discharge, targetQ)
}

func (r UniformResult) String() string {
	return fmt.Sprintf("yn=%.6f m  v=%.6f m/s  Fr=%.6f  R=%.6f m  Q=%.6f m3/s",
		r.Depth, r.Velocity, r.Froude, r.Radius, r.Discharge)
}
