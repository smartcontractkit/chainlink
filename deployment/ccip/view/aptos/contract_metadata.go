package aptos

type ContractMetaData struct {
	Address        string `json:"address,omitempty"`
	Owner          string `json:"owner,omitempty"`
	TypeAndVersion string `json:"typeAndVersion,omitempty"`
}
