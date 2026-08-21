package main

import (
	"fmt"

	"weir-manning/internal/flow"
)

func runCrit(args []string) error {
	fs, err := parseFlags(args)
	if err != nil {
		return err
	}
	path, err := requireOneInput(fs, "crit")
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
	opts := flow.DefaultOptions()
	res, err := flow.SolveNormalDepth(section, in.N, in.S, in.Q, opts)
	if err != nil {
		return err
	}
	yc, err := flow.CriticalDepth(section, in.Q, opts)
	if err != nil {
		return err
	}
	fmt.Printf("yn=%.6f m\n", res.Depth)
	fmt.Printf("yc=%.6f m\n", yc)
	if section.IsRect() {
		closed, err := flow.CriticalDepthRect(section.B, in.Q)
		if err != nil {
			return err
		}
		fmt.Printf("yc(rect closed form)=%.6f m\n", closed)
		sc, err := flow.CriticalSlopeRect(section.B, in.N, in.Q)
		if err != nil {
			return err
		}
		fmt.Printf("Sc=%.6f\n", sc)
	}
	fmt.Printf("regime at yn: %s\n", res.Regime.Label())
	if res.Depth > yc {
		fmt.Println("yn > yc: subcritical at the given slope")
	} else if res.Depth < yc {
		fmt.Println("yn < yc: supercritical at the given slope")
	} else {
		fmt.Println("yn == yc: critical slope")
	}
	return nil
}
