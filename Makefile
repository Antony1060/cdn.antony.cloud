.DEFAULT_GOAL := build

.PHONY: build
clean:
	rm -rf build/*

.PHONY: build
build: clean
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags="-extldflags=-static" -v -o build/cdn main.go