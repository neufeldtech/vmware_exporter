package util

import (
	"github.com/prometheus/client_golang/prometheus"
)

func NewVmwareDesc(objectType string, metricName string, description string, otherLabels ...string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName("vmware", objectType, metricName),
		description,
		otherLabels,
		nil)
}

func NewVmwareGauge(desc *prometheus.Desc, value float64, id int, name string, otherLabels ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value, otherLabels...)
}
