package server

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	pb "hw12_13_14_15_calendar/internal/protobuf/api"
	common "hw12_13_14_15_calendar/internal/protobuf/common"
)

type RPCServer struct {
	pb.UnimplementedApplicationServer
	server *grpc.Server
	App    interfaces.Applicationer
}

func (s *RPCServer) CreateEvent(ctx context.Context, pbEvent *pb.Event) (*pb.Event, error) {
	event := common.PBEvent2Event(pbEvent)
	createdEvent, err := s.App.CreateEvent(event)
	if err != nil {
		return nil, err
	}
	return common.Event2PBEvent(createdEvent), nil
}

func (s *RPCServer) ReadEvent(ctx context.Context, ident *pb.Id) (*pb.Event, error) {
	event, err := s.App.ReadEvent(int(ident.Pk))
	if err != nil {
		return nil, err
	}
	return common.Event2PBEvent(event), nil
}

func (s *RPCServer) UpdateEvent(ctx context.Context, pbEvent *pb.Event) (*pb.Event, error) {
	event := common.PBEvent2Event(pbEvent)
	updatedEvent, err := s.App.UpdateEvent(event)
	if err != nil {
		return nil, err
	}
	return common.Event2PBEvent(updatedEvent), nil
}

func (s *RPCServer) DeleteEvent(ctx context.Context, pbEvent *pb.Event) (*pb.Event, error) {
	event := common.PBEvent2Event(pbEvent)
	deletedEvent, err := s.App.DeleteEvent(event)
	if err != nil {
		return nil, err
	}
	return common.Event2PBEvent(deletedEvent), nil
}

func (s *RPCServer) ListEvents(_ *emptypb.Empty, stream pb.Application_ListEventsServer) error {
	events, err := s.App.ListEvents()
	if err != nil {
		return err
	}
	for _, event := range events {
		pbEvent := common.Event2PBEvent(&event)
		if err := stream.Send(pbEvent); err != nil {
			return err
		}
	}
	return nil
}

func (s *RPCServer) Start(ctx context.Context, address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	gRPCServer := grpc.NewServer()
	pb.RegisterApplicationServer(gRPCServer, s)
	s.server = gRPCServer
	if err := s.server.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *RPCServer) Stop(ctx context.Context) {
	s.server.Stop()
}
