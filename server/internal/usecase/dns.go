package usecase

import (
	"context"
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

	err := u.repo.Update(ctx, func(list []domain.DNS) ([]domain.DNS, error) {
		for _, d := range list {
			if d.IP == ip {
				return nil, domain.ErrAlreadyExists
			}
		}

		return append(list, domain.DNS{IP: ip}), nil
	})
	if err != nil {
		return domain.DNS{}, err
	}

	return domain.DNS{IP: ip}, nil
}

func (u *DNSUseCase) Remove(ctx context.Context, ip string) (domain.DNS, error) {
	u.log.Info("Removing DNS entry", "ip", ip)
	if net.ParseIP(ip) == nil {
		return domain.DNS{}, domain.ErrInvalidIP
	}

	err := u.repo.Update(ctx, func(list []domain.DNS) ([]domain.DNS, error) {
		newList := []domain.DNS{}
		found := false
		for _, d := range list {
			if d.IP == ip {
				found = true
				continue
			}
			newList = append(newList, d)
		}

		if !found {
			return []domain.DNS{}, domain.ErrNotFound
		}

		return newList, nil
	})
	if err != nil {
		return domain.DNS{}, err
	}

	return domain.DNS{IP: ip}, nil
}

func (u *DNSUseCase) List(ctx context.Context) ([]domain.DNS, error) {
	u.log.Info("Listing DNS entries")
	return u.repo.Get(ctx)
}
