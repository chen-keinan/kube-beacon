SHELL := /bin/bash

GOCMD=go
MOVESANDBOX=mv ~/vms/kube-beacon/kube-beacon ~/vms-local/kube-beacon
GOPACKR=$(GOCMD) get -u github.com/gobuffalo/packr/packr && packr
GOMOD=$(GOCMD) mod
GOMOCKS=$(GOCMD) generate ./...
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=kube-beacon
GOCOPY=cp kube-beacon ~/vagrant_file/.

all:test lint build

fmt:
	$(GOCMD) fmt ./...
lint:
	./scripts/lint.sh
tidy:
	$(GOMOD) tidy -v
test:
	$(GOCMD) get github.com/golang/mock/mockgen@latest
	$(GOCMD) install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	$(GOMOCKS)
	$(GOTEST) ./... -coverprofile coverage.md fmt
	$(GOCMD) tool cover -html=coverage.md -o coverage.html
	$(GOCMD) tool cover  -func coverage.md
build:
	$(GOPACKR)
	export PATH=$GOPATH/bin:$PATH;
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v cmd/kube/kube-beacon.go;
install:build_travis
	cp $(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
test_travis:
	$(GOCMD) get github.com/golang/mock/mockgen@latest
	$(GOCMD) install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	$(GOMOCKS)
	$(GOTEST) ./... -coverprofile coverage.md fmt
	$(GOCMD) tool cover -html=coverage.md -o coverage.html
build_travis:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v cmd/kube/kube-beacon.go;
build_remote:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v -gcflags='-N -l' cmd/kube/kube-beacon.go
	$(MOVESANDBOX)
build_docker:
	docker build -t kbeacon.jfrog.io/docker-local/kube-beacon .
	docker push kbeacon.jfrog.io/docker-local/kube-beacon
build_docker_local:
	docker build -t beacon:8082/docker-local/kube-beacon .
	docker push beacon:8082/docker-local/kube-beacon:latest

build_beb:
	$(GOPACKR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v -gcflags='-N -l' cmd/kube/kube-beacon.go
	scripts/deb.sh
.PHONY: all build install test
