package arbiter

// ReleaseStatus represents the status of a Helm release.
type ReleaseStatus int

const (
	// StatusInstalling indicates that helm release is being installed.
	StatusInstalling ReleaseStatus = iota
	// StatusInstallSucceeded indicates that helm release installation succeeded.
	StatusInstallSucceeded
	// StatusInstallFailed indicates that helm release installation failed.
	StatusInstallFailed
	// StatusReady indicates that deployed service is ready and operational.
	StatusReady
	// StatusDegraded indicates that deployed service is degraded.
	StatusDegraded
	// StatusUnknown indicates that the status of the helm release is unknown.
	StatusUnknown
)

// String returns the string representation of ReleaseStatus.
func (s ReleaseStatus) String() string {
	switch s {
	case StatusInstalling:
		return "INSTALLING"
	case StatusInstallSucceeded:
		return "INSTALL_SUCCEEDED"
	case StatusInstallFailed:
		return "INSTALL_FAILED"
	case StatusReady:
		return "READY"
	case StatusDegraded:
		return "DEGRADED"
	default:
		return "UNKNOWN"
	}
}

// ShardReplica represents the status and message of a single shard replica.
type ShardReplica struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Metrics map[string]float64 `json:"metrics,omitempty"`
}

// ScaleIntentRequest represents the request body for scale-intents endpoint.
type ScaleIntentRequest struct {
	CurrentReplicas     map[string]ShardReplica `json:"currentReplicas"`
	Reason              string                  `json:"reason"`
	DesiredReplicaCount int                     `json:"desiredReplicaCount"`
}

// ScalingSpecResponse represents the response body for scaling-spec endpoint.
type ScalingSpecResponse struct {
	LastScalingReason    string `json:"lastScalingReason"`
	CurrentReplicaCount  int    `json:"currentReplicaCount"`
	DesiredReplicaCount  int    `json:"desiredReplicaCount"`
	ApprovedReplicaCount int    `json:"approvedReplicaCount"`
}

// StatusResponse represents a simple status response.
type StatusResponse struct {
	Status string `json:"status"`
}
