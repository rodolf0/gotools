package main

import (
	"column"
	"fmt"
	"io"
	"os"
)

func main() {
	r := column.NewReader(os.Stdin, ',')
	for {
		if cols, err := r.ReadLine(); err != io.EOF {
			fmt.Printf("%v, %v\n", string(cols[0]), string(cols[1]))
		} else {
			break
		}
	}
}
