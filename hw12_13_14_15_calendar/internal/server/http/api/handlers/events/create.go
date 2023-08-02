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

type ApiEventsCreateHandler struct {
	Logger interfaces.Logger
	App    interfaces.Applicationer
}

func (h ApiEventsCreateHandler) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	ApiMethod := "api.events.create"
	if rr.Method != "POST" {
		commonHandlers.InvalidHTTPMethod{ApiMethod: ApiMethod}.ServeHTTP(rw, rr)
		return
	}
	var event models.Event
	err := json.NewDecoder(rr.Body).Decode(&event)
	if err != nil {
		commonHandlers.CustomErrorHandler{ApiMethod: ApiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	createdEvent, err := h.App.CreateEvent(&event)
	if err != nil {
		commonHandlers.CustomErrorHandler{ApiMethod: ApiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	apiResponse := responses.NewAPIResponse(ApiMethod)
	apiResponse.Data = responses.DataItem{Item: createdEvent}
	apiResponseJSON, err := apiResponse.MarshalJSON()
	if err != nil {
		commonHandlers.CustomErrorHandler{ApiMethod: ApiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(apiResponseJSON)
}
