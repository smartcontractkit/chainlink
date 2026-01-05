package grpcsourcemock

import (
	"context"
	"encoding/hex"
	"log/slog"
	"os"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/privateregistry"
)

var registryLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
})).With("logger", "grpc_source_mock.PrivateRegistryService")

// PrivateRegistryService implements the WorkflowDeploymentAction interface for managing workflows
type PrivateRegistryService struct {
	store *WorkflowStore
}

// NewPrivateRegistryService creates a new PrivateRegistryService
func NewPrivateRegistryService(store *WorkflowStore) *PrivateRegistryService {
	return &PrivateRegistryService{
		store: store,
	}
}

// Ensure PrivateRegistryService implements WorkflowDeploymentAction
var _ privateregistry.WorkflowDeploymentAction = (*PrivateRegistryService)(nil)

// AddWorkflow registers a new workflow with the source
func (s *PrivateRegistryService) AddWorkflow(ctx context.Context, workflow *privateregistry.WorkflowRegistration) error {
	registryLogger.Debug("AddWorkflow called",
		"workflowID", hex.EncodeToString(workflow.WorkflowID[:]),
		"workflowName", workflow.WorkflowName,
		"donFamily", workflow.DonFamily,
		"binaryURL", workflow.BinaryURL,
	)
	return s.store.Add(workflow)
}

// UpdateWorkflow updates the workflow's status configuration
func (s *PrivateRegistryService) UpdateWorkflow(ctx context.Context, workflowID [32]byte, config *privateregistry.WorkflowStatusConfig) error {
	return s.store.Update(workflowID, config)
}

// DeleteWorkflow removes the workflow from the source
func (s *PrivateRegistryService) DeleteWorkflow(ctx context.Context, workflowID [32]byte) error {
	return s.store.Delete(workflowID)
}
