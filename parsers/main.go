package main

import (
	"fmt"
	"parsers/opprec"
	"flag"
)

func init() {
	flag.Parse()
}

func main() {
	expr := flag.Arg(0)
	rpn, err := opprec.Parse(expr)
	if err != nil {
		panic(err)
	}
	res, err := opprec.Eval(rpn)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", res)
}
