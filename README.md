# GoX

[Experimental] React-like library for Golang in WASM

## Features

* Virtual DOM and HTML DSL with function calls.
* CSS in Go like [Emotion JS](https://github.com/emotion-js/emotion) did.
* `Fragment` && `Portal` supports.
* Component support as `interface { Render(ctx context.Context, childen ...interface{}) interface{}}`.
* Basic hooks support `UseState`, `UseEffect`, `UseMemo`, `UseRef`
    * `UseContext` not needed in Go, the `context.Context` will pass into Component
* Request HTTP in web worker by XHR

## Known Issues

* Compiled to wasm
    * bundle size is so big, but could be compressed by `gzip`/`brotli` (brotli is better)
    