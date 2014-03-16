package util

import (
	"bytes"
	"errors"
	"strconv"
	"bufio"
	"io"
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

// RowReader pumps rows built by reading from 'in' to 'out'
func RowReader(in io.Reader,  out chan<- Row, delim []byte, done <-chan struct{}) {
	for scanner := bufio.NewScanner(in); scanner.Scan(); {
		select {
			case out <-MakeRow(scanner.Bytes(), delim):
			case <-done:
				return
		}
	}
}

// markFields finds the indexes where 'delim' marks fields separation.
// Field indexes are 0-based
func (r *Row) markFields(n int) error {
	if len(r.delims) > n {
		return nil
	}
	// figure out where to start searching for the next delimiter
	start := 0
	if len(r.delims) > 0 {
		start = r.delims[len(r.delims)-1] + len(r.delim)
	} else {
		r.delims = make([]int, 0, n+1)
	}
	// search for as many delimiters as needed to reach one past field n
	if start+len(r.delim) < len(r.Data) {
		for len(r.delims) <= n {
			if idx := bytes.Index(r.Data[start:], r.delim); idx != -1 {
				r.delims = append(r.delims, start+idx)
				start += idx + len(r.delim)
			} else {
				r.delims = append(r.delims, len(r.Data))
				break
			}
		}
	} else {
		// next delim doesn't fit in data: we're at the end
		r.delims = append(r.delims, len(r.Data))
	}
	// error if we still have fewer delims than needed
	if len(r.delims) <= n {
		return errors.New("Field not found")
	}
	return nil
}


func (r *Row) Bytes(idx int) ([]byte, error) {
	if idx < 0 {
		return nil, errors.New("Negative index")
	}
	if err := r.markFields(idx); err != nil {
		return nil, err
	}
	if idx > 0 {
		return r.Data[r.delims[idx-1]+len(r.delim) : r.delims[idx]], nil
	}
	return r.Data[:r.delims[idx]], nil
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
