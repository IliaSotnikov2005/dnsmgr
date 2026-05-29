package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	"github.com/IliaSotnikov2005/dnsmgr/proto"
	"github.com/IliaSotnikov2005/dnsmgr/server/internal/adapter/repository/file"
	"github.com/IliaSotnikov2005/dnsmgr/server/internal/config"
	"github.com/IliaSotnikov2005/dnsmgr/server/internal/controller"
	"github.com/IliaSotnikov2005/dnsmgr/server/internal/usecase"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func NewApp(log *slog.Logger, storagePath string, grpcCfg config.GRPCConfig) *App {
	repo := file.NewFileRepository(log, storagePath)
	uc := usecase.NewDNSUseCase(log, repo)

	handler := controller.NewDNSHandler(log, uc)

	grpcServer := grpc.NewServer()
	proto.RegisterDNSServiceServer(grpcServer, handler)

	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       grpcCfg.Port,
	}
}

func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", "app.Run", err)
	}

	a.log.Info("gRPC server is running", slog.String("port", a.port))

	errCh := make(chan error, 1)

	go func() {
		if err := a.gRPCServer.Serve(l); err != nil {
			if err != grpc.ErrServerStopped {
				a.log.Error("gRPC server failed", slog.String("error", err.Error()))
				errCh <- err
			}
		}
	}()

	select {
	case <-ctx.Done():
		a.log.Info("stopping application due to signal")
	case err := <-errCh:
		a.log.Error("gRPC server error", slog.String("error", err.Error()))
		return err
	}

	a.gRPCServer.GracefulStop()
	a.log.Info("application stopped")

	return nil
}
