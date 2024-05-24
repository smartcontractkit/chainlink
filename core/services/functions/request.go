package functions

const (
	LocationInline     = 0
	LocationRemote     = 1
	LocationDONHosted  = 2
	LanguageJavaScript = 0

	RequestStatePending       = 1
	RequestStateComplete      = 2
	RequestStateInternalError = 3
)

type RequestFlags [32]byte

type OffchainRequest struct {
	RequestId         []byte      `json:"requestId"`
	RequestInitiator  []byte      `json:"requestInitiator"`
	SubscriptionId    uint64      `json:"subscriptionId"`
	SubscriptionOwner []byte      `json:"subscriptionOwner"`
	Timestamp         uint64      `json:"timestamp"`
	Data              RequestData `json:"data"`
}

type RequestData struct {
	Source          string   `json:"source" cbor:"source"`
	Language        int      `json:"language" cbor:"language"`
	CodeLocation    int      `json:"codeLocation" cbor:"codeLocation"`
	Secrets         []byte   `json:"secrets,omitempty" cbor:"secrets"`
	SecretsLocation int      `json:"secretsLocation" cbor:"secretsLocation"`
	Args            []string `json:"args,omitempty" cbor:"args"`
	BytesArgs       [][]byte `json:"bytesArgs,omitempty" cbor:"bytesArgs"`
}

// NOTE: to be extended with raw report and signatures when needed
type OffchainResponse struct {
	RequestId []byte `json:"requestId"`
	Result    []byte `json:"result,omitempty"`
	Error     []byte `json:"error,omitempty"`
}

type HeartbeatResponse struct {
	Status        int               `json:"status"`
	InternalError string            `json:"internalError,omitempty"`
	ReceivedTs    uint64            `json:"receivedTs"`
	CompletedTs   uint64            `json:"completedTs"`
	Response      *OffchainResponse `json:"response,omitempty"`
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
