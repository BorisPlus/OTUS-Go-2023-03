package api

import (
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	commonHandlers "hw12_13_14_15_calendar/internal/server/http/api/handlers/common"
	apiEvents "hw12_13_14_15_calendar/internal/server/http/api/handlers/events"
	middleware "hw12_13_14_15_calendar/internal/server/http/middleware"
	regexped "hw12_13_14_15_calendar/internal/server/http/regexphandlers"
)

type DefaultHandler struct{}

func (h DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/version", http.StatusTemporaryRedirect)
}

var (
	none = regexped.Params{}
	id   = regexped.Params{"id"}
)

func Handlers(logger interfaces.Logger, app interfaces.Applicationer) regexped.RegexpHandlers {
	return regexped.NewRegexpHandlers(
		middleware.Instance().Listen(DefaultHandler{}),
		logger,
		app,
		*regexped.NewRegexpHandler(
			`/api/version`,
			none,
			middleware.Instance().Listen(commonHandlers.VersionHandler{}),
		),
		*regexped.NewRegexpHandler(
			`/api/events/`,
			none,
			middleware.Instance().Listen(apiEvents.EventsListHandler{
				Logger: logger, App: app, APIMethod: "api.events.list",
				Listing: app.ListEvents,
			}),
		),
		*regexped.NewRegexpHandler(
			`/api/events/notsheduled`,
			none,
			middleware.Instance().Listen(apiEvents.EventsListHandler{
				Logger: logger, App: app, APIMethod: "api.events.listnotsheduled",
				Listing: app.ListNotSheduledEvents,
			}),
		),
		*regexped.NewRegexpHandler(
			`/api/events/create`,
			none,
			middleware.Instance().Listen(apiEvents.EventsCreateHandler{Logger: logger, App: app}),
		),
		*regexped.NewRegexpHandler(
			`/api/events/{numeric}`,
			id,
			apiEvents.EventsGetHandler{Logger: logger, App: app},
		),
		*regexped.NewRegexpHandler(
			`/api/events/{numeric}/update`,
			id,
			apiEvents.EventsUpdateHandler{Logger: logger, App: app},
		),
		*regexped.NewRegexpHandler(
			`/api/events/{numeric}/delete`,
			id,
			apiEvents.EventsDeleteHandler{Logger: logger, App: app},
		),
	)
}
