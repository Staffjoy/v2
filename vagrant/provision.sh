#!/bin/bash

set -e
set -u
set -x

export DEBIAN_FRONTEND=noninteractive
export VHOME=/home/vagrant
export GOPATH=$VHOME/golang
export PROJECT_ROOT=$GOPATH/src/v2.staffjoy.com

sudo apt-get update -y -q
sudo apt-get install -y -q build-essential git curl ca-certificates bash-completion autoconf unison mysql-client

sudo mkdir -p $PROJECT_ROOT
sudo chown -R vagrant $GOPATH
sudo chgrp -R vagrant $GOPATH

source /vagrant/vagrant/golang.sh
source /vagrant/vagrant/bazel.sh
source /vagrant/vagrant/docker.sh
source /vagrant/vagrant/k8s.sh
source /vagrant/vagrant/npm.sh
source /vagrant/vagrant/docker.sh
source /vagrant/vagrant/nginx.sh
source /vagrant/vagrant/grpc.sh
source /vagrant/vagrant/mysql.sh

sudo apt-get autoremove -y -q
echo "export STAFFJOY=/home/vagrant/golang/src/v2.staffjoy.com/" >> "$VHOME/.profile"
echo "alias k=\"kubectl --namespace=development\"" >> "$VHOME/.profile"
echo "export ACCOUNT_MYSQL_CONFIG=\"mysql://root:SHIBBOLETH@tcp(10.0.0.100:3306)/account\"" >> "$VHOME/.profile"
echo "export COMPANY_MYSQL_CONFIG=\"mysql://root:SHIBBOLETH@tcp(10.0.0.100:3306)/company\"" >> "$VHOME/.profile"
echo "192.168.69.69 suite.local" >> "/etc/hosts"
