# Testing

## Unit tests

```sh
make unit-test
```

Offline Perses schema validation (no cluster) remains in CI via
`make build-dashboards-local` + `percli validate` in the `check-dashboards` job.

## Dashboard E2E tests

Validates mixin-built dashboards against a live Perses on a local kind cluster:

1. **Operator format** — `PersesDashboard` CRs reconcile through
   [perses-operator](https://github.com/perses/perses-operator)
   (`Available=True`, `Degraded=False`).
2. **Perses format** — native YAML is applied with `percli` to project
   `community-mixins` and listed via the Perses HTTP API.

**Prerequisites:** Docker, [kind](https://kind.sigs.k8s.io/), Helm 3, kubectl, Go,
[percli](https://github.com/perses/perses/releases) (for Perses-format apply).

The default kind cluster name is `community-mixins-e2e` (override with `E2E_KIND_CLUSTER`).

### Local run

```sh
make e2e        # create cluster, install stack, run operator + Perses-format tests
make e2e-down   # delete the kind cluster
```

`make e2e` runs `test/e2e/setup.sh up`, which creates the cluster (if needed),
installs the operator via Helm, applies Perses fixtures, builds and applies
operator CRs (`perses-dev`), then builds and applies Perses-native YAML to
project `community-mixins` via `percli`.

### Install stack without running tests

If you already have the kind cluster but need to reinstall the operator, Perses, and dashboards:

```sh
make e2e-install
```

### Re-run tests only

If the cluster is already up from a previous `make e2e`:

```sh
make test-e2e-operator   # OPERATOR_E2E=true
make test-e2e-perses     # PERSES_E2E=true
```

Tests use the same default kubeconfig as `kubectl` after `make e2e` / kind cluster creation.
