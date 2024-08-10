package module

import (
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/structtag"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

type mod struct {
	*pgs.ModuleBase
	pgsgo.Context
}

func New() pgs.Module {
	return &mod{ModuleBase: &pgs.ModuleBase{}}
}

func (m *mod) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.Context = pgsgo.InitContext(c.Parameters())
}

func (mod) Name() string {
	return "gotag"
}

func (m mod) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	xtv := m.Parameters().Str("xxx")

	xtv = strings.Replace(xtv, "+", ":", -1)

	xt, err := structtag.Parse(xtv)
	m.CheckErr(err)

	autoTag := m.Parameters().Str("auto")
	var autoTags []string
	if autoTag != "" {
		autoTags = strings.Split(autoTag, "+")
	}

	module := m.Parameters().Str("module")
	outdir := m.Parameters().Str("outdir")
	outAbsDir := m.Parameters().Str("outabsdir")

	extractor := newTagExtractor(m, m.Context, autoTags)

	for _, f := range targets {
		tags := extractor.Extract(f)

		tags.AddTagsToXXXFields(xt)

		gfname := m.Context.OutputPath(f).SetExt(".go").String()

		filename := gfname
		if len(outAbsDir) > 0 {
			filename = filepath.Join(outAbsDir, gfname)
		} else if outdir != "" {
			filename = filepath.Join(outdir, gfname)
		}

		if module != "" {
			filename = strings.ReplaceAll(filename, string(filepath.Separator), "/")
			trim := module + "/"
			if !strings.HasPrefix(filename, trim) {
				m.Debug(fmt.Sprintf("%v: generated file does not match prefix %q", filename, module))
				m.Exit(1)
			}
			filename = strings.TrimPrefix(filename, trim)
		}

		fs := token.NewFileSet()
		fn, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
		m.CheckErr(err)

		m.CheckErr(Retag(fn, tags))

		var buf strings.Builder
		m.CheckErr(printer.Fprint(&buf, fs, fn))

		if len(outAbsDir) > 0 {
			m.OverwriteCustomFile(filename, buf.String(), os.ModePerm)
		} else {
			m.OverwriteGeneratorFile(filename, buf.String())
		}
	}

	return m.Artifacts()
}
