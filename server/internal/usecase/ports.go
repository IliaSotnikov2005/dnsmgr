package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/dnsmgr/server/internal/domain"
)

type DNSRepository interface {
	Get(ctx context.Context) ([]domain.DNS, error)
	Save(ctx context.Context, dns []domain.DNS) error
	Update(ctx context.Context, fn func([]domain.DNS) ([]domain.DNS, error)) error
}
