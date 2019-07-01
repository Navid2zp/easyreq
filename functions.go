package easyreq

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

// Returns request method.
// Default is 'get' if you don't set one.
func (r *Request) getMethod() string {
	if r.Method == "" {
		return "get"
	}
	return r.Method
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

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	response := RequestResponse{Response: res}

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

// Unmarshals json response and saves into the given pointer
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