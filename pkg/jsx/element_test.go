package jsx

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-courier/jsx/pkg/css"
)

func TestJSX(t *testing.T) {
	c := css.NewCache()

	e := Provide(func(ctx context.Context) context.Context {
		return css.ContextWithCache(ctx, c)
	})(Div(
		CSS{
			"color": "red",
		},
		Div(
			Attr("style", "color:blue;"),
		),
	))

	_ = Render(context.Background(), e, os.Stdout)
	_, _ = fmt.Fprintf(os.Stdout, "\n")
}
