package goloader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-courier/gox/pkg/pwasm/core"
	"github.com/pkg/errors"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type GoLoader struct {
}

func (GoLoader) Transform(ctx context.Context, filename string, src io.Reader) (core.Output, error) {
	p := core.ProcessorFromContext(ctx)

	dir := filepath.Dir(filename)
	parts := strings.Split(dir, "/")
	appName := parts[len(parts)-1]
	outputWASM := filepath.Join(filepath.Dir(filename), "bin/"+appName+".wasm")

	err := RunScript(fmt.Sprintf("GOOS=js GOARCH=wasm go build -o %s %s", outputWASM, filename))
	if err != nil {
		return nil, errors.Wrapf(err, "go build wasm failed")
	}

	resolved, err := p.Resolve(ctx, outputWASM, filename, false)
	if err != nil {
		return nil, err
	}

	wasmExec, err := os.ReadFile(filepath.Join(runtime.GOROOT(), "./misc/wasm/wasm_exec.js"))
	if err != nil {
		return nil, err
	}

	// mark to dead code to disable commonjs
	wasmExec = bytes.ReplaceAll(wasmExec, []byte(`global.require = require;`), []byte(""))
	wasmExec = bytes.ReplaceAll(wasmExec, []byte(`global.require`), []byte("false"))
	wasmExec = bytes.ReplaceAll(wasmExec, []byte(`require`), []byte("__require"))

	sf := core.NewSourceFile(filename+".ts", append(wasmExec, []byte(`
export const main = () => {
    const go = new Go();

    WebAssembly
        .instantiateStreaming(
            fetch("`+resolved+`"),
            go.importObject
        )
        .then((result) => go.run(result.instance));
}`)...), false)

	o, err := p.RoundTrip(ctx, sf)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func RunScript(script string) (e error) {
	sh, err := syntax.NewParser().Parse(strings.NewReader(script), "")
	if err != nil {
		return err
	}

	runner, err := interp.New(interp.StdIO(os.Stdin, os.Stdout, os.Stderr))
	if err != nil {
		return err
	}

	return runner.Run(context.Background(), sh)
}
