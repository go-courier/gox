package css

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/davecgh/go-spew/spew"
)

func TestCache(t *testing.T) {
	c := CSSCache{}

	s := c.SerializeStyles(
		context.Background(),
		CSS{
			"color": "red",
			"&:focus": CSS{
				"color": "blue",
			},
		},
		CSS{
			"fontSize": "10px",
		},
	)

	spew.Dump(s)
}
