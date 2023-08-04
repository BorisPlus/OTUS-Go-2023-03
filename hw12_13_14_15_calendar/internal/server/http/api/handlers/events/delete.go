package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
	responses "hw12_13_14_15_calendar/internal/server/http/api/api_response"
	commonHandlers "hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

type EventsDeleteHandler struct {
	Logger interfaces.Logger
	App    interfaces.Applicationer
}

func (h EventsDeleteHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	apiMethod := "api.events.delete"
	if request.Method != "DELETE" {
		commonHandlers.InvalidHTTPMethod{APIMethod: apiMethod}.ServeHTTP(response, request)
		return
	}
	pkString := request.Form.Get("id")
	pk, err := strconv.Atoi(pkString)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(response, request)
		return
	}
	event := models.Event{}
	event.PK = pk
	deletedEvent, err := h.App.DeleteEvent(&event)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(response, request)
		return
	}
	apiResponse := responses.NewAPIResponse(apiMethod)
	apiResponse.Data = responses.DataItem{Item: deletedEvent}
	apiResponseJSON, err := json.Marshal(apiResponse)
	if err != nil {
		commonHandlers.CustomErrorHandler{APIMethod: apiMethod, Error: err}.ServeHTTP(response, request)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(apiResponseJSON)
}
