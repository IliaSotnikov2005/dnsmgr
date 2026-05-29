package usecase

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/IliaSotnikov2005/dnsmgr/server/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type DNSRepositoryMock struct {
	mock.Mock
}

func (m *DNSRepositoryMock) Get(ctx context.Context) ([]domain.DNS, error) {
	args := m.Called(ctx)

	result := args.Get(0)

	if result == nil {
		return nil, args.Error(1)
	}

	return result.([]domain.DNS), args.Error(1)
}

func (m *DNSRepositoryMock) Save(ctx context.Context, dns []domain.DNS) error {
	args := m.Called(ctx, dns)
	return args.Error(0)
}

func (m *DNSRepositoryMock) Update(ctx context.Context, fn func([]domain.DNS) ([]domain.DNS, error)) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func TestDNSUseCase_Add(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	type testCase struct {
		name         string
		inputIP      string
		mockBehavior func(m *DNSRepositoryMock)
		wantErr      error
	}

	tests := []testCase{
		{
			name:    "Success: Valid IP added",
			inputIP: "1.1.1.1",
			mockBehavior: func(m *DNSRepositoryMock) {
				m.On("Update", mock.Anything, mock.MatchedBy(func(fn func([]domain.DNS) ([]domain.DNS, error)) bool {
					result, err := fn([]domain.DNS{{IP: "8.8.8.8"}})
					return err == nil && len(result) == 2 && result[1].IP == "1.1.1.1"
				})).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:    "Fail: Invalid IP format",
			inputIP: "999.999.999.999",
			mockBehavior: func(m *DNSRepositoryMock) {
			},
			wantErr: domain.ErrInvalidIP,
		},
		{
			name:    "Fail: IP already exists",
			inputIP: "8.8.8.8",
			mockBehavior: func(m *DNSRepositoryMock) {
				m.On("Update", mock.Anything, mock.MatchedBy(func(fn func([]domain.DNS) ([]domain.DNS, error)) bool {
					_, err := fn([]domain.DNS{{IP: "8.8.8.8"}})
					return errors.Is(err, domain.ErrAlreadyExists)
				})).Return(domain.ErrAlreadyExists)
			},
			wantErr: domain.ErrAlreadyExists,
		},
		{
			name:    "Fail: Repository Update error",
			inputIP: "1.2.3.4",
			mockBehavior: func(m *DNSRepositoryMock) {
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("disk failure"))
			},
			wantErr: errors.New("disk failure"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(DNSRepositoryMock)
			tt.mockBehavior(repo)
			uc := NewDNSUseCase(logger, repo)

			_, err := uc.Add(context.Background(), tt.inputIP)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestDNSUseCase_Remove(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name         string
		ipToRemove   string
		mockBehavior func(m *DNSRepositoryMock)
		wantErr      error
	}{
		{
			name:       "Success: IP removed",
			ipToRemove: "8.8.8.8",
			mockBehavior: func(m *DNSRepositoryMock) {
				m.On("Update", mock.Anything, mock.MatchedBy(func(fn func([]domain.DNS) ([]domain.DNS, error)) bool {
					result, err := fn([]domain.DNS{{IP: "8.8.8.8"}, {IP: "1.1.1.1"}})
					return err == nil && len(result) == 1 && result[0].IP == "1.1.1.1"
				})).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:       "Fail: IP not found",
			ipToRemove: "1.2.3.4",
			mockBehavior: func(m *DNSRepositoryMock) {
				m.On("Update", mock.Anything, mock.MatchedBy(func(fn func([]domain.DNS) ([]domain.DNS, error)) bool {
					_, err := fn([]domain.DNS{{IP: "8.8.8.8"}})
					return errors.Is(err, domain.ErrNotFound)
				})).Return(domain.ErrNotFound)
			},
			wantErr: domain.ErrNotFound,
		},
		{
			name:       "Fail: Repository Update error",
			ipToRemove: "8.8.8.8",
			mockBehavior: func(m *DNSRepositoryMock) {
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("write error"))
			},
			wantErr: errors.New("write error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(DNSRepositoryMock)
			tt.mockBehavior(repo)
			uc := NewDNSUseCase(logger, repo)

			_, err := uc.Remove(context.Background(), tt.ipToRemove)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestDNSUseCase_List(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name         string
		mockBehavior func(m *DNSRepositoryMock)
		wantLen      int
		wantErr      error
	}{
		{
			name: "Success: Multiple items",
			mockBehavior: func(m *DNSRepositoryMock) {
				m.On("Get", mock.Anything).Return([]domain.DNS{{IP: "1.1.1.1"}, {IP: "8.8.8.8"}}, nil)
			},
			wantLen: 2,
			wantErr: nil,
		},
		{
			name: "Success: Empty list",
			mockBehavior: func(m *DNSRepositoryMock) {
				m.On("Get", mock.Anything).Return([]domain.DNS{}, nil)
			},
			wantLen: 0,
			wantErr: nil,
		},
		{
			name: "Fail: Repository error",
			mockBehavior: func(m *DNSRepositoryMock) {
				m.On("Get", mock.Anything).Return([]domain.DNS(nil), errors.New("read error"))
			},
			wantLen: 0,
			wantErr: errors.New("read error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(DNSRepositoryMock)
			tt.mockBehavior(repo)
			uc := NewDNSUseCase(logger, repo)

			res, err := uc.List(context.Background())

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Len(t, res, tt.wantLen)
			}
			repo.AssertExpectations(t)
		})
	}
}
