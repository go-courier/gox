package esloader

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-courier/gox/pkg/pwasm/esloader/nodejs"
	"io"
	"path/filepath"
	"strconv"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/go-courier/gox/pkg/pwasm/core"
)

type ESLoader struct {
}

var loaders = map[string]api.Loader{
	".ts":  api.LoaderTS,
	".tsx": api.LoaderTSX,
	".js":  api.LoaderJS,
	".jsx": api.LoaderJSX,
	//".css":  api.LoaderCSS,
	//".wasm": api.LoaderBinary,
	//".go":   api.LoaderBinary,
}

func (ESLoader) Transform(ctx context.Context, filename string, r io.Reader) (core.Output, error) {
	code, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	p := core.ProcessorFromContext(ctx)

	const minify = true

	defines := map[string]string{
		"__filename": strconv.Quote(filename),
		"__dirname":  strconv.Quote(filepath.Dir(filename)),
		"__manifests": func() string {
			manifestOrigin := p.Manifests()
			manifests := make([]core.Manifest, 0, len(manifestOrigin))
			for i := range manifestOrigin {
				m := manifestOrigin[i]
				if !m.SourceFile.Async {
					manifests = append(manifests, m)
				}
			}
			d, _ := json.Marshal(manifests)
			return string(d)
		}(),
	}

	opt := api.BuildOptions{
		Write:             false,
		Bundle:            true,
		Target:            api.ES2017,
		Format:            api.FormatESModule,
		JSXMode:           api.JSXModeTransform,
		MinifyWhitespace:  minify,
		MinifyIdentifiers: minify,
		MinifySyntax:      minify,
		Define:            defines,
	}

	opt.Loader = loaders

	opt.Stdin = &api.StdinOptions{
		Sourcefile: filename,
		Loader:     api.LoaderTS,
		Contents:   string(code),
		ResolveDir: filepath.Dir(filename),
	}

	nodeResolver := &nodejs.NodeResolver{}
	nodeResolver.SetDefaults()

	opt.Plugins = []api.Plugin{
		{
			Name: "module-resolver",
			Setup: func(build api.PluginBuild) {
				build.OnResolve(
					api.OnResolveOptions{Filter: ".*"},

					func(args api.OnResolveArgs) (api.OnResolveResult, error) {
						if nodejs.IsBuiltinNodeModules(args.Path) {
							return api.OnResolveResult{External: true}, nil
						}

						srcFile := args.Path

						if resolvedSrcPath, ok := nodeResolver.Resolve(args.Path, args.ResolveDir+"/package.json"); ok {
							srcFile = resolvedSrcPath

							return api.OnResolveResult{
								External: true,
								//Path:     resolvedPath,
							}, nil
						}

						fmt.Println(srcFile)

						resolvedPath, err := p.Resolve(ctx, srcFile, filename, args.Kind == api.ResolveJSDynamicImport)
						if err != nil {
							return api.OnResolveResult{}, err
						}

						return api.OnResolveResult{
							External: true,
							Path:     resolvedPath,
						}, nil
					},
				)
			},
		},
	}

	tr := api.Build(opt)

	if len(tr.OutputFiles) > 0 {
		return &outputES{
			filename: filename,
			contents: tr.OutputFiles[0].Contents,
		}, nil
	}

	return nil, errors.New(fmt.Sprintf("%#v", tr.Errors))

}

type outputES struct {
	filename string
	contents []byte
}

func (o *outputES) Filename() string {
	switch ext := filepath.Ext(o.filename); ext {
	case ".ts", ".tsx":
		b := filepath.Base(o.filename)
		return filepath.Join(filepath.Dir(o.filename), b[0:len(b)-len(ext)]+".js")
	}
	return o.filename
}

func (o *outputES) Code() []byte {
	return o.contents
}
