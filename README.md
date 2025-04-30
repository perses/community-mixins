# Perses Community Dashboards

Welcome to the **Perses Community Dashboards** repository! This project is designed to provide Prometheus mixins tailored for the Perses platform. Developed with the **Perses Go SDK**, these dashboards are modular, reusable, and simple to integrate into various observability setups.

## Overview of Available Dashboards

### Prometheus Dashboards
- Prometheus Overview
- Prometheus Remote Write

### Node Exporter Dashboards
- Nodes
- Cluster USE Method

### AlertManager Dashboards
- AlertManager Overview

### Thanos Dashboards
- Receive Overview
- Query Overview
- Query Frontend Overview
- Store Gateway Overview
- Compact Overview
- Ruler Overview

## Library Panels

In addition to the community dashboards, this repository also offers a **library of reusable panels**. These panels can be used as building blocks for custom dashboard creation, enabling you to craft tailored setups to suit specific observability needs.

## Rendering Dashboards

To render and generate the dashboards, run the following command:

```bash
make build-dashboards
```

The generated dashboard files will be stored as **YAML files** in the `examples/dashboards/` directory by default (both in native Perses and Perses Operator format). You can then import these files into your Perses instances.

## Local Development Guide

For local development, you can quickly spin up a Perses environment with the following command:

```bash
make start-demo
```

This command initializes a local Perses instance that includes predefined resources such as Projects and DataSources. Once the instance is running, you can access the Perses UI at [http://localhost:8080](http://localhost:8080).

### Applying Dashboards with `percli`

To apply the dashboards to your Perses instance, use the [percli](https://pkg.go.dev/github.com/perses/perses/cmd/percli) tool with the following command:

```bash
percli apply -d examples/dashboards/perses/
```

This will deploy the dashboards from the `examples/dashboards/perses` directory to your local Perses instance.

### Applying Dashboards with Perses Operator

Once you have installed Perses Operator, with a `Perses` instance and `PersesDatasource` object in your cluster (following the instructions [here](https://github.com/perses/perses-operator?tab=readme-ov-file#running-on-the-cluster)), you can apply the generated `PersesDashboard` objects using:

```bash
kubectl apply -f examples/dashboards/operator/
```
