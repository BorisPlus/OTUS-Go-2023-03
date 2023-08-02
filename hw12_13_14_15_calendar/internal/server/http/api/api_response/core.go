package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

type DataItem struct {
	Item any `json:"item"`
}
type DataItems struct {
	Items any `json:"items"`
}

type APIResponse struct {
	ApiMethod string `json:"method"`
	Error     error  `json:"error"`
	Data      any    `json:"data"`
}

// https://github.com/golang/go/issues/5161

func (r APIResponse) MarshalJSON() ([]byte, error) {
	var errorMsg string
	if r.Error != nil {
		errorMsg = r.Error.Error()
	}
	anon := struct {
		ApiMethod string `json:"method"`
		Error     string `json:"error"`
		Data      any    `json:"data"`
	}{
		ApiMethod: r.ApiMethod,
		Error:     errorMsg,
		Data:      r.Data,
	}
	return json.Marshal(anon)
}

func NewAPIResponse(ApiMethod string) *APIResponse {
	r := new(APIResponse)
	r.ApiMethod = ApiMethod
	return r
}

func NewErrorAPIResponse(ApiMethod string, Error error) *APIResponse {
	r := new(APIResponse)
	r.ApiMethod = ApiMethod
	r.Error = Error
	return r
}

func NewErrorAPIResponseByString(ApiMethod string, format string, a ...any) *APIResponse {
	r := new(APIResponse)
	r.ApiMethod = ApiMethod
	r.Error = fmt.Errorf(format, a...)
	return r
}

func InternalServerError(ApiMethod string) *APIResponse {
	r := new(APIResponse)
	r.ApiMethod = ApiMethod
	r.Error = fmt.Errorf("internal server error")
	return r
}

func InvalidHTTPMethod(ApiMethod string) *APIResponse {
	r := new(APIResponse)
	r.ApiMethod = ApiMethod
	// r.Error = fmt.Errorf("invalid HTTP method")
	r.Error = errors.New("invalid HTTP method")
	return r
}

func InvalidRequestBody(ApiMethod string) *APIResponse {
	r := new(APIResponse)
	r.ApiMethod = ApiMethod
	r.Error = fmt.Errorf("invalid HTTP request body")
	return r
}
