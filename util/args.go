package util

import (
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

//func parse_fields()
