package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq" // a blank import should be justifying.

	"golang.org/x/exp/slices"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/icrowley/fake"

	config "hw12_13_14_15_calendar/internal/config"
	models "hw12_13_14_15_calendar/internal/models"
	// amqp "github.com/rabbitmq/amqp091-go"
)

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "", "Path to configuration file")
}

func main() {
	pflag.Parse()
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
	mainConfig := config.NewCheckerConfig()
	err = viper.Unmarshal(mainConfig)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	datasetSize := 10
	// CREATE DATASET SERIA
	client := &http.Client{}
	requestOfCreate := fmt.Sprintf("http://%s:%d/api/events/create", mainConfig.HTTP.Host, mainConfig.HTTP.Port)
	now := time.Now()
	titles := []string{}
	// MUST BY NOTIFIED
	notifiedTitles := []string{}
	for i := 1; i <= datasetSize; i++ {
		event := models.Event{
			Title:       fake.Title(),
			StartAt:     now.Add(3000 * time.Second),
			Duration:    1800,
			Description: fake.EmailSubject(),
			Owner:       fake.EmailAddress(),
			NotifyEarly: 3600,
		}
		notifiedTitles = append(notifiedTitles, event.Title)
		titles = append(titles, event.Title)
		payloadOfCreateRaw, _ := json.Marshal(event)
		payloadOfCreate := strings.NewReader(string(payloadOfCreateRaw))
		request, err := http.NewRequestWithContext(context.Background(), "POST", requestOfCreate, payloadOfCreate)
		if err != nil {
			log.Printf("FAIL: error prepare http request: %s\n", requestOfCreate)
			return
		}
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		if err != nil {
			log.Printf("FAIL: error decode event http request: %s\n", err)
			return
		}
		response.Body.Close()
		log.Printf("Put event: %+v\n", event)
	}
	_ = notifiedTitles
	// MUST BE ARCHIVE
	archivedTitles := []string{}
	for i := 1; i <= datasetSize; i++ {
		event := models.Event{
			Title:       fake.Title(),
			StartAt:     now.Add(-18000 * time.Minute),
			Duration:    1800,
			Description: fake.EmailSubject(),
			Owner:       fake.EmailAddress(),
			NotifyEarly: 60,
		}
		archivedTitles = append(archivedTitles, event.Title)
		titles = append(titles, event.Title)
		payloadOfCreateRaw, _ := json.Marshal(event)
		payloadOfCreate := strings.NewReader(string(payloadOfCreateRaw))
		request, err := http.NewRequestWithContext(context.Background(), "POST", requestOfCreate, payloadOfCreate)
		if err != nil {
			log.Printf("FAIL: error prepare http request: %s\n", requestOfCreate)
			return
		}
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		if err != nil {
			log.Printf("FAIL: error decode event http request: %s\n", err)
			return
		}
		response.Body.Close()
		log.Printf("Put event: %+v\n", event)
	}
	_ = archivedTitles
	// WAIT FOR NOTIFY
	for i := 1; i <= datasetSize; i++ {
		event := models.Event{
			Title:       fake.Title(),
			StartAt:     now.Add(36000 * time.Second),
			Duration:    1800,
			Description: fake.EmailSubject(),
			Owner:       fake.EmailAddress(),
			NotifyEarly: 1,
		}
		titles = append(titles, event.Title)
		payloadOfCreateRaw, _ := json.Marshal(event)
		payloadOfCreate := strings.NewReader(string(payloadOfCreateRaw))
		request, err := http.NewRequestWithContext(context.Background(), "POST", requestOfCreate, payloadOfCreate)
		if err != nil {
			log.Printf("FAIL: error prepare http request: %s\n", requestOfCreate)
			return
		}
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		if err != nil {
			log.Printf("FAIL: error decode event http request: %s\n", err)
			return
		}
		response.Body.Close()
		log.Printf("Put event: %+v\n", event)
	}
	//
	fmt.Printf("%+v", mainConfig)
	// connectionSended, err := amqp.Dial(mainConfig.Sended.DSN)
	// if err != nil {
	// 	log.Print(err.Error())
	// 	os.Exit(1)
	// }
	// defer connectionSended.Close()
	// channelSended, err := connectionSended.Channel()
	// if err != nil {
	// 	log.Print(err.Error())
	// 	os.Exit(2)
	// }
	// defer channelSended.Close()
	// q, err := channelSended.QueueDeclare(
	// 	mainConfig.Sended.QueueName, // name
	// 	false,                       // durable
	// 	false,                       // delete when unused
	// 	false,                       // exclusive
	// 	false,                       // no-wait
	// 	nil,                         // arguments
	// )
	// if err != nil {
	// 	log.Print(err.Error())
	// 	os.Exit(12)
	// }
	// msgs, err := channelSended.Consume(
	// 	q.Name, // queue
	// 	"",     // consumer
	// 	true,   // auto-ack
	// 	false,  // exclusive
	// 	false,  // no-local
	// 	false,  // no-wait
	// 	nil,    // args
	// )
	// var notice models.Notice
	// count := 0
	// for d := range msgs {
	// 	log.Printf("Received a message: %s", d.Body)
	// 	json.Unmarshal(d.Body, &notice)
	// 	if slices.Contains(notifiedTitles, notice.Title) {
	// 		count += 1
	// 	} else {
	// 		os.Exit(3)
	// 	}
	// 	if count == 10 {
	// 		break
	// 	}
	// }
	// channelSended.Close()
	// connectionSended.Close()
	// //
	// connectionArchived, err := amqp.Dial(mainConfig.Archived.DSN)
	// if err != nil {
	// 	log.Print(err.Error())
	// 	os.Exit(4)
	// }
	// defer connectionArchived.Close()
	// channelArchived, err := connectionArchived.Channel()
	// if err != nil {
	// 	log.Print(err.Error())
	// 	os.Exit(5)
	// }
	// defer channelArchived.Close()
	// qArchived, err := channelArchived.QueueDeclare(
	// 	mainConfig.Archived.QueueName, // name
	// 	false,                         // durable
	// 	false,                         // delete when unused
	// 	false,                         // exclusive
	// 	false,                         // no-wait
	// 	nil,                           // arguments
	// )
	// if err != nil {
	// 	log.Print(err.Error())
	// 	os.Exit(6)
	// }
	// msgsArchived, err := channelArchived.Consume(
	// 	qArchived.Name, // queue
	// 	"",             // consumer
	// 	true,           // auto-ack
	// 	false,          // exclusive
	// 	false,          // no-local
	// 	false,          // no-wait
	// 	nil,            // args
	// )
	// countArchived := 0
	// for d := range msgsArchived {
	// 	log.Printf("Received a message: %s", d.Body)
	// 	json.Unmarshal(d.Body, &notice)
	// 	if slices.Contains(notifiedTitles, notice.Title) {
	// 		countArchived += 1
	// 	} else {
	// 		os.Exit(7)
	// 	}
	// 	if count == 10 {
	// 		break
	// 	}
	// }
	// channelArchived.Close()
	// connectionArchived.Close()
	//
	db, err := sql.Open("postgres", mainConfig.Storage.DSN)
	if err != nil {
		log.Print(err.Error())
		os.Exit(8)
	}
	var e models.Event
	sqlStatement := `
		SELECT 
			"pk", "title", "description", "startat", "durationseconds", "owner", "notifyearlyseconds", "sheduled"
		FROM 
			hw12calendar.events`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Print(err.Error())
		os.Exit(9)
	}
	defer rows.Close()
	countTitled := 0
	for rows.Next() {
		countTitled += 1
		err = rows.Scan(&e.PK, &e.Title, &e.Description, &e.StartAt, &e.Duration, &e.Owner, &e.NotifyEarly, &e.Sheduled)
		if err != nil {
			log.Print(err.Error())
			os.Exit(10)
		}
		if slices.Contains(titles, e.Title) {
			countTitled += 1
		} else {
			os.Exit(11)
		}
	}
	// TODO: Count, not any other
	count := 0
	row, err := db.Query(`SELECT COUNT(*) FROM hw12calendar.events`)
	if err != nil {
		log.Print(err.Error())
		os.Exit(13)
	}
	defer row.Close()
	err = row.Scan(&count)
	if err != nil {
		log.Print(err.Error())
		os.Exit(14)
	}
	if count != len(titles) {
		os.Exit(15)
	}
}
