#!/bin/bash

set -e
set -u
set -x

export DEBIAN_FRONTEND=noninteractive
export VHOME=/home/${USER}
export GOPATH=${VHOME}/golang
export STAFFJOY=${GOPATH}/src/v2.staffjoy.com

## apt-fast
sudo add-apt-repository ppa:apt-fast/stable < /dev/null
echo debconf apt-fast/maxdownloads string 16 | sudo debconf-set-selections
echo debconf apt-fast/dlflag boolean true | sudo debconf-set-selections
echo debconf apt-fast/aptmanager string apt | sudo debconf-set-selections
sudo apt install -y -q  apt-fast

sudo apt update -y -q
sudo apt install -y -q  build-essential bash-completion autoconf git curl unison mc
sudo apt install -y -q  apt-transport-https ca-certificates gnupg-agent software-properties-common debconf-utils

sudo mkdir -p ${STAFFJOY}
sudo chown -R ${USER} ${GOPATH}
sudo chgrp -R ${USER} ${GOPATH}

source ${STAFFJOY}/vagrant/golang.sh
source ${STAFFJOY}/vagrant/bazel.sh
source ${STAFFJOY}/vagrant/npm.sh
source ${STAFFJOY}/vagrant/grpc.sh
source ${STAFFJOY}/vagrant/nginx.sh
source ${STAFFJOY}/vagrant/docker.sh
source ${STAFFJOY}/vagrant/minikube.sh
source ${STAFFJOY}/vagrant/mysql.sh

sudo apt autoremove -y -q && sudo apt clean

echo "export STAFFJOY=${STAFFJOY}" | tee -a ${VHOME}/.profile
echo "export ACCOUNT_MYSQL_CONFIG=\"mysql://root:SHIBBOLETH@tcp(10.0.0.100:3306)/account\"" | tee -a ${VHOME}/.profile
echo "export COMPANY_MYSQL_CONFIG=\"mysql://root:SHIBBOLETH@tcp(10.0.0.100:3306)/company\"" | tee -a ${VHOME}/.profile

echo "192.168.69.69 suite.local" | sudo tee -a /etc/hosts

echo "alias bazel=\"${VHOME}/.bazel/bin/bazel\"" | tee -a ${VHOME}/.bash_aliases
echo "alias k=\"kubectl --namespace=development\"" | tee -a ${VHOME}/.bash_aliases

#echo "alias minikube-kill = `docker rm $(docker kill $(docker ps -a --filter=\"name=k8s_\" --format=\"{{.ID}}\"))`" | tee -a ${VHOME}/.bash_aliases
#echo "alias minikube-stop = `docker stop $(docker ps -a --filter=\"name=k8s_\" --format=\"{{.ID}}\")`" | tee -a ${VHOME}/.bash_aliases
