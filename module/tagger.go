package module

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"

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
	extractor := newTagExtractor(m)
	for _, f := range target.Files() {
		tags := extractor.Extract(f)

		gfname := f.OutputPath().SetExt(".go").String()
		
		fs := token.NewFileSet()
		fn, err := parser.ParseFile(fs, gfname, nil, parser.ParseComments)
		m.CheckErr(err)

		ast.Walk(replacer{m, tags}, fn)

		var buf strings.Builder
		m.CheckErr(printer.Fprint(&buf, fs, fn))

		m.OverwriteGeneratorFile(gfname, buf.String())
	}

	return m.Artifacts()
}
