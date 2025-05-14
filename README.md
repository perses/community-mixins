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

### Kubernetes-mixin Dashboards
- API Server
- Compute Resources / Multi-Cluster
- Compute Resources / Cluster
- Compute Resources / Node
- Compute Resources / Namespace (Pods)
- Compute Resources / Namespace (Workloads)
- Compute Resources / Workloads
- Compute Resources / Pod
- Networking / Cluster
- Networking / Namespace (Pods)
- Networking / Namespace (Workload)
- Networking / Pod
- Networking / Workload
- Persistent Volume
- Kubelet
- Controller Manager
- Scheduler
- Proxy

### etcd-mixin Dashboards
- etcd Overview

### Blackbox Exporter
- Blackbox Exporter Overview

### Perses
- Perses Overview

## Library Panels

In addition to the community dashboards, this repository also offers a **library of reusable panels**. These panels can be used as building blocks for custom dashboard creation, enabling you to craft tailored setups to suit specific observability needs.

## Rendering Dashboards

To render and generate the dashboards, run the following command:

```bash
make build-dashboards
```

The generated dashboard files will be stored as **YAML files** in the `examples/dashboards/` directory by default and split by component (both in native Perses and Perses Operator format). You can then import these files into your Perses instances.

## Local Development Guide

For local development, you can quickly spin up a Perses environment with the following command:

```bash
make start-demo
```

This command initializes a local Perses instance that includes predefined resources such as Projects and DataSources. Once the instance is running, you can access the Perses UI at [http://localhost:8080](http://localhost:8080).

### Applying Dashboards with `percli`

You can use `percli` to apply the dashboards to your Perses instance. Use the [percli](https://pkg.go.dev/github.com/perses/perses/cmd/percli) tool with the following commands:

```bash
# Connect to your perses instance, for example if you want to use the demo perses
percli login http://localhost:8080/

# Add the dashboards to your instance (component-wise)
percli apply -d examples/dashboards/perses/prometheus/

# Add all the dashboards to your instance
find ./examples/dashboards/perses -type d -mindepth 1 -exec percli apply -d {} \;
```

This will deploy the dashboards from the `examples/dashboards/perses` directory to your local Perses instance.

### Applying Dashboards with Perses Operator

Once you have installed Perses Operator, with a `Perses` instance and `PersesDatasource` object in your cluster (following the instructions [here](https://github.com/perses/perses-operator?tab=readme-ov-file#running-on-the-cluster)), you can apply the generated `PersesDashboard` objects using:

```bash
# Add the dashboard CRDs to Perses (component-wise)
kubectl apply -f examples/dashboards/operator/kubernetes/

# Add all the dashboard CRDs to Perses
kubectl apply -f examples/dashboards/operator/ -R
```

### Dashboard Design Recommendations

You can find some recomendation for dashboard design [here](docs/README.md)
