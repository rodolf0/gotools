package lexer

import (
	"bufio"
	"io"
	"bytes"
	"strings"
)

const (
	EOF = -(1 + iota)
	BOF
)

// A Scanner buffers runes as the user advances the scanner.
// Once reached an interesting point the user can Accept or
// Ignore which will return the results and reset the buffer.
type Scanner struct {
	rdr *bufio.Reader
	buf []rune
	pos int
}

// Create a new scanner
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		rdr: bufio.NewReader(r),
		buf: make([]rune, 0, 32),
		pos: -1,
	}
}

// Create a new scanner
func NewScannerString(s string) *Scanner {
	return &Scanner{
		rdr: bufio.NewReader(bytes.NewBufferString(s)),
		buf: make([]rune, 0, 32),
		pos: -1,
	}
}

// Check if we've hit EOF
func (s *Scanner) EOF() bool {
	return s.pos >= len(s.buf)
}

// Check if we've hit BOF while going backwards
func (s *Scanner) BOF() bool {
	return s.pos < 0
}

// Get the current rune the scanner is on
func (s *Scanner) Curr() rune {
	if s.pos < 0 {
		return BOF
	} else if s.pos >= len(s.buf) {
		return EOF
	}
	return s.buf[s.pos]
}

// Go back and return previous rune (or BOF if at start of buffer)
func (s *Scanner) Prev() rune {
	if s.pos >= 0 {
		s.pos -= 1
	}
	return s.Curr()
}

// Rewind scanner, only up to non-accepted runes
func (s *Scanner) Rewind() {
	s.pos = -1
}

// Advance the scanner by one run and return it
func (s *Scanner) Next() rune {
	s.pos += 1
	// get more runes if we've hit the end of the buffer
	if s.pos >= len(s.buf) {
		r, _, e := s.rdr.ReadRune()
		// actually reached EOF adjust pos to equal buffer len
		if e == io.EOF {
			s.pos = len(s.buf)
			return EOF
		} else if e != nil {
			panic(e)
		}
		s.buf = append(s.buf, r)
	}
	return s.Curr()
}

// Peek the next rune without actually advancing the scanner
func (s *Scanner) Peek() rune {
	p := s.pos
	n := s.Next()
	s.pos = p
	return n
}

// Get a view of what the scanner has currently gone over
// Accept and Ignore methods will invalidate the returned slice
func (s *Scanner) View() []rune {
	if s.pos < 0 {
		return nil
	}
	n := s.pos + 1
	return s.buf[:n]
}

// Ignore drops currently scanned runes and is ready
// to start scanning from the next position
func (s *Scanner) Ignore() {
	if s.pos < 0 {
		return
	}
	n := s.pos + 1 // num chars
	if len(s.buf) >= n {
		copy(s.buf, s.buf[n:])
		s.buf = s.buf[:len(s.buf)-n]
	}
	s.pos = -1
}

// Return the known scanned slice of runes and resets
// internal buffers and positions for the next run
func (s *Scanner) Extract() string {
	r := string(s.View())
	s.Ignore()
	return r
}

// Advance the scanner if the next rune matches the
// set 'any'. Return true if advanced
func (s *Scanner) Accept(any string) bool {
	if strings.IndexRune(any, s.Peek()) >= 0 {
		s.Next()
		return true
	}
	return false
}

// Advance the scanner while input matches the 'over' set
// The scanner will be positioned on the last matching rune
// or will stay in the same place if it cannot match anything
// Will return true if the scanner advanced because of matches
func (s *Scanner) Skip(over string) bool {
	advanced := false
	for strings.IndexRune(over, s.Peek()) >= 0 {
		advanced = true
		s.Next()
	}
	return advanced
}

// Advance the scanner until we reach a rune in the 'find' set or EOF
// The scanner will be positioned on the last non-matching rune
// or will stay in the same place if the next rune is a match
// Will return true if the scanner advanced because of matches
func (s *Scanner) Until(find string) bool {
	advanced := false
	for r := s.Peek(); r != EOF &&
		strings.IndexRune(find, r) == -1; r = s.Peek() {
		advanced = true
		s.Next()
	}
	return advanced
}

// what we consider whitespace
var ws = " \n\r\t"

// Skip over white-space
func (s *Scanner) SkipWs() bool {
	return s.Skip(ws)
}

// Advance until we find whitespace
func (s *Scanner) UntilWs() bool {
	return s.Until(ws)
}
