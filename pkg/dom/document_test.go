package dom

import (
	"bytes"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNode(t *testing.T) {
	p := Document.CreateElement("parent")
	div1 := Document.CreateElement("div1")
	div2 := Document.CreateElement("div2")

	t.Run("CRUD", func(t *testing.T) {

		buf := bytes.NewBuffer(nil)

		{
			p.AppendChild(div1)

			buf.Reset()
			RenderToHTML(buf, p)
			NewWithT(t).Expect(buf.String()).To(Equal("<parent><div1></div1></parent>"))
		}

		{
			p.InsertBefore(div2, div1)

			buf.Reset()
			RenderToHTML(buf, p)
			NewWithT(t).Expect(buf.String()).To(Equal("<parent><div2></div2><div1></div1></parent>"))
		}

		{
			p.RemoveChild(div2)

			buf.Reset()
			RenderToHTML(buf, p)
			NewWithT(t).Expect(buf.String()).To(Equal("<parent><div1></div1></parent>"))
		}

		{
			p.InsertBefore(div2, nil)

			buf.Reset()
			RenderToHTML(buf, p)
			NewWithT(t).Expect(buf.String()).To(Equal("<parent><div1></div1><div2></div2></parent>"))
		}
	})
}
