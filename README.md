# Prometheus exporter for VMware vSphere

This exporter for VMware vSphere makes request against the vSphere SDK to
retrieve basic metrics about datastores, hosts, and VMs, using the [govmomi
library](https://github.com/vmware/govmomi/).

## Building this exporter

The exporter can be built using Bazel:

    bazel build //...

## Using this exporter

The exporter retrieves credentials and URL of vSphere from the `VSPHERE_URL` environment.

Once that variable set, the following command will start the exporter, causing
it to listen on TCP port 9536:

    ./vmware_exporter
