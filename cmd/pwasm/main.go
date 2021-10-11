package main

import (
	"flag"
	"github.com/go-courier/gox/pkg/pwasm"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	input := flag.Arg(0)

	o := pwasm.Option{}

	o.Cwd, _ = os.Getwd()
	o.RootDir = input
	o.PublicPath = "/__built__/"

	if filepath.Ext(input) != "" {
		o.RootDir = filepath.Dir(input)
	}

	o.OutputDir = filepath.Join(o.RootDir, "public")

	inst := pwasm.NewInstance(o)

	_ = inst.Build()

	//cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" $(WEBAPP_PUBLIC)/__built__/wasm_exec.js
	//GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o $(WEBAPP_PUBLIC)/__built__/genshin-level-up.wasm $(WEBAPP_SRC_ROOT)
}
