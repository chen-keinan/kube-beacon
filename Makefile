SHELL := /bin/bash

GOCMD=go
MOVESANDBOX=mv beacon ~/vagrant_file/.
GOPACKR=$(GOCMD) get -u github.com/gobuffalo/packr/packr && packr
GOMOD=$(GOCMD) mod
GOMOCKS=$(GOCMD) generate ./...
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=beacon
GOCOPY=cp beacon ~/vagrant_file/.

all:test lint build

fmt:
	$(GOCMD) fmt ./...
lint:
	$(GOMOCKS)
	./lint.sh
tidy:
	$(GOMOD) tidy -v
test:
	$(GOCMD) get github.com/golang/mock/mockgen@latest
	$(GOCMD) install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	mv coverage_badge.png ./pkg/images/coverage_badge.png
	$(GOMOCKS)
	$(GOTEST) ./... -coverprofile coverage.md fmt
	$(GOCMD) tool cover -html=coverage.md -o coverage.html
build:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v cmd/kube/beacon.go;
	$(MOVESANDBOX)
install:build
	cp $(BINARY_NAME) $(GOPATH)/bin/kube-beacon
test_travis:
	$(GOCMD) get github.com/golang/mock/mockgen@latest
	$(GOCMD) install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	$(GOMOCKS)
	$(GOTEST) ./... -coverprofile coverage.md fmt
	$(GOCMD) tool cover -html=coverage.md -o coverage.html
build_travis:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v cmd/kube/beacon.go;
build_remote:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v -gcflags='-N -l' cmd/kube/beacon.go
	$(MOVESANDBOX)
setup:
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: all build install test
