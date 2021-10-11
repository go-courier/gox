package pwasm

import (
	"bytes"
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-courier/gox/pkg/pwasm/core"
)

func NewInstance(o Option) *Instance {
	opts := &o
	opts.SetDefaults()

	return &Instance{
		Option: opts,
	}
}

type Instance struct {
	*Option
	core.SourceFiles
}

func (inst *Instance) Build() error {
	rootDir := filepath.Join(inst.Cwd, inst.Option.RootDir)
	outputDir := filepath.Join(inst.Cwd, inst.Option.OutputDir)

	_, err := inst.Resolve(context.Background(), inst.EntryPoint, filepath.Join(rootDir, ".root"), false)
	if err != nil {
		return err
	}

	filesToWrite := map[string]bool{}

	for _, f := range inst.Manifests() {
		if err := writeFile(filepath.Join(outputDir, f.URL), f.SourceFile.OutputCode()); err != nil {
			return err
		}
		filesToWrite[f.URL] = true
	}

	return nil
}

func (inst *Instance) Manifests() []core.Manifest {
	entryPoint := filepath.Join(inst.Cwd, inst.Option.RootDir, inst.EntryPoint)
	return inst.ListManifests(entryPoint, inst)
}

func (inst *Instance) Load(filepath string, atFile string) (string, []byte, error) {
	filename := resolve(filepath, atFile)

	data, err := os.ReadFile(filename)
	if err != nil {
		inst.Logger.Error(err, "load file failed")
		return "", nil, err
	}

	return filename, data, err
}

func resolve(fp string, atFile string) string {
	if fp == "" || fp[0] == '/' {
		return fp
	}
	return filepath.Join(filepath.Dir(atFile), fp)
}

func (inst *Instance) Resolve(ctx context.Context, filepath string, atFile string, async bool) (string, error) {
	filename, data, err := inst.Load(filepath, atFile)
	if err != nil {
		return "", err
	}

	ctx = core.ContextWithProcessor(ctx, inst)

	sf := core.NewSourceFile(filename, data, async)

	output, err := inst.RoundTrip(ctx, sf)
	if err != nil {
		return "", err
	}
	if _, ok := output.(*core.SourceFile); !ok {
		sf.Transformed = output
	}

	sf.OutputFilename = core.OutputFilename(output)

	inst.PutSourceFile(sf)

	return inst.OutputRequestURL(sf.OutputFilename), nil
}

func (inst *Instance) RoundTrip(ctx context.Context, output core.Output) (core.Output, error) {
	for i := range inst.Rules {
		rule := inst.Rules[i]
		f := output.Filename()

		if rule.Test.MatchString(f) {
			for i := range rule.Loaders {
				r, err := rule.Loaders[i].Transform(ctx, f, bytes.NewReader(output.Code()))
				if err != nil {
					inst.Logger.Error(err, "transform failed")
					return nil, err
				}
				output = r
			}
			break
		}
	}

	return output, nil
}

func (inst *Instance) OutputRequestURL(f string) string {
	rootDir := filepath.Join(inst.Cwd, inst.RootDir)
	var rel string

	// under root dir
	if strings.HasPrefix(f, rootDir) {
		rel, _ = filepath.Rel(filepath.Join(inst.Cwd, inst.RootDir), f)
	} else {
		rel = path.Join("vendor", filepath.Base(f))
	}

	if inst.PublicPath != "" {
		return path.Join(inst.PublicPath, rel)
	}

	return rel

}
