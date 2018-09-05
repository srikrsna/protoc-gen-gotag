package module

import (
	"go/ast"
	"strings"

	"github.com/fatih/structtag"
	"github.com/lyft/protoc-gen-star"
)

type replacer struct {
	pgs.DebuggerCommon
	tags map[string]map[string]*structtag.Tags
}

func (v replacer) Visit(n ast.Node) ast.Visitor {
	if tp, ok := n.(*ast.TypeSpec); ok {
		if _, ok := tp.Type.(*ast.StructType); ok {
			ast.Walk(retagger{DebuggerCommon: v, tags: v.tags[tp.Name.String()]}, n)
			return nil
		}
	}
	return v
}

type retagger struct {
	pgs.DebuggerCommon
	tags map[string]*structtag.Tags
}

func (v retagger) Visit(n ast.Node) ast.Visitor {
	if f, ok := n.(*ast.Field); ok {
		newTags := v.tags[f.Names[0].String()]
		if newTags == nil {
			return nil
		}

		oldTags, err := structtag.Parse(strings.Trim(f.Tag.Value, "`"))
		v.CheckErr(err)

		for _, t := range newTags.Tags() {
			oldTags.Set(t)
		}

		f.Tag.Value = "`" + oldTags.String() + "`"

		return nil
	}

	return v
}
