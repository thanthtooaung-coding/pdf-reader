package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/app"
	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/config"
)

func main() {
	cfg := config.Load()
	log := newLogger(cfg.LogLevel)

	application, err := app.Bootstrap(cfg)
	if err != nil {
		log.WithError(err).Fatal("bootstrap")
	}
	defer application.Close()

	addr := ":" + cfg.Port
	go func() {
		log.Infof("pdf-reader-backend listening on %s", addr)
		if err := application.App.Listen(addr); err != nil && !errors.Is(err, fiber.ErrServiceUnavailable) {
			log.WithError(err).Fatal("server stopped")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := application.App.ShutdownWithContext(ctx); err != nil {
		log.WithError(err).Error("graceful shutdown failed")
	}
	log.Info("bye")
}

func newLogger(level string) *logrus.Logger {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	l.SetLevel(lvl)
	return l
}
