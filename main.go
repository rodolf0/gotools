package main

import (
	"aggregate"
	"column"
	"flag"
	"os"
)

var Delim = flag.String("d", ",", "Field delimiter")

var Keys = flag.String("k", "", "Key fields")
var Counts = flag.String("c", "", "Count fields")
var Sums = flag.String("s", "", "Sum fields")
var Averages = flag.String("a", "", "Average fields")
var Mins = flag.String("m", "", "Minimum fields")
var Maxs = flag.String("x", "", "Maximum fields")
var Firsts = flag.String("f", "", "First fields")
var Lasts = flag.String("l", "", "Last fields")
var Concats = flag.String("t", "", "Concat fields")

func main() {
	flag.Parse()

	var r = column.NewReader(os.Stdin, []byte(*Delim))
	var header, _ = r.ReadLine()
	var ks, as = aggregate.Configure(header, Keys, Counts, Sums,
		Averages, Mins, Maxs, Firsts, Lasts, Concats)

	var aggs = aggregate.Aggregate(r, ks, as)

	aggregate.String(aggs, []byte(*Delim), os.Stdout)
}
