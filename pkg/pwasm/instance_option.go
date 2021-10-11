package pwasm

import (
	"log"
	"os"
	"regexp"

	"github.com/go-courier/gox/pkg/pwasm/core"
	"github.com/go-courier/gox/pkg/pwasm/esloader"
	"github.com/go-courier/gox/pkg/pwasm/goloader"
	"github.com/go-courier/gox/pkg/pwasm/htmlloader"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

type Option struct {
	Cwd        string
	RootDir    string
	EntryPoint string
	OutputDir  string
	PublicPath string
	Rules      []*Rule
	Logger     logr.Logger
}

func (opt *Option) SetDefaults() {
	if opt.EntryPoint == "" {
		opt.EntryPoint = "index.html"
	}

	if opt.Rules == nil {
		opt.Rules = []*Rule{
			{
				Test: regexp.MustCompile("(\\.tsx?|\\.jsx?|\\.mjsx?)$"),
				Loaders: []core.Loader{
					esloader.ESLoader{},
				},
			},
			{
				Test: regexp.MustCompile("\\.html$"),
				Loaders: []core.Loader{
					htmlloader.HTMLLoader{},
				},
			},
			{
				Test: regexp.MustCompile("\\.go$"),
				Loaders: []core.Loader{
					goloader.GoLoader{},
				},
			},
		}
	}

	if opt.Logger.GetSink() == nil {
		opt.Logger = opt.Logger.WithSink(stdr.New(log.New(os.Stderr, "[pwasm] ", log.LstdFlags|log.Lshortfile)).GetSink())
	}
}

type Rule struct {
	Test    *regexp.Regexp
	Loaders []core.Loader
}
