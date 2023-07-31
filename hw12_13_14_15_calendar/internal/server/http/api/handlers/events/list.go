package api

import (
	"net/http"
	"encoding/json"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	commonHandlers "hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

type ApiEventsListHandler struct {
	Logger interfaces.Logger
	App    interfaces.Applicationer
}

func (h ApiEventsListHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	_ = h.Logger
	if request.Method != "GET" {
		commonHandlers.InvalidHTTPMethodForUrlPathHandler{}.ServeHTTP(response, request)
		return
	}
	events, err := h.App.ListEvents()
	if err != nil {
		commonHandlers.CustomErrorHandler{Error: err}.ServeHTTP(response, request)
		return
	}
	eventsJSON, err := json.Marshal(events)
	if err != nil {
		commonHandlers.CustomErrorHandler{Error: err}.ServeHTTP(response, request)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(eventsJSON)
}
