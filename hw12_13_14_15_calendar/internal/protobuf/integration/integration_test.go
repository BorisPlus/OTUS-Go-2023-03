package integration

import (
	"context"
	"log"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"

	app "hw12_13_14_15_calendar/internal/app"
	logger "hw12_13_14_15_calendar/internal/logger"
	pb "hw12_13_14_15_calendar/internal/protobuf/api"
	grpcClient "hw12_13_14_15_calendar/internal/protobuf/client"
	grpcServer "hw12_13_14_15_calendar/internal/protobuf/server"
	storage "hw12_13_14_15_calendar/internal/storage"
)

func TestLogger(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var once sync.Once
	defer once.Do(cancel)
	lis, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gRPCServer := grpc.NewServer()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		gRPCServer.Stop()
	}()
	mainLogger := logger.NewLogger(logger.INFO, os.Stdout)
	storage := storage.NewStorageByType(storage.GOMEMORY_STORAGE)
	calendar := app.NewApp(mainLogger, storage)
	pb.RegisterApplicationServer(gRPCServer, &grpcServer.RPCServer{App: calendar})
	mainLogger.Info("server listening at %v", lis.Addr())
	//
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := gRPCServer.Serve(lis); err != nil {
			mainLogger.Error("failed to serve: %v", err)
		}
	}()
	//
	time.Sleep(5 * time.Second)
	//
	// var dialOpts []grpc.DialOption
	client := grpcClient.Client{}
	mainLogger.Info("grpcClient.Client{}")
	// grpc.WithInsecure()
	client.Connect("localhost:5000")
	mainLogger.Info("localhost:5000")
	//
	pbEvent1 := pb.Event{}
	pbEvent1.Title = "Title 1"
	createdEvent1, err := client.CreateEvent(ctx, &pbEvent1)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	mainLogger.Info("createdEvent1.PK = %d", createdEvent1.PK)
	if createdEvent1.PK != 1 {
		t.Errorf("createdEvent1.PK = %d, expected %d", createdEvent1.PK, 2)
	}
	//
	pbEvent2 := pb.Event{}
	pbEvent2.Title = "Title 2"
	createdEvent2, err := client.CreateEvent(ctx, &pbEvent2)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	mainLogger.Info("createdEvent2.PK = %d", createdEvent2.PK)
	if createdEvent2.PK != 2 {
		t.Errorf("createdEvent2.PK = %d, expected %d", createdEvent2.PK, 2)
	}
	//
	deletedEvent2, err := client.DeleteEvent(ctx, createdEvent2)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	mainLogger.Info("deletedEvent2.PK = %d", deletedEvent2.PK)
	//
	pbEvent3 := pb.Event{}
	pbEvent3.Title = "Title 3"
	createdEvent3, err := client.CreateEvent(ctx, &pbEvent3)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	mainLogger.Info("createdEvent3.PK = %d", createdEvent3.PK)
	if createdEvent3.PK != 3 {
		t.Errorf("createdEvent3.PK = %d, expected %d", createdEvent3.PK, 3)
	}
	//
	ident := pb.Id{}
	ident.Pk = 3
	pbEvent3Copy, err := client.ReadEvent(ctx, &ident)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	mainLogger.Info("pbEvent3Copy.Title = %s", pbEvent3Copy.Title)
	if pbEvent3Copy.Title != pbEvent3.Title {
		t.Errorf("pbEvent3Copy.Title = %s, expected %s", pbEvent3Copy.Title, pbEvent3.Title)
	}
	//
	events, err := client.ListEvents(ctx)
	if err != nil {
		mainLogger.Error(err.Error())
	}
	mainLogger.Info("%+v", events)
	//
	cancel()
	wg.Wait()
}
