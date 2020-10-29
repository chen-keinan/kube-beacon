##Remote Debug
###Install dlv
$ git clone https://github.com/go-delve/delve.git $GOPATH/src/github.com/go-delve/delve
$ cd $GOPATH/src/github.com/go-delve/delve
    $ make install

### export dlv bin path
export PATH=$PATH:/home/vagrant/go/bin

### compile binary with debug params
GOOS=linux GOARCH=amd64 go build -v -gcflags='-N -l' cmd/beacon/beacon.go

### run on remote machine
dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./beacon
