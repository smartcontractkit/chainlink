package pipeline_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/stretchr/testify/require"
)

func TestResolveMap(t *testing.T) {
	t.Parallel()

	vars := pipeline.Vars{
		"foo": map[string]interface{}{
			"abc": "def",
		},
		"bar": "123",
	}

	// input := map[string]interface{}{
	// 	"chain": "$(foo)",
	// 	"link": map[string]interface{}{
	// 		"sergey": "$(foo.abc)",
	// 		"$(bar)": "satoshi",
	// 	},
	// }

	input := `
    {
        "chain": $(foo),
        "link": {
            $(bar): "satoshi",
            "sergey": $(foo.abc)
        }
    }
    `

	expected := pipeline.MapParam{
		"chain": map[string]interface{}{
			"abc": "def",
		},
		"link": map[string]interface{}{
			"sergey": "def",
			"123":    "satoshi",
		},
	}

	var got pipeline.MapParam
	err := got.UnmarshalPipelineParam(input, vars)

	// got, err := pipeline.ExportedResolveMap(input, vars)
	require.NoError(t, err)

	bs, _ := json.MarshalIndent(got, "", "    ")
	fmt.Println(string(bs))

	require.Equal(t, expected, got)
}
