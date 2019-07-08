#!/bin/bash
PATH=$PATH:$GOPATH/bin:/usr/local/go/bin

if [ ! -d /usr/local/go ]; then
    sudo curl -O https://storage.googleapis.com/golang/go1.12.6.linux-amd64.tar.gz
    sudo tar -xvf go1.12.6.linux-amd64.tar.gz
    sudo mv go /usr/local
    sudo rm go1.12.6.linux-amd64.tar.gz
    echo "export GOPATH=$GOPATH" >> "$VHOME/.profile"
    echo "export PATH=\$PATH:\$GOPATH/bin:/usr/local/go/bin" >> "$VHOME/.profile"
    echo "export GO111MODULE=on" >> "$VHOME/.profile"
fi

sudo -u vagrant -H bash -c "
id
source ~/.profile

if ! command -V golint ; then
    go get -u golang.org/x/lint/golint
    go get -u golang.org/x/tools/cmd/cover
    go get -u golang.org/x/tools/cmd/goimports
fi

if ! command -V protoc-gen-go ; then 
    go get -u github.com/golang/protobuf/protoc-gen-go
    go get -u golang.org/x/tools/cmd/cover
    go get -u golang.org/x/tools/cmd/goimports
    go get -u github.com/grpc-ecosystem/grpc-gateway/...
fi

if ! command -V glide ; then
    curl https://glide.sh/get | sh
fi

if ! command -V migrate ; then
    # https://github.com/golang-migrate/migrate
    curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | sudo apt-key add -
    echo 'deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ bionic main' | sudo tee /etc/apt/sources.list.d/migrate.list
    sudo apt update -y -q
    sudo apt install -y -q  migrate
fi

if ! command -V buildifier ; then
    go get github.com/bazelbuild/buildtools/buildifier
fi

if ! command -V go-bindata ; then
    go get -u github.com/jteeuwen/go-bindata/...
fi

go get -u github.com/gogo/protobuf/...

# used for local filesystem watching
if ! command -V modd ; then
    go get github.com/cortesi/modd/cmd/modd
fi
"
