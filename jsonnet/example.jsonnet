local dashboards = import 'dashboards.libsonnet';

local config = {
    datasource: 'custom-datasource',
};

dashboards(config)
