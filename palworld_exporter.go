package main

import (
	"net/http"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jimmysharp/palworld_exporter/collector"
	"github.com/jimmysharp/palworld_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	app           = kingpin.New("palworld_exporter", "Prometheus exporter for Palworld")
	listenAddress = app.
			Flag("web.listen-address", "Address to expose metrics.").
			Default(":18212").
			Envar("LISTEN_ADDRESS").
			String()
	scrapeURI = app.
			Flag("scrape_uri", "URI to Palworld REST API metrics endpoint.").
			Default("http://localhost:8212/v1/api/metrics").
			Envar("SCRAPE_URI").
			String()
	httpUsername = app.
			Flag("http_user", "Username for Palworld REST API basic authentication.").
			Envar("HTTP_USER").
			Default("admin").
			String()
	httpPassword = app.
			Flag("http_password", "Password for Palworld REST API basic authentication.").
			Envar("HTTP_PASSWORD").
			Required().
			String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	config := &config.Config{
		ListenAddress: *listenAddress,
		ScrapeURI:     *scrapeURI,
		HttpUsername:  *httpUsername,
		HttpPassword:  *httpPassword,
	}
	exporter := collector.NewExporter(config)
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(config.ListenAddress, nil); err != nil {
		os.Exit(1)
	}
}
