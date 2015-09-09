package opprec

import "fmt"

func Example1() {
	rpn, err := Parse("3+4*2/-sin(1-5)^2^3")
	if err != nil {
		panic(err)
	}
	for _, t := range rpn {
		fmt.Println(t.tok)
	}
	// Output:
	// Number: "3"
	// Number: "4"
	// Number: "2"
	// Times: "*"
	// Number: "1"
	// Number: "5"
	// Minus: "-"
	// Function: "sin("
	// Number: "2"
	// Number: "3"
	// Power: "^"
	// Power: "^"
	// UMinus: "-"
	// Divide: "/"
	// Plus: "+"
}


func Example2() {
	rpn, err := Parse("3+4*2/-(1-5)^2^3")
	if err != nil {
		panic(err)
	}
	for _, t := range rpn {
		fmt.Println(t.tok)
	}
	// Output:
	// Number: "3"
	// Number: "4"
	// Number: "2"
	// Times: "*"
	// Number: "1"
	// Number: "5"
	// Minus: "-"
	// Number: "2"
	// Number: "3"
	// Power: "^"
	// Power: "^"
	// UMinus: "-"
	// Divide: "/"
	// Plus: "+"
}


func Example3() {
	rpn, err := Parse("3.2/4!")
	if err != nil {
		panic(err)
	}
	for _, t := range rpn {
		fmt.Println(t.tok)
	}
	// Output:
	// Number: "3.2"
	// Number: "4"
	// Factorial: "!"
	// Divide: "/"
}
