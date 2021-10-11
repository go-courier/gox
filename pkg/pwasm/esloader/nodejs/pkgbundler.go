package nodejs

import (
	"github.com/davecgh/go-spew/spew"
	"os"
)

type PkgBundler struct {
	NodeResolver
}


func (b *PkgBundler) Bundle(name string) {
	cwd, _ := os.Getwd()
	p, err := b.Resolve(name, cwd)
	spew.Dump(p, err)
}
