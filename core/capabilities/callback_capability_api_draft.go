package capabilities

import (
	"context"
	"fmt"
)

// To flesh out alternative API for CallbackCapabilites, the exising use of channel complicates the remote target implementation, but arguably (IMO) requires Capability developers to correctly
// handle the returned channel, this could lead to potential issues
func main() {

	// Current API
	currentCBC := CurrentCallbackCapability{}
	chResp, err := currentCBC.Execute(context.Background(), CapabilityRequest{})
	if err != nil {
		panic("doom")
	}

	result := <-chResp
	if result.Err != nil {
		panic("doom")
	}
	fmt.Printf("Result %v", result)

	// Proposed API 1
	api1CBC := ProposedAPI1CallbackCapability{}
	err = api1CBC.Execute(context.Background(), CapabilityRequest{}, func(response CapabilityResponse) {
		if response.Err != nil {
			panic("doom")
		}
		fmt.Printf("API 1 response %v", response)
	})
	if err != nil {
		panic("doom")
	}

	// Or just sync API and let the client call it async as needed
	syncCBC := SyncExecuteCallbackCapability{}
	resp, err := syncCBC.Execute(context.Background(), CapabilityRequest{})
	if err != nil {
		panic("doom")
	}
	fmt.Printf("Sync response %v", resp)

	//  async like this or.....
	go func() {
		resp, err := syncCBC.Execute(context.Background(), CapabilityRequest{})
		if err != nil {
			panic("doom")
		}
		fmt.Printf("Async response %v", resp)
	}()

	// ....if there is client side code that requires to execute async and return  a channel that is shut after one response
	// (as per current API) then a simple helper method would be enough?
	respCh := executeAsyncWithChannel(context.Background(), CapabilityRequest{}, syncCBC.Execute)
	resp2 := <-respCh
	fmt.Printf("Async response %v", resp2)

}

func executeAsyncWithChannel(ctx context.Context, request CapabilityRequest, toExecute func(ctx context.Context, request CapabilityRequest) (SyncCapabilityResponse, error)) chan AsyncCapabilityResponse {
	respCh := make(chan AsyncCapabilityResponse, 1)
	go func() {
		resp, err := toExecute(ctx, request)
		respCh <- AsyncCapabilityResponse{Value: resp.Value, Err: err}
		close(respCh)
	}()

	return respCh
}

type AsyncCapabilityResponse struct {
	Value string // not string in practice
	Err   error
}

type CapabilityRequest struct {
}

type CapabilityResponse struct {
	Value string // not string in practice
	Err   error
}

type SyncCapabilityResponse struct {
	Value string // not string in practice
}

type SyncExecuteCallbackCapability struct {
}

func (s SyncExecuteCallbackCapability) Execute(ctx context.Context, request CapabilityRequest) (SyncCapabilityResponse, error) {

	return SyncCapabilityResponse{}, nil
}

type ProposedAPI1CallbackCapability struct {
}

func (p1 ProposedAPI1CallbackCapability) Execute(ctx context.Context, request CapabilityRequest, callback func(response CapabilityResponse)) error {

	go func() {
		// Do stuff
		callback(CapabilityResponse{Value: "result of stuff", Err: nil})
	}()

	return nil
}

type CurrentCallbackCapability struct {
}

func (c CurrentCallbackCapability) Execute(ctx context.Context, request CapabilityRequest) (<-chan CapabilityResponse, error) {
	result := make(chan CapabilityResponse, 10)

	// Capability developer has to know to only return 1 result and make sure to close the channel, not doing this
	// will currently cause issues with remote targets, possibly other issues in areas of code that rely on the developer

	result <- CapabilityResponse{}
	close(result)

	return result, nil
}
