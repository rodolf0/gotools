package lexer

import (
	"bytes"
	"strings"
	"testing"
)

func TestZero(t *testing.T) {
	s := NewScanner(bytes.NewBufferString(""))
	if r := s.Curr(); r != BOF {
		t.Errorf("Expected BOF, got %q\n", r)
	}
	if r := s.Prev(); r != BOF {
		t.Errorf("Expected BOF, got %q\n", r)
	}
	if r := s.Peek(); r != EOF {
		t.Errorf("Expected EOF, got %q\n", r)
	}
	if r := s.Next(); r != EOF {
		t.Errorf("Expected EOF, got %q\n", r)
	}
	if s.pos != 0 {
		t.FailNow()
	}
}

func TestOne(t *testing.T) {
	s := NewScanner(bytes.NewBuffer([]byte("1")))
	for i := 0; i < 5; i++ { // check multiple times for carried errors
		if r := s.Next(); r != '1' || s.pos != 0 {
			t.Errorf("Expected '1' at 0, got %q at %q\n", r, s.pos)
		}
		if r := s.Next(); r != EOF || s.pos != 1 {
			t.Errorf("Expected EOF at 1, got %q at %q\n", r, s.pos)
		}
		if r := s.Peek(); r != EOF || s.pos != 1 {
			t.Errorf("Expected EOF at 1, got %q at %q\n", r, s.pos)
		}
		if r := s.Prev(); r != '1' || s.pos != 0 {
			t.Errorf("Expected '1' at 0, got %q at %q\n", r, s.pos)
		}
		if r := s.Prev(); r != BOF || s.pos != -1 {
			t.Errorf("Expected BOF at -1, got %q at %q\n", r, s.pos)
		}
		if r := s.Peek(); r != '1' || s.pos != -1 {
			t.Errorf("Expected pos -1 and value '1', got %q at %q\n", r, s.pos)
		}
	}
}

var samples = []string{
	"1",
	"one",
	"just a sample line",
	"this is a tokenization test that should be fed to the scanner",
}

func TestStepping(t *testing.T) {
	for _, in := range samples {
		b := bytes.NewBufferString(in)
		s := NewScanner(b)
		for k := 0; k < 5; k++ {
			for i := 0; i < len(in); i++ {
				if r := s.Next(); rune(in[i]) != r {
					t.Errorf("Input %q, expected %q at %v, got %q\n", in, in[i], i, r)
				} else if r != s.Curr() {
					t.Errorf("Curr() doesn't match rune\n")
				}
			}
			if s.Next() != EOF {
				t.Errorf("Expected EOF\n")
			}
			for i := len(in) - 1; i >= 0; i-- {
				if r := s.Prev(); rune(in[i]) != r {
					t.Errorf("Input %q, expected %q at %v, got %q\n", in, in[i], i, r)
				} else if r != s.Curr() {
					t.Errorf("Curr() doesn't match rune\n")
				}
			}
			if s.Prev() != BOF {
				t.Errorf("Expected BOF\n")
			}
			if r := s.Peek(); r != rune(in[0]) {
				t.Errorf("Expected %q, got %q\n", in[0], r)
			}
		}
	}
}

func TestIgnore(t *testing.T) {
	b := bytes.NewBufferString("the word")
	s := NewScanner(b)
	s.Next()
	s.Next()
	s.Next()
	s.Next()
	s.Ignore()
	if s.pos != -1 || len(s.buf) != 0 {
		t.Fail()
	}
	if r := s.Prev(); r != BOF {
		t.Errorf("Expected BOF\n")
	}
	if r := s.Peek(); r != 'w' {
		t.Errorf("Expected 'w'\n")
	}
}

func TestExtract(t *testing.T) {
	for _, in := range samples {
		words := strings.Split(in, " ")
		b := bytes.NewBufferString(in)
		s := NewScanner(b)
		for w := 0; s.Peek() != EOF; w++ {
			if s.SkipWs() {
				s.Ignore()
			}
			s.UntilWs()
			if e := s.Extract(); e != words[w] {
				t.Errorf("Extract failed, got %q %v expected %q\n", e, len(e), words[w])
			}
		}
	}
}
