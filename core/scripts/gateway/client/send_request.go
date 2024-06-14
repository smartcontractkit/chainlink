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
	"github.com/joho/godotenv"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
)

func main() {
	gatewayURL := flag.String("gateway_url", "", "Gateway URL")
	privateKey := flag.String("private_key", "", "Private key to sign the message with")
	messageId := flag.String("message_id", "", "Request ID")
	methodName := flag.String("method", "", "Method name")
	donId := flag.String("don_id", "", "DON ID")
	s4SetSlotId := flag.Uint("s4_set_slot_id", 0, "S4 set slot ID")
	s4SetVersion := flag.Uint64("s4_set_version", 0, "S4 set version")
	s4SetExpirationPeriod := flag.Int64("s4_set_expiration_period", 60*60*1000, "S4 how long until the entry expires from now (in milliseconds)")
	s4SetPayload := flag.String("s4_set_payload", "", "S4 set payload")
	repeat := flag.Bool("repeat", false, "Repeat sending the request every 10 seconds")
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
	address := crypto.PubkeyToAddress(key.PublicKey)

	// build payload (if relevant)
	var payloadJSON []byte
	if *methodName == functions.MethodSecretsSet {
		envelope := s4.Envelope{
			Address:    address.Bytes(),
			SlotID:     *s4SetSlotId,
			Version:    *s4SetVersion,
			Payload:    []byte(*s4SetPayload),
			Expiration: time.Now().UnixMilli() + *s4SetExpirationPeriod,
		}
		signature, err := envelope.Sign(key)
		if err != nil {
			fmt.Println("error signing S4 envelope", err)
			return
		}

		s4SetPayload := functions.SecretsSetRequest{
			SlotID:     envelope.SlotID,
			Version:    envelope.Version,
			Expiration: envelope.Expiration,
			Payload:    []byte(*s4SetPayload),
			Signature:  signature,
		}

		payloadJSON, err = json.Marshal(s4SetPayload)
		if err != nil {
			fmt.Println("error marshaling S4 payload", err)
			return
		}
	}

	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: *messageId,
			Method:    *methodName,
			DonId:     *donId,
			Payload:   json.RawMessage(payloadJSON),
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
		req, err := createRequest()
		if err != nil {
			fmt.Println("error creating a request", err)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("error sending a request", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error sending a request", err)
			return
		}

		fmt.Println(string(body))
	}

	sendRequest()

	for *repeat {
		time.Sleep(10 * time.Second)
		sendRequest()
	}
}
