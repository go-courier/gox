package stylis

import (
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func withUser(r string) string {
	return ".user{" + r + "}"
}

func TestParse(t *testing.T) {
	t.Run("unnested", func(t *testing.T) {
		NewWithT(t).Expect(Parse("--foo:none;@supports{--bar:none;}")).To(BeCSS(`
--foo:none;
@supports{--bar:none;}
`))
	})

	t.Run("escape", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser("height:calc(\\))\t!important"))).To(BeCSS(`
.user{height:calc(\))!important;}
`))
	})

	t.Run("calc", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(` 
height:calc( 100vh - 1px );
height:calc(
  100vh -
  1px
);
     `))).To(BeCSS(`
.user{height:calc( 100vh - 1px );height:calc(
  100vh -
  1px
);}
`))
	})

	t.Run("at-rules", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(` 
@-ms-viewport {
  width:device-width;
}
@viewport {
  width:device-width;
}
@page & {
  invalid:true;
}
@page {
  size:A4 landscape;
}
@document url(://www.w3.org/),url-prefix(//www.w3.org/),domain(mozilla.org),regexp("https:.*") {
  body {
	color: red;
  }
}
@viewport {
  min-width:640px;
  max-width:800px;
}
@counter-style list {
  system:fixed;
  symbols:url();
  suffix:" ";
}
@-moz-document url-prefix() {
  .selector {
	color:lime;
  }
}
@page {
  color:red;
  @bottom-right {
	content: counter(pages);
	margin-right: 1cm;
  }
  width: none;
}
     `))).To(BeCSS(`
@-ms-viewport{width:device-width;}
@viewport{width:device-width;}
@page &{invalid:true;}
@page{size:A4 landscape;}
@document url(://www.w3.org/),url-prefix(//www.w3.org/),domain(mozilla.org),regexp("https:.*"){.user body{color:red;}}
@viewport{min-width:640px;max-width:800px;}
@counter-style list{system:fixed;symbols:url();suffix:" ";}
@-moz-document url-prefix(){.user .selector{color:lime;}}
@page{color:red;@bottom-right{content:counter(pages);margin-right:1cm;}width:none;}
`))
	})

	t.Run("universal selector", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
* {
  color:red;
}
svg {
  &, & * {
	fill: currentColor;
  }
}
* * {
	color:hotpink;
}
`))).To(BeCSS(`
.user *{color:red;}
.user svg,.user svg *{fill:currentColor;}
.user * *{color:hotpink;}
`))
	})

	t.Run("flat", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
color:20px;
font-size:20px
`))).To(BeCSS(`
.user{color:20px;font-size:20px;}
`))
	})

	t.Run("namespace", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
& {
	color:red;
}
`))).To(BeCSS(`
.user{color:red;}
`))
	})

	t.Run("& in a string", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
& [href="https://css-tricks.com?a=1&b=2"] {
  color:red;
}
`))).To(BeCSS(`
.user [href="https://css-tricks.com?a=1&b=2"]{color:red;}
`))
	})

	t.Run("no trailing semi-colons", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
h2 {
display:none
}
div:hover
{
color:red
}
`))).To(BeCSS(`
.user h2{display:none;}
.user div:hover{color:red;}
`))
	})

	t.Run("edge cases", func(t *testing.T) {

	})

	t.Run("whitespace", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
div {
  width:0;    
}
.foo {
  color   :   hotpink;
}
`))).To(BeCSS(`
.user div{width:0;}
.user .foo{color:hotpink;}
`))
	})

	t.Run("multiline declaration", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
html {
background:
	linear-gradient(0deg, rgba(255, 255, 255, 0.8), rgba(255, 255, 255, 0.8)),
	url(/static/background.svg);
}
`))).To(BeCSS(`
.user html{background:linear-gradient(0deg, rgba(255, 255, 255, 0.8), rgba(255, 255, 255, 0.8)),url(/static/background.svg);}
`))
	})

	t.Run("nesting selector multiple levels", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
 a {
          a {
            a {
              a {
                a {
                  a {
                    a {
                      a {
                        a {
                          a {
                            a {
                              a {
                                color:red;
                              }
                            }
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
`))).To(BeCSS(`
.user a a a a a a a a a a a a{color:red;}
`))
	})

	t.Run("nesting @media multiple levels", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
  div {
          @media {
            a {
              color:red;
        @media {
                h1 {
                  color:hotpink;
                }
              }
            }
          }
        }
`))).To(BeCSS(`
@media{.user div a{color:red;}@media{.user div a h1{color:hotpink;}}}
`))
	})

	t.Run("context character", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
.a {background: url("img}.png");}
.a {background: url(img}.png);}
.a {background: url[img}.png];}
`))).To(BeCSS(`
.user .a{background:url("img}.png");}
.user .a{background:url(img}.png);}
.user .a{background:url[img}.png];}
`))
	})

	t.Run("`--` in an identifier", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
.block--modifier {
  color: hotpink;
}
.card {
  color: black;
}
`))).To(BeCSS(`
.user .block--modifier{color:hotpink;}
.user .card{color:black;}
`))
	})

	t.Run("parenthesis in string literal", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
@media only screen and (max-width: 320px){
	background: url("image_(1).jpg");
}
@media only screen and (min-width:321px) {
	background: url("image_(1).jpg");
}
.a {
  background: url("image_(1).jpg)");
}
.b {
  background: url("abc");
}
.icon {
  background:url("data:image/svg+xml;charset=utf-8,%3Csvg width='12' height='12' viewBox='0 0 16 16' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M14.117.323L8.044 6.398 2.595.323a1.105 1.105 0 0 0-1.562 1.562L6.482 7.96.323 14.119a1.105 1.105 0 0 0 1.562 1.562L7.96 9.608l5.449 6.073a1.103 1.103 0 1 0 1.56-1.562L9.517 8.046l6.159-6.161a1.103 1.103 0 1 0-1.56-1.562z' fill='rgba(85, 85, 85, 0.5)'/%3E%3C/svg%3E");
}
div {
  background: cyan;
}
.a .b .c {
  width: calc(100% / "func()");
}
.d {
  background: yellow;
}
`))).To(BeCSS(`
@media only screen and (max-width: 320px){.user{background:url("image_(1).jpg");}}
@media only screen and (min-width:321px){.user{background:url("image_(1).jpg");}}
.user .a{background:url("image_(1).jpg)");}
.user .b{background:url("abc");}
.user .icon{background:url("data:image/svg+xml;charset=utf-8,%3Csvg width='12' height='12' viewBox='0 0 16 16' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M14.117.323L8.044 6.398 2.595.323a1.105 1.105 0 0 0-1.562 1.562L6.482 7.96.323 14.119a1.105 1.105 0 0 0 1.562 1.562L7.96 9.608l5.449 6.073a1.103 1.103 0 1 0 1.56-1.562L9.517 8.046l6.159-6.161a1.103 1.103 0 1 0-1.56-1.562z' fill='rgba(85, 85, 85, 0.5)'/%3E%3C/svg%3E");}
.user div{background:cyan;}
.user .a .b .c{width:calc(100% / "func()");}
.user .d{background:yellow;}
`))
	})

	t.Run("nested parenthesis", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
width: calc(calc(1) + 10);
`))).To(BeCSS(`
.user{width:calc(calc(1) + 10);}
`))
	})

	t.Run("css variables edge cases", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
--braces: { };
--at-keyword-unknown-block: @foobar {};
--at-keyword-known-block: @media {};
--cdo-at-top-level: <!--;
--cdc-at-top-level: -->;
--semicolon-not-top-level: (;);
--cdo-not-top-level: (<!--);
--cdc-not-top-level: (-->);
--ampersand-preserved: foo & bar;
`))).To(BeCSS(`.user{` + strings.Join([]string{
			`--braces:{ };`,
			`--at-keyword-unknown-block:@foobar {};`,
			`--at-keyword-known-block:@media {};`,
			`--cdo-at-top-level:<!--;`,
			`--cdc-at-top-level:-->;`,
			`--semicolon-not-top-level:(;);`,
			`--cdo-not-top-level:(<!--);`,
			`--cdc-not-top-level:(-->);`,
			`--ampersand-preserved:foo & bar;`,
		}, "") + `}`))
	})

	t.Run("handles single `/` in a value", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
font: 12px/14px serif;
`))).To(BeCSS(`
.user{font:12px/14px serif;}
`))
	})

	t.Run("nested", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
  div {
	h2 {
	  color:red;
	  h3 {
		color:blue;
	  }
	}
  }
  .foo & {
	  width:1px;
	  &:hover {
		color:black;
	  }
	  li {
		color:white;
	  }
  }
  h1, div {
	color:red;
	h2, 
	&:before {
	  color:red;
	}
	color:blue;
	header {
	  font-size:12px;
	}
	@media {
	  color:red;
	}
	@media {
	  color:blue;
	}
  }
  &.foo {
	&.bar {
	  color:orange
	}
  }
  &.foo {
	&.bar {
	  &.barbar {
		color:orange
	  }
	}
  }
`))).To(BeCSS(`
.user div h2{color:red;}
.user div h2 h3{color:blue;}
.foo .user{width:1px;}
.foo .user:hover{color:black;}
.foo .user li{color:white;}
.user h1,.user div{color:red;color:blue;}
.user h1 h2,.user div h2,.user h1:before,.user div:before{color:red;}
.user h1 header,.user div header{font-size:12px;}
@media{.user h1,.user div{color:red;}}
@media{.user h1,.user div{color:blue;}}
.user.foo.bar{color:orange;}
.user.foo.bar.barbar{color:orange;}
`))
	})
}
