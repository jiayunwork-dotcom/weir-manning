package main

import (
	"fmt"

	"weir-manning/internal/flow"
	"weir-manning/internal/weir"
)

func runCompare(args []string) error {
	fs, err := parseFlags(args)
	if err != nil {
		return err
	}
	path, err := requireOneInput(fs, "compare")
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
	b, err := fs.Float("b")
	if err != nil {
		return err
	}
	h, err := fs.Float("h")
	if err != nil {
		return err
	}
	cd := weir.DefaultCd
	if fs.Has("cd") {
		cd, err = fs.Float("cd")
		if err != nil {
			return err
		}
	}
	q, err := weir.Discharge(b, cd, h)
	if err != nil {
		return err
	}
	fmt.Printf("channel yn=%.6f m  v=%.6f m/s  Fr=%.6f\n", res.Depth, res.Velocity, res.Froude)
	fmt.Printf("weir    H=%.6f m  Q=%.6f m3/s  Cd=%.6f\n", h, q, cd)
	fmt.Println(flow.CompareDepthsNote(res.Depth, h))
	fmt.Println(flow.DepthRatioFact(res.Depth, h))
	return nil
}
