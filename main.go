package main

import (
	"aggregate"
	"column"
	"flag"
	"os"
)

var Delim = flag.String("d", ",", "Field delimiter")

var Keys = flag.String("K", "", "Aggregation key fields")
var Sums = flag.String("S", "", "Aggregation sum fields")
var Average = flag.String("A", "", "Aggregation average fields")
var Count = flag.String("C", "", "Aggregation count fields")

func main() {
	flag.Parse()

	var r = column.NewReader(os.Stdin, []byte(*Delim))
	var header, _ = r.ReadLine()
	var key_fields = header.ParseFields(Keys)
	var agg_fields []aggregate.AggSpec

	for _, field := range header.ParseFields(Sums) {
		agg_fields = append(agg_fields, aggregate.AggSpec{field,
			func() aggregate.Aggregator { return new(aggregate.Adder) }})
	}
	for _, field := range header.ParseFields(Average) {
		agg_fields = append(agg_fields, aggregate.AggSpec{field,
			func() aggregate.Aggregator { return new(aggregate.Averager) }})
	}
	for _, field := range header.ParseFields(Count) {
		agg_fields = append(agg_fields, aggregate.AggSpec{field,
			func() aggregate.Aggregator { return new(aggregate.Counter) }})
	}

	aggregate.Aggregate(r, key_fields, agg_fields)
}
