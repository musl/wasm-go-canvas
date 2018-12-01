BIN := main.wasm
VERSION := $(shell git describe --abbrev=0 --tags || echo '0.0.0')
REVISION := $(shell git rev-list -1 HEAD || echo 'NONE')

.PHONY: all build clean run

all: serve

clean:
	go clean .
	rm -f *.png *.out *.wasm
	rm -fr build

run:
	GOOS=js GOARCH=wasm go run -exec="$(shell go env GOROOT)/misc/wasm/go_js_wasm_exec" .

$(BIN):
	GOOS=js GOARCH=wasm go build -o $@ -ldflags="-X main.Version=$(VERSION) -X main.Revision=$(REVISION)" .

build: clean $(BIN)
	mkdir -p build
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" build/
	cp static/* build/
	cp $(BIN) build/

serve: build
	go run tools/serve.go

