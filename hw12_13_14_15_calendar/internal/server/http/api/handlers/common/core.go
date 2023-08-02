package api

import (
	"net/http"

	responses "hw12_13_14_15_calendar/internal/server/http/api/api_response"
)

type InvalidHTTPMethod struct {
	ApiMethod string
}

func (h InvalidHTTPMethod) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	apiResponse := responses.InvalidHTTPMethod(h.ApiMethod)
	apiResponseJSON, _ := apiResponse.MarshalJSON()
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusMethodNotAllowed)
	rw.Write(apiResponseJSON)
}

type InvalidRequestBodyHandler struct {
	ApiMethod string
}

func (h InvalidRequestBodyHandler) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	apiResponse := responses.InvalidRequestBody(h.ApiMethod)
	apiResponseJSON, _ := apiResponse.MarshalJSON()
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write(apiResponseJSON)
}

type InternalServerErrorHandler struct {
	ApiMethod string
}

func (h InternalServerErrorHandler) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	apiResponse := responses.InternalServerError(h.ApiMethod)
	apiResponseJSON, _ := apiResponse.MarshalJSON()
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write(apiResponseJSON)
}

type CustomErrorHandler struct {
	ApiMethod string
	Error     error
}

func (h CustomErrorHandler) ServeHTTP(rw http.ResponseWriter, rr *http.Request) {
	apiResponse := responses.NewErrorAPIResponse(h.ApiMethod, h.Error)
	apiResponseJSON, err := apiResponse.MarshalJSON()
	if err != nil {
		InternalServerErrorHandler{h.ApiMethod}.ServeHTTP(rw, rr)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write(apiResponseJSON)
}
