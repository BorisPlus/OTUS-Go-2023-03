package api

import (
	"encoding/json"
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
	responses "hw12_13_14_15_calendar/internal/server/http/api/api_response"
	commonHandlers "hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

type EventsListHandler struct {
	Logger    interfaces.Logger
	App       interfaces.Applicationer
	APIMethod string
	Listing   func() ([]models.Event, error)
}

func (h EventsListHandler) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	_ = h.Logger
	if rr.Method != commonHandlers.GET {
		commonHandlers.InvalidHTTPMethod{APIMethod: h.APIMethod}.ServeHTTP(rw, rr)
		return
	}
	events, err := h.App.ListEvents()
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: h.APIMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	apiResponse := responses.NewAPIResponse(h.APIMethod)
	apiResponse.Data = responses.DataItems{Items: events}
	apiResponseJSON, err := json.Marshal(apiResponse)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: h.APIMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(apiResponseJSON)
}
