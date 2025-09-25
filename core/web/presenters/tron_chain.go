package presenters

// TronChainResource is an Tron chain JSONAPI resource.
type TronChainResource struct {
	ChainResource
}

// GetName implements the api2go EntityNamer interface
func (r TronChainResource) GetName() string {
	return "tron_chain"
}

// TronNodeResource is a Tron node JSONAPI resource.
type TronNodeResource struct {
	NodeResource
}

// GetName implements the api2go EntityNamer interface
func (r TronNodeResource) GetName() string {
	return "tron_node"
}
