package main

import (
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":18212", nil); err != nil {
		os.Exit(1)
	}
}
