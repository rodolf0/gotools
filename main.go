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

	var aggregations = make(map[string]map[string]interface{})

	var line, err = r.ReadLine()
	for err != io.EOF {
		var key = string(line.JoinFields(keys, []byte(*Delim)))
		if _, ok := aggregations[key]; !ok {
			aggregations[key] = make(map[string]interface{})
		}
		for i, s := range sums {
			var num, _ = strconv.ParseFloat(string(line[s]), 64)
			switch acc := aggregations[key][sumaggs[i]].(type) {
			case float64:
				aggregations[key][sumaggs[i]] = acc + num
			case nil:
				aggregations[key][sumaggs[i]] = num
			}

		}
		line, err = r.ReadLine()
	}

	for k, v := range aggregations {
		os.Stdout.Write([]byte(k))

		for _, s := range sumaggs {
			os.Stderr.Write([]byte(*Delim))
			os.Stderr.Write([]byte(strconv.FormatFloat(v[s].(float64), 'g', -1, 64)))
		}
		os.Stderr.Write([]byte("\n"))
	}
}
