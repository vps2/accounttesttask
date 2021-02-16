package client

import (
	"accounts/internal/api"
	"context"

	"google.golang.org/grpc"
)

type StatisticsServiceClient struct {
	addr string
}

func NewStatisticsServiceClient(addr string) *StatisticsServiceClient {
	return &StatisticsServiceClient{
		addr: addr,
	}
}

func (c *StatisticsServiceClient) ResetStatistics(ctx context.Context) error {
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := api.NewStatisticsServiceClient(conn)

	_, err = client.Reset(ctx, &api.Empty{})
	if err != nil {
		return err
	}

	return nil
}
