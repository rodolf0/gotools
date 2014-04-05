package util

import (
	"bufio"
	"os"
)

func OpenFiles(files []string, done <-chan struct{}) <-chan *os.File {
	out := make(chan *os.File)
	go func() {
		defer close(out)
		for _, file := range files {
			in := os.Stdin
			if file != "-" {
				var err error
				if in, err = os.Open(file); err != nil {
					panic(err)
				}
			}
			select {
			case out <- in:
			case <-done:
				return
			}
		}
	}()
	return out
}

func Files2Rows(files []string, delim []byte, done <-chan struct{}) <-chan Row {
	out := make(chan Row)
	go func() {
		defer close(out)
		for _, file := range files {
			in := os.Stdin
			if file != "-" {
				var err error
				if in, err = os.Open(file); err != nil {
					panic(err)
				}
			}
			for scanner := bufio.NewScanner(in); scanner.Scan(); {
				select {
				case out <- MakeRow(scanner.Bytes(), delim):
				case <-done:
					return
				}
			}
		}
	}()
	return out
}

// HeaderMap returns a map field-name -> field-index
func HeaderMap(header Row) map[string]int {
	hm := make(map[string]int)
	for i := 0; ; i++ {
		if name, err := header.Bytes(i); err == nil {
			hm[string(name)] = i
		} else {
			return hm
		}
	}
}
