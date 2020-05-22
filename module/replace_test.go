package module_test

import (
	"bytes"
	"flag"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/fatih/structtag"

	"github.com/srikrsna/protoc-gen-gotag/module"
)

var replaceOut = flag.Bool("tag-rep", false, "")

func TestRetag(t *testing.T) {
	fs := token.NewFileSet()

	n, err := parser.ParseFile(fs, "./test/input.txt", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	module.Retag(n, map[string]map[string]*structtag.Tags{
		"Simple": map[string]*structtag.Tags{
			"Single":   tagMust(structtag.Parse(`sql:"-,omitempty"`)),
			"Multiple": tagMust(structtag.Parse(`xml:"-,omitempty"`)),
			"None":     tagMust(structtag.Parse(`json:"none,omitempty"`)),
		},
	})

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fs, n); err != nil {
		t.Fatal(err)
	}

	if *replaceOut {
		f, err := os.Create("./test/golden.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		if _, err := io.Copy(f, &buf); err != nil {
			t.Fatal(err)
		}

		return
	}

	out, err := ioutil.ReadFile("./test/golden.txt")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(out, buf.Bytes()) {
		t.Error("output doesnot match golden file")
	}
}

func tagMust(t *structtag.Tags, err error) *structtag.Tags {
	if err != nil {
		panic(err)
	}
	return t
}
