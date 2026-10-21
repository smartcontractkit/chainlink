package types

const (
	// ShardNameMarker is the marker a workflow DON's name is suffixed with to
	// encode its shard index (e.g. "workflow-1-zone-a-shard-1" -> shard 1). A
	// name with no such suffix (e.g. "workflow-1-zone-a") is shard index 0.
	ShardNameMarker = "shard-"

	// WorkflowFamilySuffix is the suffix appended to a DON family name to
	// denote its dedicated workflow-fetching family. A DON belonging to a
	// family named "<name>WorkflowFamilySuffix" (e.g. "zone-a_workflows")
	// should also fetch workflows registered under the base family with the
	// suffix stripped (e.g. "zone-a").
	WorkflowFamilySuffix = "_workflows"
)
