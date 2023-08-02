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

type ApiEventsUpdateHandler struct {
	Logger interfaces.Logger
	App    interfaces.Applicationer
}

func (h ApiEventsUpdateHandler) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	ApiMethod := "api.events.update"
	h.Logger.Info("%+v", rr.Form)
	if rr.Method != "PUTCH" {
		commonHandlers.InvalidHTTPMethod{ApiMethod: ApiMethod}.ServeHTTP(rw, rr)
		return
	}
	var event models.Event
	err := json.NewDecoder(rr.Body).Decode(&event)
	if err != nil {
		commonHandlers.CustomErrorHandler{ApiMethod: ApiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	pkString := rr.Form.Get("id")
	pk, err := strconv.Atoi(pkString)
	if err != nil {
		commonHandlers.CustomErrorHandler{ApiMethod: ApiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	event.PK = pk
	_, err = h.App.UpdateEvent(&event)
	if err != nil {
		commonHandlers.CustomErrorHandler{ApiMethod: ApiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	apiResponse := responses.NewAPIResponse(ApiMethod)
	apiResponse.Data = responses.DataItem{Item: event}
	apiResponseJSON, err := apiResponse.MarshalJSON()
	if err != nil {
		commonHandlers.CustomErrorHandler{ApiMethod: ApiMethod, Error: err}.ServeHTTP(rw, rr)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(apiResponseJSON)
}
