# easyreq
Make HTTP requests in go as easy as possible!

This package uses builtin ```net/http``` to make requests.

`json` and `xml` package will be used to unmarshal these types.

### Install
```
go get github.com/Navid2zp/easyreq
```

### Usage

```go
import (
	"github.com/Navid2zp/easyreq"
)


func main() {
    ereq = easyreq.NewRequest("GET", "https://site.com/api")
    resp, err := ereq.Make()
    if err != nil {
        fmt.Println(err)
    }
    defer resp.CloseBody()
}
```

##### Unmarshaling
You can unmarshal the response directly:

```go
type MyData struct {
	Name     string `json:"name" xml:"name"`
	LastName string `json:"last_name" xml:"last_name"`
	Github   string `json:"github" xml:"github"`
}
var stringResult string
var result MyData

err := resp.ToJson(&result)
// or xml:
err := resp.ToXML(&result)
// Or string:
err := resp.ToString(&stringResult)
```

##### Posting Data
You pass any type of data to be posted with the request.

###### Byte:
```go
ereq.SetData([]byte("this is my data"))
```

###### Reader:
```go
ereq.SetDataReader(strings.NewReader("this is my data"))
```


###### String:
```go
ereq.SetStringData("this is my data")
```

###### JSON/XML:
```go
sendData := MyData{
      Name: "Navid",
      LastName: "Zarepak",
      Github: "Navid2zp",
    }

err := ereq.SetJsonData(sendData)
err := ereq.SetXMLData(sendData)
```

##### Request header

```go
myHeaders := map[string]string{"myHeader": "Header Value"}
ereq.SetHeaders(myHeaders)

// Add a header:
ereq.AddHeader("NewHeaderKey", "Value")
```

##### Download Files
You can download files directly by calling ```DownloadAsFile``` method on a request response and providing a path to save the file.
```go
result, err := resp.DownloadAsFile("myfile.zip")

fmt.Println("Bytes copied:", result.BytesCopied)
fmt.Println("Download time:", result.DownloadTime)
```

##### Proxy


```go
ereq.SetHttpProxy("http://<PROXY_ADDRESS>:<PROXY_PORT>")
```


##### Request Shortcuts

```go
response, err := easyreq.Get("http://site.com")
response, err := easyreq.Post("http://site.com", []byte("my data"))
response, err := easyreq.Put("http://site.com", []byte("my data"))
response, err := easyreq.Patch("http://site.com", []byte("my data"))
response, err := easyreq.Delete("http://site.com")
```

##### Response Methods

```go
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

// Close response body
resp.CloseBody()

```

You can also access the original response returned by ```http``` package by calling ```resp.Response```. (Will be a pointer to original response)


License
----

MIT
