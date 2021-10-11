package stylis

import (
	"bytes"
	"io"
	"strings"
	"text/scanner"
)

type Tokenizer struct {
	buf strings.Builder
}

func (t *Tokenizer) Tokenize(src []byte) (tokens []string) {
	s := &scanner.Scanner{}
	s.Init(bytes.NewBuffer(src))

	t.buf.Reset()

	collect := func() {
		text := strings.TrimSpace(t.buf.String())
		if len(text) > 0 {
			tokens = append(tokens, text)
		}
		t.buf.Reset()
	}

	for tok := t.scan(s); tok != scanner.EOF; tok = t.scan(s) {
		if tok == scanner.Comment {
			continue
		}
		collect()
	}
	collect()

	return
}

func (t *Tokenizer) Rules(rule []byte) (rules []string) {
	s := &scanner.Scanner{}
	s.Init(bytes.NewBuffer(rule))

	t.buf.Reset()

	collect := func() {
		text := strings.TrimSpace(t.buf.String())
		if len(text) > 0 && text != "," {
			rules = append(rules, text)
		}
		t.buf.Reset()
	}

	for tok := t.scanRule(s); tok != scanner.EOF; tok = t.scanRule(s) {
		if tok == scanner.Comment {
			continue
		}
		collect()
	}

	collect()

	return
}

func (t *Tokenizer) scanRule(s *scanner.Scanner) rune {
	tok := s.Peek()
	switch tok {
	case ',':
		t.buf.WriteRune(s.Next())
		return tok
	default:
		return t.nextUntil(s, &t.buf, ',')
	}
}

func (t *Tokenizer) scan(s *scanner.Scanner) rune {
	tok := s.Peek()
	switch tok {
	case '\t', '\n', '\r':
		s.Next()
		return t.scan(s)
	case '-':
		t.buf.WriteRune(s.Next())
		if tok := s.Peek(); tok == '-' {
			t.buf.WriteRune(s.Next())
			return t.nextUntil(s, &t.buf, ';', '\n')
		}
		return t.nextUntil(s, &t.buf, ';', '{', '}')
	case ';', '{', '}':
		t.buf.WriteRune(s.Next())
		return tok
	default:
		return t.nextUntil(s, &t.buf, ';', '{', '}')
	}
}

func (t *Tokenizer) nextUntil(s *scanner.Scanner, b *strings.Builder, breakpoints ...rune) rune {
	tok := s.Peek()

	for i := range breakpoints {
		if tok == breakpoints[i] {
			return tok
		}
	}

	switch tok {
	case scanner.EOF:
		return tok
	case '\t':
		// skip \t
		tok = s.Next()
	case ',':
		tok = s.Next()
		b.WriteRune(tok)

		tok = s.Next()
		for {
			switch tok {
			case '\t', '\n', '\r':
				tok = s.Next()
				continue
			}
			break
		}
		b.WriteRune(tok)
	case '/':
		tok = s.Next()

		if ch := s.Peek(); ch == '*' {
			ch = s.Next() // read character after "/*"

			for {
				if ch == scanner.EOF {
					break
				}

				ch0 := ch
				ch = s.Next()
				if ch0 == '*' && ch == '/' {
					ch = s.Next()
					break
				}
			}

			return s.Peek()
		}
		b.WriteRune(tok)
	case '\\':
		b.WriteRune(s.Next())
		b.WriteRune(s.Next())
	case '\'', '"':
		quote := tok
		buf := &strings.Builder{}

		ch := s.Next() // read character after quote
		buf.WriteRune(ch)
		ch = s.Next()
		buf.WriteRune(ch)

		for ch != quote {
			if ch == '\\' {
				buf.WriteRune(s.Next())

				ch = s.Next()
				buf.WriteRune(ch)
			} else {
				ch = s.Next()
				buf.WriteRune(ch)
			}
		}

		b.WriteString(buf.String())
	case '(':
		buf := &strings.Builder{}
		buf.WriteRune(s.Next())
		tok = t.nextUntil(s, buf, ')')
		buf.WriteRune(s.Next())

		b.WriteString(buf.String())
	case '[':
		buf := &strings.Builder{}
		buf.WriteRune(s.Next())
		tok = t.nextUntil(s, buf, ']')
		buf.WriteRune(s.Next())

		b.WriteString(buf.String())
	default:
		_, _ = io.WriteString(b, string(s.Next()))
	}

	return t.nextUntil(s, b, breakpoints...)
}

func SplitRules(s string) []string {
	return (&Tokenizer{}).Rules([]byte(s))
}

func ReplaceHolder(rule string, p string) string {
	if p == "" {
		return rule
	}

	s := &scanner.Scanner{}

	s.Init(bytes.NewBufferString(rule))
	s.Whitespace = 1<<'\t' | 1<<'\r' | 1<<'\n'
	s.Mode = scanner.ScanChars | scanner.ScanStrings | scanner.ScanComments

	b := bytes.NewBuffer(nil)
	hasHolder := false

	for t := s.Scan(); t != scanner.EOF; t = s.Scan() {
		switch t {
		case '\\':
			b.WriteString(s.TokenText())
			t = s.Scan()
			if t != scanner.EOF {
				b.WriteString(s.TokenText())
			}
		case '&':
			hasHolder = true
			b.WriteString(p)
		default:
			b.WriteString(s.TokenText())
		}
	}

	if hasHolder {
		return b.String()
	}

	return p + " " + b.String()
}
