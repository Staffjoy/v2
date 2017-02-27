#!/bin/bash

# docker deps
sudo apt-get install -y -q btrfs-tools libsystemd-journal-dev apparmor debhelper dh-apparmor dh-systemd libapparmor-dev libdevmapper-dev libltdl-dev libsqlite3-dev pkg-config "linux-image-extra-$(uname -r)"

# docker key
if [ ! -f /etc/apt/sources.list.d/docker.list ]; then
    sudo apt-key adv --keyserver hkp://pgp.mit.edu:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
    echo "deb https://apt.dockerproject.org/repo ubuntu-trusty main" | sudo tee /etc/apt/sources.list.d/docker.list
    sudo apt-get update -y -q
    sudo apt-get install -y -q docker-engine
fi

# docker-machine
if [ ! -f /usr/local/bin/docker-machine ]; then
    curl -L "https://github.com/docker/machine/releases/download/v0.7.0/docker-machine-$(uname -s)-$(uname -m)" > docker-machine
    chmod +x docker-machine
    sudo mv docker-machine /usr/local/bin/docker-machine
fi

# add vagrant to docker for dockering
sudo usermod -aG docker vagrant
