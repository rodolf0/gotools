package aggregate

import (
	"stream"
)

type AggSpec struct {
	Keys    []int
	Pivots  []int
	Aggs    []int
	AggCtor []func() Aggregator
}

func Configure(Keys, Pivots *string, Aggs map[string]*string, IdxMap map[string]int, SubDelim *string) AggSpec {

	var aggspec AggSpec

	for _, k := range stream.Line(*Keys).SplitFields([]byte{','}) {
		aggspec.Keys = append(aggspec.Keys, IdxMap[string(k)])
	}

	for _, p := range stream.Line(*Pivots).SplitFields([]byte{','}) {
		aggspec.Pivots = append(aggspec.Pivots, IdxMap[string(p)])
	}

	for agg_type, agg_cols := range Aggs {
		for _, col := range stream.Line(*agg_cols).SplitFields([]byte{','}) {

			aggspec.Aggs = append(aggspec.Aggs, IdxMap[string(col)])
			switch agg_type {
			case "Counter":
				aggspec.AggCtor = append(aggspec.AggCtor, func() Aggregator { return new(Counter) })
			case "Adder":
				aggspec.AggCtor = append(aggspec.AggCtor, func() Aggregator { return new(Adder) })
			case "Averager":
				aggspec.AggCtor = append(aggspec.AggCtor, func() Aggregator { return &Averager{0.0, 0} })
			case "Miner":
				aggspec.AggCtor = append(aggspec.AggCtor, func() Aggregator { return &Miner{0.0, false} })
			case "Maxer":
				aggspec.AggCtor = append(aggspec.AggCtor, func() Aggregator { return &Maxer{0.0, false} })
			case "Firster":
				aggspec.AggCtor = append(aggspec.AggCtor, func() Aggregator { return new(Firster) })
			case "Laster":
				aggspec.AggCtor = append(aggspec.AggCtor, func() Aggregator { return new(Laster) })
			case "Concater":
				aggspec.AggCtor = append(aggspec.AggCtor, func() Aggregator { return &Concater{nil, []byte(*SubDelim)} })
			}
		}
	}

	return aggspec
}
