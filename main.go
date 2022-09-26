package main

import (
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/srikrsna/protoc-gen-gotag/module"
)

func main() {
	var supportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	pgs.Init(pgs.DebugEnv("GOTAG_DEBUG"), pgs.SupportedFeatures(&supportedFeatures)).
		RegisterModule(module.New()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
