package lexer

import (
	"fmt"
)

func Example1() {
	ml := NewMathLexer("3+4*2/-(1-5)^2^3")
	for {
		if t := ml.Next(); t != nil {
			fmt.Println(t)
		} else {
			break
		}
	}
	// Output:
	// Number: "3"
	// Plus: "+"
	// Number: "4"
	// Power: "*"
	// Number: "2"
	// Divide: "/"
	// UMinus: "-"
	// OParen: "("
	// Number: "1"
	// Minus: "-"
	// Number: "5"
	// CParen: ")"
	// Power: "^"
	// Number: "2"
	// Power: "^"
	// Number: "3"
}

func Example2() {
	ml := NewMathLexer("3.4e-2 * sin(x)/(7! % -4) * max(2, x)")
	for {
		if t := ml.Next(); t != nil {
			fmt.Println(t)
		} else {
			break
		}
	}
	// Output:
	// Number: "3.4e-2"
	// Times: "*"
	// Function: "sin("
	// Variable: "x"
	// CParen: ")"
	// Divide: "/"
	// OParen: "("
	// Number: "7"
	// Factorial: "!"
	// Modulo: "%"
	// UMinus: "-"
	// Number: "4"
	// CParen: ")"
	// Times: "*"
	// Function: "max("
	// Number: "2"
	// Comma: ","
	// Variable: "x"
	// CParen: ")"
}

func Example3() {
	ml := NewMathLexer("x---y")
	for {
		if t := ml.Next(); t != nil {
			fmt.Println(t)
		} else {
			break
		}
	}
	// Output:
	// Variable: "x"
	// Minus: "-"
	// UMinus: "-"
	// UMinus: "-"
	// Variable: "y"
}

func Example4() {
	ml := NewMathLexer("sqrt(-(1i-x^2) / (1 + x^2))")
	for {
		if t := ml.Next(); t != nil {
			fmt.Println(t)
		} else {
			break
		}
	}
	// Output:
	// Function: "sqrt("
	// UMinus: "-"
	// OParen: "("
	// Number: "1i"
	// Minus: "-"
	// Variable: "x"
	// Power: "^"
	// Number: "2"
	// CParen: ")"
	// Divide: "/"
	// OParen: "("
	// Number: "1"
	// Plus: "+"
	// Variable: "x"
	// Power: "^"
	// Number: "2"
	// CParen: ")"
	// CParen: ")"
}
