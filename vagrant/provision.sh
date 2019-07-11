#!/bin/bash

set -e
set -u
set -x

export DEBIAN_FRONTEND=noninteractive
export VHOME=/home/vagrant
export GOPATH=$VHOME/golang
export PROJECT_ROOT=$GOPATH/src/v2.staffjoy.com

sudo apt update -y -q
sudo apt install -y -q  build-essential git curl mc bash-completion autoconf unison mysql-client
sudo apt install -y -q  apt-transport-https ca-certificates gnupg-agent software-properties-common debconf-utils

sudo mkdir -p $PROJECT_ROOT
sudo chown -R vagrant $GOPATH
sudo chgrp -R vagrant $GOPATH

source $PROJECT_ROOT/golang.sh
source $PROJECT_ROOT/bazel.sh
source $PROJECT_ROOT/docker.sh
source $PROJECT_ROOT/k8s.sh
source $PROJECT_ROOT/npm.sh
source $PROJECT_ROOT/docker.sh
source $PROJECT_ROOT/nginx.sh
source $PROJECT_ROOT/grpc.sh
source $PROJECT_ROOT/mysql.sh

sudo apt autoremove -y -q && sudo apt clean

echo "export STAFFJOY=${PROJECT_ROOT}" | tee -a $VHOME/.profile
echo "export ACCOUNT_MYSQL_CONFIG=\"mysql://root:SHIBBOLETH@tcp(10.0.0.100:3306)/account\"" | tee -a $VHOME/.profile
echo "export COMPANY_MYSQL_CONFIG=\"mysql://root:SHIBBOLETH@tcp(10.0.0.100:3306)/company\"" | tee -a $VHOME/.profile

echo "alias k=\"kubectl --namespace=development\"" | tee -a $VHOME/.bash_aliases
echo "alias bazel=\"${VHOME}/.bazel/bin/bazel\"" | tee -a $VHOME/.bash_aliases

echo "192.168.69.69 suite.local" | sudo tee -a /etc/hosts
