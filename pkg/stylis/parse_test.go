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

	t.Run("escaped chars in selector identifiers", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
&.B\&W{color:red;}
&.\@example\.com{color:blue;}
&.owner\/founder{color:green;}
&.discount\%  {color:purple;}
`))).To(BeCSS(`
.user.B\&W{color:red;}
.user.\@example\.com{color:blue;}
.user.owner\/founder{color:green;}
.user.discount\%{color:purple;}
`))
	})

	t.Run("escaped hex codes in selector identifiers", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
&.B\\26W{color:red;}
&.B\\000026W{color:green;}
&.B\\26 W{color:blue;}
`))).To(BeCSS(`
.user.B\\26W{color:red;}
.user.B\\000026W{color:green;}
.user.B\\26 W{color:blue;}
`))
	})

	t.Run("double spaces after escaped hex codes in selector identifiers", func(t *testing.T) {
		// TODO
		// https://github.com/thysultan/stylis.js/blob/master/test/Parser.js#L181
		// https://www.w3.org/International/questions/qa-escapes
		NewWithT(t).Expect(Parse(withUser(`
&.endsWith\\0000A9  a.childNode{color:green;}
&.endsWith\\AE  a.childNode{color:yellow;}
&.Q\\000026A  a.childNode{color:purple;}
`))).To(BeCSS(`
.user.endsWith\\0000A9  a.childNode{color:green;}
.user.endsWith\\AE  a.childNode{color:yellow;}
.user.Q\\000026A  a.childNode{color:purple;}
`))
	})

	t.Run("comments", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
/**
 * removes block comments and line comments,
 * there't a fire in the house // there is
 */
button /*
  // what't
  xxx
  */
{color: blue;}
button /* 1 */
{
  color: red; /* 2 */
}
/*! 1 */
color: red;
/*! 2 */
h1 {
  /*! 1 */
  color: red;
  /*! 2 */
  color: red;
  /*! 3 */
}
`))).To(BeCSS(
			`.user{color:red;}`,
			`.user button{color:blue;}`,
			`.user button{color:red;}`,
			`.user h1{color:red;color:red;}`,
		))
	})

	t.Run("&", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
& {
  color:blue;
}
&&& {
  color:red;
}
&+& {
  color:red;
}
.wrapper button& {
  color:red;
}
&:hover & {
  color: green;
}
div:hover & {
  color: green;
}
div:hover & {
  h1 & {
	color:red;
  }
}
`))).To(BeCSS(`
.user{color:blue;}
.user.user.user{color:red;}
.user+.user{color:red;}
.wrapper button.user{color:red;}
.user:hover .user{color:green;}
div:hover .user{color:green;}
h1 div:hover .user{color:red;}
`))
	})

	t.Run("&:before", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
&:before{
  color:blue;
}
`))).To(BeCSS(`
.user:before{color:blue;}
`))
	})

	t.Run("& :hover", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
& :hover{
  color:blue;
}
`))).To(BeCSS(`
.user :hover{color:blue;}
`))
	})

	t.Run("div :hover", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
div :hover{
  color:blue;
}
`))).To(BeCSS(`
.user div :hover{color:blue;}
`))
	})

	t.Run("@supports", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
@supports (display:block) {
  color:red;
  h1 {
	color:red;
	h2 {
	  color:blue;
	}
  }
  display:none;
}
@supports (appearance: none) {
  color:red;
}
@supports (backdrop-filter: blur(10px)) {
  backdrop-filter: blur(10px);
}
`))).To(BeCSS(`
@supports (display:block){.user{color:red;display:none;}.user h1{color:red;}.user h1 h2{color:blue;}}
@supports (appearance: none){.user{color:red;}}
@supports (backdrop-filter: blur(10px)){.user{backdrop-filter:blur(10px);}}
`))
	})

	t.Run("@media", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
@media (max-width:600px) {
  color:red;
  h1 {
	color:red;
	h2 {
	  color:blue;
	}
  }
  display:none;
}
@media (min-width:576px) {
  &.card-deck {
	.card {
	  &:not(:first-child) {
		margin-left:15px;
	  }
	  &:not(:last-child) {
		margin-right:15px;
	  }
	}
  }
}
@supports (display:block) {
  @media (min-width:10px) {
	background-color:seagreen;
  }
}
@media (max-width:600px) {
  & { color:red }
}
&:hover {
  color:orange
}
`))).To(BeCSS(`
@media (max-width:600px){.user{color:red;display:none;}.user h1{color:red;}.user h1 h2{color:blue;}}
@media (min-width:576px){.user.card-deck .card:not(:first-child){margin-left:15px;}.user.card-deck .card:not(:last-child){margin-right:15px;}}
@supports (display:block){@media (min-width:10px){.user{background-color:seagreen;}}}
@media (max-width:600px){.user{color:red;}}
.user:hover{color:orange;}
`))
	})

	t.Run("@media specifity", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
> #box-not-working {
  background:red;
  padding-left:8px;
  width:10px;
  @media only screen and (min-width:10px) {
	width: calc(10px + 90px * (100vw - 10px) / 90);
  }
  @media only screen and (min-width:90px) {
	width: 90px;
  }
  height: 10px;
  @media only screen and (min-width:10px) {
	height: calc(10px + 90px * (100vw - 10px) / 90);
  }
  @media only screen and (min-width:90px) {
	height: 90px;
  }
}
`))).To(BeCSS(`
.user > #box-not-working{background:red;padding-left:8px;width:10px;height:10px;}
@media only screen and (min-width:10px){.user > #box-not-working{width:calc(10px + 90px * (100vw - 10px) / 90);}}
@media only screen and (min-width:90px){.user > #box-not-working{width:90px;}}
@media only screen and (min-width:10px){.user > #box-not-working{height:calc(10px + 90px * (100vw - 10px) / 90);}}
@media only screen and (min-width:90px){.user > #box-not-working{height:90px;}}
`))
	})

	t.Run("@font-face", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
@font-face {
  font-family:Pangolin;
  src:url('Pangolin-Regular.ttf') format('truetype');
}
`))).To(BeCSS(`
@font-face{font-family:Pangolin;src:url('Pangolin-Regular.ttf') format('truetype');}
`))
	})

	t.Run("multiple selectors", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
span, h1 {
  color:red;
}
h1, &:after, &:before {
  color:red;
}
`))).To(BeCSS(`
.user span,.user h1{color:red;}
.user h1,.user:after,.user:before{color:red;}
`))
	})

	t.Run("[title=\"a,b\"] and :matches(a,b)", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
.test:matches(a,b,c), .test {
  color:blue;
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
`))).To(BeCSS(`
.user .test:matches(a,b,c),.user .test{color:blue;}
.user .test[title=","]{color:red;}
.user [title="a,b,c, something"],.user h1,.user [title="a,b,c"]{color:red;}
.user [title="a"],.user [title="b"]{color:red;}
`))
	})

	t.Run("quotes", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
.foo:before {
  content:".hello {world}";
  content:".hello {world} ' ";
  content:'.hello {world} " ';
}
.foo:before {
  content:"\"";
  content:"\\\"";
  content:'\'';
  content:'\\\'';
}
`))).To(BeCSS(`
.user .foo:before{content:".hello {world}";content:".hello {world} ' ";content:'.hello {world} " ';}
.user .foo:before{content:"\"";content:"\\\"";content:'\'';content:'\\\'';}
`))
	})

	t.Run("remove empty css", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`& {   }`))).To(BeCSS(``))
	})

	t.Run("remove empty declarations", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
width:;
height:;
max-width:     ;
max-height:     ;
`))).To(BeCSS(``))
	})

	t.Run("urls", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
background:url(http://url.com/});
background:url(http://url.com//1234) '(';
background:url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEYAAAABCAIAAADsEU8HAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAIklEQVQI12P8//8/Aw4wbdq0rKysAZG1trbGJXv06FH8sgDIJBbBfp+hFAAAAABJRU5ErkJggg==");
`))).To(BeCSS(`
.user{background:url(http://url.com/});background:url(http://url.com//1234) '(';background:url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEYAAAABCAIAAADsEU8HAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAIklEQVQI12P8//8/Aw4wbdq0rKysAZG1trbGJXv06FH8sgDIJBbBfp+hFAAAAABJRU5ErkJggg==");}
`))
	})

	t.Run("last semicolon omission", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
.content {
  color:red
}
.content {
  color:blue
}
`))).To(BeCSS(`
.user .content{color:red;}
.user .content{color:blue;}
`))
	})

	t.Run(":matches(:not())", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
h1:matches(.a,.b,:not(.c)) {
  display: none
}
`))).To(BeCSS(`
.user h1:matches(.a,.b,:not(.c)){display:none;}
`))
	})

	t.Run("@keyframes", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
@-webkit-keyframes slidein {
	to { transform:translate(20px); }
}
  @keyframes slidein {
			to { transform:translate(20px); }
  }
  @keyframes hahaha {
		  0%,1%{t:0}
  }
  @keyframes infinite-spinning {
	from {
	  transform: rotate(0deg);
	}
	to {
	  transform: rotate(360deg);
	}
  }
`))).To(BeCSS(
			`@-webkit-keyframes slidein{to{transform:translate(20px);}}`,
			`@keyframes slidein{to{transform:translate(20px);}}`,
			`@keyframes hahaha{0%,1%{t:0;}}`,
			`@keyframes infinite-spinning{from{transform:rotate(0deg);}to{transform:rotate(360deg);}}`,
		))
	})

	t.Run("edge cases", func(t *testing.T) {
		NewWithT(t).Expect(Parse(withUser(`
@media (min-width:537px) {
  border-bottom:4px solid red;
}
&::placeholder {
  color:pink;
}
.a {color:'red'}
.b {color:"red"}
.a {color:red;}[role=button]{color:red;}
.b {padding:30 3}
.c {v-text-anchor: middle;}
`))).To(BeCSS(`
@media (min-width:537px){.user{border-bottom:4px solid red;}}
.user::placeholder{color:pink;}
.user .a{color:'red';}
.user .b{color:"red";}
.user .a{color:red;}
.user [role=button]{color:red;}
.user .b{padding:30 3;}
.user .c{v-text-anchor:middle;}
`))
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
