#!/bin/bash

ARCH=amd64
export K8S_VERSION="v1.4.0" # should match google cloud deployed version


if [ ! -f /usr/local/bin/kubectl ]; then
    # export K8S_VERSION=$(curl -sS https://storage.googleapis.com/kubernetes-release/release/stable.txt)
    curl -O https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/linux/amd64/kubectl
    chmod +x kubectl
    sudo mv kubectl /usr/local/bin/kubectl
fi

# make the kubelet shared dir
# see https://github.com/kubernetes/kubernetes/issues/18239
if [ ! -d /var/lib/kubelet ]; then
    mkdir -p /var/lib/kubelet
    sudo mount -o bind /var/lib/kubelet /var/lib/kubelet
    sudo mount --make-shared /var/lib/kubelet
fi

if ! pgrep -c hyperkube >/dev/null 2>&1 ; then
    docker run -d \
        --volume=/:/rootfs:ro \
        --volume=/sys:/sys:ro \
        --volume=/dev:/dev \
        --volume=/var/lib/docker/:/var/lib/docker:rw \
        --volume=/var/lib/kubelet/:/var/lib/kubelet:shared \
        --volume=/var/run:/var/run:rw \
        --net=host \
        --pid=host \
        --privileged=true \
        --restart=always \
        gcr.io/google_containers/hyperkube-${ARCH}:${K8S_VERSION} \
        /hyperkube kubelet \
            --containerized \
            --address="0.0.0.0" \
            --hostname-override=127.0.0.1 \
            --api-servers=http://localhost:8080 \
            --config=/etc/kubernetes/manifests \
            --cluster-dns=10.0.0.10 \
            --cluster-domain=cluster.local \
            --allow-privileged=true \
            --v=2
fi

# We need to run a local registry - k8s cannot just pull locally
if ! pgrep -c registry >/dev/null 2>&1 ; then
    docker run -d \
        -p 5000:5000 \
        --restart=always \
        --name registry \
        registry:2
fi

# above may fail, wipe and re-run
# `docker rm -f $(docker ps -aq)`


# setup cluster config
kubectl config set-cluster staffjoy-dev --server=http://localhost:8080
kubectl config set-context staffjoy-dev --cluster=staffjoy-dev
kubectl config use-context staffjoy-dev

echo "Waiting for api to finish booting"
until curl 127.0.0.1:8080 &>/dev/null;
do
    echo ...
    sleep 1
done

kubectl create namespace development

# kick off account-mysql
kubectl --namespace=development create -R -f /home/vagrant/golang/src/v2.staffjoy.com/ci/k8s/development/infrastructure/app-mysql
