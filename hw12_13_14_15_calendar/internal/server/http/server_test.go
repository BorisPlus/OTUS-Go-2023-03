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

	logger "hw12_13_14_15_calendar/internal/logger"
	middleware "hw12_13_14_15_calendar/internal/server/http/middleware"
)

func TestServerStopNotRun(t *testing.T) {
	log := logger.NewLogger(logger.INFO, os.Stdout)
	middleware.Init(log)
	httpServer := NewHTTPServer(
		"localhost",
		8080,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		log,
		nil)
	err := httpServer.Stop()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestServerStartStopNormally(t *testing.T) {
	ctx := context.Background()
	log := logger.NewLogger(logger.INFO, os.Stdout)
	middleware.Init(log)
	httpServer := NewHTTPServer(
		"localhost",
		8080,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		log,
		nil)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Start(ctx); err != nil {
			log.Error("http server goroutine: " + err.Error())
		}
	}()
	time.Sleep(1 * time.Second)
	err := httpServer.Stop()
	if err != nil {
		t.Errorf(err.Error())
	}
	wg.Wait()
}

func WhyTestServerStartStopByContextIsNotWork(t *testing.T) {
	ctx, ctxCancel := context.WithCancel(context.Background())
	log := logger.NewLogger(logger.INFO, os.Stdout)
	middleware.Init(log)
	httpServer := NewHTTPServer(
		"localhost",
		8080,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		log,
		nil)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Start(ctx); err != nil {
			log.Error("http server goroutine: " + err.Error())
		}
	}()
	ctxCancel() // Graceful Shutdown не срабатывает внутри httpServer.Start(&ctx)?
	wg.Wait()
}

func TestServerStartStopByContext(t *testing.T) {
	ctx, ctxCancel := signal.NotifyContext(context.Background(), syscall.SIGHUP)
	defer ctxCancel()
	log := logger.NewLogger(logger.INFO, os.Stdout)
	middleware.Init(log)
	httpServer := NewHTTPServer(
		"localhost",
		8080,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		log,
		nil)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Start(ctx); err != nil {
			log.Error("http server goroutine: " + err.Error())
		}
	}()
	pid, _, _ := syscall.Syscall(syscall.SYS_GETPID, 0, 0, 0)
	process, _ := os.FindProcess(int(pid))
	process.Signal(syscall.SIGHUP)
	wg.Wait()
}

func TestServerCode(t *testing.T) {
	host := "localhost"
	var port uint16 = 8080
	// Server
	httpOutput := &bytes.Buffer{}
	log := logger.NewLogger(logger.INFO, httpOutput)
	middleware.Init(log)
	httpServer := NewHTTPServer(
		host,
		port,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		log,
		nil)
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Start(ctx); err != nil {
			log.Error("http server goroutine: " + err.Error())
		}
	}()
	// Client
	url := fmt.Sprintf("http://%s:%d", host, port)
	request, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	httpResponse := httpOutput.String()
	if !strings.Contains(httpResponse, "StatusCode:418") {
		t.Errorf("Server must contain 'StatusCode:418', but get %s\n", httpResponse)
	} else {
		fmt.Printf("OK. Middleware catch status code '418':\n%s\n", httpResponse)
	}
	//
	ctxCancel()
	httpServer.Stop()
	wg.Wait()
}
