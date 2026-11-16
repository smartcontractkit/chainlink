package v2

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"sync"
	"testing"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/workflowkey"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	artifacts "github.com/smartcontractkit/chainlink/v2/core/services/workflows/artifacts/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

// singleRowORM stands in for workflow_specs_v2: one row keyed by workflow ID,
// shared by every source, with the production pause semantics (keep the row,
// clear the artifact payload).
type singleRowORM struct {
	mu   sync.Mutex
	spec *job.WorkflowSpec
}

func (o *singleRowORM) UpsertWorkflowSpec(_ context.Context, spec *job.WorkflowSpec) (int64, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	cp := *spec
	o.spec = &cp
	return 1, nil
}

func (o *singleRowORM) GetWorkflowSpec(_ context.Context, id string) (*job.WorkflowSpec, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.spec == nil || o.spec.WorkflowID != id {
		return nil, sql.ErrNoRows
	}
	cp := *o.spec
	return &cp, nil
}

func (o *singleRowORM) ListWorkflowSpecs(context.Context) ([]*job.WorkflowSpec, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.spec == nil {
		return nil, nil
	}
	cp := *o.spec
	return []*job.WorkflowSpec{&cp}, nil
}

func (o *singleRowORM) DeleteWorkflowSpec(_ context.Context, id string) (*job.WorkflowSpec, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.spec == nil || o.spec.WorkflowID != id {
		return nil, nil
	}
	cp := *o.spec
	o.spec = nil
	return &cp, nil
}

func (o *singleRowORM) PauseWorkflowSpec(_ context.Context, id string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.spec != nil && o.spec.WorkflowID == id {
		o.spec.Status = job.WorkflowSpecStatusPaused
		o.spec.Workflow = ""
		o.spec.Config = ""
	}
	return nil
}

func (o *singleRowORM) DeleteWorkflowSpecs(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if _, err := o.DeleteWorkflowSpec(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (o *singleRowORM) SaveTriggerSubscriptions(context.Context, string, []byte) error {
	return nil
}

func (o *singleRowORM) row() job.WorkflowSpec {
	o.mu.Lock()
	defer o.mu.Unlock()
	return *o.spec
}

// A public registration claiming a private workflow's ID, with another owner,
// name and artifacts, must not leave its artifacts behind for the private
// workflow to pick up when it is reactivated after a pause.
func Test_workflowRegisteredEvent_ForeignClaimDoesNotPoisonPausedWorkflow(t *testing.T) {
	t.Parallel()

	const (
		victimBinaryURL   = "https://private.example/workflow.wasm"
		victimConfigURL   = "https://private.example/config.json"
		attackerBinaryURL = "https://attacker.example/workflow.wasm"
	)
	victimOwner := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	attackerOwner := []byte{0x11, 0x22, 0x33, 0x44}
	victimBinary := []byte("victim-binary")
	victimConfig := []byte("victim-config")
	attackerBinary := []byte("attacker-binary")

	rawID, err := pkgworkflows.GenerateWorkflowID(victimOwner, "victim-wf", victimBinary, victimConfig, "")
	require.NoError(t, err)
	victimID := types.WorkflowID(rawID)

	var fetchMu sync.Mutex
	fetches := map[string]int{}
	fetcher := func(_ context.Context, _ string, req ghcapabilities.Request) ([]byte, error) {
		fetchMu.Lock()
		fetches[req.URL]++
		fetchMu.Unlock()
		switch req.URL {
		case victimBinaryURL:
			return []byte(base64.StdEncoding.EncodeToString(victimBinary)), nil
		case victimConfigURL:
			return victimConfig, nil
		case attackerBinaryURL:
			return []byte(base64.StdEncoding.EncodeToString(attackerBinary)), nil
		}
		return nil, fmt.Errorf("unexpected fetch: %s", req.URL)
	}
	fetchCount := func(url string) int {
		fetchMu.Lock()
		defer fetchMu.Unlock()
		return fetches[url]
	}

	lggr := logger.TestLogger(t)
	orm := &singleRowORM{}
	store, err := artifacts.NewStore(lggr, orm, fetcher, nil, clockwork.NewFakeClock(), workflowkey.Key{}, custmsg.NewLabeler(), limits.Factory{Logger: lggr})
	require.NoError(t, err)
	h := newMeteringTestHandler(t, store, nil)

	victim := WorkflowRegisteredEvent{
		Status:        WorkflowStatusActive,
		WorkflowID:    victimID,
		WorkflowOwner: victimOwner,
		WorkflowName:  "victim-wf",
		BinaryURL:     victimBinaryURL,
		ConfigURL:     victimConfigURL,
		Source:        "grpc:private-registry",
	}
	attacker := WorkflowRegisteredEvent{
		Status:        WorkflowStatusActive,
		WorkflowID:    victimID,
		WorkflowOwner: attackerOwner,
		WorkflowName:  "attacker-wf",
		BinaryURL:     attackerBinaryURL,
		Source:        "contract:public",
	}

	require.NoError(t, h.workflowRegisteredEvent(t.Context(), victim))
	_, ok := h.engineRegistry.Get(victimID)
	require.True(t, ok)

	require.NoError(t, h.workflowPausedEvent(t.Context(), WorkflowPausedEvent{WorkflowID: victimID}))
	_, ok = h.engineRegistry.Get(victimID)
	require.False(t, ok)
	require.Empty(t, orm.row().Workflow)

	err = h.workflowRegisteredEvent(t.Context(), attacker)
	require.ErrorContains(t, err, "workflowID mismatch")
	assert.Equal(t, ActivationNonRetryable, classifyActivationError(err))
	assert.Equal(t, 1, fetchCount(attackerBinaryURL))

	row := orm.row()
	assert.Equal(t, hex.EncodeToString(victimOwner), row.WorkflowOwner, "rejected claim must not overwrite the row")
	assert.Equal(t, "victim-wf", row.WorkflowName)
	assert.Empty(t, row.Workflow, "rejected claim must not persist its artifacts")

	require.NoError(t, h.workflowRegisteredEvent(t.Context(), victim))
	_, ok = h.engineRegistry.Get(victimID)
	assert.True(t, ok, "victim must reactivate")
	assert.Equal(t, 2, fetchCount(victimBinaryURL))
	assert.Equal(t, hex.EncodeToString(victimBinary), orm.row().Workflow)
}

// When the stored row belongs to another identity, its artifacts must be
// replaced by a fresh fetch rather than reused, even if they are present.
func Test_workflowRegisteredEvent_IdentityMismatchRefetchesArtifacts(t *testing.T) {
	t.Parallel()

	owner := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	rawID, err := pkgworkflows.GenerateWorkflowID(owner, "wf-name", []byte("binary"), []byte("config"), "")
	require.NoError(t, err)
	wfID := types.WorkflowID(rawID)

	stub := &stubWorkflowArtifactsStore{
		persistUpserts: true,
		spec: &job.WorkflowSpec{
			WorkflowID:    wfID.Hex(),
			WorkflowOwner: "11223344",
			WorkflowName:  "other-name",
			Workflow:      hex.EncodeToString([]byte("foreign-binary")),
			Status:        job.WorkflowSpecStatusPaused,
		},
	}
	h := newMeteringTestHandler(t, stub, nil)

	require.NoError(t, h.workflowRegisteredEvent(t.Context(), WorkflowRegisteredEvent{
		Status:        WorkflowStatusActive,
		WorkflowID:    wfID,
		WorkflowOwner: owner,
		WorkflowName:  "wf-name",
	}))
	assert.Equal(t, int32(1), stub.fetchCalls.Load())
	assert.Equal(t, hex.EncodeToString([]byte("binary")), stub.spec.Workflow)
	assert.Equal(t, hex.EncodeToString(owner), stub.spec.WorkflowOwner)
	_, ok := h.engineRegistry.Get(wfID)
	assert.True(t, ok)
}
