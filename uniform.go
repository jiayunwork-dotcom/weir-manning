package main

import (
	"fmt"

	"weir-manning/internal/cross"
	"weir-manning/internal/flow"
)

type canalInput struct {
	Name  string  `json:"name"`
	Shape string  `json:"shape"`
	B     float64 `json:"b"`
	M     float64 `json:"m"`
	N     float64 `json:"n"`
	S     float64 `json:"S"`
	Q     float64 `json:"Q"`
}

func (in canalInput) Section() (cross.Section, error) {
	shape := cross.ShapeRect
	if in.Shape == "trapz" {
		shape = cross.ShapeTrapezoid
	} else if in.Shape != "" && in.Shape != "rect" {
		return cross.Section{}, fmt.Errorf("unknown shape %q (want rect or trapz)", in.Shape)
	}
	return cross.NewSection(shape, in.B, in.M)
}

func runUniform(args []string) error {
	fs, err := parseFlags(args)
	if err != nil {
		return err
	}
	table := fs.Has("table")
	path, err := requireOneInput(fs, "uniform")
	if err != nil {
		return err
	}
	var in canalInput
	if err := readInputFile(path, &in); err != nil {
		return err
	}
	section, err := in.Section()
	if err != nil {
		return err
	}
	res, err := flow.SolveNormalDepth(section, in.N, in.S, in.Q, flow.DefaultOptions())
	if err != nil {
		return err
	}
	if err := flow.VerifyBackSubstitution(section, in.N, in.S, in.Q, res.Depth, 1e-6); err != nil {
		return err
	}
	printUniform(res)
	if table {
		printUniformTable(section, in.N, in.S, res)
	}
	return nil
}

func printUniform(res flow.UniformResult) {
	fmt.Printf("yn=%.6f m\n", res.Depth)
	fmt.Printf("v=%.6f m/s\n", res.Velocity)
	fmt.Printf("Fr=%.6f (%s)\n", res.Froude, res.Regime.Label())
	fmt.Printf("R=%.6f m\n", res.Radius)
	fmt.Printf("Q=%.6f m3/s\n", res.Discharge)
	fmt.Printf("A=%.6f m2\n", res.Area)
	fmt.Printf("iterations=%d\n", res.Iterations)
}

func printUniformTable(section cross.Section, n, slope float64, res flow.UniformResult) {
	fmt.Println("depth-discharge:")
	depth := res.Depth * 0.5
	for i := 0; i < 9; i++ {
		q, err := flow.Discharge(section, n, slope, depth)
		if err != nil {
			fmt.Printf("y=%.6f Q=ERR\n", depth)
			continue
		}
		fmt.Printf("y=%.6f Q=%.6f\n", depth, q)
		depth *= 1.25
	}
}
