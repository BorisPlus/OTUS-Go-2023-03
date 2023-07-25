package internalhttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	logger "hw12_13_14_15_calendar/internal/logger"
)

// After fix code https://github.com/sonatard/noctx

func Send(body io.Reader) error {
	return SendWithContext(context.Background(), body)
}

func SendWithContext(ctx context.Context, body io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "localhost:8081", body)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func NotForGitHubTestServerCode(t *testing.T) {
	t.Helper()
	host := "localhost"
	var port uint16 = 8081
	outputInto := &bytes.Buffer{}
	mainLogger := logger.NewLogger(logger.INFO, outputInto)
	httpServer := NewServer(host, port, mainLogger, nil)
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Start(ctx); err != nil {
			mainLogger.Error("failed to start http server: " + err.Error())
		}
	}()

	url := fmt.Sprintf("http://%s:%d", host, port)
	client := &http.Client{}
	request, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	timeoutCtx, cancelByTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelByTimeout()
	if err := httpServer.Stop(timeoutCtx); err != nil {
		fmt.Println(err)
	}

	outputted := outputInto.String()
	if !strings.Contains(outputted, "StatusCode:418") {
		t.Errorf("Server must contain 'StatusCode:418', but get %s\n", outputted)
	} else {
		fmt.Printf("OK. Middleware catch status code '418':\n%s\n", outputted)
	}
}
