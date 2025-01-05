package grpcapp

import (
	authgrpc "auth/internal/grpc/auth"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	log *slog.Logger
	gRPCServer *grpc.Server
	port int
}

func New(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer)

	return &App{
		log: log,
		gRPCServer: gRPCServer,
		port: port,
	}
}

func (a *App) MustRun() {
	const op = "grpcapp.Run"

	a.log.With(slog.String("op", op))

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		panic(err.Error())
	}

	a.log.Info("gRPC server is running", slog.String("addr", listen.Addr().String()))

	err = a.gRPCServer.Serve(listen)
	if err != nil {
		panic(err.Error())
	}
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}