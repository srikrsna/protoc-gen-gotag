package main

import (
	"github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	"github.com/srikrsna/protoc-gen-gotag/module"
)

func main() {
	pgs.Init(pgs.DebugMode()).RegisterModule(module.New()).RegisterPostProcessor(pgsgo.GoFmt()).Render()
}
