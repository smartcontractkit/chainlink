[![Go Report Card](https://goreportcard.com/badge/github.com/ybbus/jsonrpc)](https://goreportcard.com/report/github.com/ybbus/jsonrpc)
[![GoDoc](https://godoc.org/github.com/ybbus/jsonrpc?status.svg)](https://godoc.org/github.com/ybbus/jsonrpc)
[![GitHub license](https://img.shields.io/github/license/mashape/apistatus.svg)]()

# JSON-RPC 2.0 Client for golang
A go implementation of an rpc client using json as data format over http.
The implementation is based on the JSON-RPC 2.0 specification: http://www.jsonrpc.org/specification

Supports:
- requests with arbitrary parameters
- convenient response retrieval
- batch requests
- custom http client (e.g. proxy, tls config)
- custom headers (e.g. basic auth)

## Installation

```sh
go get -u github.com/ybbus/jsonrpc
```

## Getting started
Let's say we want to retrieve a person struct with a specific id using rpc-json over http.
Then we want to save this person after we changed a property.
(Error handling is omitted here)

```go
type Person struct {
    Id   int `json:"id"`
    Name string `json:"name"`
    Age  int `json:"age"`
}

func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")

    var person *Person
    rpcClient.CallFor(&person, "getPersonById", 4711)

    person.Age = 33
    rpcClient.Call("updatePerson", person)
}
```

## In detail

### Generating rpc-json requests

Let's start by executing a simple json-rpc http call:
In production code: Always make sure to check err != nil first!

This calls generate and send a valid rpc-json object. (see: http://www.jsonrpc.org/specification#request_object)

```go
func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")
    rpcClient.Call("getDate")
    // generates body: {"method":"getDate","id":1,"jsonrpc":"2.0"}
}
```

Call a function with parameter:

```go
func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")
    rpcClient.Call("addNumbers", 1, 2)
    // generates body: {"method":"addNumbers","params":[1,2],"id":1,"jsonrpc":"2.0"}
}
```

Call a function with arbitrary parameters:

```go
func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")
    rpcClient.Call("createPerson", "Alex", 33, "Germany")
    // generates body: {"method":"createPerson","params":["Alex",33,"Germany"],"id":1,"jsonrpc":"2.0"}
}
```

Call a function providing custom data structures as parameters:

```go
type Person struct {
  Name    string `json:"name"`
  Age     int `json:"age"`
  Country string `json:"country"`
}
func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")
    rpcClient.Call("createPerson", &Person{"Alex", 33, "Germany"})
    // generates body: {"jsonrpc":"2.0","method":"createPerson","params":{"name":"Alex","age":33,"country":"Germany"},"id":1}
}
```

Complex example:

```go
type Person struct {
  Name    string `json:"name"`
  Age     int `json:"age"`
  Country string `json:"country"`
}
func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")
    rpcClient.Call("createPersonsWithRole", &Person{"Alex", 33, "Germany"}, &Person{"Barney", 38, "Germany"}, []string{"Admin", "User"})
    // generates body: {"jsonrpc":"2.0","method":"createPersonsWithRole","params":[{"name":"Alex","age":33,"country":"Germany"},{"name":"Barney","age":38,"country":"Germany"},["Admin","User"]],"id":1}
}
```

Some examples and resulting JSON-RPC objects:

```go
rpcClient.Call("missingParam")
{"method":"missingParam"}

rpcClient.Call("nullParam", nil)
{"method":"nullParam","params":[null]}

rpcClient.Call("boolParam", true)
{"method":"boolParam","params":[true]}

rpcClient.Call("boolParams", true, false, true)
{"method":"boolParams","params":[true,false,true]}

rpcClient.Call("stringParam", "Alex")
{"method":"stringParam","params":["Alex"]}

rpcClient.Call("stringParams", "JSON", "RPC")
{"method":"stringParams","params":["JSON","RPC"]}

rpcClient.Call("numberParam", 123)
{"method":"numberParam","params":[123]}

rpcClient.Call("numberParams", 123, 321)
{"method":"numberParams","params":[123,321]}

rpcClient.Call("floatParam", 1.23)
{"method":"floatParam","params":[1.23]}

rpcClient.Call("floatParams", 1.23, 3.21)
{"method":"floatParams","params":[1.23,3.21]}

rpcClient.Call("manyParams", "Alex", 35, true, nil, 2.34)
{"method":"manyParams","params":["Alex",35,true,null,2.34]}

rpcClient.Call("singlePointerToStruct", &person)
{"method":"singlePointerToStruct","params":{"name":"Alex","age":35,"country":"Germany"}}

rpcClient.Call("multipleStructs", &person, &drink)
{"method":"multipleStructs","params":[{"name":"Alex","age":35,"country":"Germany"},{"name":"Cuba Libre","ingredients":["rum","cola"]}]}

rpcClient.Call("singleStructInArray", []*Person{&person})
{"method":"singleStructInArray","params":[{"name":"Alex","age":35,"country":"Germany"}]}

rpcClient.Call("namedParameters", map[string]interface{}{
	"name": "Alex",
	"age":  35,
})
{"method":"namedParameters","params":{"age":35,"name":"Alex"}}

rpcClient.Call("anonymousStruct", struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}{"Alex", 33})
{"method":"anonymousStructWithTags","params":{"name":"Alex","age":33}}

rpcClient.Call("structWithNullField", struct {
	Name    string  `json:"name"`
	Address *string `json:"address"`
}{"Alex", nil})
{"method":"structWithNullField","params":{"name":"Alex","address":null}}
```

### Working with rpc-json responses


Before working with the response object, make sure to check err != nil.
Also keep in mind that the json-rpc result field can be nil even on success.

```go
func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")
    response, err := rpcClient.Call("addNumbers", 1, 2)
    if err != nil {
      // error handling goes here e.g. network / http error
    }
}
```

If an http error occurred, maybe you are interested in the error code (403 etc.)
```go
func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")
    response, err := rpcClient.Call("addNumbers", 1, 2)

    switch e := err.(type) {
      case nil: // if error is nil, do nothing
      case *HTTPError:
        // use e.Code here
        return
      default:
        // any other error
        return
    }

    // no error, go on...
}
```

The next thing you have to check is if an rpc-json protocol error occurred. This is done by checking if the Error field in the rpc-response != nil:
(see: http://www.jsonrpc.org/specification#error_object)

```go
func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")
    response, err := rpcClient.Call("addNumbers", 1, 2)
    if err != nil {
        //error handling goes here
    }

    if response.Error != nil {
        // rpc error handling goes here
        // check response.Error.Code, response.Error.Message and optional response.Error.Data
    }
}
```

After making sure that no errors occurred you can now examine the RPCResponse object.
When executing a json-rpc request, most of the time you will be interested in the "result"-property of the returned json-rpc response object.
(see: http://www.jsonrpc.org/specification#response_object)
The library provides some helper functions to retrieve the result in the data format you are interested in.
Again: check for err != nil here to be sure the expected type was provided in the response and could be parsed.

```go
func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")
    response, _ := rpcClient.Call("addNumbers", 1, 2)

    result, err := response.GetInt()
    if err != nil {
        // result cannot be unmarshalled as integer
    }

    // helpers provided for all primitive types:
    response.GetInt()
    response.GetFloat()
    response.GetString()
    response.GetBool()
}
```

Retrieving arrays and objects is also very simple:

```go
// json annotations are only required to transform the structure back to json
type Person struct {
    Id   int `json:"id"`
    Name string `json:"name"`
    Age  int `json:"age"`
}

func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")
    response, _ := rpcClient.Call("getPersonById", 123)

    var person *Person
    err := response.GetObject(&person) // expects a rpc-object result value like: {"id": 123, "name": "alex", "age": 33}
    if err != nil || person == nil {
        // some error on json unmarshal level or json result field was null
    }

    fmt.Println(person.Name)

    // we can also set default values if they are missing from the result, or result == null:
    person2 := &Person{
        Id: 0,
        Name: "<empty>",
        Age: -1,
    }
    err := response.GetObject(&person2) // expects a rpc-object result value like: {"id": 123, "name": "alex", "age": 33}
    if err != nil || person2 == nil {
        // some error on json unmarshal level or json result field was null
    }

    fmt.Println(person2.Name) // prints "<empty>" if "name" field was missing in result-json
}
```

Retrieving arrays:

```go
func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")
    response, _ := rpcClient.Call("getRandomNumbers", 10)

    rndNumbers := []int{}
    err := response.GetObject(&rndNumbers) // expects a rpc-object result value like: [10, 188, 14, 3]
    if err != nil {
        // do error handling
    }

    for _, num := range rndNumbers {
        fmt.Printf("%v\n", num)
    }
}
```

### Using convenient function CallFor()
A very handy way to quickly invoke methods and retrieve results is by using CallFor()

You can directly provide an object where the result should be stored. Be sure to provide it be reference.
An error is returned if:
- there was an network / http error
- RPCError object is not nil (err can be casted to this object)
- rpc result could not be parsed into provided object

One of te above examples could look like this:

```go
// json annotations are only required to transform the structure back to json
type Person struct {
    Id   int `json:"id"`
    Name string `json:"name"`
    Age  int `json:"age"`
}

func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")

    var person *Person
    err := rpcClient.CallFor(&person, "getPersonById", 123)

    if err != nil || person == nil {
      // handle error
    }

    fmt.Println(person.Name)
}
```

Most of the time it is ok to check if a struct field is 0, empty string "" etc. to check if it was provided by the json rpc response.
But if you want to be sure that a JSON-RPC response field was missing or not, you should use pointers to the fields.
This is just a single example since all this Unmarshaling is standard go json functionality, exactly as if you would call json.Unmarshal(rpcResponse.ResultAsByteArray, &objectToStoreResult)

```
type Person struct {
    Id   *int    `json:"id"`
    Name *string `json:"name"`
    Age  *int    `json:"age"`
}

func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")

    var person *Person
    err := rpcClient.CallFor(&person, "getPersonById", 123)

    if err != nil || person == nil {
      // handle error
    }

    if person.Name == nil {
      // json rpc response did not provide a field "name" in the result object
    }
}
```

### Using RPC Batch Requests

You can send multiple RPC-Requests in one single HTTP request using RPC Batch Requests.

```
func main() {
    rpcClient := jsonrpc.NewClient("http://my-rpc-service:8080/rpc")

    response, _ := rpcClient.CallBatch(RPCRequests{
      NewRequest("myMethod1", 1, 2, 3),
      NewRequest("anotherMethod", "Alex", 35, true),
      NewRequest("myMethod2", &Person{
        Name: "Emmy",
        Age: 4,
      }),
    })
}
```

Keep the following in mind:
- the request / response id's are important to map the requests to the responses. CallBatch() automatically sets the ids to requests[i].ID == i
- the response can be provided in an unordered and maybe incomplete form
- when you want to set the id yourself use, CallRaw()

There are some helper methods for batch request results:
```
func main() {
    // [...]

    result.HasErrors() // returns true if one of the rpc response objects has Error field != nil
    resultMap := result.AsMap() // returns a map for easier retrieval of requests

    if response123, ok := resultMap[123]; ok {
      // response object with id 123 exists, use it here
      // response123.ID == 123
      response123.GetObjectAs(&person)
      // ...
    }

}
```

### Raw functions
There are also Raw function calls. Consider the non Raw functions first, unless you know what you are doing.
You can create invalid json rpc requests and have to take care of id's etc. yourself.
Also check documentation of Params() for raw requests.

### Custom Headers, Basic authentication

If the rpc-service is running behind a basic authentication you can easily set the Authorization header:

```go
func main() {
    rpcClient := jsonrpc.NewClientWithOpts("http://my-rpc-service:8080/rpc", &jsonrpc.RPCClientOpts{
   		CustomHeaders: map[string]string{
   			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte("myUser"+":"+"mySecret")),
   		},
   	})
    response, _ := rpcClient.Call("addNumbers", 1, 2) // send with Authorization-Header
}
```

### Using oauth

Using oauth is also easy, e.g. with clientID and clientSecret authentication

```go
func main() {
		credentials := clientcredentials.Config{
    		ClientID:     "myID",
    		ClientSecret: "mySecret",
    		TokenURL:     "http://mytokenurl",
    	}

    	rpcClient := jsonrpc.NewClientWithOpts("http://my-rpc-service:8080/rpc", &jsonrpc.RPCClientOpts{
    		HTTPClient: credentials.Client(context.Background()),
    	})

	// requests now retrieve and use an oauth token
}
```

### Set a custom httpClient

If you have some special needs on the http.Client of the standard go library, just provide your own one.
For example to use a proxy when executing json-rpc calls:

```go
func main() {
	proxyURL, _ := url.Parse("http://proxy:8080")
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}

	httpClient := &http.Client{
		Transport: transport,
	}

	rpcClient := jsonrpc.NewClientWithOpts("http://my-rpc-service:8080/rpc", &jsonrpc.RPCClientOpts{
		HTTPClient: httpClient,
	})

	// requests now use proxy
}
```
