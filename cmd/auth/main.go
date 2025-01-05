package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main(){
	cfg := config.MustLoad()

	fmt.Println(cfg)
	// TODO: logger

	log := setupLogger(cfg.Env)

	log.Info("starting app", slog.Any("cfg", cfg))


	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL.String())

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()

	log.Info("gRPC server stopped")
	// TODO: sso

	// TODO: gRPC-server start
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}