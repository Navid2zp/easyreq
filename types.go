package easyreq

import (
	"net/http"
	"time"
)

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
	Header   http.Header
}

type DownloadResult struct {
	BytesCopied  int64
	DownloadTime time.Duration
}
