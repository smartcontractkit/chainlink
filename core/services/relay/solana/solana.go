package solana

// CL Core OCR2 job spec RelayConfig member for Solana
type RelayConfig struct {
	// network data
	NodeEndpointHTTP string `json:"nodeEndpointHTTP"`

	// state account passed as the ContractID in main job spec
	// on-chain program + transmissions account + store programID
	OCR2ProgramID   string `json:"ocr2ProgramID"`
	TransmissionsID string `json:"transmissionsID"`
	StoreProgramID  string `json:"storeProgramID"`

	// transaction + state parameters [OPTIONAL]
	UsePreflight bool   `json:"usePreflight"`
	Commitment   string `json:"commitment"`

	// polling parameters [OPTIONAL]
	PollingInterval   string `json:"pollingInterval"`
	PollingCtxTimeout string `json:"pollingCtxTimeout"`
	StaleTimeout      string `json:"staleTimeout"`
}
