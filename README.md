[![Go Report Card](https://goreportcard.com/badge/github.com/chen-keinan/beacon)](https://goreportcard.com/report/github.com/chen-keinan/beacon)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/chen-keinan/beacon/blob/main/LICENSE)
[![Build Status](https://travis-ci.com/chen-keinan/kube-beacon.svg?branch=main)](https://travis-ci.com/chen-keinan/kube-beacon)
[![Coverage Status](https://coveralls.io/repos/github/chen-keinan/kube-beacon/badge.svg?branch=main)](https://coveralls.io/github/chen-keinan/kube-beacon?branch=main)
<br><img src="./pkg/images/beacon-gopher.png" width="200" alt="kube-beacon logo"><br>
# Kube-Beacon Project
###  Scan your kubernetes runtime !!
Kube-Beacon is an open source audit scanner who perform audit check on a deployed kubernetes cluster and output a security report.

The audit tests are the full implementation of [CIS Kubernetes Benchmark specification](https://www.cisecurity.org/benchmark/kubernetes/) <br>

#### Audit checks are performed  on master and worker nodes and the output audit report include :
* root cause of the security issue
* proposed remediation for security issue

#### kubernetes cluster audit scan output: 
![k8s audit](./pkg/images/beacon.gif) 

* [Installation](#installation)
* [Quick Start](#quick-start)


## Installation

```sh
git clone https://github.com/chen-keinan/kube-beacon
cd kube-beacon
make install
```


## Quick Start

Execute kube-eacon without any flags , execute all tests

```
 ./kube-beacon 

```

Execute kube-beacon  with flags , execute test on demand

```
Usage: kube-Beacon [--version] [--help] <command> [<args>]

Available commands are:
  -i , --include, execute specific tests only,  example -i=1.2.3,1.4.5
  -e , --exclude, exclude tests from execution, example -e=1.2.3,1.4.5
  -n , --node,    execute specific node test,   example -n=master,-n=worker
```

Execute tests and generate failure tests report

```
./kube-beacon -r
```