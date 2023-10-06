package module_test

import (
	"bytes"
	"os"
	"testing"

	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/spf13/afero"
	"github.com/srikrsna/protoc-gen-gotag/module"
)

func TestExtract(t *testing.T) {
	req, err := os.Open("../debug/code_generator_request.pb.bin")
	if err != nil {
		t.Fatal(err)
	}

	fs := afero.NewMemMapFs()
	res := &bytes.Buffer{}

	pgs.Init(
		pgs.ProtocInput(req),
		pgs.ProtocOutput(res),
		pgs.FileSystem(fs),
	).RegisterModule(module.New()).Render()
}
