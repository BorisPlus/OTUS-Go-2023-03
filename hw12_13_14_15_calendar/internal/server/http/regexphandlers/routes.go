package regexhandlers

import (
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
)

type RegexpHandler struct {
	QueryPathPattern QueryPathPattern
	Handler          http.Handler
}

func NewRegexpHandler(Pattern string, ParamsNaming []string, Handler http.Handler) RegexpHandler {
	return RegexpHandler{
		QueryPathPattern: QueryPathPattern{
			Pattern:      Pattern,
			ParamsNaming: ParamsNaming,
		},
		Handler: Handler,
	}
}

type RegexpHandlers struct {
	Default   func(w http.ResponseWriter, r *http.Request)
	logger    interfaces.Logger
	app       interfaces.Applicationer
	Сrossroad []RegexpHandler
}

func (routes RegexpHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routeWasNotFind := true
	for _, route := range routes.Сrossroad {
		if route.QueryPathPattern.Match(r.URL.Path) {
			routeWasNotFind = false
			r.Form = route.QueryPathPattern.Fetch(r.URL.Path)
			route.Handler.ServeHTTP(w, r)
			break
		}
	}
	if routeWasNotFind && routes.Default != nil {
		routes.Default(w, r)
	}
}

func NewRoutes(DefaultRoute func(w http.ResponseWriter, r *http.Request), logger interfaces.Logger, app interfaces.Applicationer, routes ...RegexpHandler) RegexpHandlers {
	return RegexpHandlers{DefaultRoute, logger, app, routes}
}
