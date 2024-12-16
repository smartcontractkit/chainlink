package changeset

// TODO: create similar struct as CCIPChainState, but with Sol specific primitives
// SolChainState holds a Go binding for all the currently deployed CCIP programs
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type SolCCIPChainState struct {
}

// TODO: create GenerateSolanaView function similar to GenerateView