package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	app "hw12_13_14_15_calendar/internal/app"
	config "hw12_13_14_15_calendar/internal/config"
	logger "hw12_13_14_15_calendar/internal/logger"
	rpc "hw12_13_14_15_calendar/internal/protobuf/server"
	http "hw12_13_14_15_calendar/internal/server/http"
	middleware "hw12_13_14_15_calendar/internal/server/http/middleware"
	storage "hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "", "Path to configuration file")
}

func main() {
	pflag.Parse()
	if pflag.Arg(0) == "version" {
		printVersion()
		return
	}
	if configFile == "" {
		fmt.Println("Please set: '--config=<Path to configuration file>'")
		return
	}
	viper.SetConfigType("yaml")
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	viper.ReadConfig(file)
	mainConfig := config.NewConfig()
	err = viper.Unmarshal(mainConfig)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	log.Printf("%+v\n", mainConfig.HTTP)
	log.Println(mainConfig.Log.Level)
	mainLogger := logger.NewLogger(mainConfig.Log.Level, os.Stdout)
	middleware.Init(mainLogger)
	storage := storage.NewStorageByType(mainConfig.Storage.Type, mainConfig.Storage.DSN)
	calendar := app.NewApp(mainLogger, storage)
	httpServer := http.NewHTTPServer(
		mainConfig.HTTP.Host,
		mainConfig.HTTP.Port,
		mainConfig.HTTP.ReadTimeout,
		mainConfig.HTTP.ReadHeaderTimeout,
		mainConfig.HTTP.WriteTimeout,
		mainConfig.HTTP.MaxHeaderBytes,
		mainLogger,
		calendar,
	)
	rpcServer := rpc.NewRPCServer(calendar, mainLogger)

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGTSTP)
	defer stop()
	wg := sync.WaitGroup{}

	var once sync.Once
	// GRASEFULL: httpServer.Stop
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Stop(); err != nil {
			fmt.Println(err)
		}
	}()

	// GRASEFULL: rpcServer.GracefulStop
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		rpcServer.GracefulStop()
	}()

	// httpServer.Start
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Start(ctx); err != nil {
			mainLogger.Error("failed to start HTTP server: " + err.Error())
			once.Do(stop)
		}
	}()

	// rpcServer.Start
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := rpcServer.Start(ctx, fmt.Sprintf("%s:%d", mainConfig.RPC.Host, mainConfig.RPC.Port)); err != nil {
			mainLogger.Error("failed to start RPC server: " + err.Error())
			once.Do(stop)
		}
	}()

	log.Println("Println Calendar is running...")
	mainLogger.Info("calendar is running...")
	<-ctx.Done()
	// stop()
	mainLogger.Info("Complex Shutting down was done gracefully by signal.")
	wg.Wait()
}
