package lexer

import (
	"bytes"
	"strconv"
	"testing"
)

func TestMatchNumberR(t *testing.T) {
	reals := map[string]float64{
		"0":           0,
		"10":          10,
		"-0":          0,
		"-7":          -7,
		"-10":         -10,
		"987654321":   987654321,
		"0.34":        0.34,
		"-2.14":       -2.14,
		"2e3":         2000,
		"-3e1":        -30,
		"5.3E2":       530,
		"1.234e2":     123.4,
		"-3.4523e+1":  -34.523,
		"256E-2":      2.56,
		"354e-4":      0.0354,
		"-3487.23e-1": -348.723,
		"0.001e+5":    100,
		"-9e-2":       -0.09,
	}

	for str, num := range reals {
		b := bytes.NewBufferString(str)
		m := &Matcher{Scanner: NewScanner(b)}

		if m.MatchNumber() {
			matched := m.Extract()
			f, e := strconv.ParseFloat(matched, 64)
			if matched != str {
				t.Errorf("Expected to match %q, got %q\n", str, matched)
			} else if e != nil {
				t.Error(e)
			} else if f != num {
				t.Errorf("Expected number %v, got %v\n", num, f)
			} else {
				t.Logf("Expected %v, got %v", str, f)
			}
		}
	}
}

func TestMatchNumberBases(t *testing.T) {
	hex := map[string]int64{
		"0x0":      0,
		"0x10":     16,
		"-0x20":    -32,
		"0xff":     255,
		"0xabcdEf": 11259375,
	}

	for str, num := range hex {
		b := bytes.NewBufferString(str)
		m := &Matcher{Scanner: NewScanner(b)}

		if m.MatchNumber() {
			matched := m.Extract()
			f, e := strconv.ParseInt(matched, 0, 64)
			if matched != str {
				t.Errorf("Expected to match %q, got %q\n", str, matched)
			} else if e != nil {
				t.Error(e)
			} else if f != num {
				t.Errorf("Expected number %v, got %v\n", num, f)
			} else {
				t.Logf("Expected %v, got %v", str, f)
			}
		}
	}
}

func TestMatchNumberW(t *testing.T) {
	hex := map[string]string{
		"0x0ugly":      "0x0",
		"0x10test":     "0x10",
		"-0x20inconv":  "-0x20",
		"0o777":        "0o777",
		"0b10101":      "0b10101",
		"-3.1415e-20i": "-3.1415e-20i",
	}

	for str, exp := range hex {
		b := bytes.NewBufferString(str)
		m := &Matcher{Scanner: NewScanner(b)}

		if m.MatchNumber() {
			matched := m.Extract()
			t.Logf("Expected %q, got %q", exp, matched)
			if matched != exp {
				t.Fail()
			}
		}
	}
}
