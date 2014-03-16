package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type SortSpec struct {
	field   int
	reverse bool
	numeric bool
}

func Config(sortspec string) ([]SortSpec, error) {
	keys := strings.Split(sortspec, ",")
	re := regexp.MustCompile("^([0-9]+)([rn]*)$")
	sp := make([]SortSpec, 0, len(keys))

	for _, key := range keys {
		matches := re.FindStringSubmatch(key)
		if len(matches) < 3 {
			return nil, errors.New("Invalid key spec")
		}
		field, err := strconv.Atoi(matches[1])
		if err != nil || field < 1 {
			return nil, errors.New("Invalid field index: " + matches[1])
		}
		reverse := false
		numeric := false
		for _, mod := range matches[2] {
			switch mod {
			case 'r':
				reverse = true
			case 'n':
				numeric = true
			}
		}
		sp = append(sp, SortSpec{field-1, reverse, numeric})
	}
	return sp, nil
}
