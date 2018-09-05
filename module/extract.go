package module

import (
	"github.com/fatih/structtag"
	"github.com/lyft/protoc-gen-star"
	"github.com/srikrsna/protoc-gen-gotag/tagger"
)

type tagExtractor struct {
	pgs.Visitor
	pgs.DebuggerCommon

	tags map[string]map[string]*structtag.Tags
}

func newTagExtractor(d pgs.DebuggerCommon) *tagExtractor {
	v := &tagExtractor{DebuggerCommon: d}
	v.Visitor = pgs.PassThroughVisitor(v)
	return v
}

func (v *tagExtractor) VisitField(f pgs.Field) (pgs.Visitor, error) {
	var tval string
	if ok, err := f.Extension(tagger.E_Tags, &tval); !ok || err != nil {
		v.Fail("tagger extension not found or malformed")
	}

	tags, err := structtag.Parse(tval)
	v.CheckErr(err)

	msgName := f.Message().Name().PGGUpperCamelCase().String()
	if v.tags[msgName] == nil {
		v.tags[msgName] = map[string]*structtag.Tags{}
	}

	v.tags[msgName][f.Name().PGGUpperCamelCase().String()] = tags

	return v, nil
}

func (v *tagExtractor) Extract(f pgs.File) map[string]map[string]*structtag.Tags {
	v.tags = map[string]map[string]*structtag.Tags{}

	v.CheckErr(pgs.Walk(v, f))

	return v.tags
}
