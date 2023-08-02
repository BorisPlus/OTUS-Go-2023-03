package models

import (
	"fmt"
	"time"
	"testing"

	"encoding/json"
)

//

type User struct {
	ID       int
	Name     string `json:"name"`
	Datetime     time.Time
}


func TestLogger(t *testing.T) {
	var event Event
	line := `{"title":"title","startat":"2021-02-18T21:54:42+02:00","duration":0,"description": "description","owner":"owner","notifyearly":0}`
	// line := `{"title":"title","startat":"2021-02-18T21:54:42+02:00","duration":0,"description": "description","owner":"owner","notifyearly":0}`
	// line := `{"title":"Event","startat":"2021-02-18T21:54:42+02:00"}`

	if err := json.Unmarshal([]byte(line), &event); err != nil {
		fmt.Println(err)
	}
	fmt.Println(event.Title)
	fmt.Println(event.PK)
	fmt.Println(event.StartAt)
}
