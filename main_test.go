package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestUniformCommandPrintsSolvedValues(t *testing.T) {
	path := writeTemp(t, "canal.json", `{"shape":"rect","b":2.0,"m":0.0,"n":0.015,"S":0.001,"Q":2.0}`)
	out := captureStdout(func() {
		if err := runUniform([]string{path}); err != nil {
			t.Fatalf("runUniform failed: %v", err)
		}
	})
	for _, key := range []string{"yn=", "v=", "Fr=", "R="} {
		if !strings.Contains(out, key) {
			t.Errorf("output must contain %q, got:\n%s", key, out)
		}
	}
	yn := parseValue(t, out, "yn=")
	if yn < 0.78 || yn > 0.84 {
		t.Errorf("printed yn=%v, want about 0.81", yn)
	}
}

func TestWeirCommandRejectsNegativeHead(t *testing.T) {
	if err := runWeir([]string{"-b", "1.5", "-cd", "0.62", "-h", "-0.1"}); err == nil {
		t.Fatal("negative weir head must be rejected")
	}
}

func TestWeirFlagFormMatchesFileForm(t *testing.T) {
	path := writeTemp(t, "weir.json", `{"b":1.5,"cd":0.62,"H":0.25}`)
	outFile := captureStdout(func() {
		if err := runWeir([]string{path}); err != nil {
			t.Fatalf("runWeir file form failed: %v", err)
		}
	})
	outFlags := captureStdout(func() {
		if err := runWeir([]string{"-b", "1.5", "-cd", "0.62", "-h", "0.25"}); err != nil {
			t.Fatalf("runWeir flag form failed: %v", err)
		}
	})
	qFile := parseValue(t, outFile, "Q=")
	qFlags := parseValue(t, outFlags, "Q=")
	if qFile != qFlags {
		t.Errorf("file form Q=%v differs from flag form Q=%v", qFile, qFlags)
	}
	if qFile <= 0 {
		t.Errorf("printed Q=%v must be positive", qFile)
	}
}

func TestUniformRejectsBadInput(t *testing.T) {
	bad := writeTemp(t, "bad.json", `{"shape":"rect","b":2.0,"m":0.0,"n":0.0,"S":0.001,"Q":2.0}`)
	if err := runUniform([]string{bad}); err == nil {
		t.Fatal("zero roughness must be rejected")
	}
	shape := writeTemp(t, "circle.json", `{"shape":"circle","b":2.0,"m":0.0,"n":0.015,"S":0.001,"Q":2.0}`)
	if err := runUniform([]string{shape}); err == nil {
		t.Fatal("unknown shape must be rejected")
	}
	if err := runUniform([]string{}); err == nil {
		t.Fatal("missing input path must be rejected")
	}
}

func TestCompareCommandPrintsHint(t *testing.T) {
	path := writeTemp(t, "canal.json", `{"shape":"rect","b":2.0,"m":0.0,"n":0.015,"S":0.001,"Q":2.0}`)
	out := captureStdout(func() {
		if err := runCompare([]string{path, "-b", "1.5", "-h", "0.25"}); err != nil {
			t.Fatalf("runCompare failed: %v", err)
		}
	})
	if !strings.Contains(out, "note:") {
		t.Errorf("compare output must carry the hint note, got:\n%s", out)
	}
	if !strings.Contains(out, "yn=") || !strings.Contains(out, "weir") {
		t.Errorf("compare output must show both channel and weir values, got:\n%s", out)
	}
}

func TestCritCommandPrintsCriticalDepth(t *testing.T) {
	path := writeTemp(t, "canal.json", `{"shape":"rect","b":2.0,"m":0.0,"n":0.015,"S":0.001,"Q":2.0}`)
	out := captureStdout(func() {
		if err := runCrit([]string{path}); err != nil {
			t.Fatalf("runCrit failed: %v", err)
		}
	})
	for _, key := range []string{"yn=", "yc=", "yc(rect closed form)="} {
		if !strings.Contains(out, key) {
			t.Errorf("crit output must contain %q, got:\n%s", key, out)
		}
	}
}

func TestWeirCommandInvertsHeadFromDischarge(t *testing.T) {
	out := captureStdout(func() {
		if err := runWeir([]string{"-b", "1.5", "-cd", "0.62", "-q", "0.5149232042"}); err != nil {
			t.Fatalf("runWeir inverse failed: %v", err)
		}
	})
	h := parseValue(t, out, "H=")
	if h < 0.249 || h > 0.251 {
		t.Errorf("inverted H=%v, want about 0.25", h)
	}
}

func TestProfileCommandPrintsTable(t *testing.T) {
	path := writeTemp(t, "canal.json", `{"shape":"rect","b":2.0,"m":0.0,"n":0.015,"S":0.001,"Q":2.0}`)
	out := captureStdout(func() {
		if err := runProfile([]string{path}); err != nil {
			t.Fatalf("runProfile failed: %v", err)
		}
	})
	if !strings.Contains(out, "regime") {
		t.Errorf("profile output must carry a table header, got:\n%s", out)
	}
	if !strings.Contains(out, "subcritical") {
		t.Errorf("profile output must show a subcritical row, got:\n%s", out)
	}
}

func parseValue(t *testing.T, out, key string) float64 {
	t.Helper()
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, key) {
			rest := line[strings.Index(line, key)+len(key):]
			var v float64
			if _, err := fmt.Sscanf(rest, "%f", &v); err != nil {
				t.Fatalf("parse %q line %q: %v", key, line, err)
			}
			return v
		}
	}
	t.Fatalf("no line containing %q in:\n%s", key, out)
	return 0
}
