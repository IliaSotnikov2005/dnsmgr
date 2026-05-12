package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/IliaSotnikov2005/dnsmgr/server/internal/domain"
)

type DNSUseCase struct {
	log  *slog.Logger
	repo DNSRepository
}

func NewDNSUseCase(log *slog.Logger, repo DNSRepository) *DNSUseCase {
	return &DNSUseCase{
		log:  log,
		repo: repo,
	}
}

func (u *DNSUseCase) Add(ctx context.Context, ip string) (domain.DNS, error) {
	u.log.Info("Adding DNS entry", "ip", ip)
	if net.ParseIP(ip) == nil {
		return domain.DNS{}, domain.ErrInvalidIP
	}

	list, err := u.repo.Get(ctx)
	if err != nil {
		return domain.DNS{}, fmt.Errorf("failed to get DNS list: %w", err)
	}

	for _, d := range list {
		if d.IP == ip {
			return domain.DNS{IP: ip}, domain.ErrAlreadyExists
		}
	}

	list = append(list, domain.DNS{IP: ip})
	if err := u.repo.Save(ctx, list); err != nil {
		return domain.DNS{}, err
	}

	return domain.DNS{IP: ip}, nil
}

func (u *DNSUseCase) Remove(ctx context.Context, ip string) (domain.DNS, error) {
	u.log.Info("Removing DNS entry", "ip", ip)
	if net.ParseIP(ip) == nil {
		return domain.DNS{}, domain.ErrInvalidIP
	}

	list, err := u.repo.Get(ctx)
	if err != nil {
		return domain.DNS{}, fmt.Errorf("failed to get DNS list: %w", err)
	}

	var newList []domain.DNS
	var found bool
	for _, d := range list {
		if d.IP == ip {
			found = true
			continue
		}
		newList = append(newList, d)
	}

	if !found {
		return domain.DNS{IP: ip}, domain.ErrNotFound
	}

	if err := u.repo.Save(ctx, newList); err != nil {
		return domain.DNS{}, err
	}

	return domain.DNS{IP: ip}, nil
}

func (u *DNSUseCase) List(ctx context.Context) ([]domain.DNS, error) {
	u.log.Info("Listing DNS entries")
	return u.repo.Get(ctx)
}
