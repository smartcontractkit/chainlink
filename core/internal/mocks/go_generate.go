package mocks

//go:generate mockery --srcpkg github.com/smartcontractkit/chainlink/core/gethwrappers/generated/flux_aggregator_wrapper --name FluxAggregatorInterface --output . --case=underscore --structname FluxAggregator --filename flux_aggregator.go
//go:generate mockery --srcpkg github.com/smartcontractkit/chainlink/core/gethwrappers/generated/flags_wrapper --name FlagsInterface --output . --case=underscore --structname Flags --filename flags.go
//go:generate mockery --srcpkg github.com/smartcontractkit/chainlink/core/gethwrappers/generated/aggregator_v3_interface --name AggregatorV3InterfaceInterface --output ../../services/vrf/mocks/ --case=underscore --structname AggregatorV3Interface --filename aggregator_v3_interface.go
//go:generate mockery --srcpkg github.com/smartcontractkit/chainlink/core/gethwrappers/generated/vrf_coordinator_v2 --name VRFCoordinatorV2Interface --output ../../services/vrf/mocks/ --case=underscore --structname VRFCoordinatorV2Interface --filename vrf_coordinator_v2.go
//go:generate mockery --srcpkg github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon --name VRFBeaconInterface --output ../../services/ocr2/plugins/ocr2vrf/coordinator/mocks --case=underscore --structname VRFBeaconInterface --filename vrf_beacon.go
//go:generate mockery --srcpkg github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_coordinator --name VRFCoordinatorInterface --output ../../services/ocr2/plugins/ocr2vrf/coordinator/mocks --case=underscore --structname VRFCoordinatorInterface --filename vrf_coordinator.go
