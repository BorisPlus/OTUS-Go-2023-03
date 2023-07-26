package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"hw12_13_14_15_calendar/internal/app"
	"hw12_13_14_15_calendar/internal/config"
	"hw12_13_14_15_calendar/internal/logger"
	internalhttp "hw12_13_14_15_calendar/internal/server/http"
	"hw12_13_14_15_calendar/internal/storage"
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

	mainLogger := logger.NewLogger(mainConfig.Log.Level, os.Stdout)
	storage := storage.NewStorageByType(mainConfig.Storage.Type, mainConfig.Storage.DSN)
	calendar := app.NewApp(mainLogger, storage)
	httpServer := internalhttp.NewServer(mainConfig.HTTP.Host, mainConfig.HTTP.Port, mainLogger, calendar)
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGTSTP)
	defer stop()
	go func() {
		if err := httpServer.Start(ctx); err != nil {
			mainLogger.Error("failed to start http server: " + err.Error())
			stop()
		}
	}()
	mainLogger.Info("calendar is running...")
	<-ctx.Done()
	stop()
	mainLogger.Info("Shutting down gracefully by signal.")
	timeoutCtx, cancelByTimeout := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelByTimeout()
	if err := httpServer.Stop(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}
