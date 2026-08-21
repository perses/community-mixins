# E2E manifests

| Path | Purpose |
| ---- | ------- |
| `perses/` | `Perses` CR and stub `PersesGlobalDatasource` |

The operator is installed via Helm; Perses fixtures and dashboards are applied by
`test/e2e/setup.sh` (`make e2e-up` / `make e2e-install`). Version pins are in
[`.github/env`](../../.github/env).
