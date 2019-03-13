package main

import (
	"github.com/ProdriveTechnologies/vmware_exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

var (
	vmwareHostCpuUsedHzDesc = util.NewVmwareDesc(
		"host",
		"cpu_used_hertz",
		"The current CPU used of the ESXi host in Hertz.",
		"host")
	vmwareHostCpuAvailableHzDesc = util.NewVmwareDesc(
		"host",
		"cpu_available_hertz",
		"The available CPU power of the ESXi host in Hertz.",
		"host")
	vmwareHostMemoryUsedBytesDesc = util.NewVmwareDesc(
		"host",
		"memory_used_bytes",
		"The memory currently used of the ESXi host in Bytes.",
		"host")
	vmwareHostMemoryTotalBytesDesc = util.NewVmwareDesc(
		"host",
		"memory_total_bytes",
		"The available memory of the ESXi host in Bytes.",
		"host")
	vmwareHostUptimeSecondsDesc = util.NewVmwareDesc(
		"host",
		"uptime_seconds",
		"The uptime of the host in seconds.",
		"host")
)

func (e *vmwareExporter) retrieveHosts(ch chan<- prometheus.Metric) error {
	manager := view.NewManager(e.Client.Client)
	view, err := manager.CreateContainerView(e.Context, e.Client.Client.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		return err
	}
	defer view.Destroy(e.Context)
	var hosts []mo.HostSystem
	err = view.Retrieve(e.Context, []string{"HostSystem"}, []string{"summary"}, &hosts)
	if err != nil {
		return err
	}
	for _, host := range hosts {
		cpuHertz := float64(int64(host.Summary.QuickStats.OverallCpuUsage)*1000000)
		ch <- prometheus.MustNewConstMetric(vmwareHostCpuUsedHzDesc, prometheus.GaugeValue, cpuHertz, host.Summary.Config.Name)
		cpuTotal := float64(int64(host.Summary.Hardware.CpuMhz)*int64(host.Summary.Hardware.NumCpuCores)*1000000)
		ch <- prometheus.MustNewConstMetric(vmwareHostCpuAvailableHzDesc, prometheus.GaugeValue, cpuTotal, host.Summary.Config.Name)
		memoryUsed := float64(int64(host.Summary.QuickStats.OverallMemoryUsage)*1000000000)
		ch <- prometheus.MustNewConstMetric(vmwareHostMemoryUsedBytesDesc, prometheus.GaugeValue, memoryUsed, host.Summary.Config.Name)
		memoryTotal := float64(host.Summary.Hardware.MemorySize)
		ch <- prometheus.MustNewConstMetric(vmwareHostMemoryTotalBytesDesc, prometheus.GaugeValue, memoryTotal, host.Summary.Config.Name)
		ch <- prometheus.MustNewConstMetric(vmwareHostUptimeSecondsDesc, prometheus.GaugeValue, float64(host.Summary.QuickStats.Uptime), host.Summary.Config.Name)
	}
	return nil
}

func describeHosts(ch chan<- *prometheus.Desc) {
	ch <- vmwareHostCpuUsedHzDesc
	ch <- vmwareHostCpuAvailableHzDesc
	ch <- vmwareHostMemoryUsedBytesDesc
	ch <- vmwareHostMemoryTotalBytesDesc
	ch <- vmwareHostUptimeSecondsDesc
}

