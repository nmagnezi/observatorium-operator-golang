#!/bin/bash

set -e

PROJECT_ROOT="$(readlink -e $(dirname "$BASH_SOURCE[0]")/../)"

CLUSTER="${CLUSTER:-OPENSHIFT}"
MARKETPLACE_NAMESPACE="${MARKETPLACE_NAMESPACE:-openshift-marketplace}"
PACKAGE="${PACKAGE:-observatorium-operator}"
APP_REGISTRY_NAMESPACE="${APP_REGISTRY_NAMESPACE:-observatorium-operator}"
TARGET_NAMESPACE="${TARGET_NAMESPACE:-observatorium}"

# Latest version from: https://quay.io/application/observatorium-operator/observatorium-operator
PACKAGE_VERSION="${PACKAGE_VERSION:-0.0.2}"

if [ "${CLUSTER}" == "KUBERNETES" ]; then
    MARKETPLACE_NAMESPACE="marketplace"
fi
if [ -z "${QUAY_USERNAME}" ]; then
    echo "QUAY_USERNAME"
    read QUAY_USERNAME
fi

if [ -z "${QUAY_PASSWORD}" ]; then
    echo "QUAY_PASSWORD"
    read -s QUAY_PASSWORD
fi

TOKEN=$("${PROJECT_ROOT}"/tools/token.sh $QUAY_USERNAME $QUAY_PASSWORD)

if [ `oc get secrets quay-registry-$APP_REGISTRY_NAMESPACE -n "${MARKETPLACE_NAMESPACE}" 2> /dev/null | wc -l` -eq 0 ]; then
cat <<EOF | oc create -f -
apiVersion: v1
kind: Secret
metadata:
  name: quay-registry-$APP_REGISTRY_NAMESPACE
  namespace: "${MARKETPLACE_NAMESPACE}"
type: Opaque
stringData:
      token: "$TOKEN"
EOF
fi

if [ `oc get operatorgroup -n "${TARGET_NAMESPACE}" 2> /dev/null | wc -l` -eq 0 ]; then
    echo "Creating OperatorGroup"
        cat <<EOF | oc create -f -
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: "${TARGET_NAMESPACE}-group"
  namespace: "${TARGET_NAMESPACE}"
spec: {}
EOF
fi

cat <<EOF | oc create -f -
apiVersion: operators.coreos.com/v1
kind: OperatorSource
metadata:
  name: "${APP_REGISTRY_NAMESPACE}"
  namespace: "${MARKETPLACE_NAMESPACE}"
spec:
  type: appregistry
  endpoint: https://quay.io/cnr
  registryNamespace: "${APP_REGISTRY_NAMESPACE}"
  displayName: "${APP_REGISTRY_NAMESPACE}"
  publisher: "Red Hat"
  authorizationToken:
    secretName: "quay-registry-${APP_REGISTRY_NAMESPACE}"
EOF
