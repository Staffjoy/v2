#!/bin/bash

if ! command -V add-apt-repository >/dev/null 2>&1; then
    sudo apt install -y -q  software-properties-common ca-certificates apt-transport-https
fi

if ! command -V java >/dev/null 2>&1; then
    sudo apt install -y -q  debconf-utils
    
    ## needs older java, @TODO: use v12
    #sudo add-apt-repository ppa:linuxuprising/java
    #sudo apt update -y -q
    #echo "oracle-java12-installer shared/accepted-oracle-license-v1-2 select true" | sudo /usr/bin/debconf-set-selections
    #sudo apt install -y -q  oracle-java12-installer

    sudo apt install -y -q  openjdk-11-jdk
fi

if ! command -V python3 >/dev/null 2>&1; then
    sudo apt install -y -q  python3
fi 

if ! command -V python >/dev/null 2>&1; then
    sudo update-alternatives --install /usr/bin/python python /usr/bin/python3 2
fi

if ! command -V bazel >/dev/null 2>&1; then
    sudo apt install -y -q  pkg-config zip g++ zlib1g-dev unzip python3

    if [ ! -f /etc/apt/sources.list.d/bazel.list ]; then
        echo "deb [arch=amd64] http://storage.googleapis.com/bazel-apt stable jdk1.8" | sudo tee /etc/apt/sources.list.d/bazel.list
        curl https://bazel.build/bazel-release.pub.gpg | sudo apt-key add -
        sudo apt update -y -q
    fi
    sudo apt install -y -q  bazel
fi
