package interfaces

import "context"

// ContractView defines the base interface for any contract view
type ContractView interface {
	// SerializeView converts the view to string format like JSON
	SerializeView() (string, error)
}

// ContractViewGenerator
// It is not strictly necessary, but it helps to have a common interface for all view generators
// Tip: The view generated should only depend on the contract state retrieved on-chain and not any external information.
// This will keep the view generation deterministic and consistent.
//
// P is an optional set of parameters needed to generate the view
// V is the type of the contract view
type ContractViewGenerator[P any, V ContractView] interface {
	Generate(ctx context.Context, params P) (V, error)
}
