package stylis

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/davecgh/go-spew/spew"
)

func TestTokenizer(t *testing.T) {
	t.Run("Rules", func(t *testing.T) {
		NewWithT(t).Expect((&Tokenizer{}).Rules([]byte(`.test:matches(a,b,c), .test`))).To(Equal([]string{
			".test:matches(a,b,c)",
			".test",
		}))
	})

	t.Run("Tokens", func(t *testing.T) {
		tokenizer := Tokenizer{}

		tokens := tokenizer.Tokenize([]byte(`
--foo:none;
@supports {
	--at-keyword-known-block:@media {};
	--cdo-not-top-level:(<!--);
	--bar:none;
	--foo: { };
	--foo2: (;);
}
@document url(://www.w3.org/),url-prefix(//www.w3.org/),domain(mozilla.org),regexp("https:.*") {
}
@counter-style list {
  system:fixed	!important;
  symbols:url(}.jpg);
  symbols:url("}.jpg");
  background:url("image_(1).jpg"/* test */);	
  suffix:" ";
  suffix:";";
  height: calc( 
				100vh - calc( 100vh - 1px ) 
			);	
  font-size:20px	
}
body /*
// test
xxx
*/ {
	padding: 2px 2px;
}

.user {
	background-color: red;

	&:hover, &:active {
		color: red;
	}
}

&.B\\26W{color:red;}

.user{
	color:20px;
	font-size: 11px
}

.test[title=","] {
  color:red;
}
[title="a,b,c, something"], h1, [title="a,b,c"] {
  color:red
}
[title="a"],
[title="b"] {
  color:red;
}
`))
		spew.Dump(tokens)
	})
}

func TestReplaceHolder(t *testing.T) {
	NewWithT(t).Expect(ReplaceHolder(`& [url="test&x"] > \&`, ".parent")).To(Equal(`.parent [url="test&x"] > \&`))
}
