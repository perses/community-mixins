# Testing

## Unit tests

```sh
make unit-test
```

## Operator E2E tests

Validates that mixin-built `PersesDashboard` CRs reconcile through
[perses-operator](https://github.com/perses/perses-operator) on a local kind cluster.
CI runs via `.github/workflows/e2e.yaml`.

**Prerequisites:** Docker, [kind](https://kind.sigs.k8s.io/), Helm 3, kubectl, Go.

### Local run

```sh
make e2e        # create cluster, install stack, run tests
make e2e-down   # delete the kind cluster
```

Version pins (operator, Perses, kind node image, etc.) live in [`.github/env`](.github/env).

### Re-run tests only

If the cluster is already up from a previous `make e2e`:

```sh
make test-e2e-operator
```

Tests require `OPERATOR_E2E=true` (set automatically by `make test-e2e-operator`).
They use the same default kubeconfig as `kubectl` after `make e2e` / kind cluster creation.
