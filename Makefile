SHELL := /bin/bash

GOCMD=go
MOVESANDBOX=mv ~/vms/kube/kube-beacon ~/vms-local/kube-beacon
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
	$(GOCMD) get -d github.com/golang/mock/mockgen@v1.6.0
	$(GOCMD) install -v github.com/golang/mock/mockgen
	export PATH=$HOME/go/bin:$PATH
	$(GOMOCKS)
	./scripts/lint.sh
tidy:
	$(GOMOD) tidy -v
test:
	$(GOCMD) get -d github.com/golang/mock/mockgen@v1.6.0
	$(GOCMD) install -v github.com/golang/mock/mockgen && export PATH=$GOPATH/bin:$PATH;
	$(GOMOCKS)
	$(GOTEST) ./... -coverprofile coverage.md fmt
	$(GOCMD) tool cover -html=coverage.md -o coverage.html
	$(GOCMD) tool cover  -func coverage.md
build:
	export PATH=$GOPATH/bin:$PATH;
	export PATH=$PATH:/home/vagrant/go/bin
	export PATH=$PATH:/home/root/go/bin
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v ./cmd/kube;
build_local:
	export PATH=$GOPATH/bin:$PATH;
	export PATH=$PATH:/home/vagrant/go/bin
	export PATH=$PATH:/home/root/go/bin
	$(GOBUILD) ./cmd/kube;

build_docker_local:
	docker build -t chenkeinan/kube-beacon:latest .
	docker push chenkeinan/kube-beacon:latest
dlv:
	dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./kube-beacon
build_beb:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -v -gcflags='-N -l' cmd/kube/kube-beacon.go
	scripts/deb.sh
.PHONY: all build install test
