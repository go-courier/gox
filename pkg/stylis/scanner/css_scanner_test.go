package scanner

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCSSScanner(t *testing.T) {
	s := NewCSSScanner(bytes.NewBuffer([]byte(`
--foo:none;
@supports {
	--bar:none;
	--foo: {};
}
@counter-style list {
  system:fixed;
  symbols:url(}.jpg);
  symbols:url("}.jpg");
  background:url("image_(1).jpg");	
  suffix:" ";
  suffix:";";
}
// test
body {
	height: calc( 100vh - 1px );
	padding: 2px 2px;
}

.user {
	background-color: red;

	&:hover, &:active {
		color: red;
	}
}
`)))

	for tok := s.Scan(); tok != TokenEOF; tok = s.Scan() {
		fmt.Printf("....%s\n", s.TokenText())
	}
}
