package main

import (
	"net/http"
	"os"

	"github.com/jimmysharp/palworld_exporter/collector"
	"github.com/jimmysharp/palworld_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config := &config.Config{
		ListenAddress: ":18212",
		ScrapeURI:     "http://admin:pass@localhost:8212/v1/api/metrics",
	}
	exporter := collector.NewExporter(config)
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(config.ListenAddress, nil); err != nil {
		os.Exit(1)
	}
}
