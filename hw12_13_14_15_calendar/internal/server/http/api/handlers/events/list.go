package api

import (
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	responses "hw12_13_14_15_calendar/internal/server/http/api/api_response"
	commonHandlers "hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

type ApiEventsListHandler struct {
	Logger interfaces.Logger
	App    interfaces.Applicationer
}

func (h ApiEventsListHandler) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	ApiMethod := "api.events.list"
	_ = h.Logger
	if rr.Method != "GET" {
		commonHandlers.InvalidHTTPMethod{ApiMethod: ApiMethod}.ServeHTTP(rw, rr)
		return
	}
	events, err := h.App.ListEvents()
	if err != nil {
		commonHandlers.CustomErrorHandler{ApiMethod: ApiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	apiResponse := responses.NewAPIResponse(ApiMethod)
	apiResponse.Data = responses.DataItems{Items: events}
	apiResponseJSON, err := apiResponse.MarshalJSON()
	if err != nil {
		commonHandlers.CustomErrorHandler{ApiMethod: ApiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(apiResponseJSON)
}
