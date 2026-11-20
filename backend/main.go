package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"cwcu/backend/config"
	"cwcu/backend/modules/academic"
	"cwcu/backend/modules/auth"
	"cwcu/backend/modules/elearning"
	"cwcu/backend/modules/notification"
	"cwcu/backend/modules/user"
	"cwcu/backend/pkg/lavinmq"
	"cwcu/backend/pkg/response"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// ── Load Config ─────────────────────────────────────────────────────────
	cfg := config.Load()

	// ── Connect Infrastructure ───────────────────────────────────────────────
	db := config.NewDatabase(cfg)
	config.RunMigrations(db, cfg)

	redisClient := config.NewRedis(cfg)
	amqpClient := lavinmq.New(cfg.AMQPUrl)

	// ── Echo Instance ────────────────────────────────────────────────────────
	e := echo.New()
	e.HideBanner = true

	// ── Global Middleware ────────────────────────────────────────────────────
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))

	// ── Health Check ─────────────────────────────────────────────────────────
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, response.Success("CWCU API is running", map[string]string{
			"status":  "ok",
			"version": "1.0.0",
		}))
	})

	// ── API Routes ───────────────────────────────────────────────────────────
	api := e.Group("/api/v1")

	auth.RegisterRoutes(api, db, redisClient, cfg)
	user.RegisterRoutes(api, db, redisClient, cfg)
	academic.RegisterRoutes(api, db, redisClient, cfg, amqpClient)
	elearning.RegisterRoutes(api, db, redisClient, cfg, amqpClient)
	notification.StartConsumers(amqpClient, cfg)

	// ── Start Server with Graceful Shutdown ──────────────────────────────────
	go func() {
		log.Printf("🚀 Server starting on port %s", cfg.AppPort)
		if err := e.Start(":" + cfg.AppPort); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	amqpClient.Close()
	log.Println("✅ Server gracefully shut down")
}
