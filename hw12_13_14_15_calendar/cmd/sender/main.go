package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"hw12_13_14_15_calendar/internal/backend/archiver"
	"hw12_13_14_15_calendar/internal/backend/sender"
	"hw12_13_14_15_calendar/internal/config"
	"hw12_13_14_15_calendar/internal/logger"
	"hw12_13_14_15_calendar/internal/models"
)

var configFile string

type SentByFmt struct{}

func (s *SentByFmt) Notify(notice models.Notice) error {
	_, err := fmt.Printf("Notice %q send to %q\n", notice.Title, notice.Owner)
	return err
}

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
	mainLogger := logger.NewLogger(cfg.Log.Level, io.Discard)
	sender := sender.NewSender(
		archiver.NewNoticesSource(cfg.Source, mainLogger),
		sender.NewNoticesTarget(cfg.Target, mainLogger, &SentByFmt{}),
		mainLogger,
		cfg.TimeoutSec,
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
		if err := sender.Stop(ctx); err != nil {
			fmt.Println(err)
		}
	}()
	if err := sender.Start(ctx); err != nil {
		once.Do(stop)
	}
	wg.Wait()
}
