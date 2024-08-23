package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/chain"
)

func main() {
	gatewayURL := flag.String("gateway_url", "http://localhost:5002", "Gateway URL")
	privateKey := flag.String("private_key", "65456ffb8af4a2b93959256a8e04f6f2fe0943579fb3c9c3350593aabb89023f", "Private key to sign the message with")
	messageID := flag.String("message_id", "12345", "Request ID")
	methodName := flag.String("method", "write", "Method name")
	donID := flag.String("don_id", "don_1", "DON ID")
	chainFamily := flag.String("chain_family", "evm", "Chain family")
	chainID := flag.String("chain_id", "11155111", "Chain ID")
	toAddress := flag.String("to_address", "0x825cb0E13f313C0943A90380ce0bE0b7e0c8AC4C", "To address")
	flag.Parse()

	if privateKey == nil || *privateKey == "" {
		panic("private key is required")
	}

	// validate key and extract address
	key, err := crypto.HexToECDSA(*privateKey)
	if err != nil {
		fmt.Println("error parsing private key", err)
		return
	}

	payload := chain.WriteRequestPayload{
		ChainFamily: *chainFamily,
		ChainID:     *chainID,
		ToAddress:   *toAddress,
	}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: *messageID,
			Method:    *methodName,
			DonId:     *donID,
			Payload:   rawPayload,
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
