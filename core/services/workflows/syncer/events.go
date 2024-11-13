package syncer

type WorkflowRegistryForceUpdateSecretsRequestedV1 struct {
	SecretsURL   string
	Owner        string
	WorkflowName string
}

var WorkflowRegistryForceUpdateSecretsRequestedV1ModifyFields = []string{"Owner"}
