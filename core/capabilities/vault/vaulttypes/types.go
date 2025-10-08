package vaulttypes

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/build"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
)

var DefaultNamespace = "main"

const (
	// MethodSecretsCreate Note: additional methods should be reflected
	// in the `Methods` list below.
	MethodSecretsCreate = "vault.secrets.create"
	MethodSecretsGet    = "vault.secrets.get"
	MethodSecretsUpdate = "vault.secrets.update"
	MethodSecretsDelete = "vault.secrets.delete"
	MethodSecretsList   = "vault.secrets.list"
	MethodPublicKeyGet  = "vault.publicKey.get"

	MaxBatchSize = 10
)

var (
	// MethodSecretsGet is intentionally omitted from this list, as it is not exposed
	// to external clients, but rather used internally by the Workflow DON.
	Methods = []string{
		MethodSecretsCreate,
		MethodSecretsUpdate,
		MethodSecretsDelete,
		MethodSecretsList,
		MethodPublicKeyGet,
	}
)

func GetSupportedMethods(lggr logger.Logger) []string {
	methods := slices.Clone(Methods)
	if build.IsDev() {
		// Allow secrets get in non-prod environments for testing purposes
		// This should never be enabled in production
		methods = append(methods, MethodSecretsGet)
		lggr.Warnw("enabling vault.secrets.get method since it is not a production build", "build-mode", build.Mode())
	}
	return methods
}

// SignedOCRResponse is the response format for OCR signed reports, as returned by the Vault DON.
// External clients should verify that the signatures match the payload and context, before trusting this response.
// Only after validating, clients should decode the payload for further processing.
// If however the Error field is non-empty, it indicates there was an error talking to the Vault DON.
type SignedOCRResponse struct {
	Error      string          `json:"error"`
	Payload    json.RawMessage `json:"payload"`
	Context    []byte          `json:"context"`
	Signatures [][]byte        `json:"signatures"`
}

func (r *SignedOCRResponse) String() string {
	return fmt.Sprintf("SignedOCRResponse { Error: %s, Payload: %s, Context: <[%d]byte blob>, Signatures: <[%d][]byte blob>}", r.Error, string(r.Payload), len(r.Context), len(r.Signatures))
}

type SecretsService interface {
	CreateSecrets(ctx context.Context, request *vaultcommon.CreateSecretsRequest) (*Response, error)
	UpdateSecrets(ctx context.Context, request *vaultcommon.UpdateSecretsRequest) (*Response, error)
	GetSecrets(ctx context.Context, requestID string, request *vaultcommon.GetSecretsRequest) (*Response, error)
	DeleteSecrets(ctx context.Context, request *vaultcommon.DeleteSecretsRequest) (*Response, error)
	ListSecretIdentifiers(ctx context.Context, request *vaultcommon.ListSecretIdentifiersRequest) (*Response, error)

	GetPublicKey(ctx context.Context, request *vaultcommon.GetPublicKeyRequest) (*vaultcommon.GetPublicKeyResponse, error)
}

func KeyFor(id *vaultcommon.SecretIdentifier) string {
	namespace := id.Namespace
	if namespace == "" {
		namespace = DefaultNamespace
	}
	return fmt.Sprintf("%s::%s::%s", id.Owner, namespace, id.Key)
}

type Request struct {
	Payload      proto.Message
	ResponseChan chan *Response

	IDVal         string
	ExpiryTimeVal time.Time
}

func (r *Request) ID() string {
	return r.IDVal
}

func (r *Request) Copy() *Request {
	newRequest := &Request{
		Payload: proto.Clone(r.Payload),

		// intentionally not copied as we want to keep the reference
		ResponseChan: r.ResponseChan,

		// copied by value
		IDVal:         r.IDVal,
		ExpiryTimeVal: r.ExpiryTimeVal,
	}
	return newRequest
}

func (r *Request) ExpiryTime() time.Time {
	return r.ExpiryTimeVal
}

func (r *Request) SendResponse(ctx context.Context, response *Response) {
	select {
	case <-ctx.Done():
		return
	case r.ResponseChan <- response:
	}
}

func (r *Request) SendTimeout(ctx context.Context) {
	r.SendResponse(ctx, &Response{
		ID:    r.IDVal,
		Error: fmt.Sprintf("timeout exceeded: could not process request %s before expiry", r.IDVal),
	})
}

type Response struct {
	ID         string
	Error      string
	Payload    []byte
	Format     string
	Context    []byte
	Signatures [][]byte
}

func (r *Response) ToJSONRPCResult() ([]byte, error) {
	return json.Marshal(SignedOCRResponse{
		Error:      r.Error,
		Payload:    r.Payload,
		Context:    r.Context,
		Signatures: r.Signatures,
	})
}

func (r *Response) RequestID() string {
	return r.ID
}

func (r *Response) String() string {
	return fmt.Sprintf("Response { ID: %s, Error: %s, Payload: %s, Format: %s }", r.ID, r.Error, string(r.Payload), r.Format)
}

func ValidateSignatures(resp *SignedOCRResponse, allowedSigners []common.Address, minRequired int) error {
	if len(resp.Context) < 64 {
		return fmt.Errorf("context too short: expected min 64 bytes, got %d bytes", len(resp.Context))
	}

	if len(resp.Signatures) < minRequired {
		return fmt.Errorf("not enough signatures: expected min %d, got %d", minRequired, len(resp.Signatures))
	}

	// The context contains:
	// 0:32 -> config digest
	// 32:64 -> epoch + round, namely:
	//   - 0:27 -> zero padding
	//   - 27:31 -> sequence number (big endian uint32)
	//   - 31:32 -> zero round value
	// 64:96 -> extra hash (not used by the vault plugin)
	cd, epochRound := resp.Context[:32], resp.Context[32:64]
	configDigest, err := ocr2types.BytesToConfigDigest(cd)
	if err != nil {
		return fmt.Errorf("invalid config digest in signature: %w", err)
	}

	epoch := binary.BigEndian.Uint32(epochRound[27:31])
	round := epochRound[31]

	fullHash := ocr2key.ReportToSigData(ocr2types.ReportContext{
		ReportTimestamp: ocr2types.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        epoch,
			Round:        round,
		},
	}, []byte(resp.Payload))

	validSigners := map[common.Address]bool{}
	for _, s := range resp.Signatures {
		signerPubkey, err := crypto.SigToPub(fullHash, s)
		if err != nil {
			return fmt.Errorf("invalid signature: %w", err)
		}
		signerAddr := crypto.PubkeyToAddress(*signerPubkey)

		for _, as := range allowedSigners {
			if as.Hex() == signerAddr.Hex() {
				validSigners[signerAddr] = true
				break
			}
		}

		if len(validSigners) >= minRequired {
			return nil
		}
	}

	return fmt.Errorf("only %d valid signatures, need at least %d", len(validSigners), minRequired)
}

func DigestForRequest(req jsonrpc.Request[json.RawMessage]) ([32]byte, error) {
	var seed any
	switch req.Method {
	case MethodSecretsCreate:
		var createSecretsRequests vaultcommon.CreateSecretsRequest
		if err := json.Unmarshal(*req.Params, &createSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling create secrets request: " + err.Error())
		}
		secrets := make([]*vaultcommon.EncryptedSecret, len(createSecretsRequests.EncryptedSecrets))
		for i, s := range createSecretsRequests.EncryptedSecrets {
			secrets[i] = &vaultcommon.EncryptedSecret{
				EncryptedValue: s.EncryptedValue,
				Id: &vaultcommon.SecretIdentifier{
					Key:       s.Id.Key,
					Namespace: s.Id.Namespace,
					Owner:     s.Id.Owner,
				},
			}
		}
		seed = vaultcommon.CreateSecretsRequest{
			EncryptedSecrets: secrets,
		}
	case MethodSecretsUpdate:
		var updateSecretsRequests vaultcommon.UpdateSecretsRequest
		if err := json.Unmarshal(*req.Params, &updateSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling update secrets request: " + err.Error())
		}
		secrets := make([]*vaultcommon.EncryptedSecret, len(updateSecretsRequests.EncryptedSecrets))
		for i, s := range updateSecretsRequests.EncryptedSecrets {
			secrets[i] = &vaultcommon.EncryptedSecret{
				EncryptedValue: s.EncryptedValue,
				Id: &vaultcommon.SecretIdentifier{
					Key:       s.Id.Key,
					Namespace: s.Id.Namespace,
					Owner:     s.Id.Owner,
				},
			}
		}
		seed = vaultcommon.CreateSecretsRequest{
			EncryptedSecrets: secrets,
		}
	case MethodSecretsList:
		var listSecretsRequests vaultcommon.ListSecretIdentifiersRequest
		if err := json.Unmarshal(*req.Params, &listSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling list secrets request: " + err.Error())
		}
		seed = vaultcommon.ListSecretIdentifiersRequest{
			Owner:     listSecretsRequests.Owner,
			Namespace: listSecretsRequests.Namespace,
		}
	case MethodSecretsDelete:
		var deleteSecretsRequests vaultcommon.DeleteSecretsRequest
		if err := json.Unmarshal(*req.Params, &deleteSecretsRequests); err != nil {
			return [32]byte{}, errors.New("error unmarshalling delete secrets request: " + err.Error())
		}
		ids := make([]*vaultcommon.SecretIdentifier, len(deleteSecretsRequests.Ids))
		for i, id := range deleteSecretsRequests.Ids {
			ids[i] = &vaultcommon.SecretIdentifier{
				Key:       id.Key,
				Namespace: id.Namespace,
				Owner:     id.Owner,
			}
		}
		seed = vaultcommon.DeleteSecretsRequest{
			Ids: ids,
		}
	default:
		return [32]byte{}, fmt.Errorf("unauthorized method: %s", req.Method)
	}

	// Critical: convert to json, to ensure consistent encoding. Otherwise, different
	// clients may generate different digests for the same logical request.
	jsonData, err := json.Marshal(seed)
	if err != nil {
		return [32]byte{}, errors.New("error marshalling request for digest: " + err.Error())
	}
	return sha256.Sum256(jsonData), nil
}
