#!/bin/bash

# We need to run a local registry - k8s cannot just pull locally
if ! pgrep -c registry >/dev/null 2>&1 ; then
    docker run -d \
        -p 5000:5000 \
        --restart=always \
        --name registry \
        registry:2
fi

# download and install kubectl ...
curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v1.15.0/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/

# ... and minikube
curl -Lo minikube https://storage.googleapis.com/minikube/releases/v1.2.0/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/

sudo -E minikube start \
    --kubernetes-version=v1.15.0 \
    --vm-driver=none \
    --dns-domain="cluster.local" \
    --service-cluster-ip-range="10.0.0.0/12" \
    --extra-config="kubelet.cluster-dns=10.0.0.10"

# set the kubectl context to minikube (this overwrites ~/.kube and ~/.minikube, but leaves files' ownership as root:root)
sudo -E minikube update-context

# enables dashboard
sudo minikube addons enable dashboard

# either use sudo on all kubectl commands, or chown/chgrp to your user
sudo chown -R ${USER}:${USER} /home/${USER}/.kube /home/${USER}/.minikube

sudo find /etc/kubernetes \
    \( -type f -exec sudo chmod +r {} \; \) , \
    \( -type d -exec sudo chmod +rx {} \; \)

# this will write over any previous configuration)
# wait for the cluster to become ready/accessible via kubectl
echo -e -n " [ ] Waiting for master components to start...";
JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}';
until sudo kubectl get nodes -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True"; do
    echo -n "."
    sleep 1
done

kubectl cluster-info

kubectl config set-cluster staffjoy-dev --server=https://10.0.2.15:8443 --certificate-authority=/home/${USER}/.minikube/ca.crt
kubectl config set-context staffjoy-dev --cluster=staffjoy-dev --user=minikube
kubectl config use-context staffjoy-dev

kubectl create namespace development

kubectl --namespace=development create -R -f ~/golang/src/v2.staffjoy.com/ci/k8s/development/infrastructure/app-mysql

kubectl --context minikube proxy &>/dev/null &