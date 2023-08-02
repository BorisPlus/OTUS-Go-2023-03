package api

import (
	"net/http"

	responses "hw12_13_14_15_calendar/internal/server/http/api/api_response"
)

type VersionHandler struct{}

func (h VersionHandler) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	apiResponse := responses.NewAPIResponse("api.version")
	apiResponse.Data = struct {
		Version string
	}{
		Version: "1.0.0",
	}
	apiResponseJSON, _ := apiResponse.MarshalJSON()
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(apiResponseJSON)
}
