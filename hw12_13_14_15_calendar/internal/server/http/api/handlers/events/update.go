package api

import (
	"net/http"
	interfaces "hw12_13_14_15_calendar/internal/interfaces"
)

type ApiEventsUpdateHandler struct{
	logger interfaces.Logger
	app interfaces.Applicationer
}

func (h ApiEventsUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = h.app
	_ = h.logger
	w.Write([]byte("ApiEventsUpdateHandler"))
}
