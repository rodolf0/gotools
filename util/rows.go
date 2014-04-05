package util

import (
	"bytes"
	"errors"
	"strconv"
)

var (
	ErrorFieldNotFound = errors.New("Field not found")
	ErrorNegativeIndex = errors.New("Negative Index")
)

type Row struct {
	Data   []byte
	delims []int
	delim  []byte
}

func MakeRow(data, delim []byte) (r Row) {
	r.Data = make([]byte, len(data))
	r.delim = delim
	copy(r.Data, data)
	return
}

func (r *Row) markdelims() {
	start := 0
	for start+len(r.delim) <= len(r.Data) {
		idx := bytes.Index(r.Data[start:], r.delim)
		if idx == -1 {
			break
		}
		r.delims = append(r.delims, start+idx)
		start += idx + len(r.delim)
	}
	r.delims = append(r.delims, len(r.Data))
}

// Bytes returns slices to the underlying row data without copying
func (r *Row) Bytes(idx int) ([]byte, error) {
	if idx < 0 {
		return nil, ErrorNegativeIndex
	}
	if len(r.delims) == 0 {
		r.markdelims()
	}
	if idx >= len(r.delims) {
		return nil, ErrorFieldNotFound
	}
	start := 0
	if idx > 0 {
		start = r.delims[idx-1] + len(r.delim)
	}
	return r.Data[start:r.delims[idx]], nil
}

func (r *Row) String(index int) (string, error) {
	field, err := r.Bytes(index)
	if err == nil {
		return string(field), nil
	}
	return "", err
}

func (r *Row) Int(index int) (int, error) {
	field, err := r.Bytes(index)
	if err == nil {
		return strconv.Atoi(string(field))
	}
	return 0, err
}

func (r *Row) JoinF(fields []int, delim []byte) ([]byte, error) {
	if len(fields) == 0 {
		return nil, nil
	}
	j := make([]byte, 0, len(r.Data)/2)
	for i, f := range fields {
		if field, err := r.Bytes(f); err == nil {
			if i > 0 && len(delim) > 0 {
				j = append(j, delim...)
			}
			j = append(j, field...)
		} else {
			return nil, err
		}
	}
	return j, nil
}
