package presenters

// SolanaMsgResource repesents a Solana message JSONAPI resource.
type SolanaMsgResource struct {
	JAID
	ChainID string
	From    string `json:"from"`
	To      string `json:"to"`
	Amount  uint64 `json:"amount"`
}

// GetName implements the api2go EntityNamer interface
func (SolanaMsgResource) GetName() string {
	return "solana_messages"
}

// NewSolanaMsgResource returns a new partial SolanaMsgResource.
func NewSolanaMsgResource(id string, chainID string) SolanaMsgResource {
	return SolanaMsgResource{
		JAID:    NewPrefixedJAID(id, chainID),
		ChainID: chainID,
	}
}
