package usecase

import (
	"context"

	"github.com/IliaSotnikov2005/dnsmgr/client/internal/domain"
)

type DNSService interface {
	Add(ctx context.Context, ip string) (domain.DNS, error)
	Remove(ctx context.Context, ip string) (domain.DNS, error)
	List(ctx context.Context) ([]string, error)
}
