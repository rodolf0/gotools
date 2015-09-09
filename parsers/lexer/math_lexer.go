package lexer

import "fmt"

const (
	TokUnknown = iota
	TokNumber
	TokVariable
	TokFunction
	TokOParen
	TokCParen
	TokComma
	TokPlus
	TokMinus
	TokTimes
	TokDivide
	TokModulo
	TokPower
	TokUMinus
	TokFactorial
)

type MathLexcomp uint

type MathToken struct {
	Lexeme  string
	Lexcomp MathLexcomp
}

func (t MathToken) String() string {
	var lcstr string
	switch t.Lexcomp {
	case TokUnknown:
		lcstr = "Unknown"
	case TokNumber:
		lcstr = "Number"
	case TokVariable:
		lcstr = "Variable"
	case TokFunction:
		lcstr = "Function"
	case TokPlus:
		lcstr = "Plus"
	case TokMinus:
		lcstr = "Minus"
	case TokTimes:
		lcstr = "Times"
	case TokDivide:
		lcstr = "Divide"
	case TokModulo:
		lcstr = "Modulo"
	case TokPower:
		lcstr = "Power"
	case TokUMinus:
		lcstr = "UMinus"
	case TokFactorial:
		lcstr = "Factorial"
	case TokOParen:
		lcstr = "OParen"
	case TokCParen:
		lcstr = "CParen"
	case TokComma:
		lcstr = "Comma"
	}
	return fmt.Sprintf("%v: %q", lcstr, t.Lexeme)
}

type MathLexer struct {
	m      *Matcher
	tokens []*MathToken
	pos    int
}

func NewMathLexer(expr string) *MathLexer {
	return &MathLexer{
		m:      NewMatcherString(expr),
		tokens: make([]*MathToken, 0, 32),
		pos:    -1,
	}
}

// check if binary or unary minus according to previous token
func makesUMinus(prev *MathToken) bool {
	// First token
	if prev == nil {
		return true
	}
	// binary cases by extension are less then unary
	switch prev.Lexcomp {
	case TokNumber, TokVariable, TokCParen:
		return false
	}
	// default to unary
	return true
}

// get more tokens
func (ml *MathLexer) readToken() *MathToken {
	// Ignore any white space
	if ml.m.SkipWs() {
		ml.m.Ignore()
	}

	var tok *MathToken

	if ml.m.MatchId() {
		// lex variable and function names
		name := ml.m.Extract()
		// skip ws between func name and paren
		if ml.m.SkipWs() {
			ml.m.Ignore()
		}
		if ml.m.Accept(`(`) {
			tok = &MathToken{Lexeme: name + ml.m.Extract(), Lexcomp: TokFunction}
		} else {
			tok = &MathToken{Lexeme: name, Lexcomp: TokVariable}
		}

	} else if ml.m.Accept(`+-*/%^!(),`) {
		// lex operators before numbers to avoid confusing `-` as part of number
		switch r := ml.m.Extract(); r {
		case `+`:
			tok = &MathToken{Lexeme: `+`, Lexcomp: TokPlus}
		case `-`:
			if len(ml.tokens) == 0 || makesUMinus(ml.tokens[len(ml.tokens)-1]) {
				tok = &MathToken{Lexeme: `-`, Lexcomp: TokUMinus}
			} else {
				tok = &MathToken{Lexeme: `-`, Lexcomp: TokMinus}
			}
		case `*`:
			tok = &MathToken{Lexeme: `*`, Lexcomp: TokTimes}
		case `/`:
			tok = &MathToken{Lexeme: `/`, Lexcomp: TokDivide}
		case `%`:
			tok = &MathToken{Lexeme: `%`, Lexcomp: TokModulo}
		case `^`:
			tok = &MathToken{Lexeme: `^`, Lexcomp: TokPower}
		case `!`:
			tok = &MathToken{Lexeme: `!`, Lexcomp: TokFactorial}
		case `(`:
			tok = &MathToken{Lexeme: `(`, Lexcomp: TokOParen}
		case `)`:
			tok = &MathToken{Lexeme: `)`, Lexcomp: TokCParen}
		case `,`:
			tok = &MathToken{Lexeme: `,`, Lexcomp: TokComma}
		}

	} else if ml.m.MatchNumber() {
		// lex numbers
		tok = &MathToken{Lexeme: ml.m.Extract(), Lexcomp: TokNumber}

	} else if unk := ml.m.Peek(); unk != EOF {
		// Don't know what this is
		ml.m.Next()
		tok = &MathToken{Lexeme: ml.m.Extract(), Lexcomp: TokUnknown}
	} // else EOF => return nil

	return tok
}

func (ml *MathLexer) Next() *MathToken {
	ml.pos += 1
	if ml.pos >= len(ml.tokens) {
		t := ml.readToken()
		if t == nil {
			ml.pos = len(ml.tokens)
			return nil
		} else if t.Lexcomp == TokUnknown {
			panic("MathLexer: Unknown token [" + t.Lexeme + "]")
		}
		ml.tokens = append(ml.tokens, t)
	}
	return ml.Curr()
}

func (ml *MathLexer) Curr() *MathToken {
	if ml.pos < 0 {
		return nil // BOF
	} else if ml.pos >= len(ml.tokens) {
		return nil // EOF
	}
	return ml.tokens[ml.pos]
}
