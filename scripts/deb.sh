#!/usr/bin/env bash

mkdir builds/kube-beacon
mkdir builds/kube-beacon/DEBIAN
{
  echo package: kube-beacon
  echo Version: 0.1
  echo Section: custom
  echo Priority: optional
  echo Architecture: all
  echo Essential: no
  echo Installed-Size: 1024
  echo Maintainer: hen.keinan@gmail.com
  echo Description: k8s audit scan tool
} >> builds/kube-beacon/DEBIAN/control
mkdir -p builds/kube-beacon/usr/bin/
mv kube-beacon  builds/kube-beacon/usr/bin/
dpkg-deb --build builds/kube-beacon
mv builds/kube-beacon.deb builds/deb
rm -rf builds/kube-beacon