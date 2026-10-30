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

PERSES_FORMAT_PROJECT="community-mixins"
PERSES_SERVICE="perses"
PERSES_LOCAL_URL="http://127.0.0.1:8080"

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

apply_operator_dashboards() {
  "${MAKE:-make}" -C "${ROOT_DIR}" build-dashboards-local \
    PROJECT=perses-dev DATASOURCE=prometheus-datasource
  kubectl apply -R -f "${ROOT_DIR}/built/dashboards/operator/"
}

wait_for_perses_url() {
  local url="$1"
  local i
  for i in $(seq 1 60); do
    if curl -sf "${url}/api/v1/health" >/dev/null; then
      return 0
    fi
    sleep 1
  done
  echo "Perses API at ${url} did not become healthy" >&2
  return 1
}

apply_perses_format_dashboards() {
  local output_dir="${ROOT_DIR}/built/dashboards/perses"
  local pf_pid=""

  if ! command -v percli >/dev/null 2>&1; then
    echo "percli is required to apply Perses-format dashboards" >&2
    exit 1
  fi

  echo "Building Perses-format dashboards for project ${PERSES_FORMAT_PROJECT}"
  rm -rf "${output_dir}"
  (
    cd "${ROOT_DIR}"
    GOOS="$(go env GOOS)" CGO_ENABLED=0 go run main.go \
      --output-dir="${output_dir}" \
      --output=yaml \
      --project="${PERSES_FORMAT_PROJECT}" \
      --datasource=prometheus-datasource \
      --loki-datasource=loki-datasource
  )

  echo "Port-forwarding Perses service ${PERSES_SERVICE}"
  kubectl port-forward -n "${E2E_NAMESPACE}" "svc/${PERSES_SERVICE}" 8080:8080 >/tmp/perses-e2e-pf.log 2>&1 &
  pf_pid=$!

  cleanup_port_forward() {
    kill "${pf_pid}" 2>/dev/null || true
    wait "${pf_pid}" 2>/dev/null || true
  }
  trap cleanup_port_forward EXIT

  wait_for_perses_url "${PERSES_LOCAL_URL}"

  echo "Applying Perses project and dashboards via percli"
  percli login "${PERSES_LOCAL_URL}"
  percli apply -f "${ROOT_DIR}/test/e2e/manifests/perses/community-mixins-project.yaml"

  local dashboard_dir
  while IFS= read -r -d '' dashboard_dir; do
    percli apply -d "${dashboard_dir}"
  done < <(find "${output_dir}" -type d -mindepth 1 -print0)

  cleanup_port_forward
  trap - EXIT
}

apply_dashboards() {
  apply_operator_dashboards
  apply_perses_format_dashboards
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
