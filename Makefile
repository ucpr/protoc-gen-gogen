.PHONY: example
example: LOCAL_PATH = $(shell pwd)
example:
	protoc -I ${LOCAL_PATH} --go_out=. --go_opt=paths=source_relative _example/example.proto
	protoc -I ${LOCAL_PATH} --gogen_out=. --gogen_opt=paths=source_relative _example/example.proto

.PHONY: install
install:
	go install .

.PHONY: test
test:
	go test -race ./...
