package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"context"
	"time"

	grpcclient "github.com/IliaSotnikov2005/dnsmgr/client/internal/adapter/grpc"
	"github.com/IliaSotnikov2005/dnsmgr/client/internal/domain"
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

	var err error
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect to the server", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Error("failed to close the connection", "error", err)
		}
	}()

	grpcClient := proto.NewDNSServiceClient(conn)
	adapter := grpcclient.NewClient(grpcClient)

	uc := usecase.NewDNSUseCase(log, adapter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch {
	case *addIP != "":
		var res domain.DNS
		res, err = uc.Add(ctx, *addIP)
		if err == nil {
			fmt.Printf("DNS %s successfully added\n", res.Ip)
		}
	case *removeIP != "":
		var res domain.DNS
		res, err = uc.Remove(ctx, *removeIP)
		if err == nil {
			fmt.Printf("DNS %s successfully removed\n", res.Ip)
		}
	case *listAll:
		var list []string
		list, err = uc.List(ctx)
		if err == nil {
			fmt.Println("DNS list:")
			for _, ip := range list {
				fmt.Println(ip)
			}
		}
	default:
		flag.Usage()
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
