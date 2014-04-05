package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"util"
)

var aggspec AggSpec
var HeaderPivots map[string]struct{}
var Aggregations map[string]map[string][]Aggregator

func init() {
	HeaderPivots = make(map[string]struct{})
	Aggregations = make(map[string]map[string][]Aggregator)
}

func main() {
	done := make(chan struct{})
	defer close(done)
	rows := util.Files2Rows(flag.Args(), Delim, done)

	var headermap map[string]int
	if !*noheader {
		if header, ok := <-rows; ok {
			headermap = util.HeaderMap(header)
		} else {
			return
		}
	}
	aggspec = Config(Keys, Pivots, Aggs, headermap)

	// Aggregate rows
	for row := range rows {
		key, err := row.JoinF(aggspec.Keys, Delim)
		if err != nil {
			panic(err)
		}
		pivot, err := row.JoinF(aggspec.Pivots, SubDelim)
		if err != nil {
			panic(err)
		}

		skey := string(key)
		pivots, p_ok := Aggregations[skey]
		// initialize pivots for this key
		if !p_ok {
			pivots = make(map[string][]Aggregator)
			Aggregations[skey] = pivots
		}
		spivot := string(pivot)
		aggs, a_ok := pivots[spivot]
		// initialize aggregators for this pivot
		if !a_ok {
			aggs = make([]Aggregator, len(aggspec.AggCtor))
			for i, ctor := range aggspec.AggCtor {
				aggs[i] = ctor()
			}
			pivots[spivot] = aggs
			HeaderPivots[spivot] = struct{}{} // collect pivots
		}

		for i, agtor := range aggs {
			field_i := aggspec.Aggs[i]
			if field, err := row.Bytes(field_i); err == nil {
				agtor.Aggregate(field)
			} else {
				panic(err)
			}
		}
	}

	var pivots []string
	for p := range HeaderPivots {
		pivots = append(pivots, p)
	}
	out := bufio.NewWriter(os.Stdout)
	if !*noheader {
		Header(out, pivots)
	}
	Print(out, pivots)
	out.Flush()
}

func Header(out io.Writer, pivots []string) {
	for i, k := range aggspec.KeysHeader {
		if i > 0 {
			out.Write(Delim)
		}
		out.Write([]byte(k))
	}
	if len(aggspec.Keys) > 0 && len(aggspec.Aggs) > 0 {
		out.Write(Delim)
	}
	// if no aggs we're just condensing keys
	if len(aggspec.Aggs) > 0 {
		if len(aggspec.Pivots) > 0 {
			// one aggregation per pivot value
			for i, pivot := range pivots {
				if i > 0 {
					out.Write(Delim)
				}
				for j, k := range aggspec.AggsHeader {
					if j > 0 {
						out.Write(Delim)
					}
					out.Write([]byte(pivot + ":" + k))
				}
			}
		} else {
			// simple aggregation, no pivots
			for i, k := range aggspec.AggsHeader {
				if i > 0 {
					out.Write(Delim)
				}
				out.Write([]byte(k))
			}
		}
	}

	out.Write([]byte{'\n'})
}

func Print(out io.Writer, pivots []string) {
	for key, agg := range Aggregations {
		if len(aggspec.Keys) > 0 {
			out.Write([]byte(key))
		}
		if len(aggspec.Keys) > 0 && len(aggspec.Aggs) > 0 {
			out.Write(Delim)
		}
		if len(aggspec.Aggs) > 0 {
			for i, pivot := range pivots {
				if i > 0 {
					out.Write(Delim)
				}
				if pivaggs, ok := agg[pivot]; ok {
					for j, val := range pivaggs {
						if j > 0 {
							out.Write(Delim)
						}
						out.Write([]byte(val.String()))
					}
				} else {
					for j := 0; j < len(aggspec.Aggs); j++ {
						if j > 0 {
							out.Write(Delim)
						}
						// out.Write(NullVal)
					}
				}
			}
		}
		out.Write([]byte{'\n'})
	}
}
