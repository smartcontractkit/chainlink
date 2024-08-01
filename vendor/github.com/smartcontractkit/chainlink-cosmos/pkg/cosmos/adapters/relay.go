package adapters

// CL Core OCR2 job spec RelayConfig member for Cosmos
type RelayConfig struct {
	ChainID  string `json:"chainID"`  // required
	NodeName string `json:"nodeName"` // optional, defaults to a random node with ChainID
}
