package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gopkg.in/guregu/null.v4"
)

// MigrateJobSpec - Does not support mixed initiator types.
func MigrateJobSpec(js models.JobSpec) (job.Job, error) {
	var jb job.Job
	if len(js.Initiators) == 0 {
		return jb, errors.New("initiator required to migrate job")
	}
	v1JobType := js.Initiators[0].Type
	switch v1JobType {
	case models.InitiatorCron:
		return migrateCronJob(js)
	default:
		return jb, errors.Wrapf(errors.New("Invalid initiator type"), "%v", v1JobType)
	}
}

func migrateCronJob(js models.JobSpec) (job.Job, error) {
	var jb job.Job
	initr := js.Initiators[0]
	jb = job.Job{
		Name: null.StringFrom(js.Name),
		CronSpec: &job.CronSpec{
			CronSchedule: string(initr.InitiatorParams.Schedule),
			CreatedAt:    js.CreatedAt,
			UpdatedAt:    js.UpdatedAt,
		},
		Type:          job.Cron,
		SchemaVersion: 1,
		ExternalJobID: uuid.NewV4(),
	}
	ps, pd, err := BuildTaskDAG(js)
	if err != nil {
		return jb, err
	}
	jb.PipelineSpec = &pipeline.Spec{
		DotDagSource: ps,
	}
	jb.Pipeline = *pd
	return jb, nil
}

func BuildTaskDAG(js models.JobSpec) (string, *pipeline.Pipeline, error) {
	replacements := make(map[string]string)
	i := 0

	dg := pipeline.NewGraph()
	for _, ts := range js.Tasks {
		switch ts.Type {
		case adapters.TaskTypeMultiply:
		case adapters.TaskTypeEthUint256, adapters.TaskTypeEthInt256:
			// Do nothing. This is implicit in FMv2 / Cron
		case adapters.TaskTypeEthTx:
			// Do nothing. This is implicit in FMV2 / Cron
		default:
			mapp := make(map[string]interface{})
			err := json.Unmarshal(ts.Params.Bytes(), &mapp)
			if err != nil {
				return "", nil, err
			}
			marshal, err := json.Marshal(&mapp)
			if err != nil {
				return "", nil, err
			}

			template := fmt.Sprintf("%%REQ_DATA_%v%%", i)
			attrs := map[string]string{
				"type":        pipeline.TaskTypeBridge.String(),
				"name":        ts.Type.String(),
				"requestData": template,
			}
			replacements["\""+template+"\""] = string(marshal)

			n := pipeline.NewGraphNode(dg.NewNode(), "send_to_bridge", attrs)
			dg.AddNode(n)
			i++
		}
	}

	s, err := dot.Marshal(dg, "", "", "")
	if err != nil {
		return "", nil, err
	}

	// Double check we can unmarshal it
	generatedDotDagSource := string(s)
	generatedDotDagSource = strings.Replace(generatedDotDagSource, "strict digraph {", "", 1)
	generatedDotDagSource = strings.Replace(generatedDotDagSource, "\n// Node definitions.\n", "", 1)
	generatedDotDagSource = strings.Replace(generatedDotDagSource, "\n", "\n\t", 100)

	for key := range replacements {
		generatedDotDagSource = strings.Replace(generatedDotDagSource, key, "<"+replacements[key]+">", 1)
	}
	generatedDotDagSource = generatedDotDagSource[:len(generatedDotDagSource)-1] // Remove final }
	p, err := pipeline.Parse(generatedDotDagSource)
	if err != nil {
		return "", nil, errors.Wrapf(err, "failed to genreate pipeline from: \n%v", generatedDotDagSource)
	}
	return generatedDotDagSource, p, err
}
