package collector

import (
	"sync"

	"github.com/jimmysharp/palworld_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "palworld"
)

type Exporter struct {
	mutex sync.Mutex

	client *PalworldClient

	up               *prometheus.Desc
	serverFps        *prometheus.Desc
	currentPlayerNum *prometheus.Desc
	serverFrameTime  *prometheus.Desc
	maxPlayerNum     *prometheus.Desc
	uptime           *prometheus.Desc
}

func NewExporter(config *config.Config) *Exporter {
	return &Exporter{
		client: NewPalworldClient(config),
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Palworld server up",
			nil,
			nil),
		serverFps: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "server_fps"),
			"The Server FPS",
			nil,
			nil),
		currentPlayerNum: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "current_player_num"),
			"The number of current players",
			nil,
			nil),
		serverFrameTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "server_frame_time"),
			"Server frame time (ms)",
			nil,
			nil),
		maxPlayerNum: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "max_player_num"),
			"The maximum number of players",
			nil,
			nil),
		uptime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "uptime"),
			"The server uptime of seconds",
			nil,
			nil),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up
	ch <- e.serverFps
	ch <- e.currentPlayerNum
	ch <- e.serverFrameTime
	ch <- e.maxPlayerNum
	ch <- e.uptime
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	metrics, err := e.client.getPalworldMetrics()

	if err != nil {
		ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(e.serverFps, prometheus.GaugeValue, float64(metrics.ServerFps))
	ch <- prometheus.MustNewConstMetric(e.currentPlayerNum, prometheus.GaugeValue, float64(metrics.CurrentPlayerNum))
	ch <- prometheus.MustNewConstMetric(e.serverFrameTime, prometheus.GaugeValue, metrics.ServerFrameTime)
	ch <- prometheus.MustNewConstMetric(e.maxPlayerNum, prometheus.GaugeValue, float64(metrics.MaxPlayerNum))
	ch <- prometheus.MustNewConstMetric(e.uptime, prometheus.CounterValue, float64(metrics.Uptime))
}
