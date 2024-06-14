package functions

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	LocationInline     = 0
	LocationRemote     = 1
	LocationDONHosted  = 2
	LanguageJavaScript = 0
)

type RequestFlags [32]byte

type RequestData struct {
	Source           string   `json:"source" cbor:"source"`
	Language         int      `json:"language" cbor:"language"`
	CodeLocation     int      `json:"codeLocation" cbor:"codeLocation"`
	Secrets          []byte   `json:"secrets" cbor:"secrets"`
	SecretsLocation  int      `json:"secretsLocation" cbor:"secretsLocation"`
	RequestSignature []byte   `json:"requestSignature,omitempty" cbor:"requestSignature"`
	Args             []string `json:"args,omitempty" cbor:"args"`
	BytesArgs        [][]byte `json:"bytesArgs,omitempty" cbor:"bytesArgs"`
}

type DONHostedSecrets struct {
	SlotID  uint   `json:"slotId" cbor:"slotId"`
	Version uint64 `json:"version" cbor:"version"`
}

type SignedRequestData struct {
	CodeLocation    int    `json:"codeLocation" cbor:"codeLocation"`
	Language        int    `json:"language" cbor:"language"`
	Secrets         []byte `json:"secrets" cbor:"secrets"`
	SecretsLocation int    `json:"secretsLocation" cbor:"secretsLocation"`
	Source          string `json:"source" cbor:"source"`
}

// The request signature should sign the keccak256 hash of the following JSON string (without extra whitespace)
// with the corresponding Request fields in the order that they appear below:
// {
//  "codeLocation": number, (0 for Location.Inline)
//  "language": number, (0 for CodeLanguage.JavaScript)
//  "secrets": string, (encryptedSecretsReference as base64 string, must be `null` if there are no secrets)
//  "secretsLocation": number, (must be `null` if there are no secrets) (1 for Location.Remote, 2 for Location.DONHosted)
//  "source": string,
// }

func VerifyRequestSignature(subscriptionOwner common.Address, requestData *RequestData) error {
	if requestData.RequestSignature == nil {
		return errors.New("missing signature")
	}
	signedRequestData := SignedRequestData{
		CodeLocation:    requestData.CodeLocation,
		Language:        requestData.Language,
		Secrets:         requestData.Secrets,
		SecretsLocation: requestData.SecretsLocation,
		Source:          requestData.Source,
	}
	js, err := json.Marshal(signedRequestData)
	if err != nil {
		return errors.New("unable to marshal request data")
	}

	signerAddr, err := utils.GetSignersEthAddress(js, requestData.RequestSignature)
	if err != nil {
		return errors.New("invalid request signature: unable to recover signer's address")
	}

	if signerAddr != subscriptionOwner {
		return errors.New("invalid request signature: signer's address does not match subscription owner")
	}

	return nil
}
