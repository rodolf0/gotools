package main

import (
	"bytes"
	"flag"
	"os"
	"sort"
	"util"
)

var Delim = flag.String("d", ",", "Field delimiter")
var Keys = flag.String("k", "", "Sort keys (eg: 2,1r,3rn,4n")
var WinSz = flag.Int("s", 2048, "Window size")
var delim []byte
var sspec []SortSpec

func init() {
	flag.Parse()
	delim = []byte(*Delim)
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

func semisort(row *util.Row, window *[]*util.Row) {
	// find insertion point (put lowest last to trim window)
	i := sort.Search(len(*window), func(j int) bool {
		return rowcmp(row, (*window)[j]) > 0
	})
	// add row at the end... use as side effect to grow slice for 'copy'
	*window = append(*window, row)
	if i < len(*window) {
		copy((*window)[i+1:], (*window)[i:])
		(*window)[i] = row
	}
}

func main() {
	win := make([]*util.Row, 0, *WinSz)
	done := make(chan struct{})
	defer close(done)

	rows := util.Files2Rows(flag.Args(), delim, done)

	// fill sort-window
	for len(win) < cap(win) {
		if row, ok := <-rows; !ok {
			break
		} else {
			semisort(&row, &win)
		}
	}
	// sort-n-flush window
	var flushrow *util.Row
	for len(win) > 0 {
		win, flushrow = win[:len(win)-1], win[len(win)-1]
		os.Stdout.Write(flushrow.Data)
		os.Stdout.Write([]byte("\n"))
		if row, ok := <-rows; ok {
			semisort(&row, &win)
		}
	}
}
