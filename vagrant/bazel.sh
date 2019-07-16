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
    sudo apt install -y -q  pkg-config zip g++ zlib1g-dev unzip

    # This release version should correspond to the version listed here:
    # https://github.com/bazelbuild/bazel/releases
    RELEASE=0.28.0

    if [[ "$OSTYPE" == "linux-gnu" ]]; then
        sudo curl -L https://github.com/bazelbuild/bazel/releases/download/${RELEASE}/bazel-${RELEASE}-installer-linux-x86_64.sh --output /usr/src/bazel-${RELEASE}-installer-linux-x86_64.sh
        sudo chmod +x /usr/src/bazel-${RELEASE}-installer-linux-x86_64.sh
        /usr/src/bazel-${RELEASE}-installer-linux-x86_64.sh --user
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        sudo curl -L https://github.com/bazelbuild/bazel/releases/download/${RELEASE}/bazel-${RELEASE}-installer-darwin-x86_64.sh --output /usr/src/bazel-${RELEASE}-installer-darwin-x86_64.sh
        sudo chmod +x /usr/src/bazel-${RELEASE}-installer-darwin-x86_64.sh
        /usr/src/bazel-${RELEASE}-installer-darwin-x86_64.sh --user
    fi

    echo "source /home/${USER}/.bazel/bin/bazel-complete.bash" | sudo tee -a ~/.bashrc
fi