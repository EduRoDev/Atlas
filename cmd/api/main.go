package main

import (
	"net/http"
	"os"

	"github.com/EduRoDev/Atlas/internal/config"
	"github.com/EduRoDev/Atlas/internal/platform/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		os.Stderr.WriteString("Error cargando configuracion: " + err.Error() + "\n")
		os.Exit(1)
	}

	log := logger.New(cfg.App.Env)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	addr := ":" + cfg.App.Port
	log.Info("Servidor arrancando", "addr", addr, "env", cfg.App.Env)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Error("El servidor se detuvo", "error", err)
		os.Exit(1)
	}
}
