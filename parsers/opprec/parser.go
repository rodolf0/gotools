package opprec

import (
	"errors"
	"parsers/lexer"
	"strconv"
)

const (
	assocLeft = iota
	assocRight
	assocNone // eg: (a := b := c) is an error
)

type opPrec struct {
	precedence uint
	assoc      uint
}

func precOf(op lexer.MathLexcomp) opPrec {
	prec := map[lexer.MathLexcomp]opPrec{
		lexer.TokFunction: {1, assocLeft},
		lexer.TokOParen:   {1, assocLeft},

		lexer.TokPlus:  {2, assocLeft},
		lexer.TokMinus: {2, assocLeft},

		lexer.TokTimes:  {3, assocLeft},
		lexer.TokDivide: {3, assocLeft},
		lexer.TokModulo: {3, assocLeft},

		lexer.TokUMinus: {4, assocRight},

		lexer.TokPower: {5, assocRight},

		lexer.TokFactorial: {6, assocLeft},
	}
	if p, ok := prec[op]; ok {
		return p
	}
	return opPrec{100, assocNone}
}

type Token struct {
	tok      *lexer.MathToken // the lexer token from which this was derived
	arity    uint             // arity if the token is an operator or a func
}

// Parse the expression and output a reverse-polish-notation slice
// assumes function names end in '('
func Parse(expr string) ([]*Token, error) {
	out := make([]*Token, 0, len(expr))
	stack := make([]*Token, 0, len(expr)/2)

	ml := lexer.NewMathLexer(expr)

	for tok := ml.Next(); tok != nil; tok = ml.Next() {
		switch tok.Lexcomp {
		case lexer.TokNumber, lexer.TokVariable:
			out = append(out, &Token{tok: tok})

		case lexer.TokFunction, lexer.TokOParen:
			stack = append(stack, &Token{tok: tok})

		case lexer.TokComma:
			// TODO: track function arity
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top.tok.Lexcomp == lexer.TokOParen ||
					top.tok.Lexcomp == lexer.TokFunction {
					break
				}
				// until we reach an OParen pop stack
				out = append(out, top)
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, errors.New("Parse: mismatched parens")
			}

		case lexer.TokCParen:
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top.tok.Lexcomp == lexer.TokOParen ||
					top.tok.Lexcomp == lexer.TokFunction {
					break
				}
				// until we reach an OParen pop stack
				out = append(out, top)
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, errors.New("Parse: mismatched parens")
			}
			// implied top of the stack is OParen/Function
			top := stack[len(stack)-1]
			if top.tok.Lexcomp == lexer.TokFunction {
				out = append(out, top)
			}
			stack = stack[:len(stack)-1] // pop OParen/Func

		case lexer.TokUnknown:
			return nil, errors.New("Parse: unknown token") // not reachable, panic at ml.Next

		default: // assuming everything left is operator (too lazy to list)
			ptok := precOf(tok.Lexcomp)
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				ptop := precOf(top.tok.Lexcomp)

				var ok = true
				if ptok.precedence < ptop.precedence {
					ok = true
				} else if ptok.precedence == ptop.precedence {
					switch ptok.assoc {
					case assocLeft:
						ok = true
					case assocRight:
						ok = false
					case assocNone:
						panic("quack!")
					}
				} else {
					ok = false
				}

				if !ok {
					break
				}

				// pop out of the stack and push onto output
				out = append(out, top)
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, &Token{tok: tok})
		}
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		if top.tok.Lexcomp == lexer.TokOParen || top.tok.Lexcomp == lexer.TokFunction {
			return nil, errors.New("Parse: mismatched parenthesis")
		}
		out = append(out, top)
		stack = stack[:len(stack)-1]
	}

	return out, nil
}

// Eval takes a reverse polish notation list of tokens and evaluates
func Eval(rpn []*Token) (float64, error) {
	var stack []float64

	for _, in := range rpn {

		if in.tok.Lexcomp == lexer.TokNumber {
			num, _ := strconv.ParseFloat(in.tok.Lexeme, 64)
			stack = append(stack, num)
			continue
		}

		switch in.tok.Lexcomp {
		case lexer.TokNumber:
			num, _ := strconv.ParseFloat(in.tok.Lexeme, 64)
			stack = append(stack, num)

		case lexer.TokPlus:
			op1, op2 := stack[len(stack)-2], stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			stack = append(stack, op1 + op2)

		case lexer.TokMinus:
			op1, op2 := stack[len(stack)-2], stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			stack = append(stack, op1 - op2)

		case lexer.TokTimes:
			op1, op2 := stack[len(stack)-2], stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			stack = append(stack, op1 * op2)

		case lexer.TokUMinus:
			op1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, -op1)

		default:
			panic("not implemented")
		}

	}
	return stack[0], nil
}
