package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jimmysharp/palworld_exporter/config"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/assert"
)

type DummyLogHandler struct{}

func (*DummyLogHandler) Enabled(context.Context, slog.Level) bool   { return true }
func (*DummyLogHandler) Handle(context.Context, slog.Record) error  { return nil }
func (h *DummyLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *DummyLogHandler) WithGroup(name string) slog.Handler       { return h }

func newDummyLogger() *slog.Logger {
	return slog.New(&DummyLogHandler{})
}

type metrics struct {
	palworldUp               float64
	palworldServerFps        float64
	palworldCurrentPlayerNum float64
	palworldServerFrameTime  float64
	palworldMaxPlayerNum     float64
	palworldUptime           float64
}

var (
	normalResponse = `{
		"currentplayernum": 0,
		"serverfps": 43,
		"serverframetime": 22.738655090332031,
		"maxplayernum": 32,
		"uptime": 66925
	}`
	normalMetrics = metrics{
		palworldUp:               1,
		palworldServerFps:        43,
		palworldCurrentPlayerNum: 0,
		palworldServerFrameTime:  22.738655090332031,
		palworldMaxPlayerNum:     32,
		palworldUptime:           66925,
	}
)

func checkExportedMetrics(t *testing.T, stubHandler http.Handler) map[string]*dto.MetricFamily {
	// Palworld REST API stub server
	stubServer := httptest.NewServer(stubHandler)
	defer stubServer.Close()

	// Test server
	config := &config.Config{
		ListenAddress: ":18212",
		ScrapeURI:     stubServer.URL + "/v1/api/metrics",
		HttpUsername:  "admin",
		HttpPassword:  "password",
		LogLevel:      "info",
		LogFormat:     "default",
	}
	logger := newDummyLogger()
	handler := createServer(config, logger).Handler
	server := httptest.NewServer(handler)
	defer server.Close()

	// Scrape metrics and parse
	resp, err := http.Get(server.URL + "/metrics")
	assert.NoError(t, err)
	defer resp.Body.Close()

	var parser expfmt.TextParser
	metrics, err := parser.TextToMetricFamilies(resp.Body)
	assert.NoError(t, err)

	return metrics
}

func TestNormalMetrics(t *testing.T) {
	stubHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(normalResponse))
		assert.NoError(t, err)
	})
	metrics := checkExportedMetrics(t, stubHandler)

	// Assert
	assert.Contains(t, metrics, "palworld_up")
	assert.Equal(t, dto.MetricType_GAUGE, metrics["palworld_up"].GetType())
	assert.Equal(t, normalMetrics.palworldUp, metrics["palworld_up"].GetMetric()[0].GetGauge().GetValue())
	assert.Contains(t, metrics, "palworld_server_fps")
	assert.Equal(t, dto.MetricType_GAUGE, metrics["palworld_server_fps"].GetType())
	assert.Equal(t, normalMetrics.palworldServerFps, metrics["palworld_server_fps"].GetMetric()[0].GetGauge().GetValue())
	assert.Contains(t, metrics, "palworld_current_player_num")
	assert.Equal(t, dto.MetricType_GAUGE, metrics["palworld_current_player_num"].GetType())
	assert.Equal(t, normalMetrics.palworldCurrentPlayerNum, metrics["palworld_current_player_num"].GetMetric()[0].GetGauge().GetValue())
	assert.Contains(t, metrics, "palworld_server_frame_time")
	assert.Equal(t, dto.MetricType_GAUGE, metrics["palworld_server_frame_time"].GetType())
	assert.Equal(t, normalMetrics.palworldServerFrameTime, metrics["palworld_server_frame_time"].GetMetric()[0].GetGauge().GetValue())
	assert.Contains(t, metrics, "palworld_max_player_num")
	assert.Equal(t, dto.MetricType_GAUGE, metrics["palworld_max_player_num"].GetType())
	assert.Equal(t, normalMetrics.palworldMaxPlayerNum, metrics["palworld_max_player_num"].GetMetric()[0].GetGauge().GetValue())
	assert.Contains(t, metrics, "palworld_uptime")
	assert.Equal(t, dto.MetricType_COUNTER, metrics["palworld_uptime"].GetType())
	assert.Equal(t, normalMetrics.palworldUptime, metrics["palworld_uptime"].GetMetric()[0].GetCounter().GetValue())
}

func TestInvalidStatusCode(t *testing.T) {
	stubHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	metrics := checkExportedMetrics(t, stubHandler)

	// Assert
	assert.Contains(t, metrics, "palworld_up")
	assert.Equal(t, dto.MetricType_GAUGE, metrics["palworld_up"].GetType())
	assert.Equal(t, float64(0), metrics["palworld_up"].GetMetric()[0].GetGauge().GetValue())
}
