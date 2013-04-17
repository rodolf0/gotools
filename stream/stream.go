package stream

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

type Reader struct {
	r *bufio.Reader
	m sync.Mutex
}

type Line []byte
type Field []byte

func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

func LineGenerator(r io.Reader) <-chan Line {
	var ch = make(chan Line, 64)
	var br = bufio.NewReader(r)
	go func() {
		line, err := br.ReadBytes('\n')
		for err != nil {
			if len(line) > 1 && line[len(line)-2] == '\r' {
				ch <- Line(line[:len(line)-2])
			} else {
				ch <- Line(line[:len(line)-1])
			}
			line, err = br.ReadBytes('\n')
		}
		ch <- Line(line)
		close(ch)
	}()
	return ch
}

// ReadLine is safe for concurrent reading. It returns a line without final \r?\n.
// on error it returns what has being read until finding the error and the error.
func (r *Reader) ReadLine() (Line, error) {
	r.m.Lock()
	line, err := r.r.ReadBytes('\n')
	r.m.Unlock()
	if err != nil {
		return Line(line), err
	} else if len(line) > 1 && line[len(line)-2] == '\r' {
		return Line(line[:len(line)-2]), nil
	}
	return Line(line[:len(line)-1]), nil
}

// SplitFields returns a slice of Fields which are a view into the line.
func (l Line) SplitFields(delim []byte) []Field {
	if len(l) == 0 {
		return nil
	}
	var fstart = l
	var fields = make([]Field, 0, 16)
	for {
		if d := bytes.Index(fstart, delim); d >= 0 {
			fields = append(fields, Field(fstart[:d]))
			fstart = fstart[d+len(delim):]
		} else {
			break
		}
	}
	return append(fields, Field(fstart))
}

// JoinFields generates a Line from a slice of fields
func JoinFields(delim []byte, fields []Field) Line {
	if len(fields) == 0 {
		return nil
	}
	var n = 0
	for _, field := range fields {
		n += len(field)
	}
	var line = make([]byte, 0, n+len(delim)*(len(fields)-1))
	n = copy(line, fields[0])
	for _, field := range fields {
		n += copy(line[n:], delim)
		n += copy(line[n:], field)
	}
	return line
}

// JoinSomeFields generates a Line from a slice of fields
func JoinSomeFields(delim []byte, fields []Field, which []int) Line {
	if len(fields) == 0 || len(which) == 0 {
		return nil
	}
	var n = 0
	for _, f := range which {
		n += len(fields[f])
	}
	var line = make([]byte, n+len(delim)*(len(which)-1))
	n = copy(line, fields[which[0]])
	for _, f := range which[1:] {
		n += copy(line[n:], delim)
		n += copy(line[n:], fields[f])
	}
	return line
}

// IndexMap returns a mapping from field content to index
func (l Line) IndexMap(delim []byte) map[string]int {
	if len(l) == 0 {
		return nil
	}
	var fields = l.SplitFields(delim)
	var idxmap = make(map[string]int, len(fields))
	for i, f := range fields {
		idxmap[string(f)] = i
	}
	return idxmap
}
