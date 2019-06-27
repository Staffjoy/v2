#!/bin/bash

# docker deps
sudo apt-get install -y -q btrfs-tools libsystemd-dev apparmor debhelper dh-apparmor dh-systemd libapparmor-dev libdevmapper-dev libltdl-dev libsqlite3-dev pkg-config
#sudo apt-get install -y -q "linux-image-extra-$(uname -r)"

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo apt-key fingerprint 0EBFCD88

sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

# docker key
if [ ! -f /etc/apt/sources.list.d/docker.list ]; then
    sudo apt-key adv \
              --keyserver hkp://ha.pool.sks-keyservers.net:80 \
              --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
    
    echo "deb https://apt.dockerproject.org/repo ubuntu-xenial main" | sudo tee /etc/apt/sources.list.d/docker.list
    
    sudo apt-get update -y -q && apt-cache policy docker-engine
    sudo apt-get install -y -q docker-engine
fi

# docker-machine
if [ ! -f /usr/local/bin/docker-machine ]; then
    curl -L "https://github.com/docker/machine/releases/download/v0.16.1/docker-machine-$(uname -s)-$(uname -m)" > docker-machine
    chmod +x docker-machine
    sudo mv docker-machine /usr/local/bin/docker-machine
fi

# add vagrant to docker for dockering
# https://stackoverflow.com/questions/48568172/docker-sock-permission-denied
sudo usermod -aG docker $(whoami)

## not perfect, but makes it work - otherwise throws permission error on docker.sock
sudo chmod 777 /var/run/docker.sock

sudo systemctl status docker
