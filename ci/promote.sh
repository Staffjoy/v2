#!/bin/bash
set -e
set -u
set -x

# TODO: we'll want to move replica count and port number into app-specific configs eventually
export VERSION="master-${BUILD_NUMBER}"
export NAMESPACE="production"
declare -a services=("www" "faraday" "accountapi" "accountserver" "emailserver" "myaccount" "whoami" "companyapi" "companyserver" "app" "smsserver" "botserver" "ical")


# Tell sentry we are deploying a new version
curl https://app.getsentry.com/api/hooks/release/SECRET/ \
  --fail \
  -X POST \
  -H 'Content-Type: application/json' \
  -d "{\"version\": \"$VERSION\"}"


## now loop through the above array
for service in "${services[@]}"
do
  export service
  # Deploy service to Kubernetes
  ./ci/deploy-service.sh
done

