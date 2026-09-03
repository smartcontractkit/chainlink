package image_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/image"
)

func TestResolve_Public(t *testing.T) {
	t.Parallel()
	res, err := image.Resolve(image.ResolveOptions{
		ECRType:        "public",
		RepositoryPath: "chainlink",
		ImageTag:       "v2.1.0",
	})
	require.NoError(t, err)
	assert.Equal(t, "public.ecr.aws/chainlink:v2.1.0", res)
}

func TestResolve_SDLC(t *testing.T) {
	t.Parallel()
	res, err := image.Resolve(image.ResolveOptions{
		ECRType:          "sdlc",
		RepositoryPath:   "chainlink-integration-tests",
		ImageTag:         "v2.1.0",
		AWSAccountNumber: "123456789012",
		AWSRegion:        "us-west-2",
	})
	require.NoError(t, err)
	assert.Equal(t, "123456789012.dkr.ecr.us-west-2.amazonaws.com/chainlink-integration-tests:v2.1.0", res)
}

func TestResolve_Normalization(t *testing.T) {
	t.Parallel()
	res, err := image.Resolve(image.ResolveOptions{
		ECRType:          "  SDLC  ",
		RepositoryPath:   "  ChainLink-Integration-Tests  ",
		ImageTag:         "  V2.1.0-customTag  ",
		AWSAccountNumber: "  123456789012  ",
		AWSRegion:        "  US-WEST-2  ",
	})
	require.NoError(t, err)
	assert.Equal(t, "123456789012.dkr.ecr.us-west-2.amazonaws.com/chainlink-integration-tests:V2.1.0-customTag", res)
}

func TestResolve_ValidationErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		opts image.ResolveOptions
		err  string
	}{
		{
			name: "missing ecr type",
			opts: image.ResolveOptions{RepositoryPath: "chainlink", ImageTag: "v2.0"},
			err:  "'ECR_TYPE' must be set and non-empty",
		},
		{
			name: "invalid ecr type",
			opts: image.ResolveOptions{ECRType: "foo", RepositoryPath: "chainlink", ImageTag: "v2.0"},
			err:  "invalid 'ECR_TYPE': 'foo'",
		},
		{
			name: "missing repo path",
			opts: image.ResolveOptions{ECRType: "public", ImageTag: "v2.0"},
			err:  "'CHAINLINK_IMAGE_REPO_PATH' must be set and non-empty",
		},
		{
			name: "missing tag",
			opts: image.ResolveOptions{ECRType: "public", RepositoryPath: "chainlink"},
			err:  "'CHAINLINK_IMAGE_TAG' must be set and non-empty",
		},
		{
			name: "missing aws credentials for sdlc",
			opts: image.ResolveOptions{ECRType: "sdlc", RepositoryPath: "chainlink", ImageTag: "v2.0"},
			err:  "both 'AWS_ACCOUNT_NUMBER' and 'AWS_REGION' must be set",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := image.Resolve(tc.opts)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.err)
		})
	}
}
