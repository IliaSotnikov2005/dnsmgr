package usecase

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type DNSUseCase struct {
	log     *slog.Logger
	service DNSService
	output  io.Writer
}

func NewDNSUseCase(log *slog.Logger, service DNSService, output io.Writer) *DNSUseCase {
	return &DNSUseCase{log: log, service: service, output: output}
}

func (uc *DNSUseCase) Add(ctx context.Context, ip string) {
	ucDns, err := uc.service.Add(ctx, ip)
	if err != nil {
		fmt.Fprintf(uc.output, "Failed to add DNS %s: %s\n", ip, err.Error())
		return
	}

	fmt.Fprintf(uc.output, "DNS %s succesfully added\n", ucDns.Ip)
}

func (uc *DNSUseCase) Remove(ctx context.Context, ip string) {
	ucDns, err := uc.service.Remove(ctx, ip)
	if err != nil {
		uc.log.Error("failed to remove dns", "ip", ucDns.Ip, "err", err)
		fmt.Fprintf(os.Stderr, "Failed to remove DNS %s: %s\n", ucDns.Ip, err.Error())
		return
	}

	fmt.Fprintf(uc.output, "DNS %s succesfully removed\n", ucDns.Ip)
}

func (uc *DNSUseCase) List(ctx context.Context) {
	dnsList, err := uc.service.List(ctx)
	if err != nil {
		uc.log.Error("failed to list dns", "err", err)
		fmt.Fprintf(os.Stderr, "Failed to list DNS: %s\n", err.Error())
		return
	}

	fmt.Fprintln(uc.output, "DNS list:")
	for _, dns := range dnsList {
		fmt.Fprintln(uc.output, dns)
	}
}
