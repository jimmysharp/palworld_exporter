package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/jimmysharp/palworld_exporter/collector"
	"github.com/jimmysharp/palworld_exporter/config"
	"github.com/jimmysharp/palworld_exporter/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
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
	logLevel = app.
			Flag("log.level", "Only log messages with the given severity or above. If log.format is set 'default', this option is ignored. Valid levels: [debug, info, warn, error]").
			Default("info").
			Envar("LOG_LEVEL").
			HintOptions("debug", "info", "warn", "error").
			String()
	logFormat = app.
			Flag("log.format", "Output format of log messages. Valid formats: [default, text, json]").
			Default("default").
			Envar("LOG_FORMAT").
			HintOptions("default", "text", "json").
			String()
)

func createServer(config *config.Config, logger *slog.Logger) *http.Server {
	exporter := collector.NewExporter(config, logger)
	prometheus.MustRegister(exporter)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:    config.ListenAddress,
		Handler: mux,
	}

	return server
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.HelpFlag.Short('h')
	app.Version(version.Print("palworld_exporter"))
	kingpin.MustParse(app.Parse(os.Args[1:]))

	config := &config.Config{
		ListenAddress: *listenAddress,
		ScrapeURI:     *scrapeURI,
		HttpUsername:  *httpUsername,
		HttpPassword:  *httpPassword,
		LogLevel:      *logLevel,
		LogFormat:     *logFormat,
	}

	logger := log.NewLogger(config)

	server := createServer(config, logger)
	go func() {
		logger.Info("Starting palworld_exporter", slog.String("version", version.Info()))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error starting server", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	stop()
	logger.Info("Caught signal, Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Error safely shutting down server", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
