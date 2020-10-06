# Beacon


### install packr utility
$ go get -u github.com/gobuffalo/packr/packr

### run packr
packr

### compile  go project to linux os/arch
GOOS=linux GOARCH=amd64 go build -v cmd/beacon/beacon.go