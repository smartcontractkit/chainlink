package workflows_test

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
)

func TestWorkflowSpecFactory_ToSpec(t *testing.T) {
	t.Parallel()

	anyData := "any data"
	anyConfig := []byte("any config")
	anySpec := sdk.WorkflowSpec{Name: "name", Owner: "owner"}

	t.Run("delegates to factory and calculates CID", func(t *testing.T) {
		factory := workflows.WorkflowSpecFactory{
			job.YamlSpec: mockSdkSpecFactory{t: t, Spec: anySpec},
		}

		results, cid, err := factory.ToSpec(anyData, anyConfig, job.YamlSpec)

		require.NoError(t, err)
		assert.Equal(t, anySpec, results)

		sha256Hash := sha256.New()
		sha256Hash.Write([]byte(anyData))
		sha256Hash.Write(anyConfig)
		expectedCid := fmt.Sprintf("%x", sha256Hash.Sum(nil))
		assert.Equal(t, expectedCid, cid)

	})

	t.Run("CID without config matches", func(t *testing.T) {
		factory := workflows.WorkflowSpecFactory{
			job.YamlSpec: mockSdkSpecFactory{t: t, noConfig: true, Spec: anySpec},
		}
		results, cid, err := factory.ToSpec(anyData, nil, job.YamlSpec)
		require.NoError(t, err)

		assert.Equal(t, anySpec, results)

		sha256Hash := sha256.New()
		sha256Hash.Write([]byte(anyData))
		expectedCid := fmt.Sprintf("%x", sha256Hash.Sum(nil))
		assert.Equal(t, expectedCid, cid)
	})

	t.Run("returns errors from sdk factory", func(t *testing.T) {
		anyErr := errors.New("nope")
		factory := workflows.WorkflowSpecFactory{
			job.YamlSpec: mockSdkSpecFactory{t: t, Err: anyErr},
		}

		_, _, err := factory.ToSpec(anyData, anyConfig, job.YamlSpec)
		assert.Equal(t, anyErr, err)
	})

	t.Run("returns an error if the type is not supported", func(t *testing.T) {
		factory := workflows.WorkflowSpecFactory{
			job.YamlSpec: mockSdkSpecFactory{t: t, Spec: anySpec},
		}

		_, _, err := factory.ToSpec(anyData, anyConfig, "unsupported")
		assert.Error(t, err)
	})
}

type mockSdkSpecFactory struct {
	t        *testing.T
	noConfig bool
	Spec     sdk.WorkflowSpec
	Err      error
}

func (f mockSdkSpecFactory) GetRawSpec(wf string) ([]byte, error) {
	return []byte(wf), nil
}

func (f mockSdkSpecFactory) GetSpec(rawSpec, config []byte) (sdk.WorkflowSpec, error) {
	assert.ElementsMatch(f.t, rawSpec, []byte("any data"))
	if f.noConfig {
		assert.Nil(f.t, config)
	} else {
		assert.ElementsMatch(f.t, config, []byte("any config"))
	}

	return f.Spec, f.Err
}
