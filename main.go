package main

import (
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/srikrsna/protoc-gen-gotag/module"
)

func main() {
	opt := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	pgs.Init(
		pgs.DebugEnv("GOTAG_DEBUG"),
		pgs.SupportedFeatures(&opt),
	).
		RegisterModule(module.New()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
