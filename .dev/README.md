# kubernetes-vagrantfile

vagrant file to be used for k8s associated  programs developments, file include :
- buntu/bionic64
- k8s cluster 
- minikube
- dlv for remote debug

## Quick Start

```
 git clone git@github.com:chen-keinan/kubernetes-vagrantfile.git
 cd kubernetes-vagrantfile
 vagrant up

```


### Compile binary with debug params
```
GOOS=linux GOARCH=amd64 go build -v -gcflags='-N -l' demo.go
```
### Run debug on remote machine
```
dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./demo
```

### Tear down
```
 vagrant destroy
