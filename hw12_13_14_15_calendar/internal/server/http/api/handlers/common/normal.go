package api

import (
	"net/http"

	responses "hw12_13_14_15_calendar/internal/server/http/api/response"
)

type ApiVersionHandler struct{}

func (h ApiVersionHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(responses.ResponseVersionJSON)
}