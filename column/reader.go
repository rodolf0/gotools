package column

import (
	"bufio"
	"bytes"
	"io"
)

type Field []byte
type Row []Field

type Reader struct {
	rd    *bufio.Reader
	delim byte
}

func NewReader(r io.Reader, delim byte) *Reader {
	return &Reader{
		rd:    bufio.NewReader(r),
		delim: delim,
	}
}

// Index returns the index of the field in Row containing field 'name'
func (r *Row) Index(name Field) int {
	for i, n := range []Field(*r) {
		if bytes.Equal(n, name) {
			return i
		}
	}
	return -1
}

// Indexes returns the a slice of uint's mapping row values to
// the indexes in the row that contain those values
func (r *Row) Indexes(names []Field) []uint {
	var indexes []uint
	for i, name := range names {
		if r.Index(name) != -1 {
			indexes = append(indexes, uint(i))
		} else {
			return nil
		}
	}
	return indexes
}

// ReadLine reads a line, finds fields boundaries and returns a slice
// of columns or nil and an error.
func (cr *Reader) ReadLine() (Row, error) {
	line, err := cr.rd.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	// remove end-of-line bytes
	if len(line) > 1 && line[len(line)-2] == '\r' {
		line = line[:len(line)-2]
	} else {
		line = line[:len(line)-1]
	}
	// find field boundaries
	var columns []Field
	for {
		if d := bytes.IndexByte(line, cr.delim); d >= 0 {
			columns = append(columns, line[:d])
			line = line[d+1:]
		} else {
			columns = append(columns, line)
			break
		}
	}
	return columns, nil
}
