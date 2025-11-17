package job_types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
)

func TestJobSpecTemplate_UnmarshalJSON(t *testing.T) {
	t.Run("cron string", func(t *testing.T) {
		var in jobs.ProposeJobSpecInput
		js := `{"environment":"e","domain":"d","don_name":"don","don_filters":[],"job_name":"j","template":"cron","inputs":{}}`
		require.NoError(t, json.Unmarshal([]byte(js), &in))
		require.Equal(t, job_types.Cron, in.Template)
	})

	t.Run("invalid string", func(t *testing.T) {
		var in jobs.ProposeJobSpecInput
		js := `{"environment":"e","domain":"d","don_name":"don","don_filters":[],"job_name":"j","template":"nope","inputs":{}}`
		err := json.Unmarshal([]byte(js), &in)
		require.Error(t, err)
	})

	t.Run("numeric legacy", func(t *testing.T) {
		var in jobs.ProposeJobSpecInput
		js := `{"environment":"e","domain":"d","don_name":"don","don_filters":[],"job_name":"j","template":0,"inputs":{}}`
		require.NoError(t, json.Unmarshal([]byte(js), &in))
		require.Equal(t, job_types.Cron, in.Template)
	})
}

func TestJobSpecTemplate_UnmarshalYAML(t *testing.T) {
	t.Run("cron string", func(t *testing.T) {
		var in jobs.ProposeJobSpecInput
		yml := "environment: e\ndomain: d\ndon_name: don\ndon_filters: []\njob_name: j\ntemplate: cron\ninputs: {}\n"
		require.NoError(t, yaml.Unmarshal([]byte(yml), &in))
		require.Equal(t, job_types.Cron, in.Template)
	})

	t.Run("invalid string", func(t *testing.T) {
		var in jobs.ProposeJobSpecInput
		yml := "environment: e\ndomain: d\ndon_name: don\ndon_filters: []\njob_name: j\ntemplate: nope\ninputs: {}\n"
		err := yaml.Unmarshal([]byte(yml), &in)
		require.Error(t, err)
	})

	t.Run("invalid type", func(t *testing.T) {
		var in jobs.ProposeJobSpecInput
		yml := "environment: e\ndomain: d\ndon_name: don\ndon_filters: []\njob_name: j\ntemplate: 0\ninputs: {}\n"
		require.Error(t, yaml.Unmarshal([]byte(yml), &in))
	})
}
