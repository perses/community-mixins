#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
ENV_FILE="${ROOT_DIR}/.github/env"
# shellcheck source=test/e2e/lib/env.sh
source "${ROOT_DIR}/test/e2e/lib/env.sh"

E2E_KIND_CLUSTER="${E2E_KIND_CLUSTER:-community-mixins-e2e}"
E2E_NAMESPACE="perses-dev"
HELM_REPO="https://perses.github.io/helm-charts"
HELM_RELEASE="perses-operator"

usage() {
  cat <<EOF
Usage: $(basename "$0") <up|install|down>

  up       Create the kind cluster (if needed) and install the e2e stack
  install  Install operator, Perses, built dashboards (cluster must exist)
  down     Delete the kind cluster
EOF
}

install_operator() {
  local operator_chart_version operator_image_tag
  operator_chart_version="$(require_env perses-operator-version)"
  operator_image_tag="v${operator_chart_version#v}"

  helm repo add perses "${HELM_REPO}" 2>/dev/null || true

  helm upgrade --install "${HELM_RELEASE}" perses/perses-operator \
    --version "${operator_chart_version}" \
    --namespace "${E2E_NAMESPACE}" \
    --create-namespace \
    --set certManager.enable=false \
    --set testFramework.enabled=false \
    --set "manager.image.tag=${operator_image_tag}" \
    --wait \
    --timeout 5m
}

apply_perses_fixtures() {
  local perses_version perses_image
  perses_version="$(require_env perses-version)"
  perses_version="${perses_version#v}"
  perses_image="persesdev/perses:v${perses_version}"

  sed "s|image: persesdev/perses:v[^[:space:]]*|image: ${perses_image}|" \
    "${ROOT_DIR}/test/e2e/manifests/perses/perses.yaml" | kubectl apply -f -

  kubectl apply -f "${ROOT_DIR}/test/e2e/manifests/perses/prometheus-global-datasource.yaml"
}

wait_for_stack() {
  kubectl wait --for=condition=Ready pods -l control-plane=controller-manager \
    -n "${E2E_NAMESPACE}" --timeout=300s
  kubectl wait --for=jsonpath='{.status.readyReplicas}'=1 statefulset/perses \
    -n "${E2E_NAMESPACE}" --timeout=600s
}

apply_dashboards() {
  "${MAKE:-make}" -C "${ROOT_DIR}" build-dashboards-local \
    PROJECT=perses-dev DATASOURCE=prometheus-datasource
  kubectl apply -R -f "${ROOT_DIR}/built/dashboards/operator/"
}

install_stack() {
  install_operator
  apply_perses_fixtures
  wait_for_stack
  apply_dashboards
}

kind_up() {
  local kind_node_image
  kind_node_image="$(require_env kind-node-image)"

  if ! kind get clusters 2>/dev/null | grep -Fxq "${E2E_KIND_CLUSTER}"; then
    kind create cluster --name "${E2E_KIND_CLUSTER}" \
      --config "${ROOT_DIR}/test/e2e/kind/config.yml" \
      --image "${kind_node_image}"
  fi

  install_stack
}

kind_down() {
  kind delete cluster --name "${E2E_KIND_CLUSTER}"
}

case "${1:-}" in
  up)
    kind_up
    ;;
  install)
    install_stack
    ;;
  down)
    kind_down
    ;;
  *)
    usage >&2
    exit 1
    ;;
esac
