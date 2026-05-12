package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"context"
	"time"

	"github.com/IliaSotnikov2005/dnsmgr/client/internal/adapter/grpc"
	"github.com/IliaSotnikov2005/dnsmgr/client/internal/usecase"
	"github.com/IliaSotnikov2005/dnsmgr/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := flag.String("addr", "localhost:50051", "server address")
	addIP := flag.String("add", "", "add DNS IP")
	removeIP := flag.String("rm", "", "remove DNS IP")
	listAll := flag.Bool("ls", false, "show all DNS")

	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	count := 0
	if *addIP != "" {
		count++
	}
	if *removeIP != "" {
		count++
	}
	if *listAll {
		count++
	}

	if count > 1 {
		fmt.Println("Error: Choose only one action.")
		flag.Usage()
		os.Exit(1)
	}

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect to the server", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	grpcClient := proto.NewDNSServiceClient(conn)
	adapter := grpcclient.NewClient(grpcClient)
	uc := usecase.NewDNSUseCase(log, adapter, os.Stdout)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	switch {
	case *addIP != "":
		uc.Add(ctx, *addIP)
	case *removeIP != "":
		uc.Remove(ctx, *removeIP)
	case *listAll:
		uc.List(ctx)
	default:
		flag.Usage()
	}
}
