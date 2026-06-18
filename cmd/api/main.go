package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/EduRoDev/Atlas/internal/config"
	"github.com/EduRoDev/Atlas/internal/platform/cache"
	"github.com/EduRoDev/Atlas/internal/platform/database"
	"github.com/EduRoDev/Atlas/internal/platform/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		os.Stderr.WriteString("Error cargando configuracion: " + err.Error() + "\n")
		os.Exit(1)
	}

	log := logger.New(cfg.App.Env)
	ctx := context.Background()

	db, err := database.NewPostgresDB(ctx, cfg.Database.DSN(), cfg.App.Env)
	if err != nil {
		log.Error("Error conectando a la base de datos", "error", err)
		os.Exit(1)
	}
	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()
	log.Info("conectando a postgres")

	redisAddr := cfg.Redis.Host + ":" + cfg.Redis.Port
	rdb, err := cache.NewRedisClient(ctx, redisAddr, cfg.Redis.Password)

	if err != nil {
		log.Error("Error conectando a redis", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()
	log.Info("conectando a redis")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		hctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		w.Header().Set("Content-Type", "application/json")

		sqlDB, err := db.DB()
		if err != nil || sqlDB.PingContext(hctx) != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"error","postgres":"down"}`))
			return
		}
		if rdb.Ping(hctx).Err() != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"error","redis":"down"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","postgres":"up","redis":"up"}`))
	})

	addr := ":" + cfg.App.Port
	log.Info("servidor arrancando", "addr", addr, "env", cfg.App.Env)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Error("el servidor se detuvo", "error", err)
		os.Exit(1)
	}
}
