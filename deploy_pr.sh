#!/bin/bash
set -euo pipefail

BONFIRE_NAMESPACE=""

on_exit() {
  echo
  echo "Releasing namespace ${BONFIRE_NAMESPACE:-<not reserved>}"
  if [[ -n "$BONFIRE_NAMESPACE" ]]; then
    bonfire namespace release "$BONFIRE_NAMESPACE"
  fi
}

echo "Prerequisites:"
echo "  1. Log in to OpenShift (oc login)"
echo "  2. Activate a Python virtual environment with bonfire installed"
echo "  3. Switch to the branch you want to deploy"
echo "  4. Push your commit and wait for the Konflux pipeline to succeed"
echo "  For setup details: https://inscope.corp.redhat.com/docs/default/component/ephemeral-environments-docs/getting-started/01-getting-started-with-ees/"
echo

QUICK=false
if [[ "${1:-}" == "--quick" ]]; then
  QUICK=true
  shift
fi

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  echo "Usage: $(basename "$0") [--quick] [COMMIT_ID]"
  echo "Deploy a PR build to an ephemeral namespace. Defaults to HEAD commit."
  echo
  echo "Options:"
  echo "  --quick   Skip the image existence check and existing reservation check"
  exit 0
fi

if ! command -v oc &>/dev/null; then
  echo "ERROR: oc (OpenShift CLI) not found. Please install and retry."
  exit 1
fi

if ! command -v skopeo &>/dev/null; then
  case "$(uname -s)" in
    Darwin*) install_hint="brew install skopeo" ;;
    *)       install_hint="dnf install skopeo" ;;
  esac
  echo "ERROR: skopeo not found. Install with '$install_hint' and retry."
  exit 1
fi

if [[ -z "${VIRTUAL_ENV:-}" ]]; then
  echo "ERROR: Not running inside a virtual environment. Activate one with bonfire installed."
  exit 1
fi

if ! python -c 'import importlib.util; exit(0 if importlib.util.find_spec("bonfire") else 1)' 2>/dev/null; then
  echo "ERROR: Bonfire is not installed in the current venv. Please install and retry."
  exit 1
fi

if ! oc whoami &>/dev/null; then
  echo "ERROR: Not logged in to OpenShift. Run 'oc login' first."
  exit 1
fi

CURRENT_BRANCH="$(git branch --show-current)"
COMMIT_ID="${1:-$(git rev-parse HEAD)}"

if ! git branch -r --contains "$COMMIT_ID" 2>/dev/null | grep -q "origin/"; then
  echo "ERROR: Commit $COMMIT_ID has not been pushed to origin."
  echo "  Push it first:  git push origin $CURRENT_BRANCH"
  exit 1
fi

trap on_exit EXIT SIGINT SIGTERM

IMAGE="quay.io/redhat-user-workloads/teamnado-konflux-tenant/composer-api"
IMAGE_TAG="on-pr-$COMMIT_ID"

echo "Deploying branch '$CURRENT_BRANCH' (commit: ${COMMIT_ID:0:12})"
echo

if [[ "$QUICK" == false ]]; then
  echo "Checking if image $IMAGE:$IMAGE_TAG exists..."
  if ! skopeo inspect "docker://$IMAGE:$IMAGE_TAG" &>/dev/null; then
    echo "ERROR: Image $IMAGE:$IMAGE_TAG not found. The Konflux pipeline may not have finished."
    exit 1
  fi
  echo "Image found!"
  echo

  EXISTING_NS="$(bonfire namespace list --mine 2>/dev/null | grep -oE 'ephemeral-[^ ]+' | head -1 || true)"

  if [[ -n "$EXISTING_NS" ]]; then
    echo "Existing namespace found: $EXISTING_NS"
    read -p "Reuse it? (y/n) " -n 1 -r
    echo
    if [[ "$REPLY" =~ ^[Yy]$ ]]; then
      BONFIRE_NAMESPACE="$EXISTING_NS"
      echo "Reusing namespace: $BONFIRE_NAMESPACE"
    else
      echo "Reserving a new bonfire namespace..."
      BONFIRE_NAMESPACE="$(bonfire namespace reserve --force)"
      echo "Namespace reserved: $BONFIRE_NAMESPACE"
    fi
  else
    echo "Reserving bonfire namespace..."
    BONFIRE_NAMESPACE="$(bonfire namespace reserve)"
    echo "Namespace reserved: $BONFIRE_NAMESPACE"
  fi
else
  echo "Reserving bonfire namespace (skipping checks)..."
  BONFIRE_NAMESPACE="$(bonfire namespace reserve --force)"
  echo "Namespace reserved: $BONFIRE_NAMESPACE"
fi

oc project "$BONFIRE_NAMESPACE" >/dev/null
echo "Switched oc project to $BONFIRE_NAMESPACE"

echo "Deploying commit $COMMIT_ID..."
if ! oc process \
  -p ENV_NAME="env-$BONFIRE_NAMESPACE" \
  -p REPLICAS=1 \
  -p IMAGE="$IMAGE" \
  -p IMAGE_TAG="$IMAGE_TAG" \
  -f deploy/clowdapp.yaml | oc apply -f - >/dev/null; then
  echo "ERROR: Deployment failed"
  exit 1
fi

echo "Deployment successful. Pod details:"
echo "  https://console-openshift-console.apps.crc-eph.r9lp.p1.openshiftapps.com/k8s/ns/$BONFIRE_NAMESPACE/core~v1~Pod"
echo
bonfire namespace describe
echo

read -p "Press any key to release namespace and exit " -n 1 -r
