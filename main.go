package main

import (
	"aggregate"
	"flag"
	"os"
	"stream"
)

var Delim = flag.String("d", ",", "Field delimiter")
var SubDelim = flag.String("b", "|", "Field delimiter")

var Keys = flag.String("k", "", "Key fields")
var Pivots = flag.String("p", "", "Pivot fields")

var Aggs = map[string]*string{
	"Counter":  flag.String("c", "", "Count fields"),
	"Adder":    flag.String("s", "", "Sum fields"),
	"Averager": flag.String("a", "", "Average fields"),
	"Miner":    flag.String("n", "", "Minimum fields"),
	"Maxer":    flag.String("x", "", "Maximum fields"),
	"Firster":  flag.String("f", "", "First fields"),
	"Laster":   flag.String("l", "", "Last fields"),
	"Concater": flag.String("t", "", "Concat fields"),
}

func main() {
	flag.Parse()

	var reader = stream.NewReader(os.Stdin)
	var header, _ = reader.ReadLine()

	var idxmap = header.IndexMap([]byte(*Delim))

	var aggspec = aggregate.Configure(Keys, Pivots, Aggs, idxmap, SubDelim)
	var aggs = aggregate.Aggregate(reader, []byte(*Delim), aggspec)

	aggregate.String(aggs, []byte(*Delim), os.Stdout)
}
