package integration

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	app "hw12_13_14_15_calendar/internal/app"
	logger "hw12_13_14_15_calendar/internal/logger"
	pb "hw12_13_14_15_calendar/internal/protobuf/api"
	client "hw12_13_14_15_calendar/internal/protobuf/client"
	server "hw12_13_14_15_calendar/internal/protobuf/server"
	storage "hw12_13_14_15_calendar/internal/storage"
)

func TestIntegration(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var once sync.Once
	defer once.Do(cancel)

	mainLogger := logger.NewLogger(logger.INFO, os.Stdout)
	storage := storage.NewStorageByType(storage.GOMEMORY_STORAGE)
	calendar := app.NewApp(mainLogger, storage)
	grpcServer := server.NewRPCServer(calendar, mainLogger)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		mainLogger.Info("GracefulStop")
		grpcServer.GracefulStop()
	}()
	//
	wg.Add(1)
	go func() {
		defer wg.Done()
		mainLogger.Info("server listening at %v", "localhost:5000")
		grpcServer.Start(ctx, "localhost:5000")
	}()
	//
	time.Sleep(5 * time.Second)
	//
	// var dialOpts []grpc.DialOption
	grpcClient := client.Client{}
	mainLogger.Info("grpcClient.Client{}")
	// grpc.WithInsecure()
	grpcClient.Connect("localhost:5000")
	mainLogger.Info("localhost:5000")
	//
	pbEvent1 := pb.Event{}
	pbEvent1.Title = "Title 1"
	createdEvent1, err := grpcClient.CreateEvent(ctx, &pbEvent1)
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
	createdEvent2, err := grpcClient.CreateEvent(ctx, &pbEvent2)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	mainLogger.Info("createdEvent2.PK = %d", createdEvent2.PK)
	if createdEvent2.PK != 2 {
		t.Errorf("createdEvent2.PK = %d, expected %d", createdEvent2.PK, 2)
	}
	//
	deletedEvent2, err := grpcClient.DeleteEvent(ctx, createdEvent2)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	mainLogger.Info("deletedEvent2.PK = %d", deletedEvent2.PK)
	//
	pbEvent3 := pb.Event{}
	pbEvent3.Title = "Title 3"
	createdEvent3, err := grpcClient.CreateEvent(ctx, &pbEvent3)
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
	pbEvent3Copy, err := grpcClient.ReadEvent(ctx, &ident)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	mainLogger.Info("pbEvent3Copy.Title = %s", pbEvent3Copy.Title)
	if pbEvent3Copy.Title != pbEvent3.Title {
		t.Errorf("pbEvent3Copy.Title = %s, expected %s", pbEvent3Copy.Title, pbEvent3.Title)
	}
	//
	events, err := grpcClient.ListEvents(ctx)
	if err != nil {
		mainLogger.Error(err.Error())
	}
	mainLogger.Info("%+v", events)
	grpcClient.Close()
	//
	once.Do(cancel) // TODO: поискать по коду подобное
	wg.Wait()
}

func TestInterceptorLogging(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var once sync.Once
	defer once.Do(cancel)

	grpcOutput := &bytes.Buffer{}
	mainLogger := logger.NewLogger(logger.INFO, grpcOutput)

	storage := storage.NewStorageByType(storage.GOMEMORY_STORAGE)
	calendar := app.NewApp(mainLogger, storage)
	grpcServer := server.NewRPCServer(calendar, mainLogger)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		mainLogger.Info("GracefulStop")
		grpcServer.GracefulStop()
	}()
	//
	wg.Add(1)
	go func() {
		defer wg.Done()
		mainLogger.Info("server listening at %v", "localhost:5000")
		grpcServer.Start(ctx, "localhost:5000")
	}()
	//
	grpcClient := client.Client{}
	grpcClient.Connect("localhost:5000")
	//
	pbEvent1 := pb.Event{}
	pbEvent1.Title = "Title 1"
	_, err := grpcClient.CreateEvent(ctx, &pbEvent1)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	if !strings.Contains(grpcOutput.String(), "/calendar.Application/CreateEvent") {
		t.Error("Log not contain interceptor data")
	} else {
		fmt.Println("It's OK.")
	}
	grpcClient.Close()
	once.Do(cancel)
	wg.Wait()
}
