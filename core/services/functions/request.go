package functions

import (
	"github.com/ethereum/go-ethereum/common"
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

// The request signature should sign the keccak256 hash of the following JSON string
// with the corresponding Request fields in the order that they appear below:
// {
//  "codeLocation": number, (0 for Location.Inline)
//  "language": number, (0 for CodeLanguage.JavaScript)
//  "secrets": string, (encryptedSecretsReference as base64 string, must be `null` if there are no secrets)
//  "secretsLocation": number, (must be `null` if there are no secrets) (1 for Location.Remote, 2 for Location.DONHosted)
//  "source": string,
// }

func (l *FunctionsListener) VerifyRequestSignature(requestID RequestID, subscriptionOwner common.Address, requestData *RequestData) error {
	return l.verifyRequestSignature(requestID, subscriptionOwner, requestData)
}

func (l *FunctionsListener) ParseCBOR(requestId RequestID, cborData []byte, maxSizeBytes uint32) (*RequestData, error) {
	return l.parseCBOR(requestId, cborData, maxSizeBytes)
}
