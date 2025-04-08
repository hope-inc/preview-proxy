package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/hope-inc/preview-proxy/pkg"
)

func main() {
	cfg := pkg.LoadConfig()
	port := cfg.Port

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	http.Handle("/proxy/healthz", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
		w.WriteHeader(http.StatusOK)
	}))
	proxy := pkg.NewReverseProxy(cfg.OriginScheme, cfg.OriginBaseDomain, cfg.OriginPort)
	http.Handle("/", proxy)
	listenHost := fmt.Sprintf(":%v", port)
	slog.Info("Listen on", "host", listenHost)
	server := &http.Server{
		Addr:              listenHost,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
