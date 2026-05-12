package grpcclient

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/IliaSotnikov2005/dnsmgr/client/internal/domain"
	"github.com/IliaSotnikov2005/dnsmgr/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	client proto.DNSServiceClient
}

func NewClient(client proto.DNSServiceClient) *Client {
	return &Client{client: client}
}

func (c *Client) Add(ctx context.Context, ip string) (domain.DNS, error) {
	resProto, err := c.client.AddDNS(ctx, &proto.DNSRequest{Ip: ip})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return domain.DNS{}, err
		}

		switch st.Code() {
		case codes.AlreadyExists:
			return domain.DNS{}, domain.ErrAlreadyExists
		case codes.InvalidArgument:
			return domain.DNS{}, domain.ErrInvalidIP
		default:
			return domain.DNS{}, err
		}
	}

	res := domain.DNS{Ip: resProto.Ip}
	return res, nil
}

func (c *Client) Remove(ctx context.Context, ip string) (domain.DNS, error) {
	resProto, err := c.client.RemoveDNS(ctx, &proto.DNSRequest{Ip: ip})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return domain.DNS{}, err
		}

		switch st.Code() {
		case codes.NotFound:
			return domain.DNS{}, domain.ErrNotFound
		case codes.InvalidArgument:
			return domain.DNS{}, domain.ErrInvalidIP
		default:
			return domain.DNS{}, err
		}
	}

	return domain.DNS{Ip: resProto.Ip}, nil
}

func (c *Client) List(ctx context.Context) ([]string, error) {
	res, err := c.client.ListDNS(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	dnsList := make([]string, 0, len(res.Ips))
	dnsList = append(dnsList, res.Ips...)

	return dnsList, nil
}
