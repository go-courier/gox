package css

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestCache(t *testing.T) {
	c := Cache{}

	s := c.SerializeStyles(
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
