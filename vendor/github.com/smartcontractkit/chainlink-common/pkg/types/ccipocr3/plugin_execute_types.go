package ccipocr3

type ExecutePluginReport struct {
	ChainReports []ExecutePluginReportSingleChain `json:"chainReports"`
}

type ExecutePluginReportSingleChain struct {
	SourceChainSelector ChainSelector `json:"sourceChainSelector"`
	Messages            []Message     `json:"messages"`
	OffchainTokenData   [][][]byte    `json:"offchainTokenData"`
	Proofs              []Bytes32     `json:"proofs"`
	ProofFlagBits       BigInt        `json:"proofFlagBits"`
}
