# easyreq
Make HTTP requests in go as easy as possible!

This package uses builtin ```net/http``` to make requests.

`json` and `xml` package will be used to unmarshal these types.

For now there is no support for third party packages.

### Install
```
go get github.com/Navid2zp/easyreq
```

### Usage

```
import (
	"github.com/Navid2zp/easyreq"
)

type MyData struct {
	Name     string `json:"name" xml:"name"`
	LastName string `json:"last_name" xml:"last_name"`
	Github   string `json:"github" xml:"github"`
}

func main() {

    sendData := MyData{
      Name: "Navid",
      LastName: "Zarepak",
      Github: "Navid2zp",
    }
    
    result := MyData{}
    
    ereq := easyreq.Request{
		URL:              "https://site.com/api",
		Method:           "post",
		Data:             []byte(sendData),
		RequestDataType:  "json",
		ResponseDataType: "json",
		SaveResponseTo:   &result,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	resp, err := ereq.Make()
}
```

You can also make a request using ```Make()``` function:

```
res, err := easyreq.Make(
	"post", // Request method
	"https://site.com/api", // URL
	[]byte(sendData), // Data to send
	"json", // Request data type
	"json", // Response data type
	&result, map[string]string{"Content-Type": "application/json",}, // Headers
)
```

By providing ```ResponseDataType``` and ```SaveResponseTo``` package will try to parse the response to the given ```ResponseDataType``` and the result will be saved to where ```SaveResponseTo``` points.


#### ```easyreq.Request``` Args

|          Arg         |                            Description                      |          Type      |
|----------------------|-------------------------------------------------------------|--------------------|
|**URL**               | URL to send the request to.                                 | string           
|**Method**            | Request method (default is get)                             | string
|**Data**              | Data to send with request                                   | []byte
|**RequestDataType**   | Data type for request (json, xml, string, ...)              | string
|**ResponseDataType**  | Data type for response (Will be used to parse the response) | string
|**SaveResponseTo**    | Where parsed data for response should be saved to           | pointer
|**Headers**           | Request headers                                             | map[string]string



#### Methods

These methods can be called on ```Make()``` response which is a ```easyreq.RequestResponse``` Type containing the original response from request sent with ```net/http``` package.

```
// A pointer to your struct to save the data to.
// Uses 'json' package to unmarshal the response body.
// Returns error if anything goes wrong.
err := resp.ToJson(&result)

// A pointer to your struct to save the data to.
// Uses 'xml' package to unmarshal the response body.
// Returns error if anything goes wrong.
err := resp.ToXML(&result)

// A pointer to a string.
// Response body will be converted to string and saved to the given pointer.
// Returns error if anything goes wrong.
err := resp.ToString(&result)
```

Other useful methods:
```
// Returns status code (200, ...)
// Type: int
resp.StatusCode()

// Returns status code (200 OK) 
// Type: string
resp.Status()

// Returns response headers
// Type: http.Header
resp.Headers()

// Retunrs response body
// Type: io.ReadCloser
resp.Body()

// Reads response body into bytes using `ioutil` and returns the result.
Type: []byte, error
resp.ReadBody()
```

You can also access the original response returned by ```http``` package by calling ```resp.Response```. (Will be a pointer to original response)


License
----

MIT
