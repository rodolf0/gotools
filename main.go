package main

import (
	"aggregate"
	"bytes"
	"column"
	"flag"
	"os"
)

var Delim = flag.String("d", ",", "Field delimiter")
var SubDelim = flag.String("b", ",", "Field delimiter")

var Keys = flag.String("k", "", "Key fields")
var Counts = flag.String("c", "", "Count fields")
var Sums = flag.String("s", "", "Sum fields")
var Averages = flag.String("a", "", "Average fields")
var Mins = flag.String("n", "", "Minimum fields")
var Maxs = flag.String("x", "", "Maximum fields")
var Firsts = flag.String("f", "", "First fields")
var Lasts = flag.String("l", "", "Last fields")
var Concats = flag.String("t", "", "Concat fields")
var Pivots = flag.String("p", "", "Pivot fields")

func main() {
	flag.Parse()

	var r = column.NewReader(os.Stdin, []byte(*Delim))
	var header, _ = r.ReadLine()

	var ks, as = aggregate.Configure(header,
		bytes.Split([]byte(*Keys), []byte{','}),
		bytes.Split([]byte(*Counts), []byte{','}),
		bytes.Split([]byte(*Sums), []byte{','}),
		bytes.Split([]byte(*Averages), []byte{','}),
		bytes.Split([]byte(*Mins), []byte{','}),
		bytes.Split([]byte(*Maxs), []byte{','}),
		bytes.Split([]byte(*Firsts), []byte{','}),
		bytes.Split([]byte(*Lasts), []byte{','}),
		bytes.Split([]byte(*Concats), []byte{','}))

	var aggs = aggregate.Aggregate(r, ks, as)

	aggregate.String(aggs, []byte(*Delim), os.Stdout)
}
