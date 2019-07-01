package easyreq

import "net/http"

type Request struct {
	URL              string
	Headers          map[string]string
	Data             []byte
	Method           string
	RequestDataType  string
	ResponseDataType string
	SaveResponseTo   interface{}
}

type RequestResponse struct {
	Response *http.Response
}
