package column

import (
	"bytes"
)

// FieldMap returns a field-content to index map
func (r Row) FieldMap() map[string]uint {
	var fm = make(map[string]uint)
	for i, f := range r {
		fm[string(f)] = uint(i)
	}
	return fm
}

// Indexes returns a slice with the index of each field or -1 if non present
func (r Row) Indexes(fields [][]byte) (indexes []int) {
	var fmap = r.FieldMap()
	for _, f := range fields {
		if i, ok := fmap[string(f)]; ok {
			indexes = append(indexes, int(i))
		} else {
			indexes = append(indexes, -1)
		}
	}
	return
}

// Fields returns a slice with the fields requested by indexes or nil
func (r Row) Fields(indexes []uint) (fields [][]byte) {
	var _r = [][]byte(r)
	for _, i := range indexes {
		if int(i) >= len(_r) {
			fields = append(fields, nil)
		} else {
			fields = append(fields, _r[i])
		}
	}
	return
}

// ParseFields returns a list of indexes matching the field names in s
func (r Row) ParseFields(s *string) (indexes []int) {
	if len(*s) > 0 {
		return r.Indexes(bytes.Split([]byte(*s), []byte{','}))
	}
	return nil
}
