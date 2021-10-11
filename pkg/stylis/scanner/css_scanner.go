package scanner

import (
	"io"
)

func NewCSSScanner(r io.Reader) *Scanner {
	s := NewScanner(r, func(ch rune, s *Scanner) (rune, rune, error) {
		tok := ch

		switch ch {
		case '-':
			// variables
			ch = s.Next()
			if ch == '-' {
				for ch != ';' {
					ch = scanCSSExpr(s)
				}
			}
		case '/':
			// comments
			ch = s.Next()
			if ch == '/' || ch == '*' {
				tok = TokenComment
				c, err := ScanComment(ch, s.Next)
				if err != nil {
					return 0, 0, err
				}
				ch = c
			}
		case ';', '{', '}':
			tok = breakpointTokens[ch]
			ch = s.Next()
		default:
			tok = TokenExpr
			ch = scanCSSExpr(s)
		}

		return ch, tok, nil
	})
	return s
}

var groups = map[rune]rune{
	'(': ')',
	'[': ']',
}

func scanCSSExpr(s *Scanner) rune {
	ch := s.Peek()

	var scan func(end rune)

	scan = func(groupEnd rune) {
		ch = s.Next() // read character after started

		for ch != TokenEOF && ch != groupEnd {
			if ch == '\\' {
				ch, _ = ScanEscape(groupEnd, s.Next)
			} else {
				if ch == '\'' || ch == '"' {
					ch, _ = ScanString(ch, s.Next)
					ch = s.Next()
				} else if end, ok := groups[ch]; ok {
					scan(end)
					ch = s.Next()
				} else {
					ch = s.Next()
				}
			}
		}
	}

	for !isBreakpointToken(ch) {
		if ch == TokenEOF {
			break
		}

		if ch == '\'' || ch == '"' {
			ch, _ = ScanString(ch, s.Next)
			ch = s.Next()
		} else if groupEnd, ok := groups[ch]; ok {
			scan(groupEnd)
			ch = s.Next()
		} else {
			ch = s.Next()
		}
	}

	return ch
}

var breakpointTokens = map[rune]rune{
	'{': TokenBlockStart,
	'}': TokenBlockEnd,
	';': TokenSemicolon,
}

func isBreakpointToken(r rune) bool {
	_, ok := breakpointTokens[r]
	return ok
}

const (
	TokenEOF = -(iota + 1)
	TokenExpr
	TokenSemicolon  // ;
	TokenBlockStart // {
	TokenBlockEnd   // }
	TokenComment
)

var TokenString = map[rune]string{
	TokenEOF:        "EOF",
	TokenExpr:       "Expr",
	TokenSemicolon:  "Semicolon",
	TokenBlockStart: "BlockStart",
	TokenBlockEnd:   "BlockEnd",
	TokenComment:    "Comment",
}
