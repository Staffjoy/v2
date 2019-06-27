#!/bin/bash

if ! command -V add-apt-repository >/dev/null 2>&1; then
    sudo apt-get install -y -q  software-properties-common ca-certificates apt-transport-https
fi

if ! command -V java >/dev/null 2>&1; then
    sudo apt-get install -y -q  python-software-properties debconf-utils
    
    ## needs older java, @TODO: use v12
    #sudo add-apt-repository ppa:linuxuprising/java
    #sudo apt-get update -y -q
    #echo "oracle-java12-installer shared/accepted-oracle-license-v1-2 select true" | sudo /usr/bin/debconf-set-selections
    #sudo apt-get install -y -q  oracle-java12-installer

    ## fallback to v1.8
    sudo apt install -y -q  openjdk-8-jdk
fi

# bazel deps
sudo apt-get install -y -q  pkg-config zip g++ zlib1g-dev unzip

if [ ! -f /etc/apt/sources.list.d/bazel.list ]; then
    echo "deb http://storage.googleapis.com/bazel-apt testing jdk1.8" | sudo tee /etc/apt/sources.list.d/bazel.list
    curl https://storage.googleapis.com/bazel-apt/doc/apt-key.pub.gpg | sudo apt-key add -
fi

sudo apt-get update -y -q
# latest bazel, 0.27rc5 - too new at that point
#sudo apt-get install -y -q bazel

# use older bazel
sudo curl https://github.com/bazelbuild/bazel/releases/download/0.6.0/bazel_0.6.0-linux-x86_64.deb --output /usr/src/bazel_0.6.0-linux-x86_64.deb -L --max-redirs 5
sudo dpkg -i /usr/src/bazel_0.6.0-linux-x86_64.deb
