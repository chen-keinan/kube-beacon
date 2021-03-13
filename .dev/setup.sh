echo "start vagrant provioning..."

sudo apt-get update

echo "install docker..."
sudo apt-get install -y docker.io

echo "set docker launch at boot..."
sudo systemctl enable docker

echo "start docker..."
sudo systemctl start docker

echo "add signing key..."
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add

echo "install curl pkg..."
sudo apt-get install -y curl

echo "Add ubuntu to default repository..."
sudo apt-add-repository "deb http://apt.kubernetes.io/ kubernetes-xenial main"

echo "install k8s tools..."
sudo apt-get install -y kubeadm kubelet kubectl
sudo apt-mark hold kubeadm kubelet kubectl

echo "set master node..."
sudo swapoff -a
sudo hostnamectl set-hostname master-node

echo "init master node..."
sudo kubeadm init --pod-network-cidr=10.244.0.0/16
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

echo "deploy pod network to cluster..."
sudo kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml

kubectl get pods --all-namespaces

kubectl get nodes

echo "install http transport..."
sudo apt-get install apt-transport-https

echo "install virtual box..."
echo virtualbox-ext-pack virtualbox-ext-pack/license select true | sudo debconf-set-selections
sudo apt install -y virtualbox virtualbox-ext-pack

echo "install minikube..." 
wget https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo cp minikube-linux-amd64 /usr/local/bin/minikube
sudo chmod 755 /usr/local/bin/minikube
minikube version

echo "install kubectl..."
curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl

chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl
kubectl version -o json

echo "start minikube..."
usermod -aG docker $USER && newgrp docker
minikube start --memory=1992mb --driver=docker --force

echo "check minikube status..."
kubectl config view
kubectl cluster-info
kubectl get nodes
kubectl get pod
minikube status

echo "install golang pkg"
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update -y
sudo apt install -y golang-go 

echo "Install dlv pkg"
 git clone https://github.com/go-delve/delve.git $GOPATH/src/github.com/go-delve/delve
 cd $GOPATH/src/github.com/go-delve/delve
 make install

### export dlv bin path
export PATH=$PATH:/home/vagrant/go/bin

echo "Finished provisioning."
