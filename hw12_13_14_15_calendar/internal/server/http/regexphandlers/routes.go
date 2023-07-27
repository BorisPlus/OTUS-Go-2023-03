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
	defaultHandler http.Handler
	logger         interfaces.Logger
	app            interfaces.Applicationer
	Сrossroad      []RegexpHandler
}

func (handlers RegexpHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlerWasNotFound := true
	for _, handler := range handlers.Сrossroad {
		if handler.QueryPathPattern.Match(r.URL.Path) {
			handlerWasNotFound = false
			r.Form = handler.QueryPathPattern.Fetch(r.URL.Path)
			handler.Handler.ServeHTTP(w, r)
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
