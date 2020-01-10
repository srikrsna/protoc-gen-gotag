package module

import (
	"github.com/fatih/structtag"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"github.com/taankyou/protoc-gen-gotag/tagger"
	"log"
)

type tagExtractor struct {
	pgs.Visitor
	pgs.DebuggerCommon
	pgsgo.Context

	tags map[string]map[string]*structtag.Tags
	autoAddTags []string
}

func newTagExtractor(d pgs.DebuggerCommon, ctx pgsgo.Context) *tagExtractor {
	v := &tagExtractor{DebuggerCommon: d, Context: ctx}
	v.Visitor = pgs.PassThroughVisitor(v)
	return v
}

func (v *tagExtractor) VisitOneOf(o pgs.OneOf) (pgs.Visitor, error) {
	var tval string
	ok, err := o.Extension(tagger.E_OneofTags, &tval)
	if err != nil {
		return nil, err
	}

	msgName := v.Context.Name(o.Message()).String()

	if v.tags[msgName] == nil {
		v.tags[msgName] = map[string]*structtag.Tags{}
	}

	if !ok {
		return v, nil
	}

	tags, err := structtag.Parse(tval)
	if err != nil {
		return nil, err
	}

	v.tags[msgName][v.Context.Name(o).String()] = tags

	return v, nil
}

func (v *tagExtractor) VisitField(f pgs.Field) (pgs.Visitor, error) {
	var tval string
	ok, err := f.Extension(tagger.E_Tags, &tval)
	if err != nil {
		return nil, err
	}

	msgName := v.Context.Name(f.Message()).String()
	if f.InOneOf() {
		msgName = f.Message().Name().UpperCamelCase().String() + "_" + f.Name().UpperCamelCase().String()
	}

	if v.tags[msgName] == nil {
		v.tags[msgName] = map[string]*structtag.Tags{}
	}

	if !ok && len(v.autoAddTags) == 0 {
		return v, nil
	}

	tags := structtag.Tags{}

	if !ok {
		val := ToSnakeCase(v.Context.Name(f).String())
		for _, tag := range v.autoAddTags {
			t := structtag.Tag{
				Key:     tag,
				Name:    val,
				Options: nil,
			}
			if err := tags.Set(&t); err != nil {
				log.Fatal(err)
			}
		}
		v.tags[msgName][v.Context.Name(f).String()] = &tags
	}

	newTags, err := structtag.Parse(tval)
	v.CheckErr(err)
	for _, tag := range newTags.Tags() {
		if err := tags.Set(tag); err != nil {
			log.Fatal(err)
		}
	}

	v.tags[msgName][v.Context.Name(f).String()] = &tags

	return v, nil
}

func (v *tagExtractor) Extract(f pgs.File) StructTags {
	v.tags = map[string]map[string]*structtag.Tags{}

	v.CheckErr(pgs.Walk(v, f))

	return v.tags
}
