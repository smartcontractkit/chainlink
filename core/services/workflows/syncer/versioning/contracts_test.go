package versioning_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/balance_reader"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/log_emitter"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/capabilities/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncer/versioning"
)

func TestContracts_TypeAndVersion(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		th := testutils.NewEVMBackendTH(t)

		// Deploy a contract that has typeAndVersion
		addr, _, _, err := balance_reader.DeployBalanceReader(th.ContractsOwner, th.Backend.Client())
		th.Backend.Commit()
		th.Backend.Commit()
		th.Backend.Commit()
		require.NoError(t, err)

		contractReaderFactory := func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
			return th.NewContractReader(ctx, t, bytes)
		}

		contractType, version, err := versioning.TypeAndVersion(t.Context(), addr.String(), contractReaderFactory)
		require.NoError(t, err)

		expectedVersion, err := semver.NewVersion("1.0.0")
		require.NoError(t, err)

		require.Equal(t, versioning.ContractType("BalanceReader"), contractType)
		require.Equal(t, expectedVersion, &version)
	})

	t.Run("contract does not have typeAndVersion", func(t *testing.T) {
		t.Parallel()

		th := testutils.NewEVMBackendTH(t)

		// Deploy a contract that does not have typeAndVersion
		addr, _, _, err :=
			log_emitter.DeployLogEmitter(th.ContractsOwner, th.Backend.Client())
		th.Backend.Commit()
		require.NoError(t, err)

		contractReaderFactory := func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
			return th.NewContractReader(ctx, t, bytes)
		}

		contractType, version, err := versioning.TypeAndVersion(t.Context(), addr.String(), contractReaderFactory)
		require.NoError(t, err)

		expectedVersion, err := semver.NewVersion("1.0.0")
		require.NoError(t, err)

		require.Equal(t, versioning.Unknown, contractType)
		require.Equal(t, expectedVersion, &version)
	})

	t.Run("errors on empty address", func(t *testing.T) {
		t.Parallel()

		th := testutils.NewEVMBackendTH(t)

		contractReaderFactory := func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
			return th.NewContractReader(ctx, t, bytes)
		}

		_, _, err := versioning.TypeAndVersion(t.Context(), "0x0000000000000000000000000000000000000000", contractReaderFactory)
		require.ErrorContains(t, err, "internal error: contract does not exist at address: 0x0000000000000000000000000000000000000000")
	})

	t.Run("errors on invalid contract reader factory", func(t *testing.T) {
		t.Parallel()

		contractReaderFactory := func(ctx context.Context, bytes []byte) (types.ContractReader, error) {
			return nil, nil
		}

		_, _, err := versioning.TypeAndVersion(t.Context(), "", contractReaderFactory)
		require.ErrorIs(t, versioning.ErrNoContractReader, err)
	})
}

func TestContracts_VerifyTypeAndVersion(t *testing.T) {
	t.Parallel()

	t.Run("incorrect type", func(t *testing.T) {
		t.Parallel()

		th := testutils.NewEVMBackendTH(t)

		addr, _, _, err := balance_reader.DeployBalanceReader(th.ContractsOwner, th.Backend.Client())
		require.NoError(t, err)
		th.Backend.Commit()
		th.Backend.Commit()
		th.Backend.Commit()

		contractReaderFactory := func(ctx context.Context, cfg []byte) (types.ContractReader, error) {
			return th.NewContractReader(ctx, t, cfg)
		}

		version, err := versioning.VerifyTypeAndVersion(t.Context(), addr.String(), contractReaderFactory, "SomeType")
		require.ErrorContains(t, err, "wrong contract type BalanceReader")
		require.Equal(t, semver.Version{}, version)
	})
}

func TestContracts_ParseTypeAndVersion(t *testing.T) {
	t.Parallel()

	t.Run("valid string", func(t *testing.T) {
		t.Parallel()
		contractType, version, err := versioning.ParseTypeAndVersion("SomeType 1.2.0")
		require.NoError(t, err)
		require.Equal(t, "SomeType", contractType)
		require.Equal(t, "1.2.0", version)
	})
	t.Run("invalid string - too short", func(t *testing.T) {
		t.Parallel()
		_, _, err := versioning.ParseTypeAndVersion("v1.2.0")
		require.ErrorContains(t, err, "invalid type and version v1.2.0")
	})
	t.Run("invalid string - too long", func(t *testing.T) {
		t.Parallel()
		_, _, err := versioning.ParseTypeAndVersion("SomeType WithMoreWords vv1.2.0")
		require.ErrorContains(t, err, "invalid type and version SomeType WithMoreWords vv1.2.0")
	})
	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		contractType, version, err := versioning.ParseTypeAndVersion("")
		require.NoError(t, err)
		require.Equal(t, versioning.Unknown, versioning.ContractType(contractType))
		require.Equal(t, "1.0.0", version)
	})
}

func TestContracts_RunWithRetries(t *testing.T) {
	t.Parallel()

	t.Run("returns an error if max retries are exceeded", func(t *testing.T) {
		t.Parallel()
		tries := 0
		fn := func() (versioning.ContractType, semver.Version, error) {
			tries++
			return "", semver.Version{}, errors.New("some error")
		}
		_, _, err := versioning.RunWithRetries(t.Context(), 10*time.Millisecond, 3, fn)
		require.ErrorContains(t, err, "max retries (3) reached, aborting")
		require.Equal(t, 4, tries)
	})

	t.Run("fails and then succeeds", func(t *testing.T) {
		t.Parallel()
		tries := 0
		fn := func() (versioning.ContractType, semver.Version, error) {
			if tries < 2 {
				tries++
				return "", semver.Version{}, errors.New("some error")
			}
			return "", semver.Version{}, nil
		}
		_, _, err := versioning.RunWithRetries(t.Context(), 10*time.Millisecond, 3, fn)
		require.NoError(t, err)
		require.Equal(t, 2, tries)
	})
}
