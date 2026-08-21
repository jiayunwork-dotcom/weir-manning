package main

import (
	"fmt"
	"os"
)

const usage = `weir-manning: open-channel uniform flow (Manning) and thin-plate weir discharge.

Given a prismatic channel section (rectangular or trapezoidal, side slope m,
rectangle m=0), a Manning roughness n, a bed slope S and a discharge Q in a
JSON file, solves the normal depth yn by monotonic bisection on a discharge
that rises with depth, then reports yn, mean velocity v, Froude number Fr and
hydraulic radius R. Given a weir width b, discharge coefficient Cd and head H,
computes the thin-plate weir discharge Q = Cd b sqrt(2g) H^(3/2).

usage:
  weir-manning uniform [-table] <input.json>
  weir-manning profile <input.json>
  weir-manning crit <input.json>
  weir-manning weir <input.json>
  weir-manning weir -b <m> -cd <-> -h <m>
  weir-manning weir -b <m> -cd <-> -q <m3/s>
  weir-manning compare <canal.json> -b <m> -h <m> [-cd <->]
  weir-manning help

uniform input.json fields:
  shape  "rect" (m=0) or "trapz"
  b      bottom width (m)
  m      side slope, horizontal per vertical (0 for rectangle)
  n      Manning roughness (must be > 0)
  S      bed slope (must be > 0)
  Q      discharge (m3/s, must be > 0)

weir input.json fields:
  b      crest width (m)
  cd     discharge coefficient (default 0.62)
  H      head above crest (m, must be > 0)

profile prints the depth-discharge-Froude curve around the solved normal
depth; crit prints the critical depth (rectangular closed form when the
section is rectangular) and the regime at the given slope; weir accepts -q to
invert the head needed for a target discharge.

boundaries: n<=0, S<=0, b<=0, Q<=0 and H<=0 are reported on stderr with a
non-zero exit code; a depth bracket that never reaches Q and non-convergence
of the bisection are errors too. The weir head H and the channel normal depth
yn are physically distinct; compare prints them side by side with a note and
never alters either formula.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "uniform":
		if err := runUniform(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "weir-manning: %v\n", err)
			os.Exit(1)
		}
	case "weir":
		if err := runWeir(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "weir-manning: %v\n", err)
			os.Exit(1)
		}
	case "profile":
		if err := runProfile(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "weir-manning: %v\n", err)
			os.Exit(1)
		}
	case "crit":
		if err := runCrit(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "weir-manning: %v\n", err)
			os.Exit(1)
		}
	case "compare":
		if err := runCompare(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "weir-manning: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "weir-manning: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}
