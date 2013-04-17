package aggregate

import (
	"io"
	"stream"
)

func Aggregate(input *stream.Reader, delim []byte, aggspec AggSpec) map[string][]Aggregator {
	var aggregations = make(map[string][]Aggregator)

	var line, err = input.ReadLine()
	for err == nil {
		var fields = line.SplitFields(delim)
		var key = string(stream.JoinSomeFields(delim, fields, aggspec.Keys))
		var agg, ok = aggregations[key]

		if !ok {
			agg = make([]Aggregator, len(aggspec.AggCtor))
			for i, ctor := range aggspec.AggCtor {
				agg[i] = ctor()
			}
			aggregations[key] = agg
		}

		for i, a := range agg {
			a.Aggregate(fields[aggspec.Aggs[i]])
		}
		line, err = input.ReadLine()
	}

	return aggregations
}

func Aggregate2(input <-chan stream.Line, delim []byte, aggspec AggSpec) map[string][]Aggregator {
	var aggregations = make(map[string][]Aggregator)

	for line := range input {
		var fields = line.SplitFields(delim)
		var key = string(stream.JoinSomeFields(delim, fields, aggspec.Keys))
		var agg, ok = aggregations[key]

		if !ok {
			agg = make([]Aggregator, len(aggspec.AggCtor))
			for i, ctor := range aggspec.AggCtor {
				agg[i] = ctor()
			}
			aggregations[key] = agg
		}

		for i, a := range agg {
			a.Aggregate(fields[aggspec.Aggs[i]])
		}
	}

	return aggregations
}

func String(aggregations map[string][]Aggregator, delim []byte, out io.Writer) {
	// output aggregation results
	for key, vals := range aggregations {
		out.Write([]byte(key))

		for _, agg := range vals {
			out.Write(delim)
			out.Write([]byte(agg.String()))
		}
		out.Write([]byte{'\n'})
	}
}
