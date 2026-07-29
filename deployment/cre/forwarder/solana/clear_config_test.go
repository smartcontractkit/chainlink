package solana

import (
	"crypto/ecdsa"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	cldftesthelpers "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils/testhelpers"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/internal/soltestutils"
)

func TestClearForwarderConfig(t *testing.T) {
	t.Parallel()

	programsPath := t.TempDir()
	te := setupForwarderTestEnv(t, programsPath, soltestutils.LoadKeystonePrograms(t, programsPath), false)

	clearTask := runtime.ChangesetTask(ClearForwarderConfigs{},
		&ClearForwarderConfigRequest{
			DonID:         te.DON.ID,
			ConfigVersion: te.DON.Version,
			Qualifier:     testQualifier,
			Version:       testVersion,
		},
	)

	t.Run("verify preconditions", func(t *testing.T) {
		valid := ClearForwarderConfigRequest{
			DonID:         te.DON.ID,
			ConfigVersion: te.DON.Version,
			Qualifier:     testQualifier,
			Version:       testVersion,
			Chains:        map[uint64]struct{}{te.Selector: {}},
		}

		tests := []struct {
			name    string
			mutate  func(*ClearForwarderConfigRequest)
			wantErr string
		}{
			{
				name:   "valid request",
				mutate: func(*ClearForwarderConfigRequest) {},
			},
			{
				name:   "no chains means all chains",
				mutate: func(req *ClearForwarderConfigRequest) { req.Chains = nil },
			},
			{
				name:    "invalid version",
				mutate:  func(req *ClearForwarderConfigRequest) { req.Version = "not-a-version" },
				wantErr: "invalid semantic version",
			},
			{
				name:    "zero DON ID",
				mutate:  func(req *ClearForwarderConfigRequest) { req.DonID = 0 },
				wantErr: "DON ID must be non-zero",
			},
			{
				name:    "zero config version",
				mutate:  func(req *ClearForwarderConfigRequest) { req.ConfigVersion = 0 },
				wantErr: "config version must be non-zero",
			},
			{
				name:    "unknown chain selector",
				mutate:  func(req *ClearForwarderConfigRequest) { req.Chains = map[uint64]struct{}{1: {}} },
				wantErr: "solana chain not found for chain selector 1",
			},
			{
				name:    "unknown qualifier",
				mutate:  func(req *ClearForwarderConfigRequest) { req.Qualifier = "does-not-exist" },
				wantErr: "failed get forwarder for chain selector",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := valid
				tt.mutate(&req)

				err := ClearForwarderConfigs{}.VerifyPreconditions(te.Runtime.Environment(), &req)
				if tt.wantErr == "" {
					require.NoError(t, err)
					return
				}
				require.ErrorContains(t, err, tt.wantErr)
			})
		}
	})

	t.Run("clearing an unconfigured don fails", func(t *testing.T) {
		err := te.Runtime.Exec(
			runtime.ChangesetTask(DeployForwarder{},
				&DeployForwarderRequest{
					ChainSel:  te.Selector,
					Qualifier: testQualifier,
					Version:   testVersion,
				},
			),
		)
		require.NoError(t, err)

		// The DON was never configured, so there is no oracles config account to close.
		require.ErrorContains(t, te.Runtime.Exec(clearTask), "no oracles config")
	})

	t.Run("clears the config of a configured don", func(t *testing.T) {
		err := te.Runtime.Exec(
			runtime.ChangesetTask(ConfigureForwarders{},
				&ConfigureForwarderRequest{
					DON:       te.DON,
					Version:   testVersion,
					Qualifier: testQualifier,
				},
			),
		)
		require.NoError(t, err)
		requireOraclesConfig(t, te.Runtime.Environment(), te.Selector, te.DON.ID, te.DON.Version, true)

		require.NoError(t, te.Runtime.Exec(clearTask))

		out := te.Runtime.State().Outputs[clearTask.ID()]
		require.NotNil(t, out, "changeset output should not be nil")
		require.Empty(t, out.MCMSTimelockProposals, "should not have MCMS proposals when not using MCMS")

		requireOraclesConfig(t, te.Runtime.Environment(), te.Selector, te.DON.ID, te.DON.Version, false)
	})
}

func TestClearForwarderConfigWithMCMS(t *testing.T) {
	t.Parallel()

	programsPath := t.TempDir()
	te := setupForwarderTestEnv(t, programsPath, soltestutils.LoadKeystonePrograms(t, programsPath), true)

	mcmsCfg := &cldfproposalutils.TimelockConfig{MinDelay: time.Second}
	signers := []*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}

	// Deploy the forwarder, hand it over to MCMS and configure the DON through a proposal.
	err := te.Runtime.Exec(
		runtime.ChangesetTask(DeployForwarder{},
			&DeployForwarderRequest{
				ChainSel:  te.Selector,
				Qualifier: testQualifier,
				Version:   testVersion,
			},
		),
		runtime.ChangesetTask(TransferOwnershipForwarder{},
			&TransferOwnershipForwarderRequest{
				ChainSel:  te.Selector,
				MCMSCfg:   *mcmsCfg,
				Qualifier: testQualifier,
				Version:   testVersion,
			},
		),
		runtime.SignAndExecuteProposalsTask(signers),
	)
	require.NoError(t, err)

	err = te.Runtime.Exec(
		runtime.ChangesetTask(ConfigureForwarders{},
			&ConfigureForwarderRequest{
				DON:       te.DON,
				Version:   testVersion,
				Qualifier: testQualifier,
				MCMS:      mcmsCfg,
			},
		),
		runtime.SignAndExecuteProposalsTask(signers),
	)
	require.NoError(t, err)
	requireOraclesConfig(t, te.Runtime.Environment(), te.Selector, te.DON.ID, te.DON.Version, true)

	clearTask := runtime.ChangesetTask(ClearForwarderConfigs{},
		&ClearForwarderConfigRequest{
			DonID:         te.DON.ID,
			ConfigVersion: te.DON.Version,
			Qualifier:     testQualifier,
			Version:       testVersion,
			MCMS:          mcmsCfg,
		},
	)
	require.NoError(t, te.Runtime.Exec(clearTask, runtime.SignAndExecuteProposalsTask(signers)))

	out := te.Runtime.State().Outputs[clearTask.ID()]
	require.NotNil(t, out, "changeset output should not be nil")
	require.Len(t, out.MCMSTimelockProposals, 1, "should have one MCMS proposal per chain")

	requireOraclesConfig(t, te.Runtime.Environment(), te.Selector, te.DON.ID, te.DON.Version, false)
}
