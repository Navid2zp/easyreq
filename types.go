package easyreq

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	URL              string
	Headers          map[string]string
	Data             []byte
	DataReader       io.Reader
	Method           string
	RequestDataType  string
	ResponseDataType string
	SaveResponseTo   interface{}
	Proxy            *url.URL
}

type RequestResponse struct {
	Response *http.Response
	Header   http.Header
}

type DownloadResult struct {
	BytesCopied  int64
	DownloadTime time.Duration
}
