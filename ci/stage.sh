#!/bin/bash
set -e
set -u
set -x

# This script is called by Jenkins
export VERSION="master-${BUILD_NUMBER}"
export NAMESPACE="staging"
declare -a targets=("www" "faraday" "account/api" "account/server" "email/server" "myaccount" "whoami" "company/api" "company/server" "sms/server" "bot/server" "app" "ical")

# Tell sentry we are deploying a new version
curl https://app.getsentry.com/api/hooks/release/SECRET/ \
  --fail \
  -X POST \
  -H 'Content-Type: application/json' \
  -d "{\"version\": \"$VERSION\"}"

## now loop through the above array
for target in "${targets[@]}"
do
    # Remove slashes to get service name
    service=$(echo $target | sed 's/\///g')
    export service
    # Run the build and upload to GKE
    bazel run //$target:docker
    docker tag bazel/$(echo $target | sed 's/\//_/g'):docker gcr.io/staffjoy-prod/$service:$VERSION
    gcloud docker push gcr.io/staffjoy-prod/$service:$VERSION

    # Deploy service to Kubernetes
    ./ci/deploy-service.sh
done

