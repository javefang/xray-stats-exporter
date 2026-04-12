package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type XrayStats struct {
	Observatory map[string]struct {
		Alive        bool   `json:"alive"`
		Delay        int    `json:"delay"`
		OutboundTag  string `json:"outbound_tag"`
		LastSeenTime int64  `json:"last_seen_time"`
		LastTryTime  int64  `json:"last_try_time"`
	} `json:"observatory"`
	Stats struct {
		Inbound  map[string]map[string]int64 `json:"inbound"`
		Outbound map[string]map[string]int64 `json:"outbound"`
		User     map[string]map[string]int64 `json:"user"`
	} `json:"stats"`
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Fetcher interface {
	Fetch(ctx context.Context, url string) (*XrayStats, error)
}

type HTTP struct {
	Client Doer
}

func (h *HTTP) Fetch(ctx context.Context, url string) (*XrayStats, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var stats XrayStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
