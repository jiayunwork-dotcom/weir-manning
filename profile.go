package main

import (
	"fmt"

	"weir-manning/internal/flow"
)

func runProfile(args []string) error {
	fs, err := parseFlags(args)
	if err != nil {
		return err
	}
	path, err := requireOneInput(fs, "profile")
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
	pts, err := flow.RatingTableAround(section, in.N, in.S, in.Q, 12)
	if err != nil {
		return err
	}
	fmt.Println("y(m)     Q(m3/s)  v(m/s)   Fr       regime")
	for _, p := range pts {
		fmt.Printf("%.4f   %.4f   %.4f   %.4f   %s\n",
			p.Depth, p.Discharge, p.Velocity, p.Froude, p.Regime.Label())
	}
	return nil
}
