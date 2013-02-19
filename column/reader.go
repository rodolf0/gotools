package column

import (
	"bufio"
	"bytes"
	"io"
)

type Row [][]byte

type Reader struct {
	rd    *bufio.Reader
	delim []byte
}

func NewReader(r io.Reader, delim string) *Reader {
	return &Reader{
		rd:    bufio.NewReader(r),
		delim: []byte(delim),
	}
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
	return bytes.Split(line, cr.delim), nil
}
