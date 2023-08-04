package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	responses "hw12_13_14_15_calendar/internal/server/http/api/api_response"
	commonHandlers "hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

type EventsGetHandler struct {
	Logger interfaces.Logger
	App    interfaces.Applicationer
}

func (h EventsGetHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	apiMethod := "api.events.get"
	h.Logger.Info("%+v", request.Form)
	if request.Method != "GET" {
		commonHandlers.InvalidHTTPMethod{APIMethod: apiMethod}.ServeHTTP(response, request)
		return
	}
	pkString := request.Form.Get("id")
	pk, err := strconv.Atoi(pkString)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(response, request)
		return
	}
	event, err := h.App.ReadEvent(pk)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(response, request)
		return
	}
	apiResponse := responses.NewAPIResponse(apiMethod)
	apiResponse.Data = responses.DataItem{Item: event}
	apiResponseJSON, err := json.Marshal(apiResponse)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(response, request)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(apiResponseJSON)
}
