package chaos

const (
	// ChaosGroupMinority a group of faulty nodes, even if they fail OCR must work
	ChaosGroupMinority = "chaosGroupMinority"
	// ChaosGroupMajority a group of nodes that are working even if minority fails
	ChaosGroupMajority = "chaosGroupMajority"
	// ChaosGroupMajorityPlus a group of nodes that are majority + 1
	ChaosGroupMajorityPlus = "chaosGroupMajorityPlus"

	PodChaosFailMercury                  = "pod-chaos-fail-mercury-server"
	PodChaosFailMinorityNodes            = "pod-chaos-fail-minority-nodes"
	PodChaosFailMajorityNodes            = "pod-chaos-fail-majority-nodes"
	PodChaosFailMajorityDB               = "pod-chaos-fail-majority-db"
	NetworkChaosFailMajorityNetwork      = "network-chaos-fail-majority-network"
	NetworkChaosFailBlockchainNode       = "network-chaos-fail-blockchain-node"
	NetworkChaosDisruptNetworkDONMercury = "network-chaos-disrupt-don-mercury"
)
