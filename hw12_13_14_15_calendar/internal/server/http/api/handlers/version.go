package api

import (
	"net/http"
)

type ApiVersionHandler struct{}

func (h ApiVersionHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("Version 1.0.0"))
}
