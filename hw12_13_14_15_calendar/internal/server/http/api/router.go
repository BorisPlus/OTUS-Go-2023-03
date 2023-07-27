package api

import (
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	apiHandlers "hw12_13_14_15_calendar/internal/server/http/api/handlers"
	apiEventsHandlers "hw12_13_14_15_calendar/internal/server/http/api/handlers/events"
	middleware "hw12_13_14_15_calendar/internal/server/http/middleware"
	regexped "hw12_13_14_15_calendar/internal/server/http/regexphandlers"
)

var defaultHandler = func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/version", http.StatusTemporaryRedirect)
}

// var a = middleware.Middleware(apiHandlers.ApiVersionHandler{}, nil)
var __ = []string{}
var id = []string{"id"}

func Handlers(logger interfaces.Logger, app interfaces.Applicationer) regexped.RegexpHandlers {
	return regexped.NewRoutes(
		defaultHandler,
		logger,
		app,
		regexped.NewRegexpHandler(`/api/version`, __, middleware.Middleware(apiHandlers.ApiVersionHandler{}, logger)),
		regexped.NewRegexpHandler(`/api/events/`, __, apiEventsHandlers.ApiEventsListHandler{Logger: logger, App: app}),
		regexped.NewRegexpHandler(`/api/events/create/`, __, apiEventsHandlers.ApiEventsListHandler{}),
		regexped.NewRegexpHandler(`/api/events/{numeric}`, id, apiEventsHandlers.ApiEventsGetHandler{}),
		regexped.NewRegexpHandler(`/api/events/{numeric}/update`, id, apiEventsHandlers.ApiEventsUpdateHandler{}),
		regexped.NewRegexpHandler(`/api/events/{numeric}/delete`, id, apiEventsHandlers.ApiEventsDeleteHandler{}),
	)
}
