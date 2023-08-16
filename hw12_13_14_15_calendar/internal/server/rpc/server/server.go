package server

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
	calendarrpcapi "hw12_13_14_15_calendar/internal/server/rpc/api"
)

func Event2PBEvent(event *models.Event) *calendarrpcapi.Event {
	pbEvent := new(calendarrpcapi.Event)
	pbEvent.PK = int32(event.PK)
	pbEvent.Title = event.Title
	pbEvent.Description = event.Description
	pbEvent.StartAt = timestamppb.New(event.StartAt)
	pbEvent.Duration = int32(event.Duration)
	pbEvent.Owner = event.Owner
	pbEvent.NotifyEarly = int32(event.NotifyEarly)
	pbEvent.Sheduled = event.Sheduled
	return pbEvent
}

func PBEvent2Event(pbEvent *calendarrpcapi.Event) *models.Event {
	event := new(models.Event)
	event.PK = int(pbEvent.PK)
	event.Title = pbEvent.Title
	event.Description = pbEvent.Description
	event.StartAt = pbEvent.StartAt.AsTime()
	event.Duration = int(pbEvent.Duration)
	event.Owner = pbEvent.Owner
	event.NotifyEarly = int(pbEvent.NotifyEarly)
	event.Sheduled = pbEvent.Sheduled
	return event
}

type RPCServer struct {
	calendarrpcapi.UnimplementedApplicationServer
	logger interfaces.Logger
	app    interfaces.Applicationer
	server *grpc.Server
}

func NewRPCServer(app interfaces.Applicationer, logger interfaces.Logger) *RPCServer {
	self := &RPCServer{}
	self.app = app
	self.logger = logger
	return self
}

func (s *RPCServer) CreateEvent(ctx context.Context, pbEvent *calendarrpcapi.Event) (*calendarrpcapi.Event, error) {
	_ = ctx // TODO: pass to s.app
	event := PBEvent2Event(pbEvent)
	createdEvent, err := s.app.CreateEvent(event)
	if err != nil {
		return nil, err
	}
	return Event2PBEvent(createdEvent), nil
}

func (s *RPCServer) ReadEvent(ctx context.Context, ident *calendarrpcapi.Id) (*calendarrpcapi.Event, error) {
	_ = ctx // TODO: pass to s.app
	event, err := s.app.ReadEvent(int(ident.Pk))
	if err != nil {
		return nil, err
	}
	return Event2PBEvent(event), nil
}

func (s *RPCServer) UpdateEvent(ctx context.Context, pbEvent *calendarrpcapi.Event) (*calendarrpcapi.Event, error) {
	_ = ctx // TODO: pass to s.app
	event := PBEvent2Event(pbEvent)
	updatedEvent, err := s.app.UpdateEvent(event)
	if err != nil {
		return nil, err
	}
	return Event2PBEvent(updatedEvent), nil
}

func (s *RPCServer) DeleteEvent(ctx context.Context, pbEvent *calendarrpcapi.Event) (*calendarrpcapi.Event, error) {
	_ = ctx // TODO: pass to s.app
	event := PBEvent2Event(pbEvent)
	deletedEvent, err := s.app.DeleteEvent(event)
	if err != nil {
		return nil, err
	}
	return Event2PBEvent(deletedEvent), nil
}

func (s *RPCServer) ListEvents(_ *emptypb.Empty, stream calendarrpcapi.Application_ListEventsServer) error {
	events, err := s.app.ListEvents()
	if err != nil {
		return err
	}
	for _, event := range events {
		tmp := event
		pbEvent := Event2PBEvent(&tmp)
		if err := stream.Send(pbEvent); err != nil {
			return err
		}
	}
	return nil
}

func (s *RPCServer) ListNotSheduledEvents(_ *emptypb.Empty,
	stream calendarrpcapi.Application_ListNotSheduledEventsServer,
) error {
	events, err := s.app.ListNotSheduledEvents()
	if err != nil {
		return err
	}
	for _, event := range events {
		tmp := event
		pbEvent := Event2PBEvent(&tmp)
		if err := stream.Send(pbEvent); err != nil {
			return err
		}
	}
	return nil
}

// type UnaryInterceptorType func(
// ctx context.Context,
// req interface{},
// info *grpc.UnaryServerInfo,
// handler grpc.UnaryHandler) (interface{}, error)

func LoggedUnaryInterceptor(logger interfaces.Logger) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Info("UnaryInterceptor: %q <-- OBJECT{%s}", info.FullMethod, req)
		return handler(ctx, req)
	}
}

// type StreamInterceptorType func(
// 	srv interface{},
// 	stream grpc.ServerStream,
// 	info *grpc.StreamServerInfo,
// 	handler grpc.StreamHandler,
// ) error

func LoggedStreamInterceptor(logger interfaces.Logger) func(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		logger.Info("StreamInterceptor: %q %s", info.FullMethod)
		return handler(srv, stream)
	}
}

func (s *RPCServer) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(LoggedUnaryInterceptor(s.logger)),
		grpc.StreamInterceptor(LoggedStreamInterceptor(s.logger)),
	)
	calendarrpcapi.RegisterApplicationServer(gRPCServer, s)
	s.server = gRPCServer
	s.logger.Info("GRPCServer.Start()")
	return s.server.Serve(lis)
}

func (s *RPCServer) Stop() {
	s.logger.Info("GRPCServer.Stop()")
	s.server.Stop()
}

func (s *RPCServer) GracefulStop() {
	s.logger.Info("GRPCServer.GracefulStop()")
	s.server.GracefulStop()
}
