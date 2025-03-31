package client

import (
	"encoding/json"
	"fmt"
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
	scrapeURI    string
	httpUsername string
	httpPassword string
}

func NewPalworldClient(config *config.Config) *PalworldClient {
	return &PalworldClient{
		scrapeURI:    config.ScrapeURI,
		httpUsername: config.HttpUsername,
		httpPassword: config.HttpPassword,
	}
}

func (c *PalworldClient) GetMetrics() (*PalworldMetricsResponse, error) {
	client := http.Client{Timeout: time.Duration(defaultTimeout)}
	req, err := http.NewRequest(http.MethodGet, c.scrapeURI, nil)
	if err != nil {
		return nil, fmt.Errorf("error building http request: %w", err)
	}
	req.SetBasicAuth(c.httpUsername, c.httpPassword)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error http response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid http status: %s (%d)", resp.Status, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	metrics := &PalworldMetricsResponse{}
	err = json.Unmarshal(body, metrics)
	if err != nil {
		return nil, fmt.Errorf("response has Invalid JSON body: %w", err)
	}

	return metrics, nil
}
