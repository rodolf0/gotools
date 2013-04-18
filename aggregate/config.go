package aggregate

import (
	"stream"
)

type Aggregation struct {
	Keys    []int
	Pivots  []int
	Aggs    []int
	AggCtor []func() Aggregator
	Data    map[string][]Aggregator
	Header  []string
	Delim   []byte
}

func Configure(Keys, Pivots *string, Aggs map[string]*string, Delim, SubDelim *string, header stream.Line) *Aggregation {

	var a = &Aggregation{
		Data:  make(map[string][]Aggregator),
		Delim: []byte(*Delim),
	}

	var IdxMap = header.IndexMap([]byte(*Delim))

	for _, k := range stream.Line(*Keys).SplitFields([]byte{','}) {
		idx, ok := IdxMap[string(k)]
		if !ok {
			panic("No column named " + string(k))
		}
		a.Keys = append(a.Keys, idx)
		a.Header = append(a.Header, string(k))
	}

	for _, p := range stream.Line(*Pivots).SplitFields([]byte{','}) {
		idx, ok := IdxMap[string(p)]
		if !ok {
			panic("No column named " + string(p))
		}
		a.Pivots = append(a.Pivots, idx)
	}

	for agg_type, agg_cols := range Aggs {
		for _, col := range stream.Line(*agg_cols).SplitFields([]byte{','}) {

			idx, ok := IdxMap[string(col)]
			if !ok {
				panic("No column named " + string(col))
			}
			a.Aggs = append(a.Aggs, idx)

			switch agg_type {
			case "Counter":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return new(Counter) })
				a.Header = append(a.Header, string(col)+"-Cnt")
			case "Adder":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return new(Adder) })
				a.Header = append(a.Header, string(col)+"-Sum")
			case "Averager":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return &Averager{0.0, 0} })
				a.Header = append(a.Header, string(col)+"-Avg")
			case "Miner":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return &Miner{0.0, false} })
				a.Header = append(a.Header, string(col)+"-Min")
			case "Maxer":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return &Maxer{0.0, false} })
				a.Header = append(a.Header, string(col)+"-Max")
			case "Firster":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return new(Firster) })
				a.Header = append(a.Header, string(col)+"-Fst")
			case "Laster":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return new(Laster) })
				a.Header = append(a.Header, string(col)+"-Lst")
			case "Concater":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return &Concater{nil, []byte(*SubDelim)} })
				a.Header = append(a.Header, string(col)+"-Cat")
			}
		}
	}

	return a
}
