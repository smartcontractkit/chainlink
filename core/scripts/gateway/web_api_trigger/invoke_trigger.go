package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/webapicap"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
)

// https://gateway-us-1.chain.link/web-api-trigger
//   {
//     jsonrpc: "2.0",
//     id: "...",
//     method: "web-api-trigger",
//     params: {
//       signature: "...",
//       body: {
//         don_id: "workflow_123",
//         payload: {
//           trigger_id: "web-api-trigger@1.0.0",
//           trigger_event_id: "action_1234567890",
//           timestamp: 1234567890,
//           sub-events: [
//             {
//               topics: ["daily_price_update"],
//               params: {
//                 bid: "101",
//                 ask: "102"
//               }
//             },
//             {
//               topics: ["daily_message", "summary"],
//               params: {
//                 message: "all good!",
//               }
//             },
//           ]
//         }
//       }
//     }
//   }

func main() {
	gatewayURL := flag.String("gateway_url", "http://localhost:5002", "Gateway URL")
	privateKey := flag.String("private_key", "65456ffb8af4a2b93959256a8e04f6f2fe0943579fb3c9c3350593aabb89023f", "Private key to sign the message with")
	messageID := flag.String("id", "12345", "Request ID")
	methodName := flag.String("method", "web_api_trigger", "Method name")
	donID := flag.String("don_id", "example_don", "DON ID")

	flag.Parse()

	if privateKey == nil || *privateKey == "" {
		if err := godotenv.Load(); err != nil {
			panic(err)
		}

		privateKeyEnvVar := os.Getenv("PRIVATE_KEY")
		privateKey = &privateKeyEnvVar
		fmt.Println("Loaded private key from .env")
	}

	// validate key and extract address
	key, err := crypto.HexToECDSA(*privateKey)
	if err != nil {
		fmt.Println("error parsing private key", err)
		return
	}

	payload := webapicap.TriggerRequestPayload{
		Timestamp: time.Now().Unix(),
		Topics:    []string{"daily_price_update"},
		Params: webapicap.TriggerRequestPayloadParams{
			"bid": "101",
			"ask": "102",
		},
		TriggerEventId: uuid.New().String(),
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error marshalling map", err)
		return
	}

	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: *messageID,
			Method:    *methodName,
			DonId:     *donID,
			Payload:   json.RawMessage(payloadJson),
		},
	}
	if err = msg.Sign(key); err != nil {
		fmt.Println("error signing message", err)
		return
	}
	codec := api.JsonRPCCodec{}
	rawMsg, err := codec.EncodeRequest(msg)
	if err != nil {
		fmt.Println("error JSON-RPC encoding", err)
		return
	}

	createRequest := func() (req *http.Request, err error) {
		req, err = http.NewRequestWithContext(context.Background(), "POST", *gatewayURL, bytes.NewBuffer(rawMsg))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
		return
	}

	client := &http.Client{}

	sendRequest := func() {
		req, err2 := createRequest()
		if err2 != nil {
			fmt.Println("error creating a request", err2)
			return
		}

		resp, err2 := client.Do(req)
		if err2 != nil {
			fmt.Println("error sending a request", err2)
			return
		}
		defer resp.Body.Close()

		body, err2 := io.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Println("error sending a request", err2)
			return
		}

		var prettyJSON bytes.Buffer
		if err2 = json.Indent(&prettyJSON, body, "", "  "); err2 != nil {
			fmt.Println(string(body))
		} else {
			fmt.Println(prettyJSON.String())
		}
	}
	sendRequest()
}
