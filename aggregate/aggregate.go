package aggregate

import (
	"column"
	"os"
)

type AggSpec struct {
	Field   int
	AggCtor func() Aggregator
}

func Aggregate(input *column.Reader, key_fields []int, agg_fields []AggSpec) {
	var aggregations = make(map[string][]Aggregator)

	var line, err = input.ReadLine()
	for err == nil {
		// build the key for the current line
		var key = string(line.JoinFields(key_fields, input.Delim))
		// instantiate aggregators if this is a new key
		if _, ok := aggregations[key]; !ok {
			for _, af := range agg_fields {
				aggregations[key] = append(aggregations[key], af.AggCtor())
			}
		}
		// feed current line to aggregators
		for i, agg := range aggregations[key] {
			agg.Aggregate(line[agg_fields[i].Field])
		}
		line, err = input.ReadLine()
	}

	// output aggregation results
	for key, aggs := range aggregations {
		os.Stdout.Write([]byte(key))

		for _, agg := range aggs {
			os.Stderr.Write(input.Delim)
			os.Stderr.Write([]byte(agg.String()))
		}
		os.Stderr.Write([]byte("\n"))
	}
}
