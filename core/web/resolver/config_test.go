package resolver

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

var (
	configEmptyEffective string
	configFull           string
	configFullEffective  string
	configMulti          string
	configMultiEffective string
)

func init() {
	d := "../../services/chainlink/testdata"

	emptyEffectivePath := filepath.Join(d, "config-empty-effective.toml")
	configEmptyEffective = mustRead(emptyEffectivePath)

	fullPath := filepath.Join(d, "config-full.toml")
	configFull = mustRead(fullPath)

	fullEffectivePath := filepath.Join(d, "config-full-effective.toml")
	configFullEffective = mustRead(fullEffectivePath)

	multiPath := filepath.Join(d, "config-multi-chain.toml")
	configMulti = mustRead(multiPath)

	multiEffectivePath := filepath.Join(d, "config-multi-chain-effective.toml")
	configMultiEffective = mustRead(multiEffectivePath)
}

func mustRead(p string) string {
	data, err := os.ReadFile(p)
	if err != nil {
		panic(fmt.Sprintf("Failed to read %s: %v", p, err))
	}
	return string(data)
}

func TestResolver_ConfigV2(t *testing.T) {
	t.Parallel()

	query := `
		query FetchConfigV2 {
			configv2 {
				user
				effective
			}
	  	}`

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "configv2"),
		{
			name:          "empty",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				opts := chainlink.GeneralConfigOpts{}
				cfg, err := opts.New()
				require.NoError(t, err)
				f.App.On("GetConfig").Return(cfg)
			},
			query:  query,
			result: fmt.Sprintf(`{"configv2":{"user":"","effective":%s}}`, mustJSONMarshal(t, configEmptyEffective)),
		},
		{
			name:          "full",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				opts := chainlink.GeneralConfigOpts{
					ConfigStrings:  []string{configFull},
					SecretsStrings: []string{},
				}
				cfg, err := opts.New()
				require.NoError(t, err)
				f.App.On("GetConfig").Return(cfg)
			},
			query:  query,
			result: fmt.Sprintf(`{"configv2":{"user":%s,"effective":%s}}`, mustJSONMarshal(t, configFull), mustJSONMarshal(t, configFullEffective)),
		},
		{
			name:          "partial",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				opts := chainlink.GeneralConfigOpts{
					ConfigStrings:  []string{configMulti},
					SecretsStrings: []string{},
				}
				cfg, err := opts.New()
				require.NoError(t, err)
				f.App.On("GetConfig").Return(cfg)
			},
			query:  query,
			result: fmt.Sprintf(`{"configv2":{"user":%s,"effective":%s}}`, mustJSONMarshal(t, configMulti), mustJSONMarshal(t, configMultiEffective)),
		},
	}

	RunGQLTests(t, testCases)
}

func mustJSONMarshal(t *testing.T, s string) string {
	b, err := json.Marshal(s)
	require.NoError(t, err)
	return string(b)
}
