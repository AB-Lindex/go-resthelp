# Simplified wrapper for calling REST-APIs

## Installation

```bash
go get github.com/AB-Lindex/go-resthelp
```

## Examples

### Create a  helper

```go
helper,err := resthelp.New(
  resthelp.WithBaseURL("https://some-custom-api.com"),
  resthelp.WithHeader("Authentication", "Bearer " + token),
)
```

### Create a Request
  
```go
// Simple GET request with a query parameter
req,err := helper.Get("/users",
  resthelp.WithQuery("page", "2"),
)

// Simple POST with a JSON body
var body myStruct
req,err := helper.Post("/users",
  resthelp.WithJSON(&body),
)
```

### Execute a Request and handle the Response

```go
// Execute the request
resp,err := req.Do()
defer resp.Close()     // always close the response (this is nil-safe)
if !resp.IsOK() {      // this checks for 'err' or anything except 2xx
  panic(resp.Error())  // will either return 'err' or the status-text
}
var data myReponseStruct
err = resp.Parse(data) // will parse the response body into 'data' according to Content-Type
```