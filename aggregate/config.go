package aggregate

import (
	"stream"
	"sync"
)

type _agg_ struct {
	d map[string][]Aggregator
	sync.Mutex
}

type Aggregation struct {
	Keys       []int
	Pivots     []int
	Aggs       []int
	AggCtor    []func() Aggregator
	Data       map[string]_agg_
	KeysHeader []string
	AggsHeader []string
	PivsHeader map[string]bool
	Delim      []byte
	SubDelim   []byte
	NullVal    []byte
}

func Configure(Keys, Pivots *string, Aggs map[string]*string, Delim, SubDelim *string, header stream.Line) *Aggregation {

	var a = &Aggregation{
		Data:       make(map[string]_agg_),
		PivsHeader: make(map[string]bool),
		Delim:      []byte(*Delim),
		SubDelim:   []byte(*SubDelim),
	}

	var IdxMap = header.IndexMap([]byte(*Delim))

	for _, k := range stream.Line(*Keys).SplitFields([]byte{','}) {
		idx, ok := IdxMap[string(k)]
		if !ok {
			panic("No column named " + string(k))
		}
		a.Keys = append(a.Keys, idx)
		a.KeysHeader = append(a.KeysHeader, string(k))
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
				a.AggsHeader = append(a.AggsHeader, string(col)+"-Cnt")
			case "Adder":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return new(Adder) })
				a.AggsHeader = append(a.AggsHeader, string(col)+"-Sum")
			case "Averager":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return &Averager{0.0, 0} })
				a.AggsHeader = append(a.AggsHeader, string(col)+"-Avg")
			case "Miner":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return &Miner{0.0, false} })
				a.AggsHeader = append(a.AggsHeader, string(col)+"-Min")
			case "Maxer":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return &Maxer{0.0, false} })
				a.AggsHeader = append(a.AggsHeader, string(col)+"-Max")
			case "Firster":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return new(Firster) })
				a.AggsHeader = append(a.AggsHeader, string(col)+"-Fst")
			case "Laster":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return new(Laster) })
				a.AggsHeader = append(a.AggsHeader, string(col)+"-Lst")
			case "Concater":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return &Concater{nil, []byte(*SubDelim)} })
				a.AggsHeader = append(a.AggsHeader, string(col)+"-Cat")
			}
		}
	}

	return a
}
