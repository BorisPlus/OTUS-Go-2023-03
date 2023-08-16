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
	httpServer "hw12_13_14_15_calendar/internal/server/http"
	middleware "hw12_13_14_15_calendar/internal/server/http/middleware"
	rpcServer "hw12_13_14_15_calendar/internal/server/rpc/server"
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
	mainConfig := config.NewCalendarConfig()
	err = viper.Unmarshal(mainConfig)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	log.Printf("HTTP Config - %+v\n", mainConfig.HTTP)
	mainLogger := logger.NewLogger(mainConfig.Log.Level, os.Stdout)
	middleware.Init(mainLogger)
	storage := storage.NewStorageByType(mainConfig.Storage.Type, mainConfig.Storage.DSN)
	calendar := app.NewApp(mainLogger, storage)
	httpServer := httpServer.NewHTTPServer(
		mainConfig.HTTP.Host,
		mainConfig.HTTP.Port,
		mainConfig.HTTP.ReadTimeout,
		mainConfig.HTTP.ReadHeaderTimeout,
		mainConfig.HTTP.WriteTimeout,
		mainConfig.HTTP.MaxHeaderBytes,
		mainLogger,
		calendar,
	)
	rpcServer := rpcServer.NewRPCServer(calendar, mainLogger)

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
		if err := httpServer.Stop(ctx); err != nil {
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
		if err := httpServer.Start(); err != nil {
			mainLogger.Error("failed to start HTTP server: " + err.Error())
			once.Do(stop)
		}
	}()

	// rpcServer.Start
	log.Printf("RPC Config - %+v\n", mainConfig.RPC)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := rpcServer.Start(fmt.Sprintf("%s:%d", mainConfig.RPC.Host, mainConfig.RPC.Port)); err != nil {
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
