package main

import (
	"context"
	"flag"
	"net/url"
	"net/http"
	"os"

	"github.com/ProdriveTechnologies/vmware_exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/vmware/govmomi"
)

type vmwareExporter struct {
	Context context.Context
	Client *govmomi.Client
}

func NewVmwareExporter(vcenterUrl string) (*vmwareExporter, error) {
	parsedUrl, err := url.Parse(vcenterUrl)
	if err != nil {
		return &vmwareExporter{}, err
	}
	ctx := context.Background()
	client, err := govmomi.NewClient(ctx, parsedUrl, true)
	if err != nil {
		return &vmwareExporter{}, err
	}
	return &vmwareExporter{
			Context: ctx,
			Client: client,
	}, nil
}

var (
	vmwareScrapeSuccessDesc = util.NewVmwareDesc(
		"",
		"scrape_success",
		"Whether scraping the VMWare environment was successful.")

)

func (e *vmwareExporter) Collect(ch chan<- prometheus.Metric) {
	if err := e.retrieveHosts(ch); err != nil {
		log.Error(err)
		ch <- prometheus.MustNewConstMetric(
			vmwareScrapeSuccessDesc, prometheus.GaugeValue, 0.0)
		return
	}
	if err := e.retrieveVms(ch); err != nil {
		log.Error(err)
		ch <- prometheus.MustNewConstMetric(
			vmwareScrapeSuccessDesc, prometheus.GaugeValue, 0.0)
		return
	}
	if err := e.retrieveDatastores(ch); err != nil {
		log.Error(err)
		ch <- prometheus.MustNewConstMetric(
			vmwareScrapeSuccessDesc, prometheus.GaugeValue, 0.0)
		return
	} else {
		ch <- prometheus.MustNewConstMetric(
			vmwareScrapeSuccessDesc, prometheus.GaugeValue, 1.0)
	}
}

func (e *vmwareExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- vmwareScrapeSuccessDesc
	describeDatastores(ch)
	describeHosts(ch)
	describeVms(ch)
}

func main() {
	var (
		listenAddress = flag.String("web.listen-address", ":9536", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	)
	flag.Parse()

	vcenterUrl := os.Getenv("VSPHERE_URL")
	log.Infof("Connecting to vCenter")
	e, err := NewVmwareExporter(vcenterUrl)
	if err != nil {
		log.Fatal(err)
	}
	prometheus.MustRegister(e)

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
			<head><title>VMWare Exporter</title></head>
			<body>
			<h1>VMWare Exporter</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
	})
	log.Info("Listening on address:port => ", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
