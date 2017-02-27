echo "Starting deployment for service $service version $VERSION"

# Create or update service
echo "Checking if k8s service for $service exists..."
kubectl get service $service-service --namespace=$NAMESPACE 2>&1 >/dev/null

if [ $? -ne 0 ]
then
  echo "K8s service for $service doesn't exist.  Creating service..."
  kubectl --namespace=$NAMESPACE create -f ./ci/k8s/$NAMESPACE/services/$service.yaml
else
  echo "K8s service for $service exists "
  # TODO - not really clear how to update a running service
  #kubectl --namespace=$NAMESPACE replace -f ./ci/k8s/$NAMESPACE/services/$service.yaml
fi

cp ./ci/k8s/$NAMESPACE/deployments/$service.yaml ./ci/k8s/$NAMESPACE/deployments/$service-copy.yaml
sed -i "s/VERSION/$VERSION/g" ./ci/k8s/$NAMESPACE/deployments/$service-copy.yaml

echo "Checking if deployment for $service exists..."
kubectl get deployment $service-deployment --namespace=$NAMESPACE 2>&1 >/dev/null
if [ $? -eq 0 ]
then
  echo "Deployment for $service exists, updating container image to version $VERSION"
  kubectl --namespace=$NAMESPACE update -f ./ci/k8s/$NAMESPACE/deployments/$service-copy.yaml
else
  echo "Deployment for $service doesn't exist, creating deployment with container image version $VERSION"
  kubectl --namespace=$NAMESPACE create -f ./ci/k8s/$NAMESPACE/deployments/$service-copy.yaml
fi
echo "Finished deploying $service, version $VERSION to $NAMESPACE."

rm ./ci/k8s/$NAMESPACE/deployments/$service-copy.yaml
