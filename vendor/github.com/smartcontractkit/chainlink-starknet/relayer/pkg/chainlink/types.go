package chainlink

// [relayConfig] member of Chainlink's job spec v2 (OCR2 only currently)
type RelayConfig struct {
	ChainID        string `json:"chainID"`
	AccountAddress string `json:"accountAddress"` // address of the account contract
	NodeName       string `json:"nodeName"`       // optional, defaults to random node with 'chainID'
}
