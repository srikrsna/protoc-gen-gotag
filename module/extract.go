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

func (v *tagExtractor) VisitOneOf(o pgs.OneOf) (pgs.Visitor, error) {
	var tval string
	ok, err := o.Extension(tagger.E_Tag, &tval)
	if err != nil {
		return nil, err
	}

	if !ok {
		return v, nil
	}

	tags, err := structtag.Parse(tval)
	if err != nil {
		return nil, err
	}

	msgName := o.Message().Name().PGGUpperCamelCase().String()

	if v.tags[msgName] == nil {
		v.tags[msgName] = map[string]*structtag.Tags{}
	}

	v.tags[msgName][o.Name().PGGUpperCamelCase().String()] = tags

	return v, nil
}

func (v *tagExtractor) VisitField(f pgs.Field) (pgs.Visitor, error) {
	var tval string
	ok, err := f.Extension(tagger.E_Tags, &tval)
	if err != nil {
		return nil, err
	}

	if !ok {
		return v, nil
	}

	tags, err := structtag.Parse(tval)
	v.CheckErr(err)

	msgName := f.Message().Name().PGGUpperCamelCase().String()

	if f.InOneOf() {
		msgName = f.Message().Name().PGGUpperCamelCase().String() + "_" + f.Name().PGGUpperCamelCase().String()
	}

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
