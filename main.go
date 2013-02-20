package main

import (
	"column"
	"flag"
	"io"
	"os"
	"strconv"
)

var Delim = flag.String("d", ",", "Field delimiter")

var Keys = flag.String("K", "", "Aggregation key fields")
var Sums = flag.String("S", "", "Aggregation sum fields")
var Average = flag.String("A", "", "Aggregation average fields")
var Count = flag.String("C", "", "Aggregation count fields")

// min, max, first, last, concat, x

type Aggregator interface {
	Aggregate(value []byte)
}

type Adder float64

func (adder *Adder) Aggregate(value []byte) {
	var f, _ = strconv.ParseFloat(string(value), 64)
	*adder += Adder(f)
}

func init() {
	flag.Parse()
}

func main() {
	r := column.NewReader(os.Stdin, []byte(*Delim))
	header, _ := r.ReadLine()
	keys := header.ParseFields(Keys)

	sums := header.ParseFields(Sums)
	var sumaggs []string
	for _, s := range sums {
		sumaggs = append(sumaggs, string(strconv.AppendInt([]byte("sum-"), int64(s), 10)))
	}

	var aggregations = make(map[string]map[string]Aggregator)

	var line, err = r.ReadLine()
	for err != io.EOF {
		var key = string(line.JoinFields(keys, []byte(*Delim)))
		if _, ok := aggregations[key]; !ok {
			aggregations[key] = make(map[string]Aggregator)
		}
		var aggs = aggregations[key]
		for i, s := range sums {
			if aggs[sumaggs[i]] == nil {
				aggs[sumaggs[i]] = new(Adder)
			}
			aggs[sumaggs[i]].Aggregate(line[s])
		}
		line, err = r.ReadLine()
	}

	for k, v := range aggregations {
		os.Stdout.Write([]byte(k))

		for _, s := range sumaggs {
			os.Stderr.Write([]byte(*Delim))
			os.Stderr.Write([]byte(strconv.FormatFloat(float64(*v[s].(*Adder)), 'g', -1, 64)))
		}
		os.Stderr.Write([]byte("\n"))
	}
}
