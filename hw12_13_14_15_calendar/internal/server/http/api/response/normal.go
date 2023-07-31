package api

import (
	"encoding/json"
)

var ResponseVersion = Response{
	nil,
	struct {
		Version string
	}{
		Version: "1.0.0",
	},
}

var ResponseVersionJSON, _ = json.Marshal(ResponseVersion)
