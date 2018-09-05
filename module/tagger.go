package module

import (
	"go/parser"
	"go/printer"
	"go/token"
	"strings"

	"github.com/fatih/structtag"
	"github.com/lyft/protoc-gen-star"
)

type mod struct {
	*pgs.ModuleBase
}

func New() pgs.Module {
	return &mod{&pgs.ModuleBase{}}
}

func (mod) Name() string {
	return "gotag"
}

func (m mod) Execute(target pgs.Package, packages map[string]pgs.Package) []pgs.Artifact {

	xtv := m.Parameters().Str("xxx")

	xtv = strings.Replace(xtv, "+", ":", -1)

	xt, err := structtag.Parse(xtv)
	m.CheckErr(err)

	extractor := newTagExtractor(m)
	for _, f := range target.Files() {
		tags := extractor.Extract(f)

		tags.AddTagsToXXXFields(xt)

		gfname := f.OutputPath().SetExt(".go").String()

		fs := token.NewFileSet()
		fn, err := parser.ParseFile(fs, gfname, nil, parser.ParseComments)
		m.CheckErr(err)

		m.CheckErr(Retag(fn, tags))

		var buf strings.Builder
		m.CheckErr(printer.Fprint(&buf, fs, fn))

		m.OverwriteGeneratorFile(gfname, buf.String())
	}

	return m.Artifacts()
}
