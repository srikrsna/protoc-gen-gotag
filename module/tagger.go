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
	xtags := map[string]*structtag.Tags{}

	xtv = strings.Replace(xtv, "+", ":", -1)

	tags, err := structtag.Parse(xtv)
	m.CheckErr(err)

	if xtv != "" {
		xtags["XXX_NoUnkeyedLiteral"] = tags
		xtags["XXX_unrecognized"] = tags
		xtags["XXX_sizecache"] = tags
	}

	extractor := newTagExtractor(m)
	for _, f := range target.Files() {
		tags := extractor.Extract(f)

		for o := range tags {
			if tags[o] == nil {
				tags[o] = map[string]*structtag.Tags{}
			}

			for k, v := range xtags {
				tags[o][k] = v
			}
		}

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
