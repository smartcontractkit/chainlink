package trigger

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/webapi/webapicap"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
)

const (
	Topic  = "sendValue"
	Method = "web_api_trigger"
)

// WebAPITriggerValue triggers a workflow with web API trigger with random value (0 - 1000000) and topic "sendValue"
// It will keep retrying until the workflow is triggered successfully or the timeout is reached
func WebAPITriggerValue(gatewayURL, sender, receiver, privateKey string, timeout time.Duration) error {
	valueInt, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return errors.Wrap(err, "error generating random value")
	}

	value := valueInt.String()
	fmt.Print(libformat.DarkYellowText("🚀 Triggering workflow with value %s\n\n", value))

	payload := webapicap.TriggerRequestPayload{
		Timestamp: time.Now().Unix(),
		Topics:    []string{Topic},
		Params: webapicap.TriggerRequestPayloadParams{
			"paymentId": uuid.New().String(),
			"sender":    sender,
			"receiver":  receiver,
			"value":     value,
		},
		TriggerEventId: uuid.New().String(),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "error marshalling json payload")
	}

	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return errors.Wrap(err, "error parsing private key")
	}

	publicAddress := crypto.PubkeyToAddress(key.PublicKey)

	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: uuid.New().String(),
			Method:    Method,
			DonId:     "1",
			Payload:   json.RawMessage(payloadJSON),
			Sender:    publicAddress.String(),
		},
	}
	if err = msg.Sign(key); err != nil {
		return errors.Wrap(err, "error signing message")
	}

	codec := api.JsonRPCCodec{}
	rawMsg, err := codec.EncodeRequest(msg)
	if err != nil {
		return errors.Wrap(err, "error JSON-RPC encoding")
	}

	createRequest := func() (req *http.Request, err error) {
		req, err = http.NewRequestWithContext(context.Background(), "POST", gatewayURL, bytes.NewBuffer(rawMsg))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
		return
	}

	client := &http.Client{}

	sendRequest := func() error {
		req, createErr := createRequest()
		if createErr != nil {
			return errors.Wrap(createErr, "error creating a request")
		}

		resp, sendErr := client.Do(req)
		if sendErr != nil {
			return errors.Wrap(sendErr, "error sending a request")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return errors.Wrap(readErr, "error reading response body")
		}

		var payloadMap map[string]any
		if unmarshalErr := json.Unmarshal(body, &payloadMap); unmarshalErr != nil {
			return errors.Wrap(unmarshalErr, "error unmarshalling response body")
		}

		if result, ok := payloadMap["result"]; ok {
			if resultMap, ok := result.(map[string]any); ok {
				if bodyMap, ok := resultMap["body"].(map[string]any); ok {
					if payloadMap, ok := bodyMap["payload"].(map[string]any); ok {
						if status, ok := payloadMap["status"].(string); ok {
							if strings.EqualFold(status, "error") {
								return fmt.Errorf("error sending request to gateway: %v", payloadMap["error_message"])
							}

							return nil
						}
					}
				}
			}
		}

		return fmt.Errorf("unexpected response body: %v", string(body))
	}

	ticker := 10 * time.Second

	for {
		select {
		case <-time.After(timeout):
			return fmt.Errorf("failed to send request to gateway within %.2f seconds", timeout.Seconds())
		case <-time.Tick(ticker):
			if sendErr := sendRequest(); sendErr == nil {
				fmt.Print(libformat.DarkYellowText("🏁 Successfully triggered workflow\n\n"))
				return nil
			} else {
				fmt.Printf("failed to send request to gateway: %v. Retrying in %.2f seconds...\n", sendErr, ticker.Seconds())
			}
		}
	}
}
