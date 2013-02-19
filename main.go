package main

import (
	"column"
	"flag"
	"fmt"
	"io"
	"os"
)

var Keys = flag.String("K", "", "Aggregation key fields")
var Sums = flag.String("S", "", "Aggregation sum fields")
var Delim = flag.String("d", ",", "Field delimiter")

func init() {
	flag.Parse()
}

func main() {
	r := column.NewReader(os.Stdin, *Delim)
	header, _ := r.ReadLine()
	keys := header.ParseFields(Keys)
	//	sums := header.ParseFields(Sums)

	for ki, fi := range keys {
		fmt.Printf("key %v is at %v\n", ki, fi)
	}

	//	var aggs = make(map[string]interface{})

	var fields, err = r.ReadLine()
	for err != io.EOF {
		fmt.Println(fields)
		fields, err = r.ReadLine()
	}
}
