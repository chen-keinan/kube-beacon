# Beacon Project
#### The beacon project is audit scanner developed in GO that run set of audit check describe by [CIS Benchmark v1.6.0](https://www.cisecurity.org/benchmark/kubernetes/) on a deployed kubernetes cluster and output a security report 

#### Audit checks are performed  on master and worker nodes and the output audit report include :
* root cause of the security issue
* proposed remediation for security issue

[![Go Report Card](https://goreportcard.com/badge/github.com/chen-keinan/beacon)](https://goreportcard.com/report/github.com/chen-keinan/beacon)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/chen-keinan/beacon/blob/main/LICENSE)
[![Build Status](https://travis-ci.org/chen-keinan/beacon.svg?branch=main)](https://travis-ci.org/chen-keinan/beacon)


#### The Following diagram describe the kubernetes services which beacon service check   
![k8s arch](./pkg/images/k8s_arch.png?raw=true)

#### kubernetes cluster audit scan output: 
![k8s audit](./pkg/images/k8s_audit.png?raw=true)

