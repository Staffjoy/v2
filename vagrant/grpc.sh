#!/bin/bash

set -e

COMPILE_FROM_SRC=false;

if $COMPILE_FROM_SRC; then
    if [ -d tmp ]; then
        rm -rf tmp
    fi

    mkdir tmp
    cd tmp

    # Subset of protobuf to have a faster setup
    curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.8.0/protobuf-cpp-3.8.0.tar.gz
    tar -xvzf protobuf-cpp-3.8.0.tar.gz
    ln -s protobuf-3.8.0 protobuf

    cd protobuf
    ./autogen.sh
    ./configure

    make
    make check
    sudo make install

    sudo ldconfig # refresh shared library cache.

    cd ../..
    rm -rf tmp
else
    sudo apt install -y -q  protobuf-compiler
fi