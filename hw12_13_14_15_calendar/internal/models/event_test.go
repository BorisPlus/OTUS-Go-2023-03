package models

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type User struct {
	ID       int
	Name     string `json:"name"`
	Datetime time.Time
}

func TestLogger(t *testing.T) {
	var event Event
	line := `{"title":"title","startat":"2021-02-18T21:54:42+02:00","duration":0}`
	if err := json.Unmarshal([]byte(line), &event); err != nil {
		fmt.Println(err)
	}
	if event.Title != "title" {
		t.Errorf("Could not Unmarshal object")
	} else {
		fmt.Println("Unmarshaled. OK.")
	}
}
