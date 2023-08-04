package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	app "hw12_13_14_15_calendar/internal/app"
	logger "hw12_13_14_15_calendar/internal/logger"
	models "hw12_13_14_15_calendar/internal/models"
	middleware "hw12_13_14_15_calendar/internal/server/http/middleware"
	storage "hw12_13_14_15_calendar/internal/storage/gomemory"
)

// After fix code https://github.com/sonatard/noctx

type APIResponseTest struct {
	APIMethod string
	Error     string
	Data      struct{ Item models.Event }
}

var (
	apiResponse APIResponseTest
	host        = "localhost"
)

func TestServerAPICreatePKSequence(t *testing.T) {
	var response *http.Response
	var err error
	mainLogger := logger.NewLogger(logger.INFO, os.Stdout)
	inmemoryDatabase := storage.NewStorage()
	calendarApp := app.NewApp(mainLogger, inmemoryDatabase)
	middleware.Init(mainLogger)
	var port uint16 = 8083
	httpServer := NewHTTPServer(
		host,
		port,
		10,
		10,
		10,
		1<<20,
		mainLogger,
		calendarApp,
	)
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeoutCancel()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer.Start(timeoutCtx)
	}()
	client := &http.Client{}
	payloadOfCreateRaw := `{
						"title": "title",
						"startat": "2023-08-05T21:54:42+02:00",
						"duration": 0,
						"description": "description",
						"owner": "owner",
						"notifyearly": 0
						}`
	requestOfCreate := fmt.Sprintf("http://%s:%d/api/events/create", host, port)
	// CREATE 1
	payloadOfCreate := strings.NewReader(payloadOfCreateRaw)
	request, err := http.NewRequestWithContext(timeoutCtx, "POST", requestOfCreate, payloadOfCreate)
	if err != nil {
		t.Errorf("FAIL: error prepare http request: %s\n", requestOfCreate)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	if err != nil {
		t.Errorf("FAIL: error decode event http request: %s\n", err)
		return
	}
	response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		t.Errorf("FAIL: error decode event http request: %s\n", err)
		return
	}
	if apiResponse.Data.Item.PK != 1 {
		t.Errorf("FAIL: get event PK %d, expected 1\n", apiResponse.Data.Item.PK)
	} else {
		fmt.Printf("OK: get event PK %d\n", apiResponse.Data.Item.PK)
	}
	// CREATE 2
	payloadOfCreate = strings.NewReader(payloadOfCreateRaw)
	request, err = http.NewRequestWithContext(timeoutCtx, "POST", requestOfCreate, payloadOfCreate)
	if err != nil {
		t.Errorf("FAIL: error prepare http request: %s\n", requestOfCreate)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	if err != nil {
		t.Errorf("FAIL: error decode event http request: %s\n", err)
		return
	}
	response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		t.Errorf("FAIL: error decode event http request: %s\n", err)
		return
	}
	if apiResponse.Data.Item.PK != 2 {
		t.Errorf("FAIL: get event PK %d, expected 1\n", apiResponse.Data.Item.PK)
	} else {
		fmt.Printf("OK: get event PK %d\n", apiResponse.Data.Item.PK)
	}
	// DELETE 3
	requestOfDelete := fmt.Sprintf("http://%s:%d/api/events/1/delete", host, port)
	payloadOfDelete := strings.NewReader(``)
	request, err = http.NewRequestWithContext(timeoutCtx, "DELETE", requestOfDelete, payloadOfDelete)
	if err != nil {
		t.Errorf("FAIL: error prepare http request: %s\n", requestOfDelete)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	if err != nil {
		t.Errorf("FAIL: error decode event http request: %s\n", err)
		return
	}
	response.Body.Close()
	// CREATE 3
	payloadOfCreate = strings.NewReader(payloadOfCreateRaw)
	request, err = http.NewRequestWithContext(timeoutCtx, "POST", requestOfCreate, payloadOfCreate)
	if err != nil {
		t.Errorf("FAIL: error prepare http request: %s\n", requestOfCreate)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	response, err = client.Do(request)
	if err != nil {
		t.Errorf("FAIL: error decode event http request: %s\n", err)
		return
	}
	response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		t.Errorf("FAIL: error decode event http request: %s\n", err)
		return
	}
	if apiResponse.Data.Item.PK != 3 {
		t.Errorf("FAIL: get event PK %d, expected 3\n", apiResponse.Data.Item.PK)
	} else {
		fmt.Printf("OK: get event PK %d\n", apiResponse.Data.Item.PK)
	}
	//
	<-timeoutCtx.Done()
	httpServer.Stop(timeoutCtx)
	wg.Wait()
}

func TestServerAPIVersion(t *testing.T) {
	mainLogger := logger.NewLogger(logger.INFO, os.Stdout)
	inmemoryDatabase := storage.NewStorage()
	calendarApp := app.NewApp(mainLogger, inmemoryDatabase)
	middleware.Init(mainLogger)
	client := &http.Client{}
	var port uint16 = 8082
	httpServer := NewHTTPServer(
		host,
		port,
		10,
		10,
		10,
		1<<20,
		mainLogger,
		calendarApp,
	)

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer timeoutCancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer.Start(timeoutCtx)
	}()

	requestOfVersion := fmt.Sprintf("http://%s:%d/api/version", host, port)
	payload := strings.NewReader(``)
	request, err := http.NewRequestWithContext(timeoutCtx, "GET", requestOfVersion, payload)
	if err != nil {
		t.Errorf("FAIL: error prepare http request: %s\n", requestOfVersion)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		t.Errorf("FAIL: error decode event http request: %s\n", err)
		return
	}
	response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("FAIL: error making http request: %s\n", err)
		return
	}
	ethalon := `{"method":"api.version","error":"","data":{"Version":"1.0.0"}}`
	if string(body) != ethalon {
		t.Errorf("FAIL: get %s\n", body)
		t.Errorf("FAIL: expected %s\n", ethalon)
	} else {
		fmt.Printf("OK: %s\n", body)
	}

	<-timeoutCtx.Done()
	httpServer.Stop(timeoutCtx)
	wg.Wait()
}
