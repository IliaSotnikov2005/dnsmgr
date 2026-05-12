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
				m.On("Get", mock.Anything).Return([]domain.DNS{{IP: "8.8.8.8"}}, nil)
				m.On("Save", mock.Anything, mock.MatchedBy(func(list []domain.DNS) bool {
					return len(list) == 2 && list[1].IP == "1.1.1.1"
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
				m.On("Get", mock.Anything).Return([]domain.DNS{{IP: "8.8.8.8"}}, nil)
			},
			wantErr: domain.ErrAlreadyExists,
		},
		{
			name:    "Fail: Repository Get error",
			inputIP: "1.2.3.4",
			mockBehavior: func(m *DNSRepositoryMock) {
				m.On("Get", mock.Anything).Return(nil, errors.New("disk failure"))
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
				existing := []domain.DNS{{IP: "8.8.8.8"}, {IP: "1.1.1.1"}}
				m.On("Get", mock.Anything).Return(existing, nil)
				m.On("Save", mock.Anything, []domain.DNS{{IP: "1.1.1.1"}}).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:       "Fail: IP not found",
			ipToRemove: "1.2.3.4",
			mockBehavior: func(m *DNSRepositoryMock) {
				existing := []domain.DNS{{IP: "8.8.8.8"}}
				m.On("Get", mock.Anything).Return(existing, nil)
			},
			wantErr: domain.ErrNotFound,
		},
		{
			name:       "Fail: Repository Save error",
			ipToRemove: "8.8.8.8",
			mockBehavior: func(m *DNSRepositoryMock) {
				m.On("Get", mock.Anything).Return([]domain.DNS{{IP: "8.8.8.8"}}, nil)
				m.On("Save", mock.Anything, []domain.DNS{}).Return(errors.New("write error"))
			},
			wantErr: errors.New("write error"),
		},
		{
			name:       "Fail Repository Get error",
			ipToRemove: "1.1.1.1",
			mockBehavior: func(m *DNSRepositoryMock) {
				m.On("Get", mock.Anything).Return(nil, errors.New("read error"))
			},
			wantErr: errors.New("read error"),
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
