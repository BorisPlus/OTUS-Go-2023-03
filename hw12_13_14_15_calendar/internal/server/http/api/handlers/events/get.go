package api

import (
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	common "hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

type ApiEventsGetHandler struct {
	Logger interfaces.Logger
	App    interfaces.Applicationer
}

func (h ApiEventsGetHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		common.InvalidHTTPMethodForUrlPathHandler{}.ServeHTTP(response, request)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	_ = h.Logger
	_ = h.App
	h.Logger.Info("%+v", request.Form)
	response.Write([]byte("ApiEventsGetHandler"))
}
