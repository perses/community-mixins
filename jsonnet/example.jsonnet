local dashboards = import 'dashboards.libsonnet';

local config = {
  datasource: 'custom-datasource',
  components: ['kubernetes', 'thanos', 'etcd', 'blackbox-exporter', 'node-exporter', 'alertmanager', 'prometheus', 'perses'],
};

dashboards(config)
