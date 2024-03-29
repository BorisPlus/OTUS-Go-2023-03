package api

// dupl-disable.

import (
	"encoding/json"
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	responses "hw12_13_14_15_calendar/internal/server/http/api/api_response"
	commonHandlers "hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

type EventsListNotSheduledHandler struct {
	Logger interfaces.Logger
	App    interfaces.Applicationer
}

func (h EventsListNotSheduledHandler) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	apiMethod := "api.events.listnotsheduled"
	_ = h.Logger
	if rr.Method != commonHandlers.GET {
		commonHandlers.InvalidHTTPMethod{APIMethod: apiMethod}.ServeHTTP(rw, rr)
		return
	}
	events, err := h.App.ListNotSheduledEvents()
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	apiResponse := responses.NewAPIResponse(apiMethod)
	apiResponse.Data = responses.DataItems{Items: events}
	apiResponseJSON, err := json.Marshal(apiResponse)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(apiResponseJSON)
}

// dupl-enable.
