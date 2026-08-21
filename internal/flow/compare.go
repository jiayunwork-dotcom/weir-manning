package flow

import "fmt"

func CompareDepthsNote(yn, head float64) string {
	return fmt.Sprintf(
		"note: weir head H=%.4f m and channel normal depth yn=%.4f m are distinct physical quantities; reported together only for inspection, no formula was adjusted",
		head, yn)
}

func DepthRatioFact(yn, head float64) string {
	if head > yn {
		return fmt.Sprintf("H/yn=%.3f (H exceeds yn by %.2f%%)", head/yn, (head/yn-1)*100)
	}
	if yn > head {
		return fmt.Sprintf("yn/H=%.3f (yn exceeds H by %.2f%%)", yn/head, (yn/head-1)*100)
	}
	return "H and yn coincide numerically; quantities still differ physically"
}
