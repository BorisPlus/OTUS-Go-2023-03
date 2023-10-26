package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq" // a blank import should be justifying.

	"golang.org/x/exp/slices"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/icrowley/fake"

	amqp "github.com/rabbitmq/amqp091-go"
	config "hw12_13_14_15_calendar/internal/config"
	logger "hw12_13_14_15_calendar/internal/logger"
	models "hw12_13_14_15_calendar/internal/models"
)

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "", "Path to configuration file")
}

func osExit(exitCode int, l logger.Logger) {
	l.Warning("ExitCode: %d", exitCode)
	os.Exit(exitCode)
}

func main() {
	pflag.Parse()
	if configFile == "" {
		fmt.Println("Please set: '--config=<Path to configuration file>'")
		os.Exit(1)
	}
	viper.SetConfigType("yaml")
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	viper.ReadConfig(file)
	cfg := config.NewCheckerConfig()
	err = viper.Unmarshal(cfg)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
		os.Exit(3)
	}
	logging := logger.NewLogger(cfg.Log.Level, os.Stdout)
	// CREATE DATASET SERIA
	client := &http.Client{}
	requestOfCreate := fmt.Sprintf("http://%s:%d/api/events/create", cfg.HTTP.Host, cfg.HTTP.Port)
	now := time.Now()
	titles := []string{}
	// cfgT BY NOTIFIED
	logging.Debug("SEND COUNT: %d", cfg.Counts.Send)
	titlesOfSend := []string{}
	for i := 0; i < cfg.Counts.Send; i++ {
		event := models.Event{
			Title:       fake.Title(),
			StartAt:     now.Add(3000 * time.Second),
			Duration:    1800,
			Description: fake.EmailSubject(),
			Owner:       fake.EmailAddress(),
			NotifyEarly: 3600,
		}
		logging.Debug("SEND: %s\n", event.Title)
		titlesOfSend = append(titlesOfSend, event.Title)
		titles = append(titles, event.Title)
		payload, err := json.Marshal(event)
		if err != nil {
			logging.Error(err.Error())
			osExit(10, *logging)
		}
		payloadOfCreate := strings.NewReader(string(payload))
		request, err := http.NewRequestWithContext(context.Background(), "POST", requestOfCreate, payloadOfCreate)
		if err != nil {
			logging.Error("FAIL: error prepare http request: %s\n", requestOfCreate)
			logging.Error(err.Error())
			osExit(10, *logging)
		}
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		if err != nil {
			logging.Error("FAIL: error decode event http request")
			logging.Error(err.Error())
			osExit(12, *logging)
		}
		if response.StatusCode != 200 {
			response.Body.Close()
			logging.Error("FAIL: HTTP-status %d\n", response.StatusCode)
			osExit(13, *logging)
		}
		response.Body.Close()
		logging.Debug("Put event: %+v\n", event)
	}
	// ARCHIVE
	titlesOfArchived := []string{}
	titlesOfArchived = append(titlesOfArchived, titlesOfSend...)
	logging.Debug("ARCHIVE COUNT: %d", cfg.Counts.Archive)
	for i := 0; i < cfg.Counts.Archive; i++ {
		event := models.Event{
			Title:       fake.Title(),
			StartAt:     now.Add(-18000 * time.Minute),
			Duration:    1800,
			Description: fake.EmailSubject(),
			Owner:       fake.EmailAddress(),
			NotifyEarly: 60,
		}
		titlesOfArchived = append(titlesOfArchived, event.Title)
		titles = append(titles, event.Title)
		payload, err := json.Marshal(event)
		if err != nil {
			logging.Error(err.Error())
			osExit(20, *logging)
		}
		payloadOfCreate := strings.NewReader(string(payload))
		request, err := http.NewRequestWithContext(context.Background(), "POST", requestOfCreate, payloadOfCreate)
		if err != nil {
			logging.Error("FAIL: error prepare http request: %s\n", requestOfCreate)
			logging.Error(err.Error())
			osExit(21, *logging)
		}
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		if err != nil {
			logging.Error("FAIL: error decode event http request")
			logging.Error(err.Error())
			osExit(22, *logging)
		}
		response.Body.Close()
		logging.Debug("Put event: %+v\n", event)
	}
	// WAIT FOR NOTIFY
	titlesOfDefer := []string{}
	logging.Debug("DEFER COUNT: %d", cfg.Counts.Defer)
	for i := 0; i < cfg.Counts.Defer; i++ {
		event := models.Event{
			Title:       fake.Title(),
			StartAt:     now.Add(36000 * time.Second),
			Duration:    1800,
			Description: fake.EmailSubject(),
			Owner:       fake.EmailAddress(),
			NotifyEarly: 1,
		}
		titlesOfDefer = append(titlesOfDefer, event.Title)
		titles = append(titles, event.Title)
		payloadOfCreateRaw, err := json.Marshal(event)
		if err != nil {
			logging.Error(err.Error())
			osExit(30, *logging)
		}
		payloadOfCreate := strings.NewReader(string(payloadOfCreateRaw))
		request, err := http.NewRequestWithContext(context.Background(), "POST", requestOfCreate, payloadOfCreate)
		if err != nil {
			logging.Error("FAIL: error prepare http request: %s\n", requestOfCreate)
			logging.Error(err.Error())
			osExit(31, *logging)
		}
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		if err != nil {
			logging.Error(err.Error())
			osExit(32, *logging)
		}
		response.Body.Close()
		logging.Debug("Put event: %+v\n", event)
	}
	// CHECKING
	connectionSended, err := amqp.Dial(cfg.Sended.DSN)
	if err != nil {
		logging.Error(err.Error())
		osExit(40, *logging)
	}
	defer connectionSended.Close()
	channelSended, err := connectionSended.Channel()
	if err != nil {
		logging.Error(err.Error())
		osExit(41, *logging)
	}
	defer channelSended.Close()
	msgs, err := channelSended.Consume(
		cfg.Sended.QueueName, // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	var notice models.Notice
	checkCountOfSend := 0
	for d := range msgs {
		logging.Debug("Received a message: %s", d.Body)
		json.Unmarshal(d.Body, &notice)
		if slices.Contains(titlesOfSend, notice.Title) {
			checkCountOfSend += 1
		} else {
			logging.Error("Get unexpected send title: %s", notice.Title)
			logging.Error("checkCountOfSend: %+v", checkCountOfSend)
			osExit(43, *logging)
		}
		if checkCountOfSend == cfg.Counts.Send {
			break
		}
	}
	channelSended.Close()
	connectionSended.Close()
	logging.Info("OK. Get all notices of sended events")
	//
	connectionArchived, err := amqp.Dial(cfg.Archived.DSN)
	if err != nil {
		logging.Error(err.Error())
		osExit(50, *logging)
	}
	defer connectionArchived.Close()
	channelArchived, err := connectionArchived.Channel()
	if err != nil {
		logging.Error(err.Error())
		osExit(51, *logging)
	}
	defer channelArchived.Close()
	msgsArchived, err := channelArchived.Consume(
		cfg.Archived.QueueName, // queue
		"",                     // consumer
		true,                   // auto-ack
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,                    // args
	)
	countArchived := 0
	for d := range msgsArchived {
		logging.Debug("Received a message: %s", d.Body)
		json.Unmarshal(d.Body, &notice)
		if slices.Contains(titlesOfArchived, notice.Title) {
			countArchived += 1
		} else {
			logging.Error("Get unexpected archived title: %s", notice.Title)
			logging.Error("titlesOfArchived: %+v", titlesOfArchived)
			osExit(53, *logging)
		}
		if countArchived == cfg.Counts.Archive+cfg.Counts.Send {
			break
		}
	}
	channelArchived.Close()
	connectionArchived.Close()
	logging.Info("OK. Get all notices of archived events (also after send).")
	//
	db, err := sql.Open("postgres", cfg.Storage.DSN)
	if err != nil {
		logging.Error(err.Error())
		osExit(300, *logging)
	}
	var e models.Event
	logging.Debug("SELECT ALL")
	sqlStatement := `
		SELECT 
			"pk", "title", "description", "startat", "durationseconds", "owner", "notifyearlyseconds", "sheduled"
		FROM 
			hw15calendar.events`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		logging.Error(err.Error())
		osExit(301, *logging)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&e.PK, &e.Title, &e.Description, &e.StartAt, &e.Duration, &e.Owner, &e.NotifyEarly, &e.Sheduled)
		if err != nil {
			logging.Error(err.Error())
			osExit(302, *logging)
		}
		if !slices.Contains(titles, e.Title) {
			logging.Error("Get unexpected title %q", e.Title)
			logging.Error("Get unexpected title %v", titles)
			osExit(303, *logging)
		}
	}
	logging.Info("OK. All selected titles of created events are correct.")
	// TODO: Count, not any other
	logging.Debug("SELECT COUNT(*)")
	count := 0
	row := db.QueryRow(`SELECT COUNT(*) FROM hw15calendar.events`)
	err = row.Scan(&count)
	if err != nil {
		logging.Error(err.Error())
		osExit(304, *logging)
	}
	if count != cfg.Counts.Send+cfg.Counts.Archive+cfg.Counts.Defer {
		logging.Error("Get unexpected title count %d, expected %d", cfg.Counts.Send+cfg.Counts.Archive+cfg.Counts.Defer, count)
		osExit(305, *logging)
	} else {
		logging.Info("OK. Get expected count of events: %d.", count)
	}
	//
	logging.Info("Everything all right.")
	// os.Exit(0)
}
