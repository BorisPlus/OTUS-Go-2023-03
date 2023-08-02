package api

import (
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	commonHandlers "hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
	apiEventsHandlers "hw12_13_14_15_calendar/internal/server/http/api/handlers/events"
	middleware "hw12_13_14_15_calendar/internal/server/http/middleware"
	regexped "hw12_13_14_15_calendar/internal/server/http/regexphandlers"
)

type DefaultHandler struct{}

func (h DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/version", http.StatusTemporaryRedirect)
}

var __ = regexped.Params{}
var id = regexped.Params{"id"}

func Handlers(logger interfaces.Logger, app interfaces.Applicationer) regexped.RegexpHandlers {
	return regexped.NewRegexpHandlers(
		middleware.Instance().Listen(DefaultHandler{}),
		logger,
		app,
		*regexped.NewRegexpHandler(`/api/version`, __, middleware.Instance().Listen(commonHandlers.VersionHandler{})),
		*regexped.NewRegexpHandler(`/api/events/`, __, middleware.Instance().Listen(apiEventsHandlers.ApiEventsListHandler{Logger: logger, App: app})),
		*regexped.NewRegexpHandler(`/api/events/create`, __, apiEventsHandlers.ApiEventsCreateHandler{Logger: logger, App: app}),
		*regexped.NewRegexpHandler(`/api/events/{numeric}`, id, apiEventsHandlers.ApiEventsGetHandler{Logger: logger, App: app}),
		*regexped.NewRegexpHandler(`/api/events/{numeric}/update`, id, apiEventsHandlers.ApiEventsUpdateHandler{Logger: logger, App: app}),
		*regexped.NewRegexpHandler(`/api/events/{numeric}/delete`, id, apiEventsHandlers.ApiEventsDeleteHandler{Logger: logger, App: app}),
	)
}
