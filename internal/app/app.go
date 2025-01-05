package app

import (
	grpcapp "auth/internal/app/grpc"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL string) *App {
	// TODO: Storage init

	// TODO: init auth
	grpcApp:= grpcapp.New(log, grpcPort)

	return &App {
		GRPCServer: grpcApp,
	}
}