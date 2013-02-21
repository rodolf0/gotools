package main

import (
	"column"
	"flag"
	"io"
	"os"
)

var Delim = flag.String("d", ",", "Field delimiter")

var Keys = flag.String("K", "", "Aggregation key fields")
var Sums = flag.String("S", "", "Aggregation sum fields")
var Average = flag.String("A", "", "Aggregation average fields")
var Count = flag.String("C", "", "Aggregation count fields")

type AggField struct {
	field   int
	AggCtor func() Aggregator
}

func main() {
	flag.Parse()

	var r = column.NewReader(os.Stdin, []byte(*Delim))
	var header, _ = r.ReadLine()
	var key_fields = header.ParseFields(Keys)
	var agg_fields []AggField

	for _, field := range header.ParseFields(Sums) {
		agg_fields = append(agg_fields, AggField{field, NewAdder})
	}
	for _, field := range header.ParseFields(Average) {
		agg_fields = append(agg_fields, AggField{field, NewAverager})
	}
	for _, field := range header.ParseFields(Count) {
		agg_fields = append(agg_fields, AggField{field, NewCounter})
	}

	var aggregations = make(map[string][]Aggregator)

	var line, err = r.ReadLine()
	for err != io.EOF {
		// build the key for the current line
		var key = string(line.JoinFields(key_fields, []byte(*Delim)))
		// instantiate aggregators if this is a new key
		if _, ok := aggregations[key]; !ok {
			for _, af := range agg_fields {
				aggregations[key] = append(aggregations[key], af.AggCtor())
			}
		}
		// feed current line to aggregators
		for i, agg := range aggregations[key] {
			agg.Aggregate(line[agg_fields[i].field])
		}
		line, err = r.ReadLine()
	}

	// output aggregation results
	for key, aggs := range aggregations {
		os.Stdout.Write([]byte(key))

		for _, agg := range aggs {
			os.Stderr.Write([]byte(*Delim))
			os.Stderr.Write([]byte(agg.String()))
		}
		os.Stderr.Write([]byte("\n"))
	}
}
