package easyreq

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"
)

// Returns request method.
// Default is 'get' if you don't set one.
func (r *Request) getMethod() string {
	if r.Method == "" {
		return "GET"
	}
	return strings.ToUpper(r.Method)
}

// getClient returns the http client used to make the request
func (r *Request) getClient() *http.Client {
	if r.Proxy != nil {
		return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(r.Proxy)}}
	}
	return http.DefaultClient
}

// Send the final request
func (r *Request) Make() (*RequestResponse, error) {
	err := r.validate()
	if err != nil {
		return nil, err
	}
	data, err := r.getData()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(r.getMethod(), r.URL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	res, err := r.getClient().Do(req)
	if err != nil {
		return nil, err
	}

	// Body won't be closed anymore
	// You should call CloseBody() on response
	//defer res.Body.Close()

	response := RequestResponse{Response: res, Header: res.Header}

	// Sends response to parseTo
	// Only when you set a SaveResponseTo
	if r.SaveResponseTo != nil {
		err := response.parseTo(r.ResponseDataType, r.SaveResponseTo)
		return &response, err
	}
	// Returns a simple RequestResponse type
	// So you can use response methods such as .json() later.
	return &response, nil
}

// NewRequest creates a new request
func NewRequest(method, url string) *Request {
	return &Request{
		URL:    url,
		Method:method,
	}
}

// Returns request data based on information provided
// It will convert data type if it doesn't match what has been set as RequestDataType
func (r *Request) getData() ([]byte, error) {
	if r.RequestDataType != "" {
		switch strings.ToLower(r.RequestDataType) {
		case "json":
			if !IsJson(r.Data) {
				jsonData, err := json.Marshal(r.Data)
				return jsonData, err
			}
			return r.Data, nil
		case "xml":
			if !IsXML(r.Data) {
				xmlData, err := xml.Marshal(r.Data)
				return xmlData, err
			}
			return r.Data, nil
		default:
			return r.Data, nil
		}
	}
	return r.Data, nil
}

func Get(url string) (*RequestResponse, error) {
	req := Request{Method: "GET", URL: url}
	return req.Make()
}

func Delete(url string) (*RequestResponse, error) {
	req := Request{Method: "DELETE", URL: url}
	return req.Make()
}

func Post(url string, data []byte) (*RequestResponse, error) {
	req := Request{Method: "POST", URL: url, Data: data}
	return req.Make()
}

func Put(url string, data []byte) (*RequestResponse, error) {
	req := Request{Method: "PUT", URL: url, Data: data}
	return req.Make()
}

func Patch(url string, data []byte) (*RequestResponse, error) {
	req := Request{Method: "PATCH", URL: url, Data: data}
	return req.Make()
}

// Converts and saves to `r.ResponseDataType`
// Using RequestResponse methods.
func (r *RequestResponse) parseTo(responseDataType string, saveResponseTo interface{}) error {
	if responseDataType != "" && saveResponseTo != nil {
		switch strings.ToLower(responseDataType) {
		case "json":
			err := r.ToJson(saveResponseTo)
			return err
		case "string", "text", "html":
			err := r.ToString(saveResponseTo)
			return err
		case "xml":
			err := r.ToXML(saveResponseTo)
			return err
		default:
			return errors.New("unknown type")
		}
	}
	return errors.New("no ResponseDataType provided ")
}

// Either call this function to create a request
// Or call Request type directly and call Make().
// Both will do the same.
func Make(method, url string, data []byte, dataType, responseType string, saveTo interface{}, headers map[string]string) (*RequestResponse, error) {
	req := Request{
		Method:           method,
		URL:              url,
		Data:             data,
		RequestDataType:  dataType,
		ResponseDataType: responseType,
		SaveResponseTo:   saveTo,
		Headers:          headers,
	}
	res, err := req.Make()
	return res, err
}

// Unmarshals json response and saves it into the given pointer
func (r *RequestResponse) ToJson(saveTo interface{}) error {
	body, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, saveTo)
	return err
}

// Unmarshals xml response and saves into the given pointer
func (r *RequestResponse) ToXML(saveTo interface{}) error {
	body, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(body, saveTo)
	return err
}

// converts response to string and saves into the given pointer
func (r *RequestResponse) ToString(saveTo interface{}) error {
	err := validateSaveToInterface(saveTo)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		return err
	}
	// If `saveTo` doesn't point to string
	_, ok := saveTo.(*string)
	if !ok {
		return errors.New("non-string pointer for 'ResponseDataType=string'")
	}
	st := reflect.ValueOf(saveTo)
	st.Elem().SetString(string(body))
	return nil
}

// Returns response status code (int: 200)
func (r *RequestResponse) StatusCode() int {
	return r.Response.StatusCode
}

// Returns status (string: 200 OK)
func (r *RequestResponse) Status() string {
	return r.Response.Status
}

// Returns response headers
func (r *RequestResponse) Headers() http.Header {
	return r.Response.Header
}

// Returns response body
func (r *RequestResponse) Body() io.ReadCloser {
	return r.Response.Body
}

// Reads response body into bytes and returns the result
func (r *RequestResponse) ReadBody() ([]byte, error) {
	return ioutil.ReadAll(r.Response.Body)
}

// Close response body
func (r *RequestResponse) CloseBody() error {
	return r.Response.Body.Close()
}

// Downloads the body and saves it into the given path
func (r *RequestResponse) DownloadAsFile(fileName string) (*DownloadResult, error) {
	st := time.Now()
	out, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	n, err := io.Copy(out, r.Response.Body)
	return &DownloadResult{
		BytesCopied:  n,
		DownloadTime: time.Now().Sub(st),
	}, err
}

// SetHeaders sets request headers
func (r *Request) SetHeaders(headers map[string]string) {
	r.Headers = headers
}

func (r *Request) SetProxy(proxy string) error {
	proxyURL, err := url.Parse(proxy)
	r.Proxy = proxyURL
	return err
}
