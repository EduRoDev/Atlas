package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/EduRoDev/Atlas/internal/auth/app"
	"github.com/EduRoDev/Atlas/internal/auth/infra/crypto"
	authhttp "github.com/EduRoDev/Atlas/internal/auth/infra/http"
	"github.com/EduRoDev/Atlas/internal/auth/infra/postgres"
	"github.com/EduRoDev/Atlas/internal/config"
	"github.com/EduRoDev/Atlas/internal/platform/cache"
	"github.com/EduRoDev/Atlas/internal/platform/database"
	"github.com/EduRoDev/Atlas/internal/platform/logger"
	"github.com/gin-gonic/gin"
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

	hasher := crypto.NewArgon2Hasher()
	userRepo := postgres.NewUserRepository(db)
	authService := app.NewService(userRepo, hasher)
	authHandler := authhttp.NewHandler(authService)

	r := gin.Default()
	authhttp.RegisterRoutes(r, authHandler)

	redisAddr := cfg.Redis.Host + ":" + cfg.Redis.Port
	rdb, err := cache.NewRedisClient(ctx, redisAddr, cfg.Redis.Password)

	if err != nil {
		log.Error("Error conectando a redis", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()
	log.Info("conectando a redis")

	r.GET("/health", func(c *gin.Context) {
		hctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		c.Header("Content-Type", "application/json")

		sqlDB, err := db.DB()
		if err != nil || sqlDB.PingContext(hctx) != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "postgres": "down"})
			return
		}
		if rdb.Ping(hctx).Err() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "redis": "down"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "postgres": "up", "redis": "up"})
	})

	addr := ":" + cfg.App.Port
	log.Info("servidor arrancando", "addr", addr, "env", cfg.App.Env)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Error("el servidor se detuvo", "error", err)
		os.Exit(1)
	}
}
