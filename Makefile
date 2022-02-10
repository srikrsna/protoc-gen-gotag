LOCAL_PATH = $(shell pwd)

.PHONY: example proto install gen-tag test

example: proto install
	protoc -I /usr/local/include \
	-I ${LOCAL_PATH} \
	--gotag_out=xxx="graphql+\"-\" bson+\"-\"":. example/example.proto

proto:
	protoc -I /usr/local/include \
	-I ${LOCAL_PATH} \
	--go_out=:. example/example.proto

install:
	go install .

gen-tag:
	buf generate
	buf generate --template=buf.gen.tag.yaml

test:
	go test ./...
