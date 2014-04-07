package main

import (
	"bytes"
	"flag"
	"os"
	"sort"
	"util"
)

var Delim = flag.String("d", ",", "Field delimiter")
var Keys = flag.String("k", "", "Sort keys (eg: 2,1r,3rn,4n)")
var WinSz = flag.Int("s", 2048, "Window size")
var delim []byte
var sspec []SortSpec
var sortwin []*util.Row

func init() {
	flag.Parse()
	delim = []byte(*Delim)
	sortwin = make([]*util.Row, 0, *WinSz)
	var err error
	sspec, err = Config(*Keys)
	if err != nil {
		panic(err)
	}
}

func rowcmp(a, b *util.Row) int {
	for _, key := range sspec {
		var c int
		if key.numeric {
			fa, _ := a.Int(key.field)
			fb, _ := b.Int(key.field)
			c = fa - fb
		} else {
			fa, _ := a.Bytes(key.field)
			fb, _ := b.Bytes(key.field)
			c = bytes.Compare(fa, fb)
		}
		switch {
		case c < 0:
			if key.reverse {
				return 1
			}
			return -1
		case c > 0:
			if key.reverse {
				return -1
			}
			return 1
		}
	}
	return 0
}

func semisort(row *util.Row) {
	// find insertion point (put lowest last to trim window)
	i := sort.Search(len(sortwin), func(j int) bool {
		return rowcmp(row, sortwin[j]) > 0
	})
	// add row at the end... use as side effect to grow slice for 'copy'
	sortwin = append(sortwin, row)
	if i < len(sortwin) {
		copy(sortwin[i+1:], sortwin[i:])
		sortwin[i] = row
	}
}

func main() {
	done := make(chan struct{})
	defer close(done)
	files := flag.Args()
	if len(files) == 0 {
		files = make([]string, 1)
		files[0] = "-"
	}
	rows := util.Files2Rows(files, delim, done)

	// fill sort-window
	for len(sortwin) < cap(sortwin) {
		if row, ok := <-rows; ok {
			semisort(&row)
		} else {
			break
		}
	}
	// sort-n-flush window
	var flushrow *util.Row
	for len(sortwin) > 0 {
		sortwin, flushrow = sortwin[:len(sortwin)-1], sortwin[len(sortwin)-1]
		os.Stdout.Write(flushrow.Data)
		os.Stdout.Write([]byte("\n"))
		if row, ok := <-rows; ok {
			semisort(&row)
		}
	}
}
