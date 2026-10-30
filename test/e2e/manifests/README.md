# E2E manifests

| Path | Purpose |
| ---- | ------- |
| `perses/` | `Perses` CR, stub `PersesGlobalDatasource`, and Perses-format e2e `Project` |

The operator is installed via Helm; Perses fixtures and dashboards are applied by
`test/e2e/setup.sh` (`make e2e-up` / `make e2e-install`). Version pins are in
[`.github/env`](../../.github/env).

Operator CRs use project `perses-dev`. Perses-native YAML is applied via `percli`
to project `community-mixins` (see `community-mixins-project.yaml`).
