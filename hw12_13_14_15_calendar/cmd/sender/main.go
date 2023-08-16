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

	config "hw12_13_14_15_calendar/internal/config"
	logger "hw12_13_14_15_calendar/internal/logger"
	rmq "hw12_13_14_15_calendar/internal/rmq"
)

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "", "Path to configuration file")
}

func main() {
	pflag.Parse()
	if pflag.Arg(0) == "version" {
		fmt.Printf("2023.08.13 v.1")
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
	cfg := config.NewSenderConfig()
	err = viper.Unmarshal(cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	mainLogger := logger.NewLogger(cfg.Log.Level, os.Stdout)
	sendTo := rmq.NewNotifier(
		cfg.SendTo,
		mainLogger,
	)
	archiveTo := rmq.NewNotifier(
		cfg.ArchiveTo,
		mainLogger,
	)
	sender := NewSender(
		cfg.Source,
		sendTo,
		archiveTo,
		mainLogger,
	)
	var once sync.Once
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGTSTP)
	defer once.Do(stop)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := sender.Stop(); err != nil {
			fmt.Println(err)
		}
	}()
	if err := sender.Start(ctx); err != nil {
		once.Do(stop)
	}
	wg.Wait()
}
