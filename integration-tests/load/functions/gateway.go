package loadfunctions

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/api"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/functions"
	"github.com/smartcontractkit/chainlink/v2/core/services/s4"
)

type RPCResponse struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Body struct {
			DonID     string `json:"don_id"`
			MessageID string `json:"message_id"`
			Method    string `json:"method"`
			Payload   struct {
				NodeResponses []struct {
					Body struct {
						DonID     string `json:"don_id"`
						MessageID string `json:"message_id"`
						Method    string `json:"method"`
						Payload   struct {
							Success bool `json:"success"`
						} `json:"payload"`
						Receiver string `json:"receiver"`
					} `json:"body"`
					Signature string `json:"signature"`
				} `json:"node_responses"`
				Success bool `json:"success"`
			} `json:"payload"`
			Receiver string `json:"receiver"`
		} `json:"body"`
		Signature string `json:"signature"`
	} `json:"result"`
}

func UploadS4Secrets(rc *resty.Client, s4Cfg *S4SecretsCfg) (uint8, uint64, error) {
	key, err := crypto.HexToECDSA(s4Cfg.PrivateKey)
	if err != nil {
		return 0, 0, err
	}
	address := crypto.PubkeyToAddress(key.PublicKey)
	var payloadJSON []byte
	envelope := s4.Envelope{
		Address:    address.Bytes(),
		SlotID:     s4Cfg.S4SetSlotID,
		Version:    s4Cfg.S4SetVersion,
		Payload:    []byte(s4Cfg.S4SetPayload),
		Expiration: time.Now().UnixMilli() + s4Cfg.S4SetExpirationPeriod,
	}
	signature, err := envelope.Sign(key)
	if err != nil {
		return 0, 0, err
	}

	s4SetPayload := functions.SecretsSetRequest{
		SlotID:     envelope.SlotID,
		Version:    envelope.Version,
		Expiration: envelope.Expiration,
		Payload:    []byte(s4Cfg.S4SetPayload),
		Signature:  signature,
	}

	payloadJSON, err = json.Marshal(s4SetPayload)
	if err != nil {
		return 0, 0, err
	}

	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: s4Cfg.MessageID,
			Method:    s4Cfg.Method,
			DonId:     s4Cfg.DonID,
			Payload:   json.RawMessage(payloadJSON),
		},
	}

	err = msg.Sign(key)
	if err != nil {
		return 0, 0, err
	}
	codec := api.JsonRPCCodec{}
	rawMsg, err := codec.EncodeRequest(msg)
	if err != nil {
		return 0, 0, err
	}
	var result *RPCResponse
	resp, err := rc.R().
		SetBody(rawMsg).
		Post(s4Cfg.GatewayURL)
	if err != nil {
		return 0, 0, err
	}
	if resp.StatusCode() != 200 {
		return 0, 0, fmt.Errorf("status code was %d, expected 200", resp.StatusCode())
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return 0, 0, err
	}
	log.Debug().Interface("Result", result).Msg("S4 secrets_set response result")
	for _, nodeResponse := range result.Result.Body.Payload.NodeResponses {
		if !nodeResponse.Body.Payload.Success {
			return 0, 0, errors.New("node response was not successful")
		}
	}
	if envelope.SlotID > math.MaxUint8 {
		return 0, 0, fmt.Errorf("slot ID overflows uint8: %d", envelope.SlotID)
	}
	return uint8(envelope.SlotID), envelope.Version, nil
}

func ListS4Secrets(rc *resty.Client, s4Cfg *S4SecretsCfg) error {
	key, err := crypto.HexToECDSA(s4Cfg.PrivateKey)
	if err != nil {
		return err
	}

	msg := &api.Message{
		Body: api.MessageBody{
			MessageId: s4Cfg.MessageID,
			Method:    s4Cfg.Method,
			DonId:     s4Cfg.DonID,
			Receiver:  s4Cfg.RecieverAddr,
		},
	}

	err = msg.Sign(key)
	if err != nil {
		return err
	}
	codec := api.JsonRPCCodec{}
	rawMsg, err := codec.EncodeRequest(msg)
	if err != nil {
		return err
	}
	msgdec, err := codec.DecodeRequest(rawMsg)
	if err != nil {
		return err
	}
	log.Debug().Interface("Request", msgdec).Msg("Sending RPC request")
	var result map[string]interface{}
	resp, err := rc.R().
		SetBody(rawMsg).
		Post(s4Cfg.GatewayURL)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return err
	}
	log.Debug().Interface("Result", result).Msg("S4 secrets_list response result")
	if resp.StatusCode() != 200 {
		return fmt.Errorf("status code was %d, expected 200", resp.StatusCode())
	}
	return nil
}

func ParseTDH2Key(data []byte) (*tdh2easy.PublicKey, error) {
	pk := &tdh2easy.PublicKey{}
	if err := pk.Unmarshal(data); err != nil {
		return nil, err
	}
	return pk, nil
}

func EncryptS4Secrets(deployerPk *ecdsa.PrivateKey, tdh2Pk *tdh2easy.PublicKey, donKey []byte, msgJSON string) (string, error) {
	// 65 bytes PublicKey format, should start with 0x04 to be processed by crypto.UnmarshalPubkey()
	b := make([]byte, 1)
	b[0] = 0x04
	donKey = bytes.Join([][]byte{b, donKey}, nil)
	donPubKey, err := crypto.UnmarshalPubkey(donKey)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal DON key: %w", err)
	}
	eciesDONPubKey := ecies.ImportECDSAPublic(donPubKey)
	signature, err := deployerPk.Sign(rand.Reader, []byte(msgJSON), nil)
	if err != nil {
		return "", fmt.Errorf("failed to sign the msg with Ethereum key: %w", err)
	}
	signedSecrets, err := json.Marshal(struct {
		Signature []byte `json:"signature"`
		Message   string `json:"message"`
	}{
		Signature: signature,
		Message:   msgJSON,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal signed secrets: %w", err)
	}
	ct, err := ecies.Encrypt(rand.Reader, eciesDONPubKey, signedSecrets, nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt with DON key: %w", err)
	}
	ct0xFormat, err := json.Marshal(map[string]interface{}{"0x0": base64.StdEncoding.EncodeToString(ct)})
	if err != nil {
		return "", fmt.Errorf("failed to marshal DON key encrypted format: %w", err)
	}
	ctTDH2Format, err := tdh2easy.Encrypt(tdh2Pk, ct0xFormat)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt with TDH2 public key: %w", err)
	}
	tdh2Message, err := ctTDH2Format.Marshal()
	if err != nil {
		return "", fmt.Errorf("failed to marshal TDH2 encrypted msg: %w", err)
	}
	finalMsg, err := json.Marshal(map[string]interface{}{
		"encryptedSecrets": "0x" + hex.EncodeToString(tdh2Message),
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal secrets msg: %w", err)
	}
	return string(finalMsg), nil
}
