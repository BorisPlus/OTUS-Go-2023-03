package api

import (
	"encoding/json"
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
	responses "hw12_13_14_15_calendar/internal/server/http/api/api_response"
	commonHandlers "hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

// curl -X POST -H 'Content-Type: application/json' -d "{\"test\": \"that\"}"

type EventsCreateHandler struct {
	Logger interfaces.Logger
	App    interfaces.Applicationer
}

func (h EventsCreateHandler) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	apiMethod := "api.events.create"
	if rr.Method != "POST" {
		commonHandlers.InvalidHTTPMethod{APIMethod: apiMethod}.ServeHTTP(rw, rr)
		return
	}
	var event models.Event
	err := json.NewDecoder(rr.Body).Decode(&event)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	createdEvent, err := h.App.CreateEvent(&event)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	apiResponse := responses.NewAPIResponse(apiMethod)
	apiResponse.Data = responses.DataItem{Item: createdEvent}
	apiResponseJSON, err := json.Marshal(apiResponse)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(apiResponseJSON)
}
