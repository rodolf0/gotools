package aggregate

import (
	"io"
	"stream"
)

func (a *Aggregation) AggregateStream(input <-chan stream.Line) {
	for line := range input {
		var fields = line.SplitFields(a.Delim)
		var key = string(stream.JoinSomeFields(a.Delim, fields, a.Keys))
		var agg, ok = a.Data[key]

		if !ok {
			agg = make([]Aggregator, len(a.AggCtor))
			for i, ctor := range a.AggCtor {
				agg[i] = ctor()
			}
			a.Data[key] = agg
		}

		for i, agtor := range agg {
			agtor.Aggregate(fields[a.Aggs[i]])
		}
	}
}

func (a *Aggregation) Print(out io.Writer) {
	if len(a.Header) > 0 {
		out.Write([]byte(a.Header[0]))
		for i := 1; i < len(a.Header); i++ {
			out.Write(a.Delim)
			out.Write([]byte(a.Header[i]))
		}
		out.Write([]byte{'\n'})
	}

	// output aggregation results
	for key, vals := range a.Data {
		if len(a.Keys) > 0 {
			out.Write([]byte(key))
		}

		if len(a.Keys) > 0 && len(a.Aggs) > 0 {
			out.Write(a.Delim)
		}

		if len(a.Aggs) > 0 {
			out.Write([]byte(vals[0].String()))
			for i := 1; i < len(a.Aggs); i++ {
				out.Write(a.Delim)
				out.Write([]byte(vals[i].String()))
			}

		}
		out.Write([]byte{'\n'})
	}
}
