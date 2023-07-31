package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
	common "hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

// curl -X POST -H 'Content-Type: application/json' -d "{\"test\": \"that\"}"

type ApiEventsCreateHandler struct {
	Logger interfaces.Logger
	App    interfaces.Applicationer
}

func (h ApiEventsCreateHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		common.InvalidHTTPMethodForUrlPathHandler{}.ServeHTTP(response, request)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var event models.Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		panic(err)
	}
	log.Println(event)
	_ = h.Logger
	_ = h.App
	h.Logger.Info("%+v", request.Form)
	response.Write([]byte("ApiEventsCreateHandler"))
}
