BINARY ?= converter
GOCACHE ?= $(PWD)/.gocache
IMAGE ?= ebook-pdf

.PHONY: build test docker clean

build:
	GOCACHE=$(GOCACHE) go build -o $(BINARY) .

test:
	GOCACHE=$(GOCACHE) go test ./...

docker:
	docker build -t $(IMAGE) .

clean:
	rm -rf $(BINARY) $(GOCACHE)
