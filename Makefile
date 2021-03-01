.PHONY: run
run: build
	./oghu
	goserve public

.PHONY: build
build:
	go build -o oghu cmd/oghu/*.go
