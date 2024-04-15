package collector

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/jimmysharp/palworld_exporter/config"
)

const (
	defaultTimeout = 5 * time.Second
)

type PalworldMetricsResponse struct {
	ServerFps        int     `json:"serverfps"`
	CurrentPlayerNum int     `json:"currentplayernum"`
	ServerFrameTime  float64 `json:"serverframetime"`
	MaxPlayerNum     int     `json:"maxplayernum"`
	Uptime           int     `json:"uptime"`
}

type PalworldClient struct {
	ScrapeURI string
}

func NewPalworldClient(config *config.Config) *PalworldClient {
	return &PalworldClient{
		ScrapeURI: config.ScrapeURI,
	}
}

func (c *PalworldClient) getPalworldMetrics() (*PalworldMetricsResponse, error) {
	client := http.Client{Timeout: time.Duration(defaultTimeout)}
	req, err := http.NewRequest(http.MethodGet, c.ScrapeURI, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	metrics := &PalworldMetricsResponse{}
	err = json.Unmarshal(body, metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
