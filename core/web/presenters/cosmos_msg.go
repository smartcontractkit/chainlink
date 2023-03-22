package presenters

// CosmosMsgResource repesents a Cosmos message JSONAPI resource.
type CosmosMsgResource struct {
	JAID
	ChainID    string
	ContractID string
	State      string
	TxHash     *string
}

// GetName implements the api2go EntityNamer interface
func (CosmosMsgResource) GetName() string {
	return "cosmos_messages"
}

// NewCosmosMsgResource returns a new partial CosmosMsgResource.
func NewCosmosMsgResource(id int64, chainID string, contractID string) CosmosMsgResource {
	return CosmosMsgResource{
		JAID:       NewJAIDInt64(id),
		ChainID:    chainID,
		ContractID: contractID,
	}
}
