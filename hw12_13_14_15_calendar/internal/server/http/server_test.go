package internalhttp

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	app "hw12_13_14_15_calendar/internal/app"
	logger "hw12_13_14_15_calendar/internal/logger"
	middleware "hw12_13_14_15_calendar/internal/server/http/middleware"
	storage "hw12_13_14_15_calendar/internal/storage/gomemory"
)

func TestServerStopNotStarted(t *testing.T) {
	log := logger.NewLogger(logger.INFO, os.Stdout)
	middleware.Init(log)
	mainLogger := logger.NewLogger(logger.INFO, os.Stdout)
	calendarApp := app.NewApp(mainLogger, storage.NewStorage())
	httpServer := NewHTTPServer(
		"localhost",
		8080,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		log,
		calendarApp,
	)
	err := httpServer.Stop(context.Background())
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestServerStopNormally(t *testing.T) {
	log := logger.NewLogger(logger.INFO, os.Stdout)
	middleware.Init(log)
	mainLogger := logger.NewLogger(logger.INFO, os.Stdout)
	calendarApp := app.NewApp(mainLogger, storage.NewStorage())
	httpServer := NewHTTPServer(
		"localhost",
		8080,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		log,
		calendarApp,
	)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(3 * time.Second)
		err := httpServer.Stop(context.Background()) // STOP
		if err != nil {
			t.Errorf(err.Error())
		}
	}()
	if err := httpServer.Start(); err != nil { // START
		fmt.Println("http server goroutine: " + err.Error())
	}
	wg.Wait()
}

func TestServerStopBySignalNoWait(_ *testing.T) {
	ctx, ctxCancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT)
	defer ctxCancel()
	log := logger.NewLogger(logger.INFO, os.Stdout)
	middleware.Init(log)
	mainLogger := logger.NewLogger(logger.INFO, os.Stdout)
	calendarApp := app.NewApp(mainLogger, storage.NewStorage())
	httpServer := NewHTTPServer(
		"localhost",
		8080,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		log,
		calendarApp,
	)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Start(); err != nil { // START
			fmt.Println("http server Start goroutine: " + err.Error())
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Stop(ctx); err != nil { // START
			fmt.Println("http server Stop goroutine: " + err.Error())
		}
	}()
	pid, _, _ := syscall.Syscall(syscall.SYS_GETPID, 0, 0, 0)
	process, _ := os.FindProcess(int(pid))
	process.Signal(syscall.SIGHUP) // STOP
	wg.Wait()
}

func TestServerStopBySignalWithWait(_ *testing.T) {
	ctx, ctxCancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT)
	defer ctxCancel()
	log := logger.NewLogger(logger.INFO, os.Stdout)
	middleware.Init(log)
	mainLogger := logger.NewLogger(logger.INFO, os.Stdout)
	calendarApp := app.NewApp(mainLogger, storage.NewStorage())
	httpServer := NewHTTPServer(
		"localhost",
		8080,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		log,
		calendarApp,
	)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Start(); err != nil { // START
			fmt.Println("http server Start goroutine: " + err.Error())
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Stop(ctx); err != nil { // START
			fmt.Println("http server Stop goroutine: " + err.Error())
		}
	}()
	time.Sleep(3 * time.Second)
	pid, _, _ := syscall.Syscall(syscall.SYS_GETPID, 0, 0, 0)
	process, _ := os.FindProcess(int(pid))
	process.Signal(syscall.SIGINT) // STOP
	wg.Wait()
}

func TestServerStopByCancel(_ *testing.T) {
	ctx, ctxCancel := context.WithCancel(context.Background())
	log := logger.NewLogger(logger.INFO, os.Stdout)
	middleware.Init(log)
	mainLogger := logger.NewLogger(logger.INFO, os.Stdout)
	calendarApp := app.NewApp(mainLogger, storage.NewStorage())
	httpServer := NewHTTPServer(
		"localhost",
		8080,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		log,
		calendarApp,
	)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Start(); err != nil { // START
			fmt.Println("http server Start goroutine: " + err.Error())
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Stop(ctx); err != nil { // START
			fmt.Println("http server Stop goroutine: " + err.Error())
		}
	}()
	time.Sleep(3 * time.Second)
	ctxCancel() // STOP
	wg.Wait()
}

func TestServerCode(t *testing.T) {
	ctx, ctxCancel := context.WithCancel(context.Background())
	host := "localhost"
	var port uint16 = 8080
	// Server
	httpOutput := &bytes.Buffer{}
	log := logger.NewLogger(logger.INFO, httpOutput)
	middleware.Init(log)
	mainLogger := logger.NewLogger(logger.INFO, os.Stdout)
	calendarApp := app.NewApp(mainLogger, storage.NewStorage())
	httpServer := NewHTTPServer(
		"localhost",
		8080,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		mainLogger,
		calendarApp,
	)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Start(); err != nil { // START
			fmt.Println("http server Start goroutine: " + err.Error())
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Stop(ctx); err != nil { // START
			fmt.Println("http server Stop goroutine: " + err.Error())
		}
	}()
	// Wait
	time.Sleep(3 * time.Second)
	// Client
	url := fmt.Sprintf("http://%s:%d/", host, port)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, strings.NewReader(``))
	if err != nil {
		t.Error(err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		t.Error(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 418 {
		t.Errorf("StatusCode must be '418', but get '%d'\n", resp.StatusCode)
	} else {
		fmt.Printf("OK. StatusCode '418'\n")
	}
	//
	ctxCancel() // STOP
	wg.Wait()
}
