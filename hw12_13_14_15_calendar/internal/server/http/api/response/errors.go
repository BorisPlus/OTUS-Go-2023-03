package api

import (
	"encoding/json"
)

var ResponseInvalidRequestBody = NewErrorResponseByString("invalid HTTP request body")

var InvalidRequestBodyJSON, _ = json.Marshal(ResponseInvalidRequestBody)

var ResponseInternalServerError = NewErrorResponseByString("internal server error")

var InternalServerErrorJSON, _ = json.Marshal(ResponseInternalServerError)

var ResponseInvalidHTTPMethodForUrlPath = NewErrorResponseByString("invalid HTTP method")

var InvalidHTTPMethodJSON, _ = json.Marshal(ResponseInvalidHTTPMethodForUrlPath)
