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
dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./kube-beacon

docker run --pid=host -v /etc:/etc:ro -v /var:/var:ro -v /*/cni/*:/*/cni/* -t  beacon

docker build ./ -t beacon -f Dockerfile

    export KUBECONFIG=/etc/kubernetes/admin.conf
mkdir -p $HOME/.kube
 sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
 sudo chown $(id -u):$(id -g) $HOME/.kube/config
https://github.com/oracle/vagrant-projects

    kubectl taint nodes master-node node-role.kubernetes.io/master-
    kubectl create clusterrolebinding default-admin --clusterrole cluster-admin --serviceaccount=default:default
