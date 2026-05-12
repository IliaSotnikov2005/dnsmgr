package controller

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/IliaSotnikov2005/dnsmgr/proto"
	"github.com/IliaSotnikov2005/dnsmgr/server/internal/domain"
	"github.com/IliaSotnikov2005/dnsmgr/server/internal/usecase"
)

type DNSHandler struct {
	proto.UnimplementedDNSServiceServer
	log     *slog.Logger
	useCase *usecase.DNSUseCase
}

func NewDNSHandler(log *slog.Logger, useCase *usecase.DNSUseCase) *DNSHandler {
	return &DNSHandler{
		log:     log,
		useCase: useCase,
	}
}

func (h *DNSHandler) AddDNS(ctx context.Context, req *proto.DNSRequest) (*proto.DNSResponse, error) {
	h.log.Info("Received AddDNS request", "ip", req.GetIp())
	res, err := h.useCase.Add(ctx, req.GetIp())
	if err != nil {
		h.log.Error("failed to add dns", "err", err)

		if errors.Is(err, domain.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		if errors.Is(err, domain.ErrInvalidIP) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	h.log.Info("DNS added successfully", "ip", req.GetIp())
	return &proto.DNSResponse{
		Ip:      res.IP,
		Message: "ok",
	}, nil
}

func (h *DNSHandler) RemoveDNS(ctx context.Context, req *proto.DNSRequest) (*proto.DNSResponse, error) {
	h.log.Info("Received RemoveDNS request", "ip", req.GetIp())
	res, err := h.useCase.Remove(ctx, req.GetIp())
	if err != nil {
		h.log.Error("failed to remove dns", "err", err)
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, domain.ErrInvalidIP) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &proto.DNSResponse{
		Ip:      res.IP,
		Message: "ok",
	}, nil
}

func (h *DNSHandler) ListDNS(ctx context.Context, _ *emptypb.Empty) (*proto.DNSListResponse, error) {
	h.log.Info("Received ListDNS request")
	res, err := h.useCase.List(ctx)
	if err != nil {
		h.log.Error("failed to list dns", "err", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	var dnsList = make([]string, 0, len(res))
	for _, dns := range res {
		dnsList = append(dnsList, dns.IP)
	}

	return &proto.DNSListResponse{
		Ips: dnsList,
	}, nil
}
