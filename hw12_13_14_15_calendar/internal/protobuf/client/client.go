package client

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "hw12_13_14_15_calendar/internal/protobuf/api"
)

type Client struct {
	grpcClient pb.ApplicationClient
	connection *grpc.ClientConn
}

var localOpts = []grpc.CallOption{}

func (c *Client) Connect(dsn string) error {
	connection, err := grpc.Dial(dsn, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed on client: %v", err)
		return err
	}
	c.connection = connection
	c.grpcClient = pb.NewApplicationClient(c.connection)
	return nil
}

func (c *Client) Close(dsn string) error {
	err := c.connection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) CreateEvent(ctx context.Context, event *pb.Event) (*pb.Event, error) {
	return c.grpcClient.CreateEvent(ctx, event, localOpts...)
}

func (c *Client) ReadEvent(ctx context.Context, event *pb.Id) (*pb.Event, error) {
	return c.grpcClient.ReadEvent(ctx, event, localOpts...)
}

func (c *Client) UpdateEvent(ctx context.Context, event *pb.Event) (*pb.Event, error) {
	return c.grpcClient.UpdateEvent(ctx, event, localOpts...)
}

func (c *Client) DeleteEvent(ctx context.Context, event *pb.Event, opts ...grpc.CallOption) (*pb.Event, error) {
	return c.grpcClient.DeleteEvent(ctx, event, localOpts...)
}

func (c *Client) ListEvents(ctx context.Context) ([]*pb.Event, error) {
	events, err := c.grpcClient.ListEvents(ctx, &emptypb.Empty{}, localOpts...)
	if err != nil {
		return nil, err
	}
	pbEvents := make([]*pb.Event, 0)
	for {
		pbEvent, err := events.Recv()
		// err := events.RecvMsg(&pbEvent)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		pbEvents = append(pbEvents, pbEvent)
	}
	return pbEvents, nil
}
