#!/bin/bash

# standard bash error handling
set -o nounset  # treat unset variables as an error and exit immediately.
set -o errexit  # exit immediately when a command fails.
set -E          # needs to be set if we want the ERR trap
set -o pipefail # prevents errors in a pipeline from being masked

VERSION=${1:-''}

# Create namespaces
kubectl create namespace kcp-system || true
kubectl create namespace kyma-system || true
kubectl create namespace istio-system || true

# Install Istio
helm repo add istio https://istio-release.storage.googleapis.com/charts
helm repo update
helm install istio-base istio/base -n istio-system --set defaultRevision=default

# Install Prometheus Operator for ServiceMonitor
kubectl create -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/master/bundle.yaml

# Install Postgres
kubectl create -f scripts/testing/yaml/postgres -n kcp-system

# Prepare fake gardener credentials
KCFG=$(kubectl config view --minify --raw)
kubectl create secret generic gardener-credentials --from-literal=kubeconfig="$KCFG" -n kcp-system

# Prepare chart for custom KEB version
if [[ -n "$VERSION" ]]; then
  if [[ "$VERSION" == PR* ]]; then
    scripts/bump_keb_chart.sh "$VERSION" "pr"
  else
    scripts/bump_keb_chart.sh "$VERSION" "release"
  fi
fi

# Deploy KEB helm chart
cd resources/keb
if [[ "$VERSION" == PR* ]]; then
  helm install keb ../keb \
    --namespace kcp-system \
    -f values.yaml \
    --set global.database.embedded.enabled=false \
    --set testConfig.kebDeployment.useAnnotations=true \
    --set global.images.container_registry.path="europe-docker.pkg.dev/kyma-project/dev" \
    --set global.secrets.mechanism=secrets \
    --debug --wait
else
  helm install keb ../keb \
    --namespace kcp-system \
    -f values.yaml \
    --set global.database.embedded.enabled=false \
    --set testConfig.kebDeployment.useAnnotations=true \
    --set global.secrets.mechanism=secrets \
    --debug --wait
fi

# Check if KEB pod is in READY state
echo "Waiting for kyma-environment-broker pod to be in READY state..."
kubectl wait --namespace kcp-system --for=condition=Ready pod -l app.kubernetes.io/name=kyma-environment-broker --timeout=60s
EXIT_CODE=$?
if [ $EXIT_CODE -ne 0 ]; then
  echo "The kyma-environment-broker pod did not become READY within the timeout."
  echo "Fetching the logs from the pod..."
  POD_NAME=$(kubectl get pod -l app.kubernetes.io/name=kyma-environment-broker -n kcp-system -o jsonpath='{.items[0].metadata.name}')
  kubectl logs $POD_NAME -n kcp-system
  exit 1
fi
