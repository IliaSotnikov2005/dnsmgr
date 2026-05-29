package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/IliaSotnikov2005/dnsmgr/server/internal/app"
	"github.com/IliaSotnikov2005/dnsmgr/server/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %s", err)
		os.Exit(1)
	}

	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	application := app.NewApp(log, cfg.StoragePath, cfg.GRPC)

	if err := application.Run(); err != nil {
		log.Error("application error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
