package mocks

//go:generate mockery --srcpkg github.com/smartcontractkit/chainlink/core/gethwrappers/generated/flux_aggregator_wrapper --name FluxAggregatorInterface --output . --case=underscore --structname FluxAggregator --filename flux_aggregator.go
//go:generate mockery --srcpkg github.com/smartcontractkit/chainlink/core/gethwrappers/generated/flags_wrapper --name FlagsInterface --output . --case=underscore --structname Flags --filename flags.go
