package main

import (
	"aggregate"
	"bufio"
	"flag"
	"os"
	"stream"
)

//var HasHeader = flag.Bool("H", true, "")
var Delim = flag.String("d", ",", "Field delimiter")
var SubDelim = flag.String("b", "|", "Field sub-delimiter")

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

	var lines = stream.LineGenerator(os.Stdin)
	var a = aggregate.Configure(Keys, Pivots, Aggs, Delim, SubDelim, <-lines)
	a.AggregateStream(lines)
	var out = bufio.NewWriter(os.Stdout)
	a.Print(out)
	out.Flush()
}
