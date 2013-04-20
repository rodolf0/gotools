package aggregate

import (
	"io"
	"stream"
	"sync"
)

func (a *Aggregation) AggregateStream(input <-chan stream.Line) {

	var wg2 sync.WaitGroup

	const partition = 4

	wg2.Add(partition)

	var aggmux [partition]chan []stream.Field

	for g := 0; g < partition; g++ {
		aggmux[g] = make(chan []stream.Field, 1024)
	}

	var m sync.Mutex

	for g := 0; g < partition; g++ {
		go func(am_i int) {
			var localagg = make(map[string]map[string][]Aggregator)
			var pivheader = make(map[string]bool)

			for fields := range aggmux[am_i] {
				var key = string(stream.JoinSomeFields(a.Delim, fields, a.Keys))
				var pivot = ""
				if len(a.Pivots) > 0 {
					pivot = string(stream.JoinSomeFields(a.SubDelim, fields, a.Pivots))
				}

				var keyagg, ok1 = localagg[key]
				if !ok1 {
					keyagg = make(map[string][]Aggregator)
					localagg[key] = keyagg
				}

				var agg, ok2 = keyagg[pivot]
				if !ok2 {
					agg = make([]Aggregator, len(a.AggCtor))
					for i, ctor := range a.AggCtor {
						agg[i] = ctor()
					}
					keyagg[pivot] = agg
					pivheader[pivot] = true
				}

				for i, agtor := range agg {
					agtor.Aggregate(fields[a.Aggs[i]])
				}
			}

			m.Lock()
			for k, v := range localagg {
				a.Data[k] = v
			}
			for k, _ := range pivheader {
				a.PivsHeader[k] = true
			}
			m.Unlock()

			wg2.Done()
		}(g)
	}

	for line := range input {
		var fields = line.SplitFields(a.Delim)
		// hash key
		var hash = int64(5381)
		for _, i := range a.Keys {
			for _, b := range fields[i] {
				hash = (hash << 5) + hash + int64(b)
			}
		}
		aggmux[hash%partition] <- fields
	}

	for g := 0; g < partition; g++ {
		close(aggmux[g])
	}
	wg2.Wait()
}

func (a *Aggregation) printHeaderKeys(out io.Writer) {
	for i, kheader := range a.KeysHeader {
		if i > 0 {
			out.Write(a.Delim)
		}
		out.Write([]byte(kheader))
	}
}

func (a *Aggregation) printHeaderPivots(out io.Writer, pivots []string) {
	for i, pivot := range pivots {
		if i > 0 {
			out.Write(a.Delim)
		}
		for j, aheader := range a.AggsHeader {
			if j > 0 {
				out.Write(a.Delim)
			}
			out.Write([]byte(pivot + ":" + aheader))
		}
	}
}

func (a *Aggregation) printHeaderAggs(out io.Writer) {
	for i, aheader := range a.AggsHeader {
		if i > 0 {
			out.Write(a.Delim)
		}
		out.Write([]byte(aheader))
	}
}

func (a *Aggregation) printHeader(out io.Writer, pivots []string) {
	if len(a.Keys) > 0 {
		a.printHeaderKeys(out)
	}
	if len(a.Keys) > 0 && len(a.Aggs) > 0 {
		out.Write(a.Delim)
	}
	if len(a.Aggs) > 0 {
		if len(a.Pivots) > 0 {
			a.printHeaderPivots(out, pivots)
		} else {
			a.printHeaderAggs(out)
		}
	}
	out.Write([]byte{'\n'})
}

func (a *Aggregation) Print(out io.Writer) {
	var pivots []string
	for p := range a.PivsHeader {
		pivots = append(pivots, p)
	}

	a.printHeader(out, pivots)

	for key, agg := range a.Data {
		if len(a.Keys) > 0 {
			out.Write([]byte(key))
		}

		if len(a.Keys) > 0 && len(a.Aggs) > 0 {
			out.Write(a.Delim)
		}

		if len(a.Aggs) > 0 {
			for i, pivot := range pivots {
				if i > 0 {
					out.Write(a.Delim)
				}
				if pivaggs, ok := agg[pivot]; ok {
					for j, val := range pivaggs {
						if j > 0 {
							out.Write(a.Delim)
						}
						out.Write([]byte(val.String()))
					}
				} else {
					for j := 0; j < len(a.Aggs); j++ {
						if j > 0 {
							out.Write(a.Delim)
						}
						out.Write(a.NullVal)
					}
				}
			}
		}
		out.Write([]byte{'\n'})
	}
}
