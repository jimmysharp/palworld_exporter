package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jimmysharp/palworld_exporter/collector"
	"github.com/jimmysharp/palworld_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	shutdownTimeout = 5 * time.Second
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	kingpin.MustParse(app.Parse(os.Args[1:]))

	config := &config.Config{
		ListenAddress: *listenAddress,
		ScrapeURI:     *scrapeURI,
		HttpUsername:  *httpUsername,
		HttpPassword:  *httpPassword,
	}
	exporter := collector.NewExporter(config)
	prometheus.MustRegister(exporter)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:    config.ListenAddress,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		os.Exit(1)
	}
}
