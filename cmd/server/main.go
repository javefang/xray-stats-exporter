package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xfang/xray-stats-exporter/internal/client"
	"github.com/xfang/xray-stats-exporter/internal/exporter"
)

func main() {
	statsURL := os.Getenv("XRAY_STATS_URL")
	if statsURL == "" {
		statsURL = "http://127.0.0.1:11111/debug/vars"
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}

	httpClient := &client.HTTP{Client: &http.Client{Timeout: 10 * time.Second}}
	exporter.MustRegister(httpClient, statsURL)

	mux := http.NewServeMux()
	exporter.NewHandler(mux)
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{Addr: listenAddr, Handler: mux}

	go func() {
		log.Printf("listening on %s", listenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	log.Println("shutting down")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
}
