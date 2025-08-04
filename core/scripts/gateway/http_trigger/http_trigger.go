package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	gateway "github.com/smartcontractkit/chainlink-common/pkg/types/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func main() {
	gatewayURL := flag.String("gateway_url", "http://localhost:5002/user", "Gateway URL")
	// ALICE Public: 0xC3Ad031A27E1A6C692cBdBafD85359b0BE1B15DD
	// ALICE Private: 08aff6010e6e21c8290a3eeaf6e6067bcbc1b6ded58c3d0c575c9d2f835b7608
	// BOB Public: 0x4b8d44A7A1302011fbc119407F8Ce3baee6Ea2FF
	// BOB Private: 4abd08efd9fdb3b4dd51e51af19243cc618da47423a26b724cb13338a885ba35
	privateKey := flag.String("private_key", "4abd08efd9fdb3b4dd51e51af19243cc618da47423a26b724cb13338a885ba35", "Private key to sign the JWT with")
	workflowID := flag.String("workflow_id", "", "Workflow ID (if not provided, uses workflow_name/owner/tag)")
	workflowName := flag.String("workflow_name", "", "Workflow Name (optional)")
	workflowOwner := flag.String("workflow_owner", "", "Workflow Owner (optional)")
	workflowTag := flag.String("workflow_tag", "", "Workflow Tag (optional)")
	dedupe := flag.Bool("dedupe", false, "Enable deduplication")
	numRequests := flag.Int("num_requests", 1, "Number of requests to send with same requestID")
	conflictTest := flag.Bool("conflict_test", false, "Send requests with same requestID but different input to test idempotency")
	uniqueID := flag.Bool("unique_id", false, "Generate a unique ID for each request (overrides num_requests)")
	toppings := flag.String("toppings", "pepperoni,mushroom,sausage", "Comma-separated list of pizza toppings")

	flag.Parse()

	if privateKey == nil || *privateKey == "" {
		if err := godotenv.Load(); err != nil {
			panic(err)
		}

		privateKeyEnvVar := os.Getenv("PRIVATE_KEY")
		privateKey = &privateKeyEnvVar
		fmt.Println("✅ Loaded private key from .env file")
	}

	// validate key and extract address
	key, err := crypto.HexToECDSA(*privateKey)
	if err != nil {
		fmt.Println("❌ Error parsing private key:", err)
		return
	}

	address := crypto.PubkeyToAddress(key.PublicKey)
	fmt.Println("🔐 Authentication Setup")
	fmt.Printf("   Public Address: %s\n", address.Hex())

	// Parse toppings from command line flag
	var toppingsSlice []string
	if *toppings != "" {
		toppingsSlice = strings.Split(*toppings, ",")
		// Trim whitespace from each topping
		for i, topping := range toppingsSlice {
			toppingsSlice[i] = strings.TrimSpace(topping)
		}
	}

	// Create sample input data with order ID (uuid) and toppings
	orderID := uuid.New().String()
	input := map[string]interface{}{
		"order_id": orderID,
		"toppings": toppingsSlice,
	}

	// Add dedupe field if requested
	if *dedupe {
		input["dedupe"] = true
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		fmt.Println("❌ Error marshalling input JSON:", err)
		return
	}

	// Create the workflow selector based on provided flags
	var workflowSelector gateway.WorkflowSelector
	if *workflowID != "" {
		// If workflowID is provided, use only workflowID
		workflowSelector = gateway.WorkflowSelector{
			WorkflowID: *workflowID,
		}
		fmt.Println("\n🔍 Workflow Selection")
		fmt.Printf("   Using Workflow ID: %s\n", *workflowID)
	} else {
		// If workflowID is not provided, use name/owner/tag
		workflowSelector = gateway.WorkflowSelector{
			WorkflowName:  *workflowName,
			WorkflowOwner: *workflowOwner,
			WorkflowTag:   *workflowTag,
		}
		fmt.Println("\n🔍 Workflow Selection")
		fmt.Printf("   Workflow Name: %s\n", *workflowName)
		fmt.Printf("   Workflow Owner: %s\n", *workflowOwner)
		fmt.Printf("   Workflow Tag: %s\n", *workflowTag)
	}

	// Create the HTTPTriggerRequest
	triggerRequest := gateway.HTTPTriggerRequest{
		Input:    json.RawMessage(inputJSON),
		Workflow: workflowSelector,
	}

	// Create the JSON-RPC request
	requestID := uuid.New().String()
	client := &http.Client{}

	sendRequest := func(requestData gateway.HTTPTriggerRequest, reqID string, requestNum int) {
		// Create request with specific data and ID
		specificRequest := jsonrpc.Request[gateway.HTTPTriggerRequest]{
			Version: jsonrpc.JsonRpcVersion,
			ID:      reqID,
			Method:  gateway.MethodWorkflowExecute,
			Params:  &requestData,
		}

		fmt.Println("\n📦 Workflow Input Data")
		fmt.Printf("   Input: %+v\n", string(requestData.Input))

		// Pretty print the JSON-RPC request
		fmt.Println("\n📋 JSON-RPC Request")

		// Create a display version without the empty key field
		displayRequest := struct {
			Version string `json:"jsonrpc"`
			ID      string `json:"id"`
			Method  string `json:"method"`
			Params  *struct {
				Input    json.RawMessage          `json:"input"`
				Workflow gateway.WorkflowSelector `json:"workflow"`
			} `json:"params"`
		}{
			Version: specificRequest.Version,
			ID:      specificRequest.ID,
			Method:  specificRequest.Method,
			Params: &struct {
				Input    json.RawMessage          `json:"input"`
				Workflow gateway.WorkflowSelector `json:"workflow"`
			}{
				Input:    requestData.Input,
				Workflow: requestData.Workflow,
			},
		}

		requestJSON, err2 := json.MarshalIndent(displayRequest, "", "    ")
		if err2 != nil {
			fmt.Printf("❌ Error marshalling JSON-RPC request for display: %v\n", err2)
		} else {
			// Add indentation to each line for better visual hierarchy
			lines := strings.Split(string(requestJSON), "\n")
			for _, line := range lines {
				if len(line) > 0 {
					fmt.Printf("   %s\n", line)
				}
			}
		}

		// Encode the JSON-RPC request
		rawRequest, err2 := jsonrpc.EncodeRequest(&specificRequest)
		if err2 != nil {
			fmt.Printf("❌ Error encoding JSON-RPC request %d: %v\n", requestNum, err2)
			return
		}

		req, err2 := http.NewRequestWithContext(context.Background(), "POST", *gatewayURL, bytes.NewBuffer(rawRequest))
		if err2 != nil {
			fmt.Printf("❌ Error creating request %d: %v\n", requestNum, err2)
			return
		}

		// Create and sign JWT
		jwtToken, err := utils.CreateRequestJWT(specificRequest, utils.WithIssuer(address.Hex()))
		if err != nil {
			fmt.Println("❌ Error creating JWT:", err)
			return
		}

		// Sign the JWT with the private key
		signedJWT, err := jwtToken.SignedString(key)
		if err != nil {
			fmt.Println("error signing JWT", err)
			return
		}

		fmt.Println("\n🔑 JWT Token Generated")
		fmt.Printf("   Token Length: %d characters\n", len(signedJWT))
		fmt.Printf("   Token: %s...\n", signedJWT)
		// Print the JWT header and payload, base64-decoded and pretty-printed
		parts := strings.Split(signedJWT, ".")
		if len(parts) >= 2 {
			fmt.Println("\n🔍 JWT Header (decoded):")
			headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
			if err != nil {
				fmt.Printf("   (error decoding header: %v)\n", err)
			} else {
				var prettyHeader bytes.Buffer
				if err := json.Indent(&prettyHeader, headerBytes, "", "    "); err != nil {
					fmt.Printf("   %s\n", string(headerBytes))
				} else {
					lines := bytes.Split(prettyHeader.Bytes(), []byte("\n"))
					for _, line := range lines {
						if len(line) > 0 {
							fmt.Printf("   %s\n", string(line))
						}
					}
				}
			}

			fmt.Println("\n🔍 JWT Payload (decoded):")
			payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
			if err != nil {
				fmt.Printf("   (error decoding payload: %v)\n", err)
			} else {
				var prettyPayload bytes.Buffer
				if err := json.Indent(&prettyPayload, payloadBytes, "", "    "); err != nil {
					fmt.Printf("   %s\n", string(payloadBytes))
				} else {
					lines := bytes.Split(prettyPayload.Bytes(), []byte("\n"))
					for _, line := range lines {
						if len(line) > 0 {
							fmt.Printf("   %s\n", string(line))
						}
					}
				}
			}
		} else {
			fmt.Println("   (JWT does not have at least 2 parts to decode)")
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+signedJWT)

		fmt.Printf("\n🚀 Sending HTTP Request %d\n", requestNum)
		fmt.Printf("   Gateway URL: %s\n", *gatewayURL)
		if *workflowID != "" {
			fmt.Printf("   Workflow ID: %s\n", *workflowID)
		} else {
			fmt.Printf("   Workflow Name: %s\n", *workflowName)
			fmt.Printf("   Workflow Owner: %s\n", *workflowOwner)
			fmt.Printf("   Workflow Tag: %s\n", *workflowTag)
		}
		fmt.Printf("   Method: %s\n", gateway.MethodWorkflowExecute)
		fmt.Printf("   Request ID: %s\n", reqID)

		resp, err2 := client.Do(req)
		if err2 != nil {
			fmt.Printf("❌ Error sending request %d: %v\n", requestNum, err2)
			return
		}
		defer resp.Body.Close()

		body, err2 := io.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Printf("❌ Error reading response %d: %v\n", requestNum, err2)
			return
		}

		fmt.Printf("\n📥 Response %d Received\n", requestNum)
		fmt.Printf("   Status: %s\n", resp.Status)
		fmt.Println("   Body:")

		var prettyJSON bytes.Buffer
		if err2 = json.Indent(&prettyJSON, body, "", "    "); err2 != nil {
			fmt.Printf("   %s\n", string(body))
		} else {
			// Add indentation to each line for better visual hierarchy
			lines := bytes.Split(prettyJSON.Bytes(), []byte("\n"))
			for _, line := range lines {
				if len(line) > 0 {
					fmt.Printf("   %s\n", string(line))
				}
			}
		}
	}

	fmt.Println("🌟 Chainlink Gateway HTTP Trigger Demo")
	fmt.Println("=====================================")

	if *conflictTest {
		// Send requests with same requestID but different input to test idempotency
		fmt.Printf("\n🔄 Conflict Test Mode: Sending %d requests with same ID but different input\n", *numRequests)

		for i := 1; i <= *numRequests; i++ {
			// Create different input for each request
			conflictOrderID := uuid.New().String()

			// Create variations of the base toppings for conflict testing
			conflictToppings := make([]string, len(toppingsSlice))
			copy(conflictToppings, toppingsSlice)
			// Add a unique topping for this request
			conflictToppings = append(conflictToppings, fmt.Sprintf("extra_%d", i))

			conflictInput := map[string]interface{}{
				"order_id":    conflictOrderID,
				"toppings":    conflictToppings,
				"request_num": i,
			}

			// Add dedupe field if requested
			if *dedupe {
				conflictInput["dedupe"] = true
			}

			conflictInputJSON, err2 := json.Marshal(conflictInput)
			if err2 != nil {
				fmt.Printf("❌ Error marshalling conflict input %d: %v\n", i, err2)
				continue
			}

			conflictRequest := gateway.HTTPTriggerRequest{
				Input:    json.RawMessage(conflictInputJSON),
				Workflow: workflowSelector,
			}

			sendRequest(conflictRequest, requestID, i) // Same requestID for all
		}
	} else if *uniqueID {
		// Generate a unique ID for each request
		fmt.Printf("\n🔄 Unique ID Mode: Sending %d requests with unique IDs\n", *numRequests)

		for i := 1; i <= *numRequests; i++ {
			uniqueRequestID := uuid.New().String()          // Generate a new unique ID for each request
			sendRequest(triggerRequest, uniqueRequestID, i) // Unique requestID for each
			time.Sleep(1000 * time.Millisecond)             // Optional: add a small delay between requests
		}
	} else {
		// Send multiple requests with same requestID and same input
		fmt.Printf("\n🔄 Sending %d requests with same ID and same input\n", *numRequests)

		for i := 1; i <= *numRequests; i++ {
			sendRequest(triggerRequest, requestID, i) // Same requestID and data for all
		}
	}

	fmt.Println("\n✅ Demo completed successfully!")
}
