package ccip

type LoadConfig struct {
	TestTimeout                     *string
	NoOfNodes                       *int
	LokiEndpoint                    *string
	MessageTypeWeights              *[]int
	SecondsPerRequestPerDestination *int
	EnabledDestionationChains       *[]uint64
}
