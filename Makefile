lint:
	golangci-lint run
.PHONY: lint

test:
	go test -race -p 8 -parallel 8 ./...
.PHONY: test

build:
	go build -o ${GOPATH}/bin/effe ./cmd/effe/...
.PHONY: build

generate:
	go generate ./...
.PHONY: generate

test-cover:
	go test -race -p 8 -parallel 8 -coverpkg ./... -coverprofile coverage.out ./...
.PHONY: test-cover

# TODO: update tools in Dockerfile.dev as well
update-deps:
	go get -u ./...
	go mod tidy
.PHONY: update-deps

check-tidy:
	cp go.mod go.check.mod
	cp go.sum go.check.sum
	go mod tidy -modfile=go.check.mod
	diff -u go.mod go.check.mod
	diff -u go.sum go.check.sum
	rm go.check.mod go.check.sum
.PHONY: check-tidy

