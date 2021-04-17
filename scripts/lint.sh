#!/usr/bin/env bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.27.0
## need to generate mock before lint
go get github.com/golang/mock/mockgen@latest
go install -v github.com/golang/mock/mockgen
export PATH=$GOPATH/bin:$PATH
export PATH=$PATH:/root/go/bin
go generate ./...
golangci-lint run -v  > lint.xml