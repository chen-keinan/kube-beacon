#!/usr/bin/env bash
# Install docker client
sudo apt-get update
sudo apt-get install docker-ce
# build docker image
docker build ./ -t kube-beacon -f Dockerfile
# login to registry
docker login beacon.jfrog.io -u $USER -p $PASSWORD
# tag image image
docker tag kube-beacon:latest beacon.jfrog.io/docker-local/kube-beacon:latest
# push image to repository
docker push beacon.jfrog.io/docker-local/kube-beacon:latest