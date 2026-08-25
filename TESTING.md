# Testing

## Unit tests

```sh
make unit-test
```

## Operator E2E tests

Validates that mixin-built `PersesDashboard` CRs reconcile through
[perses-operator](https://github.com/perses/perses-operator) on a local kind cluster.

**Prerequisites:** Docker, [kind](https://kind.sigs.k8s.io/), Helm 3, kubectl, Go.

The default kind cluster name is `community-mixins-e2e` (override with `E2E_KIND_CLUSTER`).

### Local run

```sh
make e2e        # create cluster, install stack, run tests
make e2e-down   # delete the kind cluster
```

`make e2e` runs `test/e2e/setup.sh up`, which creates the cluster (if needed), installs the operator via Helm, applies Perses fixtures, builds dashboards to `built/dashboards/operator/`, and applies them.

### Install stack without running tests

If you already have the kind cluster but need to reinstall the operator, Perses, and dashboards:

```sh
make e2e-install
```

### Re-run tests only

If the cluster is already up from a previous `make e2e`:

```sh
make test-e2e-operator
```

Tests require `OPERATOR_E2E=true` (set automatically by `make test-e2e-operator`).
They use the same default kubeconfig as `kubectl` after `make e2e` / kind cluster creation.
