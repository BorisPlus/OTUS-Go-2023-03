package regexhandlers

import (
	"net/http"
	// "net/url"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
)


// type ParamsNamedHandler interface {
// 	http.Handler
// 	GetParamsNames() []string
// }

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

// func (rh RegexpHandler) GetValues(url string) url.Values { // TODO: url.URL
// 	return rh.qpp.GetValues(url, rh.handler.GetParamsNames())
// }

type RegexpHandlers struct {
	defaultHandler http.Handler
	logger         interfaces.Logger
	app            interfaces.Applicationer
	crossroad      []RegexpHandler
}

func NewRegexpHandlers(defaultHandler http.Handler, logger interfaces.Logger, app interfaces.Applicationer, rh ...RegexpHandler) RegexpHandlers {
	return RegexpHandlers{defaultHandler, logger, app, rh}
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
