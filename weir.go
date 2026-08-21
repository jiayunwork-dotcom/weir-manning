package main

import (
	"fmt"
	"strings"

	"weir-manning/internal/weir"
)

type weirInput struct {
	Name string  `json:"name"`
	B    float64 `json:"b"`
	Cd   float64 `json:"cd"`
	H    float64 `json:"H"`
	Q    float64 `json:"Q"`
}

func runWeir(args []string) error {
	fs, err := parseFlags(args)
	if err != nil {
		return err
	}
	var b, cd, h, q float64
	haveCD := false
	haveH := false
	haveQ := false
	if len(fs.rest) == 1 && strings.HasSuffix(fs.rest[0], ".json") {
		var in weirInput
		if err := readInputFile(fs.rest[0], &in); err != nil {
			return err
		}
		b, cd, h, q = in.B, in.Cd, in.H, in.Q
		haveCD = in.Cd != 0
		haveH = in.H != 0
		haveQ = in.Q != 0
	} else {
		if len(fs.rest) != 0 {
			return fmt.Errorf("weir accepts one JSON file or -b/-cd/-h flags")
		}
		b, err = fs.Float("b")
		if err != nil {
			return err
		}
		if fs.Has("h") {
			h, err = fs.Float("h")
			if err != nil {
				return err
			}
			haveH = true
		}
		if fs.Has("q") {
			q, err = fs.Float("q")
			if err != nil {
				return err
			}
			haveQ = true
		}
		if fs.Has("cd") {
			cd, err = fs.Float("cd")
			if err != nil {
				return err
			}
			haveCD = true
		}
	}
	if !haveCD {
		cd = weir.DefaultCd
	}
	if haveH {
		q, err = weir.Discharge(b, cd, h)
		if err != nil {
			return err
		}
	} else if haveQ {
		h, err = weir.HeadForDischarge(b, cd, q)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("weir needs a head H (or a target discharge -q)")
	}
	fmt.Printf("Q=%.6f m3/s\n", q)
	fmt.Printf("b=%.6f m  Cd=%.6f  H=%.6f m\n", b, cd, h)
	return nil
}
