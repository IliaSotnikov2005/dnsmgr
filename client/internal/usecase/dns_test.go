package usecase

import (
	"context"
	"log/slog"
	"testing"

	"github.com/IliaSotnikov2005/dnsmgr/client/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type DNSServiceMock struct {
	mock.Mock
}

func (m *DNSServiceMock) Add(ctx context.Context, ip string) (domain.DNS, error) {
	args := m.Called(ctx, ip)
	return args.Get(0).(domain.DNS), args.Error(1)
}

func (m *DNSServiceMock) Remove(ctx context.Context, ip string) (domain.DNS, error) {
	args := m.Called(ctx, ip)
	return args.Get(0).(domain.DNS), args.Error(1)
}

func (m *DNSServiceMock) List(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	var res []string
	if args.Get(0) != nil {
		res = args.Get(0).([]string)
	}

	return res, args.Error(1)
}

func TestDNSUseCase_Add(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(DNSServiceMock)
		uc := NewDNSUseCase(slog.Default(), mockSvc)
		ip := "1.1.1.1"

		mockSvc.On("Add", ctx, ip).Return(domain.DNS{Ip: ip}, nil)

		res, err := uc.Add(ctx, ip)

		assert.NoError(t, err)
		assert.Equal(t, ip, res.Ip)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Validation_Error", func(t *testing.T) {
		mockSvc := new(DNSServiceMock)
		uc := NewDNSUseCase(slog.Default(), mockSvc)

		res, err := uc.Add(ctx, "invalid-ip-format")

		assert.Error(t, err)
		assert.Empty(t, res.Ip)
		mockSvc.AssertNotCalled(t, "Add", mock.Anything, mock.Anything)
	})
}

func TestDNSUseCase_Remove(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(DNSServiceMock)
		uc := NewDNSUseCase(slog.Default(), mockSvc)
		ip := "8.8.8.8"

		mockSvc.On("Remove", ctx, ip).Return(domain.DNS{Ip: ip}, nil)

		res, err := uc.Remove(ctx, ip)

		assert.NoError(t, err)
		assert.Equal(t, ip, res.Ip)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Not_Found", func(t *testing.T) {
		mockSvc := new(DNSServiceMock)
		uc := NewDNSUseCase(slog.Default(), mockSvc)
		ip := "1.2.3.4"

		mockSvc.On("Remove", ctx, ip).Return(domain.DNS{}, domain.ErrNotFound)

		_, err := uc.Remove(ctx, ip)

		assert.ErrorIs(t, err, domain.ErrNotFound)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Validation_Error", func(t *testing.T) {
		mockSvc := new(DNSServiceMock)
		uc := NewDNSUseCase(slog.Default(), mockSvc)

		_, err := uc.Remove(ctx, "  ")

		assert.Error(t, err)
		mockSvc.AssertNotCalled(t, "Remove", mock.Anything, mock.Anything)
	})
}

func TestDNSUseCase_List(t *testing.T) {
	mockSvc := new(DNSServiceMock)
	uc := NewDNSUseCase(slog.Default(), mockSvc)

	list := []string{"1.1.1.1", "8.8.8.8"}
	mockSvc.On("List", mock.Anything).Return(list, nil)

	res, err := uc.List(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, list, res)
}
