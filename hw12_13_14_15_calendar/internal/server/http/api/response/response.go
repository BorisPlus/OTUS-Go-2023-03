package api

import "fmt"

// import (
// 	models "hw12_13_14_15_calendar/internal/models"
// )

type Response struct {
	Error error `json:"error"`
	Data  any   `json:"data"`
}

func NewErrorResponse(Error error) *Response {
	r := new(Response)
	r.Error = Error
	return r
}

func NewErrorResponseByString(Error string, a... any) *Response {
	r := new(Response)
	r.Error = fmt.Errorf(Error, a...)
	return r
}

