package column

import (
	"bufio"
	"bytes"
	"io"
)

type Row struct {
	Delim []byte
	Line  []byte
}

type SplitRow struct {
	Row    *Row
	Fields [][]byte
}

type Reader struct {
	rd    *bufio.Reader
	Delim []byte
}

func NewReader(r io.Reader, delim []byte) *Reader {
	var nr = new(Reader)
	nr.rd = bufio.NewReader(r) // auto-detects if already buffered
	nr.Delim = make([]byte, len(delim))
	copy(nr.Delim, delim)
	return nr
}

// ReadLine reads a line, finds fields boundaries and returns a slice
// of columns or nil and an error.
func (cr *Reader) ReadLine() (Row, error) {
	line, err := cr.rd.ReadBytes('\n')
	if err != nil {
		return Row{cr.Delim, line}, err
	}
	// remove end-of-line bytes
	if len(line) > 1 && line[len(line)-2] == '\r' {
		line = line[:len(line)-2]
	} else {
		line = line[:len(line)-1]
	}
	return Row{cr.Delim, line}, nil
}

// Split returns slices to the underlying fields of the row. The
// underlying data is only valid until a new line is read.
func (r *Row) Split(delim []byte) SplitRow {
	var fstart = r.Line
	if delim == nil {
		delim = r.Delim
	}
	var fields = make([][]byte, 0, 16)
	for {
		if s := bytes.Index(fstart, delim); s >= 0 {
			fields = append(fields, fstart[:s])
			fstart = fstart[s+len(delim):]
		} else {
			break
		}
	}
	return SplitRow{r, append(fields, fstart)}
}

// ParseFields returns a the indexes corresponding to the field names
// from the row which the function is called from. It is useful for
// indexing columns by name.
func (r Row) ParseFields(fields [][]byte) (indexes []int) {
	var sr = r.Split(r.Delim)
	for _, rf := range fields {
		for ai, af := range sr.Fields {
			if bytes.Equal(rf, af) {
				indexes = append(indexes, ai)
			}
		}
	}
	return
}

func (sr SplitRow) JoinFields(indexes []int) Row {
	var line = make([]byte, 0, len(sr.Row.Line)) // estimate size to avoid allocations
	line = append(line, sr.Fields[indexes[0]]...)
	for _, i := range indexes[1:] {
		line = append(line, sr.Row.Delim...)
		line = append(line, sr.Fields[i]...)
	}
	return Row{sr.Row.Delim, line}
}
