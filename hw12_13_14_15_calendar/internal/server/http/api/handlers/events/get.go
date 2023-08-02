package api

import (
	"net/http"
	"strconv"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	responses "hw12_13_14_15_calendar/internal/server/http/api/api_response"
	commonHandlers "hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

type ApiEventsGetHandler struct {
	Logger interfaces.Logger
	App    interfaces.Applicationer
}

func (h ApiEventsGetHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ApiMethod := "api.events.get"
	h.Logger.Info("%+v", request.Form)
	if request.Method != "GET" {
		commonHandlers.InvalidHTTPMethod{ApiMethod: ApiMethod}.ServeHTTP(response, request)
		return
	}
	pkString := request.Form.Get("id")
	pk, err := strconv.Atoi(pkString)
	if err != nil {
		commonHandlers.CustomErrorHandler{ApiMethod: ApiMethod, Error: err}.ServeHTTP(response, request)
		return
	}
	event, err := h.App.ReadEvent(pk)
	if err != nil {
		commonHandlers.CustomErrorHandler{ApiMethod: ApiMethod, Error: err}.ServeHTTP(response, request)
		return
	}
	apiResponse := responses.NewAPIResponse(ApiMethod)
	apiResponse.Data = responses.DataItem{Item: event}
	apiResponseJSON, err := apiResponse.MarshalJSON()
	if err != nil {
		commonHandlers.CustomErrorHandler{ApiMethod: ApiMethod, Error: err}.ServeHTTP(response, request)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(apiResponseJSON)
}
