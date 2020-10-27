# Beacon Project
###  Scan your kubernetes runtime !!
[![Go Report Card](https://goreportcard.com/badge/github.com/chen-keinan/beacon)](https://goreportcard.com/report/github.com/chen-keinan/beacon)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/chen-keinan/beacon/blob/main/LICENSE)
[![Build Status](https://travis-ci.org/chen-keinan/beacon.svg?branch=main)](https://travis-ci.org/chen-keinan/beacon)
![Go Coverage](./pkg/images/coverage_badge.png?raw=true)

Beacon is a GO base audit scanner who perform audit check on a deployed kubernetes cluster and output a security report.
The audit tests are the full implementation of [CIS Kubernetes Benchmark specification](https://www.cisecurity.org/benchmark/kubernetes/) <br>

#### Audit checks are performed  on master and worker nodes and the output audit report include :
* root cause of the security issue
* proposed remediation for security issue


#### The Following diagram describe the kubernetes services which take part in beacon audit   
![k8s arch](./pkg/images/k8s_arch.png?raw=true)

#### kubernetes cluster audit scan output: 
![k8s audit](./pkg/images/beacon.gif) 

