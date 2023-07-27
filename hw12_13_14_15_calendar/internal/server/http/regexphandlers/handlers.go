package regexhandlers

import (
	"net/http"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
)

type RegexpHandler struct {
	QueryPathPattern QueryPathPattern
	handler          http.Handler
}

func NewRegexpHandler(pattern string, ParamsNaming []string, handler http.Handler) RegexpHandler {
	return RegexpHandler{
		QueryPathPattern: QueryPathPattern{
			pattern:      pattern,
			ParamsNaming: ParamsNaming,
		},
		handler: handler,
	}
}

type RegexpHandlers struct {
	defaultHandler http.Handler
	logger         interfaces.Logger
	app            interfaces.Applicationer
	Сrossroad      []RegexpHandler
}

func (handlers RegexpHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlerWasNotFound := true
	for _, handler := range handlers.Сrossroad {
		if handler.QueryPathPattern.match(r.URL.Path) {
			handlerWasNotFound = false
			r.Form = handler.QueryPathPattern.fetch(r.URL.Path)
			handler.handler.ServeHTTP(w, r)
			break
		}
	}
	if handlerWasNotFound && handlers.defaultHandler != nil {
		handlers.defaultHandler.ServeHTTP(w, r)
	}
}

func NewRegexpHandlers(defaultHandler http.Handler, logger interfaces.Logger, app interfaces.Applicationer, handlers ...RegexpHandler) RegexpHandlers {
	return RegexpHandlers{defaultHandler, logger, app, handlers}
}
