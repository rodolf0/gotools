package main

import (
	"flag"
	"strconv"
	"strings"
)

var delim = flag.String("d", ",", "Field delimiter")
var subdelim = flag.String("b", "|", "Field sub-delimiter")
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
var noheader = flag.Bool("H", false, "No header row")

var Delim []byte
var SubDelim []byte

func init() {
	flag.Parse()
	Delim = []byte(*delim)
	SubDelim = []byte(*subdelim)
}

type AggSpec struct {
	Keys       []int
	Pivots     []int
	Aggs       []int
	KeysHeader []string
	AggsHeader []string
	AggCtor    []func() Aggregator
}

func Config(keys, pivots *string, aggs map[string]*string, headermap map[string]int) AggSpec {
	var a AggSpec

	if len(*keys) > 0 {
		for _, k := range strings.Split(*keys, ",") {
			if i, ok := headermap[k]; ok {
				a.Keys = append(a.Keys, i)
			} else if i, err := strconv.Atoi(k); err == nil {
				a.Keys = append(a.Keys, i)
			} else {
				panic("Invalid key: " + k)
			}
			a.KeysHeader = append(a.KeysHeader, k)
		}
	}

	if len(*pivots) > 0 {
		for _, k := range strings.Split(*pivots, ",") {
			if i, ok := headermap[k]; ok {
				a.Pivots = append(a.Pivots, i)
			} else if i, err := strconv.Atoi(k); err == nil {
				a.Pivots = append(a.Pivots, i)
			} else {
				panic("Invalid pivot: " + k)
			}
		}
	}

	for aggtype, agg := range aggs {
		if len(*agg) == 0 {
			continue
		}
		for _, k := range strings.Split(*agg, ",") {
			if i, ok := headermap[k]; ok {
				a.Aggs = append(a.Aggs, i)
			} else if i, err := strconv.Atoi(k); err == nil {
				a.Aggs = append(a.Aggs, i)
			} else {
				panic("Invalid agg field: " + k)
			}
			// TODO: reduce to one function call only
			switch aggtype {
			case "Counter":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return new(Counter) })
				a.AggsHeader = append(a.AggsHeader, k+"-Cnt")
			case "Adder":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return new(Adder) })
				a.AggsHeader = append(a.AggsHeader, k+"-Sum")
			case "Averager":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return &Averager{0.0, 0} })
				a.AggsHeader = append(a.AggsHeader, k+"-Avg")
			case "Miner":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return &Miner{0.0, false} })
				a.AggsHeader = append(a.AggsHeader, k+"-Min")
			case "Maxer":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return &Maxer{0.0, false} })
				a.AggsHeader = append(a.AggsHeader, k+"-Max")
			case "Firster":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return new(Firster) })
				a.AggsHeader = append(a.AggsHeader, k+"-Fst")
			case "Laster":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return new(Laster) })
				a.AggsHeader = append(a.AggsHeader, k+"-Lst")
			case "Concater":
				a.AggCtor = append(a.AggCtor, func() Aggregator { return &Concater{nil, SubDelim} })
				a.AggsHeader = append(a.AggsHeader, k+"-Cat")
			}
		}
	}

	return a
}
