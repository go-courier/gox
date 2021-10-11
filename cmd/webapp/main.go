package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	. "github.com/go-courier/gox/pkg/css"
	. "github.com/go-courier/gox/pkg/dom"
	. "github.com/go-courier/gox/pkg/gox"
	"github.com/go-courier/gox/pkg/gox/renderer"
	"github.com/go-courier/gox/pkg/httputil"
)

func withCSSCache(c *CSSCache) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return ContextWithCSSCache(ctx, c)
	}
}

var css = NewCSSCache("app", Document.QuerySelector("head"))

var spanList = make([]interface{}, 0)

func init() {
	for i := 0; i < 10000; i++ {
		spanList = append(spanList, Span(fmt.Sprintf(" %d ", i)))
	}
}

func UseOn(ctx context.Context, eventName string, handle func(event Event)) *Ref {
	elRef := UseRef(ctx, nil)
	handleRef := UseRef(ctx, nil)
	handleRef.Current = handle

	UseEffect(ctx, func() func() {
		if n, ok := elRef.Current.(Node); ok {
			h := func(event Event) {
				handleRef.Current.(func(event Event))(event)
			}

			n.AddEventListener(eventName, h)
			return func() {
				n.RemoveEventListener(eventName, h)
			}
		}
		return nil
	}, []interface{}{elRef.Current})

	return elRef
}

type SubComp struct {
}

func (a SubComp) Render(ctx context.Context, children ...interface{}) interface{} {
	refInput := UseOn(ctx, "click", func(event Event) {
		fmt.Println("click")
	})

	return Button(refInput, "button")
}

type App struct {
}

func (a App) Render(ctx context.Context, children ...interface{}) interface{} {
	s := time.Now()
	defer func() {
		fmt.Println("render app cost", time.Since(s))
	}()

	UseEffect(ctx, func() func() {
		fmt.Println("render app full cost", time.Since(s))
		return nil
	}, nil)

	value, setValue := UseState(ctx, "")
	v, setValue2 := UseState(ctx, &[]string{})
	emojis := *v.(*[]string)

	hello := UseMemo(ctx, func() interface{} {
		return "Hello"
	}, []interface{}{}).(string)

	UseEffect(ctx, func() func() {
		go func() {
			c := httputil.GetShortConnClientContext(context.Background(), 5*time.Second)
			req, _ := http.NewRequest("GET", "https://dog.ceo/api/breeds/list/all", nil)
			resp, _ := c.Do(req)
			m := struct {
				Message map[string]interface{} `json:"message"`
			}{}

			_ = json.NewDecoder(resp.Body).Decode(&m)

			keys := make([]string, 0, len(m.Message))

			for i := range m.Message {
				keys = append(keys, i)
			}

			setValue2(&keys)
		}()

		return nil
	}, []interface{}{})

	refInput := UseOn(ctx, "input", func(event Event) {
		setValue(func(v interface{}) interface{} {
			return event.Target().(Element).Get("value")
		})
	})

	return Provider(withCSSCache(css))(
		Main(
			CSS{
				"backgroundColor": "#ddd",
			},
			Div(
				Div(hello),
			),
			Div(
				CSS{
					"overflow": "auto",
				},
				Span("Input:"),
				Span(value),
			),
			Div(
				Input(
					refInput,
					Attr("value", value),
				),
			),
			H(SubComp{})(),
			Div(
				CSS{
					"overflow": "scroll",
					"height":   "2em",
				},
				RangeSlice(emojis, func(i int) interface{} {
					return Div(
						Key(emojis[i]),
						emojis[i],
					)
				}),
			),
			Div(
				CSS{
					"overflow": "scroll",
					"height":   "20em",
					"fontSize": "8px",
				},
				Fragment(spanList),
			),
		),
		renderer.Portal(Document.QuerySelector("#portal"))(
			Div(
				CSS{
					"overflow": "auto",
				},
				"in portal", value,
			),
		),
	)
}

func main() {
	r := renderer.CreateRoot(Document.QuerySelector("#root"))
	_ = r.Render(context.Background(), H(App{})())
	fmt.Println("App ready")
	<-make(chan struct{})
}
