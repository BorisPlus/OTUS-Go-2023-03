package internalhttp

import (
	// "bytes"
	"context"
	// "fmt"
	"io"
	"net/http"
	// "strings"
	// "sync"
	"testing"
	// "time"

	// logger "hw12_13_14_15_calendar/internal/logger"
	// "hw12_13_14_15_calendar/internal/server/http/middleware"
)

// After fix code https://github.com/sonatard/noctx

func SendData(body io.Reader) error {
	return SendWithContext(context.Background(), body)
}

func SendDataWithContext(ctx context.Context, body io.Reader) error {
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

func TestServerAPICode(t *testing.T) {
	// host := "localhost"
	// var port uint16 = 8081
	// outputInto := &bytes.Buffer{}
	// mainLogger := logger.NewLogger(logger.INFO, outputInto)
	// middleware.Init(mainLogger)
	// httpServer := NewHTTPServer(
	// 	host,
	// 	port,
	// 	10,
	// 	10,
	// 	10,
	// 	1<<20,
	// 	mainLogger,
	// 	nil)
	// ctx := context.Background()
	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	httpServer.Start(ctx)
	// 	// if err := httpServer.Start(ctx); err != nil {
	// 	// 	mainLogger.Error("http server goroutine: " + err.Error())
	// 	// }
	// }()

	// createBody := `{"title":"title","startat":"2021-02-18T21:54:42+02:00","duration":0,"description": "description","owner":"owner","notifyearly":0}"`
	// createBodyReader := strings.NewReader(createBody)
	// createUrl := fmt.Sprintf("http://%s:%d/api/events/create", host, port)
	// client := &http.Client{}

	// request, _ := http.NewRequestWithContext(ctx, http.MethodPost, createUrl, createBodyReader)
	// resp, err := client.Do(request)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(resp.Body)
	// resp.Body.Close()

	// timeoutCtx, cancelByTimeout := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancelByTimeout()
	// if err := httpServer.Stop(timeoutCtx); err != nil {
	// 	fmt.Println("httpServer.Stop", err)
	// }

	// // outputted := outputInto.String()
	// // if !strings.Contains(outputted, "StatusCode:418") {
	// // 	t.Errorf("Server must contain 'StatusCode:418', but get %s\n", outputted)
	// // } else {
	// // 	fmt.Printf("OK. Middleware catch status code '418':\n%s\n", outputted)
	// // }
	// wg.Wait()
}
