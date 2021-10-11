package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/adapters"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/tidwall/gjson"
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
	case models.InitiatorRunLog:
		return migrateRunLogJob(js)
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
	ps, pd, err := BuildTaskDAG(js, job.Cron)
	if err != nil {
		return jb, err
	}
	jb.PipelineSpec = &pipeline.Spec{
		DotDagSource: ps,
	}
	jb.Pipeline = *pd
	return jb, nil
}

func migrateRunLogJob(js models.JobSpec) (job.Job, error) {
	var jb job.Job
	initr := js.Initiators[0]
	jb = job.Job{
		Name: null.StringFrom(js.Name),
		DirectRequestSpec: &job.DirectRequestSpec{
			ContractAddress:          ethkey.EIP55AddressFromAddress(initr.InitiatorParams.Address),
			MinIncomingConfirmations: clnull.Uint32From(10),
			Requesters:               requesterWhitelist,
			CreatedAt:                js.CreatedAt,
			UpdatedAt:                js.UpdatedAt,
		},
		Type:          job.DirectRequest,
		SchemaVersion: 1,
		ExternalJobID: uuid.NewV4(),
	}
	ps, pd, err := BuildTaskDAG(js, job.DirectRequest)
	if err != nil {
		return jb, err
	}
	jb.PipelineSpec = &pipeline.Spec{
		DotDagSource: ps,
	}
	jb.Pipeline = *pd
	return jb, nil
}

func BuildTaskDAG(js models.JobSpec, tpe job.Type) (string, *pipeline.Pipeline, error) {
	replacements := make(map[string]string)
	dg := pipeline.NewGraph()
	var foundEthTx = false
	var last *pipeline.GraphNode

	if tpe == job.DirectRequest {
		attrs := map[string]string{
			"type":   "ethabidecodelog",
			"abi":    "OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes32 data)",
			"data":   "$(jobRun.logData)",
			"topics": "$(jobRun.logTopics)",
		}
		n := pipeline.NewGraphNode(dg.NewNode(), "decode_log", attrs)
		dg.AddNode(n)
		last = n

		/*
		   decode_log   [type=ethabidecodelog
		                 abi="OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes32 data)"
		                 data="$(jobRun.logData)"
		                 topics="$(jobRun.logTopics)"]
		*/
	}

	for i, ts := range js.Tasks {
		var n *pipeline.GraphNode
		switch ts.Type {
		case adapters.TaskTypeHTTPGet:
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
			replacements["\""+template+"\""] = string(marshal)
			attrs := map[string]string{
				"type":        pipeline.TaskTypeHTTP.String(),
				"method":      "GET",
				"requestData": template,
			}
			n = pipeline.NewGraphNode(dg.NewNode(), fmt.Sprintf("http_%d", i), attrs)

		case adapters.TaskTypeJSONParse:
			attrs := map[string]string{
				"type": pipeline.TaskTypeJSONParse.String(),
			}
			if ts.Params.Get("path").Exists() {

				path := ts.Params.Get("path")
				pathString := path.String()

				if path.IsArray() {
					var pathSegments []string
					path.ForEach(func(key, value gjson.Result) bool {
						pathSegments = append(pathSegments, value.String())
						return true
					})

					pathString = strings.Join(pathSegments, ",")
				}

				attrs["path"] = pathString
			} else {
				return "", nil, errors.New("no path param on jsonparse task")
			}
			n = pipeline.NewGraphNode(dg.NewNode(), fmt.Sprintf("jsonparse_%d", i), attrs)

		case adapters.TaskTypeMultiply:
			attrs := map[string]string{
				"type": pipeline.TaskTypeMultiply.String(),
			}
			if ts.Params.Get("times").Exists() {
				attrs["times"] = ts.Params.Get("times").String()
			} else {
				return "", nil, errors.New("no times param on multiply task")
			}
			n = pipeline.NewGraphNode(dg.NewNode(), fmt.Sprintf("multiply_%d", i), attrs)
		case adapters.TaskTypeEthUint256, adapters.TaskTypeEthInt256:
			// Do nothing. This is implicit in FMv2 / DR
		case adapters.TaskTypeEthTx:
			if tpe == job.DirectRequest {
				attrs := map[string]string{
					"type": "ethabiencode",
					"abi":  "(uint256 value)",
					//"data": <{ "value": $(multiply) }>,
				}
				n = pipeline.NewGraphNode(dg.NewNode(), fmt.Sprintf("encode_data_%d", i), attrs)
				dg.AddNode(n)
				if last != nil {
					dg.SetEdge(dg.NewEdge(last, n))
				}
				last = n
			}
			if tpe == job.DirectRequest {

				template := fmt.Sprintf("%%REQ_DATA_%v%%", i)
				attrs := map[string]string{
					"type": "ethabiencode",
					"abi":  "fulfillOracleRequest(bytes32 requestId, uint256 payment, address callbackAddress, bytes4 callbackFunctionId, uint256 expiration, bytes32 calldata data)",
					"data": template,
				}
				replacements["\""+template+"\""] = `{
"requestId":          $(decode_log.requestId),
"payment":            $(decode_log.payment),
"callbackAddress":    $(decode_log.callbackAddr),
"callbackFunctionId": $(decode_log.callbackFunctionId),
"expiration":         $(decode_log.cancelExpiration),
"data":               $(encode_data)
}
`

				n = pipeline.NewGraphNode(dg.NewNode(), fmt.Sprintf("encode_tx_%d", i), attrs)
				dg.AddNode(n)
				if last != nil {
					dg.SetEdge(dg.NewEdge(last, n))
				}
				last = n
			}
			attrs := map[string]string{
				"type": pipeline.TaskTypeETHTx.String(),
				"to":   js.Initiators[0].Address.String(),
				"data": fmt.Sprintf("$(%v)", last.DOTID()),
			}
			n = pipeline.NewGraphNode(dg.NewNode(), fmt.Sprintf("send_tx_%d", i), attrs)
			foundEthTx = true
		default:
			// assume it's a bridge task

			encodedValue1, err := encodeTemplate(ts.Params.Bytes())
			if err != nil {
				return "", nil, err
			}
			template1 := fmt.Sprintf("%%REQ_DATA_%v%%", i)
			i++
			replacements["\""+template1+"\""] = encodedValue1

			attrs1 := map[string]string{
				"type":  "merge",
				"right": template1,
				//"data": <{ "value": $(multiply) }>,
			}
			n = pipeline.NewGraphNode(dg.NewNode(), fmt.Sprintf("merge_%d", i), attrs1)
			dg.AddNode(n)
			if last != nil {
				dg.SetEdge(dg.NewEdge(last, n))
			}
			last = n
			template := fmt.Sprintf("%%REQ_DATA_%v%%", i)

			attrs := map[string]string{
				"type":        pipeline.TaskTypeBridge.String(),
				"name":        ts.Type.String(),
				"requestData": template,
			}
			replacements["\""+template+"\""] = fmt.Sprintf("{ \"data\": $(%v) }", last.DOTID())

			n = pipeline.NewGraphNode(dg.NewNode(), fmt.Sprintf("send_to_bridge_%d", i), attrs)
			i++
		}
		if n != nil {
			dg.AddNode(n)
			if last != nil {
				dg.SetEdge(dg.NewEdge(last, n))
			}
			last = n
		}
	}
	if !foundEthTx && tpe == job.DirectRequest {
		return "", nil, errors.New("expected ethtx in FM v1 / Runlog job spec")
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
		return "", nil, errors.Wrapf(err, "failed to generate pipeline from: \n%v", generatedDotDagSource)
	}
	return generatedDotDagSource, p, err
}

func encodeTemplate(bytes []byte) (string, error) {
	mapp := make(map[string]interface{})
	err := json.Unmarshal(bytes, &mapp)
	if err != nil {
		return "", err
	}
	marshal, err := json.Marshal(&mapp)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
