package chainlink

// [relayConfig] member of Chainlink's job spec v2 (OCR2 only currently)
type RelayConfig struct {
	ChainID  string `json:"chainID"`
	NodeName string `json:"nodeName"` // optional, defaults to random node with 'chainID'
}
