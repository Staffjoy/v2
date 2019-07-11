#!/bin/bash

set -e
set -u
set -x

export DEBIAN_FRONTEND=noninteractive
export VHOME=/home/vagrant
export GOPATH=$VHOME/golang
export STAFFJOY=$GOPATH/src/v2.staffjoy.com

sudo apt update -y -q
sudo apt install -y -q  build-essential git curl mc bash-completion autoconf unison mysql-client
sudo apt install -y -q  apt-transport-https ca-certificates gnupg-agent software-properties-common debconf-utils

sudo mkdir -p $STAFFJOY
sudo chown -R vagrant $GOPATH
sudo chgrp -R vagrant $GOPATH

source $STAFFJOY/vagrant/golang.sh
source $STAFFJOY/vagrant/bazel.sh
source $STAFFJOY/vagrant/docker.sh
source $STAFFJOY/vagrant/k8s.sh
source $STAFFJOY/vagrant/npm.sh
source $STAFFJOY/vagrant/docker.sh
source $STAFFJOY/vagrant/nginx.sh
source $STAFFJOY/vagrant/grpc.sh
source $STAFFJOY/vagrant/mysql.sh

sudo apt autoremove -y -q && sudo apt clean

echo "export STAFFJOY=${STAFFJOY}" | tee -a $VHOME/.profile
echo "export ACCOUNT_MYSQL_CONFIG=\"mysql://root:SHIBBOLETH@tcp(10.0.0.100:3306)/account\"" | tee -a $VHOME/.profile
echo "export COMPANY_MYSQL_CONFIG=\"mysql://root:SHIBBOLETH@tcp(10.0.0.100:3306)/company\"" | tee -a $VHOME/.profile

echo "alias k=\"kubectl --namespace=development\"" | tee -a $VHOME/.bash_aliases
echo "alias bazel=\"${VHOME}/.bazel/bin/bazel\"" | tee -a $VHOME/.bash_aliases

echo "192.168.69.69 suite.local" | sudo tee -a /etc/hosts
