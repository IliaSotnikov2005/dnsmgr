package app

import (
	"fmt"
	"log/slog"
	"net"
	"os"
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
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", "app.Run", err)
	}

	a.log.Info("gRPC server is running", slog.String("port", a.port))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := a.gRPCServer.Serve(l); err != nil {
			a.log.Error("gRPC server failed", slog.String("error", err.Error()))
		}
	}()

	sign := <-stop
	a.log.Info("stopping application", slog.String("signal", sign.String()))

	a.gRPCServer.GracefulStop()
	a.log.Info("application stopped")

	return nil
}
