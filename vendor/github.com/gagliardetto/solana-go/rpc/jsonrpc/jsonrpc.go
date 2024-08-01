// Package jsonrpc provides a JSON-RPC 2.0 client that sends JSON-RPC requests and receives JSON-RPC responses using HTTP.
package jsonrpc

import (
	"bytes"
	"context"
	stdjson "encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	jsonrpcVersion = "2.0"
)

// RPCClient sends JSON-RPC requests over HTTP to the provided JSON-RPC backend.
//
// RPCClient is created using the factory function NewClient().
type RPCClient interface {
	// Call is used to send a JSON-RPC request to the server endpoint.
	//
	// The spec states, that params can only be an array or an object, no primitive values.
	// So there are a few simple rules to notice:
	//
	// 1. no params: params field is omitted. e.g. Call("getinfo")
	//
	// 2. single params primitive value: value is wrapped in array. e.g. Call("getByID", 1423)
	//
	// 3. single params value array or object: value is unchanged. e.g. Call("storePerson", &Person{Name: "Alex"})
	//
	// 4. multiple params values: always wrapped in array. e.g. Call("setDetails", "Alex, 35, "Germany", true)
	//
	// Examples:
	//   Call("getinfo") -> {"method": "getinfo"}
	//   Call("getPersonId", 123) -> {"method": "getPersonId", "params": [123]}
	//   Call("setName", "Alex") -> {"method": "setName", "params": ["Alex"]}
	//   Call("setMale", true) -> {"method": "setMale", "params": [true]}
	//   Call("setNumbers", []int{1, 2, 3}) -> {"method": "setNumbers", "params": [1, 2, 3]}
	//   Call("setNumbers", 1, 2, 3) -> {"method": "setNumbers", "params": [1, 2, 3]}
	//   Call("savePerson", &Person{Name: "Alex", Age: 35}) -> {"method": "savePerson", "params": {"name": "Alex", "age": 35}}
	//   Call("setPersonDetails", "Alex", 35, "Germany") -> {"method": "setPersonDetails", "params": ["Alex", 35, "Germany"}}
	//
	// for more information, see the examples or the unit tests
	Call(ctx context.Context, method string, params ...interface{}) (*RPCResponse, error)

	// CallRaw is like Call() but without magic in the requests.Params field.
	// The RPCRequest object is sent exactly as you provide it.
	// See docs: NewRequest, RPCRequest, Params()
	//
	// It is recommended to first consider Call() and CallFor()
	CallRaw(ctx context.Context, request *RPCRequest) (*RPCResponse, error)

	// CallFor is a very handy function to send a JSON-RPC request to the server endpoint
	// and directly specify an object to store the response.
	//
	// out: will store the unmarshaled object, if request was successful.
	// should always be provided by references. can be nil even on success.
	// the behaviour is the same as expected from json.Unmarshal()
	//
	// method and params: see Call() function
	//
	// if the request was not successful (network, http error) or the rpc response returns an error,
	// an error is returned. if it was an JSON-RPC error it can be casted
	// to *RPCError.
	//
	CallFor(ctx context.Context, out interface{}, method string, params ...interface{}) error

	// CallBatch invokes a list of RPCRequests in a single batch request.
	//
	// Most convenient is to use the following form:
	// CallBatch(RPCRequests{
	//   Batch("myMethod1", 1, 2, 3),
	//   Batch("myMethod2), "Test"),
	// })
	//
	// You can create the []*RPCRequest array yourself, but it is not recommended and you should notice the following:
	// - field Params is sent as provided, so Params: 2 forms an invalid json (correct would be Params: []int{2})
	// - you can use the helper function Params(1, 2, 3) to use the same format as in Call()
	// - field JSONRPC is overwritten and set to value: "2.0"
	// - field ID is overwritten and set incrementally and maps to the array position (e.g. requests[5].ID == 5)
	//
	//
	// Returns RPCResponses that is of type []*RPCResponse
	// - note that a list of RPCResponses can be received unordered so it can happen that: responses[i] != responses[i].ID
	// - RPCPersponses is enriched with helper functions e.g.: responses.HasError() returns  true if one of the responses holds an RPCError
	CallBatch(ctx context.Context, requests RPCRequests) (RPCResponses, error)

	// CallBatchRaw invokes a list of RPCRequests in a single batch request.
	// It sends the RPCRequests parameter is it passed (no magic, no id autoincrement).
	//
	// Consider to use CallBatch() instead except you have some good reason not to.
	//
	// CallBatchRaw(RPCRequests{
	//   &RPCRequest{
	//     ID: 123,            // this won't be replaced in CallBatchRaw
	//     JSONRPC: "wrong",   // this won't be replaced in CallBatchRaw
	//     Method: "myMethod1",
	//     Params: []int{1},   // there is no magic, be sure to only use array or object
	//   },
	//   &RPCRequest{
	//     ID: 612,
	//     JSONRPC: "2.0",
	//     Method: "myMethod2",
	//     Params: Params("Alex", 35, true), // you can use helper function Params() (see doc)
	//   },
	// })
	//
	// Returns RPCResponses that is of type []*RPCResponse
	// - note that a list of RPCResponses can be received unordered
	// - the id's must be mapped against the id's you provided
	// - RPCPersponses is enriched with helper functions e.g.: responses.HasError() returns  true if one of the responses holds an RPCError
	CallBatchRaw(ctx context.Context, requests RPCRequests) (RPCResponses, error)

	CallForInto(ctx context.Context, out interface{}, method string, params []interface{}) error
	CallWithCallback(ctx context.Context, method string, params []interface{}, callback func(*http.Request, *http.Response) error) error
	Close() error
}

// RPCRequest represents a JSON-RPC request object.
//
// Method: string containing the method to be invoked
//
// Params: can be nil. if not must be an json array or object
//
// ID: may always set to 1 for single requests. Should be unique for every request in one batch request.
//
// JSONRPC: must always be set to "2.0" for JSON-RPC version 2.0
//
// See: http://www.jsonrpc.org/specification#request_object
//
// Most of the time you shouldn't create the RPCRequest object yourself.
// The following functions do that for you:
// Call(), CallFor(), NewRequest()
//
// If you want to create it yourself (e.g. in batch or CallRaw()), consider using Params().
// Params() is a helper function that uses the same parameter syntax as Call().
//
// e.g. to manually create an RPCRequest object:
//
//	request := &RPCRequest{
//	  Method: "myMethod",
//	  Params: Params("Alex", 35, true),
//	}
//
// If you know what you are doing you can omit the Params() call to avoid some reflection but potentially create incorrect rpc requests:
//
//	request := &RPCRequest{
//	  Method: "myMethod",
//	  Params: 2, <-- invalid since a single primitive value must be wrapped in an array --> no magic without Params()
//	}
//
// correct:
//
//	request := &RPCRequest{
//	  Method: "myMethod",
//	  Params: []int{2}, <-- invalid since a single primitive value must be wrapped in an array
//	}
type RPCRequest struct {
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      int         `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
}

// NewRequest returns a new RPCRequest that can be created using the same convenient parameter syntax as Call()
//
// e.g. NewRequest("myMethod", "Alex", 35, true)
func NewRequest(method string, params ...interface{}) *RPCRequest {
	request := &RPCRequest{
		Method:  method,
		Params:  Params(params...),
		JSONRPC: jsonrpcVersion,
	}

	return request
}

// RPCResponse represents a JSON-RPC response object.
//
// Result: holds the result of the rpc call if no error occurred, nil otherwise. can be nil even on success.
//
// Error: holds an RPCError object if an error occurred. must be nil on success.
//
// ID: may always be 0 for single requests. is unique for each request in a batch call (see CallBatch())
//
// JSONRPC: must always be set to "2.0" for JSON-RPC version 2.0
//
// See: http://www.jsonrpc.org/specification#response_object
type RPCResponse struct {
	JSONRPC string             `json:"jsonrpc"`
	Result  stdjson.RawMessage `json:"result,omitempty"`
	Error   *RPCError          `json:"error,omitempty"`
	ID      int                `json:"id"`
}

// RPCError represents a JSON-RPC error object if an RPC error occurred.
//
// Code: holds the error code
//
// Message: holds a short error message
//
// Data: holds additional error data, may be nil
//
// See: http://www.jsonrpc.org/specification#error_object
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var spewConf = spew.ConfigState{
	Indent:                " ",
	DisableMethods:        true,
	DisablePointerMethods: true,
	SortKeys:              true,
}

// Error function is provided to be used as error object.
func (e *RPCError) Error() string {
	return spewConf.Sdump(e)
}

// HTTPError represents a error that occurred on HTTP level.
//
// An error of type HTTPError is returned when a HTTP error occurred (status code)
// and the body could not be parsed to a valid RPCResponse object that holds a RPCError.
//
// Otherwise a RPCResponse object is returned with a RPCError field that is not nil.
type HTTPError struct {
	Code int
	err  error
}

// HTTPClient is an abstraction for a HTTP client
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
	CloseIdleConnections()
}

func NewHTTPError(code int, err error) *HTTPError {
	return &HTTPError{
		Code: code,
		err:  err,
	}
}

// Error function is provided to be used as error object.
func (e *HTTPError) Error() string {
	return e.err.Error()
}

type rpcClient struct {
	endpoint      string
	httpClient    HTTPClient
	customHeaders map[string]string
}

// RPCClientOpts can be provided to NewClientWithOpts() to change configuration of RPCClient.
//
// HTTPClient: provide a custom http.Client (e.g. to set a proxy, or tls options)
//
// CustomHeaders: provide custom headers, e.g. to set BasicAuth
type RPCClientOpts struct {
	HTTPClient    HTTPClient
	CustomHeaders map[string]string
}

// RPCResponses is of type []*RPCResponse.
// This type is used to provide helper functions on the result list
type RPCResponses []*RPCResponse

// AsMap returns the responses as map with response id as key.
func (res RPCResponses) AsMap() map[int]*RPCResponse {
	resMap := make(map[int]*RPCResponse, 0)
	for _, r := range res {
		resMap[r.ID] = r
	}

	return resMap
}

// GetByID returns the response object of the given id, nil if it does not exist.
func (res RPCResponses) GetByID(id int) *RPCResponse {
	for _, r := range res {
		if r.ID == id {
			return r
		}
	}

	return nil
}

// HasError returns true if one of the response objects has Error field != nil
func (res RPCResponses) HasError() bool {
	for _, res := range res {
		if res.Error != nil {
			return true
		}
	}
	return false
}

// RPCRequests is of type []*RPCRequest.
// This type is used to provide helper functions on the request list
type RPCRequests []*RPCRequest

// NewClient returns a new RPCClient instance with default configuration.
//
// endpoint: JSON-RPC service URL to which JSON-RPC requests are sent.
func NewClient(endpoint string) RPCClient {
	return NewClientWithOpts(endpoint, nil)
}

// NewClientWithOpts returns a new RPCClient instance with custom configuration.
//
// endpoint: JSON-RPC service URL to which JSON-RPC requests are sent.
//
// opts: RPCClientOpts provide custom configuration
func NewClientWithOpts(endpoint string, opts *RPCClientOpts) RPCClient {
	rpcClient := &rpcClient{
		endpoint:      endpoint,
		httpClient:    &http.Client{},
		customHeaders: make(map[string]string),
	}

	if opts == nil {
		return rpcClient
	}

	if opts.HTTPClient != nil {
		rpcClient.httpClient = opts.HTTPClient
	}

	if opts.CustomHeaders != nil {
		for k, v := range opts.CustomHeaders {
			rpcClient.customHeaders[k] = v
		}
	}

	return rpcClient
}

func (client *rpcClient) Call(ctx context.Context, method string, params ...interface{}) (*RPCResponse, error) {
	request := &RPCRequest{
		Method:  method,
		Params:  Params(params...),
		JSONRPC: jsonrpcVersion,
	}

	return client.doCall(ctx, request)
}

func (client *rpcClient) Close() error {
	if client.httpClient != nil {
		client.httpClient.CloseIdleConnections()
	}
	return nil
}

func (client *rpcClient) CallForInto(
	ctx context.Context,
	out interface{},
	method string,
	params []interface{},
) error {
	request := &RPCRequest{
		Method:  method,
		JSONRPC: jsonrpcVersion,
	}

	if params != nil {
		request.Params = params
	}

	rpcResponse, err := client.doCall(ctx, request)
	if err != nil {
		return err
	}

	if rpcResponse.Error != nil {
		return rpcResponse.Error
	}

	return rpcResponse.GetObject(out)
}

func (client *rpcClient) CallWithCallback(
	ctx context.Context,
	method string,
	params []interface{},
	callback func(*http.Request, *http.Response) error,
) error {
	request := &RPCRequest{
		Method:  method,
		JSONRPC: jsonrpcVersion,
	}

	if params != nil {
		request.Params = params
	}

	return client.doCallWithCallbackOnHTTPResponse(
		ctx,
		request,
		callback,
	)
}

func (client *rpcClient) CallRaw(ctx context.Context, request *RPCRequest) (*RPCResponse, error) {
	return client.doCall(ctx, request)
}

func (client *rpcClient) CallFor(ctx context.Context, out interface{}, method string, params ...interface{}) error {
	rpcResponse, err := client.Call(ctx, method, params...)
	if err != nil {
		return err
	}

	if rpcResponse.Error != nil {
		return rpcResponse.Error
	}

	return rpcResponse.GetObject(out)
}

func (client *rpcClient) CallBatch(ctx context.Context, requests RPCRequests) (RPCResponses, error) {
	if len(requests) == 0 {
		return nil, errors.New("empty request list")
	}

	for i, req := range requests {
		req.ID = i
		req.JSONRPC = jsonrpcVersion
	}

	return client.doBatchCall(ctx, requests)
}

func (client *rpcClient) CallBatchRaw(ctx context.Context, requests RPCRequests) (RPCResponses, error) {
	if len(requests) == 0 {
		return nil, errors.New("empty request list")
	}

	return client.doBatchCall(ctx, requests)
}

func (client *rpcClient) newRequest(ctx context.Context, req interface{}) (*http.Request, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "POST", client.endpoint, bytes.NewReader(body))
	if err != nil {
		return request, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	// set default headers first, so that even content type and accept can be overwritten
	for k, v := range client.customHeaders {
		request.Header.Set(k, v)
	}

	return request, nil
}

func (client *rpcClient) doCall(
	ctx context.Context,
	RPCRequest *RPCRequest,
) (*RPCResponse, error) {
	var rpcResponse *RPCResponse
	err := client.doCallWithCallbackOnHTTPResponse(
		ctx,
		RPCRequest,
		func(httpRequest *http.Request, httpResponse *http.Response) error {
			decoder := json.NewDecoder(httpResponse.Body)
			decoder.DisallowUnknownFields()
			decoder.UseNumber()
			err := decoder.Decode(&rpcResponse)
			// parsing error
			if err != nil {
				// if we have some http error, return it
				if httpResponse.StatusCode >= 400 {
					return &HTTPError{
						Code: httpResponse.StatusCode,
						err:  fmt.Errorf("rpc call %v() on %v status code: %v. could not decode body to rpc response: %w", RPCRequest.Method, httpRequest.URL.String(), httpResponse.StatusCode, err),
					}
				}
				return fmt.Errorf("rpc call %v() on %v status code: %v. could not decode body to rpc response: %w", RPCRequest.Method, httpRequest.URL.String(), httpResponse.StatusCode, err)
			}

			// response body empty
			if rpcResponse == nil {
				// if we have some http error, return it
				if httpResponse.StatusCode >= 400 {
					return &HTTPError{
						Code: httpResponse.StatusCode,
						err:  fmt.Errorf("rpc call %v() on %v status code: %v. rpc response missing", RPCRequest.Method, httpRequest.URL.String(), httpResponse.StatusCode),
					}
				}
				return fmt.Errorf("rpc call %v() on %v status code: %v. rpc response missing", RPCRequest.Method, httpRequest.URL.String(), httpResponse.StatusCode)
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return rpcResponse, nil
}

func (client *rpcClient) doCallWithCallbackOnHTTPResponse(
	ctx context.Context,
	RPCRequest *RPCRequest,
	callback func(*http.Request, *http.Response) error,
) error {
	httpRequest, err := client.newRequest(ctx, RPCRequest)
	if err != nil {
		if httpRequest != nil {
			return fmt.Errorf("rpc call %v() on %v: %w", RPCRequest.Method, httpRequest.URL.String(), err)
		}
		return fmt.Errorf("rpc call %v(): %w", RPCRequest.Method, err)
	}
	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("rpc call %v() on %v: %w", RPCRequest.Method, httpRequest.URL.String(), err)
	}
	defer httpResponse.Body.Close()

	return callback(httpRequest, httpResponse)
}

func (client *rpcClient) doBatchCall(ctx context.Context, rpcRequest []*RPCRequest) ([]*RPCResponse, error) {
	httpRequest, err := client.newRequest(ctx, rpcRequest)
	if err != nil {
		if httpRequest != nil {
			return nil, fmt.Errorf("rpc batch call on %v: %w", httpRequest.URL.String(), err)
		}
		return nil, fmt.Errorf("rpc batch call: %w", err)
	}
	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("rpc batch call on %v: %w", httpRequest.URL.String(), err)
	}
	defer httpResponse.Body.Close()

	var rpcResponse RPCResponses
	decoder := json.NewDecoder(httpResponse.Body)
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	err = decoder.Decode(&rpcResponse)

	// parsing error
	if err != nil {
		// if we have some http error, return it
		if httpResponse.StatusCode >= 400 {
			return nil, &HTTPError{
				Code: httpResponse.StatusCode,
				err:  fmt.Errorf("rpc batch call on %v status code: %v. could not decode body to rpc response: %w", httpRequest.URL.String(), httpResponse.StatusCode, err),
			}
		}
		return nil, fmt.Errorf("rpc batch call on %v status code: %v. could not decode body to rpc response: %w", httpRequest.URL.String(), httpResponse.StatusCode, err)
	}

	// response body empty
	if rpcResponse == nil || len(rpcResponse) == 0 {
		// if we have some http error, return it
		if httpResponse.StatusCode >= 400 {
			return nil, &HTTPError{
				Code: httpResponse.StatusCode,
				err:  fmt.Errorf("rpc batch call on %v status code: %v. rpc response missing", httpRequest.URL.String(), httpResponse.StatusCode),
			}
		}
		return nil, fmt.Errorf("rpc batch call on %v status code: %v. rpc response missing", httpRequest.URL.String(), httpResponse.StatusCode)
	}

	return rpcResponse, nil
}

// Params is a helper function that uses the same parameter syntax as Call().
// But you should consider to always use NewRequest() instead.
//
// e.g. to manually create an RPCRequest object:
//
//	request := &RPCRequest{
//	  Method: "myMethod",
//	  Params: Params("Alex", 35, true),
//	}
//
// same with new request:
// request := NewRequest("myMethod", "Alex", 35, true)
//
// If you know what you are doing you can omit the Params() call but potentially create incorrect rpc requests:
//
//	request := &RPCRequest{
//	  Method: "myMethod",
//	  Params: 2, <-- invalid since a single primitive value must be wrapped in an array --> no magic without Params()
//	}
//
// correct:
//
//	request := &RPCRequest{
//	  Method: "myMethod",
//	  Params: []int{2}, <-- invalid since a single primitive value must be wrapped in an array
//	}
func Params(params ...interface{}) interface{} {
	var finalParams interface{}

	// if params was nil skip this and p stays nil
	if params != nil {
		switch len(params) {
		case 0: // no parameters were provided, do nothing so finalParam is nil and will be omitted
		case 1: // one param was provided, use it directly as is, or wrap primitive types in array
			if params[0] != nil {
				var typeOf reflect.Type

				// traverse until nil or not a pointer type
				for typeOf = reflect.TypeOf(params[0]); typeOf != nil && typeOf.Kind() == reflect.Ptr; typeOf = typeOf.Elem() {
				}

				if typeOf != nil {
					// now check if we can directly marshal the type or if it must be wrapped in an array
					switch typeOf.Kind() {
					// for these types we just do nothing, since value of p is already unwrapped from the array params
					case reflect.Struct:
						finalParams = params[0]
					case reflect.Array:
						finalParams = params[0]
					case reflect.Slice:
						finalParams = params[0]
					case reflect.Interface:
						finalParams = params[0]
					case reflect.Map:
						finalParams = params[0]
					default: // everything else must stay in an array (int, string, etc)
						finalParams = params
					}
				}
			} else {
				finalParams = params
			}
		default: // if more than one parameter was provided it should be treated as an array
			finalParams = params
		}
	}

	return finalParams
}

// GetObject converts the rpc response to an arbitrary type.
//
// The function works as you would expect it from json.Unmarshal()
func (RPCResponse *RPCResponse) GetObject(toType interface{}) error {
	rv := reflect.ValueOf(toType)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("expected a pointer, got a value: %s", reflect.TypeOf(toType))
	}
	if RPCResponse.Result == nil {
		RPCResponse.Result = []byte(`null`)
	}
	return json.Unmarshal(RPCResponse.Result, toType)
}
