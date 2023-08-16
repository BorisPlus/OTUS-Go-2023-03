package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/icrowley/fake"

	config "hw12_13_14_15_calendar/internal/config"
	models "hw12_13_14_15_calendar/internal/models"
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
	mainConfig := config.NewCalendarConfig()
	err = viper.Unmarshal(mainConfig)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	datasetSize := 20
	// CREATE DATASET SERIA
	client := &http.Client{}
	requestOfCreate := fmt.Sprintf("http://%s:%d/api/events/create", mainConfig.HTTP.Host, mainConfig.HTTP.Port)
	for i := 1; i <= datasetSize; i++ {
		event := models.Event{
			Title:       fake.Title(),
			StartAt:     time.Now().Add(15 * time.Second),
			Duration:    1800,
			Description: fake.EmailSubject(),
			Owner:       fake.EmailAddress(),
			NotifyEarly: 30,
		}
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
}
