package scanner

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode/utf8"
)

func NewScanner(r io.Reader, scan func(ch rune, s *Scanner) (nextCh rune, tok rune, err error)) *Scanner {
	s := &Scanner{}
	s.src = r
	s.scan = scan

	// initialize source buffer
	// (the first call to Next() will fill it by calling src.Read)
	s.srcBuf[0] = utf8.RuneSelf // sentinel
	s.srcPos = 0
	s.srcEnd = 0

	// initialize source position
	s.srcBufOffset = 0
	s.line = 1
	s.column = 0
	s.lastLineLen = 0
	s.lastCharLen = 0

	// initialize token text buffer
	// (required for first call to Next()).
	s.tokPos = -1

	// initialize one character look-ahead
	s.ch = -2 // no char read yet, not EOF

	return s
}

const Whitespace = 1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' '
const bufLen = 1024 // at least utf8.UTFMax

type Scanner struct {
	// Input
	src  io.Reader
	scan func(ch rune, s *Scanner) (nextCh rune, tok rune, err error)

	// Source buffer
	srcBuf [bufLen + 1]byte // +1 for sentinel for common case of s.Next()
	srcPos int              // reading position (srcBuf index)
	srcEnd int              // source end (srcBuf index)

	// Source position
	srcBufOffset int // byte offset of srcBuf[0] in source
	line         int // line count
	column       int // character count
	lastLineLen  int // length of last line in characters (for correct column reporting)
	lastCharLen  int // length of last character in bytes

	// Token text buffer
	// Typically, token text is stored completely in srcBuf, but in general
	// the token text's head may be buffered in tokBuf while the token text's
	// tail is stored in srcBuf.
	tokBuf bytes.Buffer // token text head that is not in srcBuf anymore
	tokPos int          // token text tail position (srcBuf index); valid if >= 0
	tokEnd int          // token text tail end (srcBuf index)

	// One character look-ahead
	ch rune // character before current srcPos
}

func (s *Scanner) TokenText() string {
	if s.tokPos < 0 {
		// no token text
		return ""
	}

	if s.tokEnd < s.tokPos {
		// if EOF was reached, s.tokEnd is set to -1 (s.srcPos == 0)
		s.tokEnd = s.tokPos
	}
	// s.tokEnd >= s.tokPos
	if s.tokBuf.Len() == 0 {
		// common case: the entire token text is still in srcBuf
		return string(s.srcBuf[s.tokPos:s.tokEnd])
	}

	// part of the token text was saved in tokBuf: save the rest in
	// tokBuf as well and return its content
	s.tokBuf.Write(s.srcBuf[s.tokPos:s.tokEnd])
	s.tokPos = s.tokEnd // ensure idempotency of TokenText() call

	return strings.TrimSpace(s.tokBuf.String())
}

func (s *Scanner) Scan() rune {
	ch := s.Peek()

	s.tokBuf.Reset()
	s.tokPos = -1

	// skip white space
	for Whitespace&(1<<uint(ch)) != 0 {
		ch = s.Next()
	}

	s.tokBuf.Reset()
	s.tokPos = s.srcPos - s.lastCharLen

	// determine token value
	tok := ch
	var err error

	if ch != TokenEOF {
		ch, tok, err = s.scan(ch, s)
		if err != nil {
			s.error(err.Error())
		}

		// end of token text
		s.tokEnd = s.srcPos - s.lastCharLen
		s.ch = ch
	}

	return tok
}

func (s *Scanner) Next() rune {
	ch, width := rune(s.srcBuf[s.srcPos]), 1

	if ch >= utf8.RuneSelf {
		// uncommon case: not ASCII or not enough bytes
		for s.srcPos+utf8.UTFMax > s.srcEnd && !utf8.FullRune(s.srcBuf[s.srcPos:s.srcEnd]) {
			// not enough bytes: read some more, but first
			// save away token text if any
			if s.tokPos >= 0 {
				s.tokBuf.Write(s.srcBuf[s.tokPos:s.srcPos])
				s.tokPos = 0
				// s.tokEnd is set by Scan()
			}
			// move unread bytes to beginning of buffer
			copy(s.srcBuf[0:], s.srcBuf[s.srcPos:s.srcEnd])
			s.srcBufOffset += s.srcPos
			// read more bytes
			// (an io.Reader must return io.EOF when it reaches
			// the end of what it is reading - simply returning
			// n == 0 will make this loop retry forever; but the
			// error is in the reader implementation in that case)
			i := s.srcEnd - s.srcPos
			n, err := s.src.Read(s.srcBuf[i:bufLen])
			s.srcPos = 0
			s.srcEnd = i + n
			s.srcBuf[s.srcEnd] = utf8.RuneSelf // sentinel
			if err != nil {
				if err != io.EOF {
					s.error(err.Error())
				}
				if s.srcEnd == 0 {
					if s.lastCharLen > 0 {
						// previous character was not EOF
						s.column++
					}
					s.lastCharLen = 0
					return TokenEOF
				}
				// If err == EOF, we won't be getting more
				// bytes; break to avoid infinite loop. If
				// err is something else, we don't know if
				// we can get more bytes; thus also break.
				break
			}
		}
		// at least one byte
		ch = rune(s.srcBuf[s.srcPos])
		if ch >= utf8.RuneSelf {
			// uncommon case: not ASCII
			ch, width = utf8.DecodeRune(s.srcBuf[s.srcPos:s.srcEnd])
			if ch == utf8.RuneError && width == 1 {
				// advance for correct error position
				s.srcPos += width
				s.lastCharLen = width
				s.column++
				s.error("invalid UTF-8 encoding")
				return ch
			}
		}
	}

	// advance
	s.srcPos += width
	s.lastCharLen = width
	s.column++

	// special situations
	switch ch {
	case 0:
		// for compatibility with other tools
		s.error("invalid character NUL")
	case '\n':
		s.line++
		s.lastLineLen = s.column
		s.column = 0
	}

	return ch
}

func (s *Scanner) error(msg string) {
}

func (s *Scanner) Peek() rune {
	if s.ch == -2 {
		// this code is only run for the very first character
		s.ch = s.Next()
		if s.ch == '\uFEFF' {
			s.ch = s.Next() // ignore BOM
		}
	}
	return s.ch
}

func ScanString(quote rune, next func() rune) (ch rune, err error) {
	ch = next() // read character after quote
	for ch != quote {
		if ch == '\n' || ch < 0 {
			err = errors.New("literal not terminated")
			return
		}
		if ch == '\\' {
			c, e := ScanEscape(quote, next)
			if e != nil {
				err = e
				return
			}
			ch = c
		} else {
			ch = next()
		}
	}
	return
}

func ScanEscape(quote rune, next func() rune) (rune, error) {
	ch := next() // read character after '/'
	switch ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		// nothing to do
		ch = next()
	default:
		return 0, errors.New("invalid char escape")
	}
	return ch, nil
}

func ScanComment(ch rune, next func() rune) (rune, error) {
	// ch == '/' || ch == '*'
	if ch == '/' {
		// line comment
		ch = next() // read character after "//"
		for ch != '\n' && ch >= 0 {
			ch = next()
		}
		return ch, nil
	}

	// general comment
	ch = next() // read character after "/*"
	for {
		if ch < 0 {
			return 0, errors.New("comment not terminated")
		}
		ch0 := ch
		ch = next()
		if ch0 == '*' && ch == '/' {
			ch = next()
			break
		}
	}
	return ch, nil
}
