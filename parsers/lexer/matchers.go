package lexer

type Matcher struct {
	*Scanner
}

// Build a new matcher
func NewMatcherString(s string) *Matcher {
	return &Matcher{
		Scanner: NewScannerString(s),
	}
}

// Try to match a number and return true if matched
func (m *Matcher) MatchNumber() bool {
	onlyinteger := false // only allow integer part of number (hex, oct, bin)
	backtrack := m.pos
	// optional number sign
	m.Accept("+-")
	// support for signed hexadecimal integers
	digits := "0123456789"

	if m.Accept("0") { // possibly another base
		onlyinteger = true
		if m.Accept("xX") { // hexadecimal
			digits = "0123456789aAbBcCdDeEfF"
		} else if m.Accept("oO") { // octal
			digits = "01234567"
		} else if m.Accept("bB") { // binary
			digits = "01"
		} else { // base 10, let the 0 be parsed as integer part
			onlyinteger = false
			m.Prev()
		}
	}
	// require some integer part (ie: not allowing .012)
	if !m.Skip(digits) {
		m.pos = backtrack
		return false
	}
	if onlyinteger {
		return true // only base 10 supports fractions/exponent/img
	}
	// take optional fractions
	backtrack = m.pos
	if m.Accept(".") {
		// require digits if fraction is present
		if !m.Skip(digits) {
			m.pos = backtrack
			return true // found an integer
		}
	}
	// take optional exponent
	backtrack = m.pos
	if m.Accept("eE") { // TODO: problem 0x3Ee+2... bad luck
		m.Accept("+-") // opt exponent sign
		// require digits if exponent present
		if !m.Skip(digits) {
			m.pos = backtrack
			return true // found a number without exponent
		}
	}
	// optional imaginary
	m.Accept("i")
	return true
}

// Try to match an identifier
func (m *Matcher) MatchId() bool {
	alfa := "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if !m.Accept(alfa) {
		return false
	}
	m.Skip(alfa + "0123456789")
	return true
}
