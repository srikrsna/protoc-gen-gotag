LOCAL_PATH = $(shell pwd)

.PHONY: install gen-tag test protobuf gotags genTags

# Generate base protobuf output files.
protobuf:
	protoc -I /usr/local/include \
	-I ${LOCAL_PATH} \
	--go_out=:. example/example.proto

# Add the gotags to the protobuf output files.
gotags:
	protoc -I /usr/local/include \
	-I ${LOCAL_PATH} \
	--gotag_out=:. example/example.proto

genTags: protobuf gotags

install:
	go install .

gen-tag:
	buf generate
	buf generate --template=buf.gen.tag.yaml
	buf generate --template=buf.gen.debug.yaml --path tagger

test:
	go test ./...
