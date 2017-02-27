#!/bin/bash
set -e

if ! command -V vagrant >/dev/null 2>&1; then
    echo "OOPS! Run this command OUTSIDE of vagrant."
    exit 1
fi

vagrant halt

vagrant up

vagrant ssh -c 'sudo docker rm -f $(docker ps -aq) || true; sudo stop docker; sleep 1; sudo start docker; sudo rm -rf /var/lib/kubelet/; cd /vagrant/; sudo bash vagrant/k8s.sh; sudo bash vagrant/mysql.sh; sudo chown -R vagrant /home/vagrant/'

echo "DONE! Now run make dev to trigger a build."
