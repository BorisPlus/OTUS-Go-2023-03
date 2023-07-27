package api

import (
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
)

type ApiEventsListHandler struct{
	Logger interfaces.Logger
	App interfaces.Applicationer
}

func (h ApiEventsListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = h.App
	_ = h.Logger
	w.Write([]byte("ApiEventsListHandler"))
}
