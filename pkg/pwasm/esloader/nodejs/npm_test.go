package nodejs

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/onsi/gomega"
)

func TestNPM(t *testing.T) {
	_ = os.Chdir("./testdata")
	cwd, _ := os.Getwd()

	_ = os.Symlink(filepath.Join(cwd, "nmodules/demo"), filepath.Join(cwd, "nmodules/demo2"))

	r := &NodeResolver{
		Modules: []string{
			"nmodules",
			"node_modules",
		},
	}
	r.SetDefaults()

	cases := []struct {
		name     string
		path     string
		resolved string
	}{
		{
			"related path without index.js",
			"./lib-without-index",
			"",
		},
		{
			"related path with index.js",
			"./lib-with-index",
			filepath.Join(cwd, "./lib-with-index/index.js"),
		},
		{
			"self with package.json",
			"./",
			filepath.Join(cwd, "./lib-without-index/main.js"),
		},
		{
			"in modules",
			"demo",
			filepath.Join(cwd, "./nmodules/demo/lib/main.js"),
		},
		{
			"in modules with symlinks",
			"demo2",
			filepath.Join(cwd, "./nmodules/demo2/lib/main.js"),
		},
		{
			"not in any modules",
			"__not_exist_module__",
			"",
		},
	}

	root := path.Join(cwd, ".root")

	for i := range cases {
		c := cases[i]
		t.Run(c.name, func(t *testing.T) {
			resolved, _ := r.Resolve(c.path, root)
			gomega.NewWithT(t).Expect(resolved).To(gomega.Equal(c.resolved))
		})
	}
}
