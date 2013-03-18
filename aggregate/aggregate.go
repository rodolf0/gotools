package aggregate

import (
	"column"
	"io"
)

type AggSpec struct {
	Field   int
	AggCtor func() Aggregator
}

func Configure(header column.Row, Keys, Counts, Sums, Averages, Mins,
	Maxs, Firsts, Lasts, Concats [][]byte) ([]int, []AggSpec) {

	var aggs []AggSpec

	for _, field := range header.ParseFields(Counts) {
		aggs = append(aggs, AggSpec{field,
			func() Aggregator { return new(Counter) }})
	}
	for _, field := range header.ParseFields(Sums) {
		aggs = append(aggs, AggSpec{field,
			func() Aggregator { return new(Adder) }})
	}
	for _, field := range header.ParseFields(Averages) {
		aggs = append(aggs, AggSpec{field,
			func() Aggregator { return new(Averager) }})
	}
	for _, field := range header.ParseFields(Mins) {
		aggs = append(aggs, AggSpec{field,
			func() Aggregator { return new(Miner) }})
	}
	for _, field := range header.ParseFields(Maxs) {
		aggs = append(aggs, AggSpec{field,
			func() Aggregator { return new(Maxer) }})
	}
	for _, field := range header.ParseFields(Firsts) {
		aggs = append(aggs, AggSpec{field,
			func() Aggregator { return new(Firster) }})
	}
	for _, field := range header.ParseFields(Lasts) {
		aggs = append(aggs, AggSpec{field,
			func() Aggregator { return new(Laster) }})
	}
	for _, field := range header.ParseFields(Concats) {
		aggs = append(aggs, AggSpec{field,
			func() Aggregator { return new(Concater) }})
	}
	return header.ParseFields(Keys), aggs
}

func Aggregate(input *column.Reader, keys []int, aggs []AggSpec) map[string][]Aggregator {
	var aggregations = make(map[string][]Aggregator)

	var line, err = input.ReadLine()
	for err == nil {
		var fields = line.Split(nil)
		var key = string(fields.JoinFields(keys).Line)
		var agg, ok = aggregations[key]

		if !ok {
			agg = make([]Aggregator, 0, len(aggs))
			for _, a := range aggs {
				agg = append(agg, a.AggCtor())
			}
			aggregations[key] = agg
		}

		for i, a := range agg {
			a.Aggregate(fields.Fields[aggs[i].Field])
		}
		line, err = input.ReadLine()
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
		out.Write([]byte("\n"))
	}
}
