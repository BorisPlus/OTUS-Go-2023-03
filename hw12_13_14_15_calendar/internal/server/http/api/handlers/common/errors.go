package api

import (
	"net/http"
	"encoding/json"

	responses "hw12_13_14_15_calendar/internal/server/http/api/response"
)

type InvalidHTTPMethodForUrlPathHandler struct{}

func (h InvalidHTTPMethodForUrlPathHandler) ServeHTTP(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusMethodNotAllowed)
	response.Write(responses.InvalidHTTPMethodJSON)
}

type InvalidRequestBodyHandler struct{}

func (h InvalidRequestBodyHandler) ServeHTTP(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusMethodNotAllowed)
	response.Write(responses.InvalidRequestBodyJSON)
}

type InternalServerErrorHandler struct{}

func (h InternalServerErrorHandler) ServeHTTP(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusInternalServerError)
	response.Write(responses.InternalServerErrorJSON)
}

type CustomErrorHandler struct{
	Error error 
}

func (h CustomErrorHandler) ServeHTTP(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusInternalServerError)
	errorResponse := responses.NewErrorResponse(h.Error)
	errorResponseJSON, err := json.Marshal(errorResponse)
	if err != nil {
		response.Write(responses.InternalServerErrorJSON)
		return
	}
	response.Write(errorResponseJSON)
}
