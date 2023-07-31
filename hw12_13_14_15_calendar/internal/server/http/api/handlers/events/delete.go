package api

import (
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	common "hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
)

type ApiEventsDeleteHandler struct {
	Logger interfaces.Logger
	App    interfaces.Applicationer
}

func (h ApiEventsDeleteHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != "PUT" {
		common.InvalidHTTPMethodForUrlPathHandler{}.ServeHTTP(response, request)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	_ = h.App
	_ = h.Logger
	response.Write([]byte("ApiEventsDeleteHandler"))
}
