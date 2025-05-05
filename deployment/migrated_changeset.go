package deployment

import (
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// ChangeSets

type (
	ChangeSet[C any]            = deployment.ChangeSet[C]
	ChangeLogic[C any]          = deployment.ChangeLogic[C]
	PreconditionVerifier[C any] = deployment.PreconditionVerifier[C]
	ChangeSetV2[C any]          = deployment.ChangeSetV2[C]
)

var (
	ErrInvalidConfig      = deployment.ErrInvalidConfig
	ErrInvalidEnvironment = deployment.ErrInvalidEnvironment
)

// Changeset Output

type (
	ChangesetOutput = deployment.ChangesetOutput
	ProposedJob     = deployment.ProposedJob
)

var (
	MergeChangesetOutput = deployment.MergeChangesetOutput
)

// ViewState

type (
	ViewState   = deployment.ViewState
	ViewStateV2 = deployment.ViewStateV2
)
