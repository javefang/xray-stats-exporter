package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch(t *testing.T) {
	fixture := `{
  "observatory": {
    "vpn-jp": {
      "alive": true,
      "delay": 886,
      "outbound_tag": "vpn-jp",
      "last_seen_time": 1776006321,
      "last_try_time": 1776006321
    }
  },
  "stats": {
    "inbound": {
      "transparent": {
        "downlink": 23108149,
        "uplink": 10664547
      }
    },
    "outbound": {
      "direct": {
        "downlink": 11876941,
        "uplink": 7439849
      }
    },
    "user": {}
  }
}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fixture))
	}))
	defer srv.Close()

	cli := &HTTP{Client: http.DefaultClient}
	stats, err := cli.Fetch(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stats.Observatory) != 1 {
		t.Errorf("expected 1 observatory probe, got %d", len(stats.Observatory))
	}
	probe, ok := stats.Observatory["vpn-jp"]
	if !ok {
		t.Fatal("expected vpn-jp probe")
	}
	if !probe.Alive {
		t.Error("expected probe to be alive")
	}
	if probe.Delay != 886 {
		t.Errorf("expected delay 886, got %d", probe.Delay)
	}

	downlink, ok := stats.Stats.Inbound["transparent"]["downlink"]
	if !ok {
		t.Fatal("expected transparent inbound downlink")
	}
	if downlink != 23108149 {
		t.Errorf("expected downlink 23108149, got %d", downlink)
	}
}

func TestFetchHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	cli := &HTTP{Client: http.DefaultClient}
	_, err := cli.Fetch(context.Background(), srv.URL)
	if err == nil {
		t.Error("expected error for HTTP 500")
	}
}

func TestXrayStatsJSON(t *testing.T) {
	raw := `{"observatory":{"vpn-jp":{"alive":true,"delay":886}},"stats":{"inbound":{},"outbound":{},"user":{}}}`
	var s XrayStats
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if !s.Observatory["vpn-jp"].Alive {
		t.Error("expected alive=true")
	}
}
