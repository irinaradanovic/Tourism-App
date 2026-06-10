// Minimal gRPC client for gateway -> tours TourCommandService calls.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	TourCommandService_CreateTour_FullMethodName  = "/tours.TourCommandService/CreateTour"
	TourCommandService_PublishTour_FullMethodName = "/tours.TourCommandService/PublishTour"
)

type TourCommandServiceClient interface {
	CreateTour(ctx context.Context, in *CreateTourRequest, opts ...grpc.CallOption) (*CreateTourResponse, error)
	PublishTour(ctx context.Context, in *PublishTourRequest, opts ...grpc.CallOption) (*TourCommandResponse, error)
}

type tourCommandServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTourCommandServiceClient(cc grpc.ClientConnInterface) TourCommandServiceClient {
	return &tourCommandServiceClient{cc}
}

func (c *tourCommandServiceClient) CreateTour(ctx context.Context, in *CreateTourRequest, opts ...grpc.CallOption) (*CreateTourResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateTourResponse)
	err := c.cc.Invoke(ctx, TourCommandService_CreateTour_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tourCommandServiceClient) PublishTour(ctx context.Context, in *PublishTourRequest, opts ...grpc.CallOption) (*TourCommandResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TourCommandResponse)
	err := c.cc.Invoke(ctx, TourCommandService_PublishTour_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
