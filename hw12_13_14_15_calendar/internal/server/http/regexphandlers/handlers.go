package regexhandlers

import (
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
)

// as HandleFunc
// http.Handle("/api/", middleware.Middleware(http.HandleFunc("/api/", api.Routers.Go), s.logger))

// as Handle
// http.HandleFunc("/api/", api.Routers(s.logger, s.app).ServeHTTP)

type RegexpHandler struct {
	qpp     QueryPathPattern
	handler http.Handler
}

func NewRegexpHandler(pattern string, params Params, handler http.Handler) *RegexpHandler {
	rh := new(RegexpHandler)
	rh.qpp = *NewQueryPathPattern(pattern, params)
	rh.handler = handler
	return rh
}

type RegexpHandlers struct {
	defaultHandler http.Handler
	logger         interfaces.Logger
	app            interfaces.Applicationer
	crossroad      []RegexpHandler
}

func NewRegexpHandlers(defaultHandler http.Handler, logger interfaces.Logger, app interfaces.Applicationer, rh ...RegexpHandler) RegexpHandlers {
	//
	// RegexpHandlers{defaultHandler, logger, app, rh}
	regexpHandlers := new(RegexpHandlers)
	regexpHandlers.defaultHandler = defaultHandler
	regexpHandlers.logger = logger
	regexpHandlers.app = app
	regexpHandlers.crossroad = rh
	return *regexpHandlers
}

func (rhs RegexpHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlerWasNotFound := true // TODO: do.Once?
	for _, rh := range rhs.crossroad {
		if rh.qpp.match(r.URL.Path) {
			handlerWasNotFound = false
			r.Form = rh.qpp.GetValues(r.URL.Path)
			rh.handler.ServeHTTP(w, r)
			break
		}
	}
	if handlerWasNotFound && rhs.defaultHandler != nil {
		rhs.defaultHandler.ServeHTTP(w, r)
	}
}
