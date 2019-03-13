package main

import (
	"github.com/ProdriveTechnologies/vmware_exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

var (
	vmwareDatastoreCapacityDesc = util.NewVmwareDesc(
		"datastore",
		"capacity_bytes",
		"The total capacity of the datastore in bytes",
		"datastore")
	vmwareDatastoreFreeSpaceDesc = util.NewVmwareDesc(
		"datastore",
		"freespace_bytes",
		"The free space available on the datastore in bytes",
		"datastore")
)

func (e *vmwareExporter) retrieveDatastores(ch chan<- prometheus.Metric) error {
	manager := view.NewManager(e.Client.Client)
	view, err := manager.CreateContainerView(e.Context, e.Client.Client.ServiceContent.RootFolder, []string{"Datastore"}, true)
	if err != nil {
		return err
	}
	defer view.Destroy(e.Context)
	var datastores []mo.Datastore
	err = view.Retrieve(e.Context, []string{"Datastore"}, []string{"summary"}, &datastores)
	if err != nil {
		return err
	}
	for _, datastore := range datastores {
		ch <- prometheus.MustNewConstMetric(vmwareDatastoreCapacityDesc, prometheus.GaugeValue, float64(datastore.Summary.Capacity), datastore.Summary.Name)
		ch <- prometheus.MustNewConstMetric(vmwareDatastoreFreeSpaceDesc, prometheus.GaugeValue, float64(datastore.Summary.FreeSpace), datastore.Summary.Name)
	}
	return nil
}

func describeDatastores(ch chan<- *prometheus.Desc) {
	ch <- vmwareDatastoreCapacityDesc
	ch <- vmwareDatastoreFreeSpaceDesc
}
