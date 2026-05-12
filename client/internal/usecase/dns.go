package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/IliaSotnikov2005/dnsmgr/client/internal/domain"
)

type DNSUseCase struct {
	log     *slog.Logger
	service DNSService
}

func NewDNSUseCase(log *slog.Logger, service DNSService) *DNSUseCase {
	return &DNSUseCase{log: log, service: service}
}

func (uc *DNSUseCase) Add(ctx context.Context, ip string) (domain.DNS, error) {
	if err := uc.validateIP(ip); err != nil {
		return domain.DNS{}, err
	}

	return uc.service.Add(ctx, ip)
}

func (uc *DNSUseCase) Remove(ctx context.Context, ip string) (domain.DNS, error) {
	if err := uc.validateIP(ip); err != nil {
		return domain.DNS{}, err
	}

	return uc.service.Remove(ctx, ip)
}

func (uc *DNSUseCase) List(ctx context.Context) ([]string, error) {
	return uc.service.List(ctx)
}

func (u *DNSUseCase) validateIP(ip string) error {
	ip = strings.TrimSpace(ip)

	if ip == "" {
		return fmt.Errorf("IP address cannot be empty")
	}

	if net.ParseIP(ip) == nil {
		return domain.ErrInvalidIP
	}

	return nil
}
