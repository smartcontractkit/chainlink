package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/store/models"
	webpresenters "github.com/smartcontractkit/chainlink/core/web/presenters"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrInvalidInitiatorType = errors.New("invalid initiator type")
)

type MigrateController struct {
	App chainlink.Application
}

// Creates a v2 job from a v1 job ID.
// Example:
//  "POST <application>/migrate/123e4567-e89b-12d3-a456-426614174000"
// Where "123e4567-e89b-12d3-a456-426614174000" is a v1 job ID.
func (mc *MigrateController) Migrate(c *gin.Context) {
	js := models.JobSpec{}
	err := js.SetID(c.Param("ID"))
	if err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}
	js, err = mc.App.GetStore().FindJobSpec(js.ID)
	if err != nil {
		jsonAPIError(c, http.StatusNotFound, err)
		return
	}
	jbV2, err := MigrateJobSpec(mc.App.GetStore().Config, js)
	if err != nil {
		if errors.Cause(err) == ErrInvalidInitiatorType {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jb, err := mc.App.AddJobV2(c, jbV2, jbV2.Name)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	logger.Infow(fmt.Sprintf("Successfully migrated job %v into %v", js.ID, jb.ID), "v1 job", js, "v2 job", jb)
	// If the migration went well, archive the v1
	if err := mc.App.ArchiveJob(js.ID); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}
	jsonAPIResponse(c, webpresenters.NewJobResource(jb), jb.Type.String())
}

// Does not support mixed initiator types.
func MigrateJobSpec(c *config.Config, js models.JobSpec) (job.Job, error) {
	var jb job.Job
	if len(js.Initiators) == 0 {
		return jb, errors.New("initiator required to migrate job")
	}
	v1JobType := js.Initiators[0].Type
	switch v1JobType {
	case models.InitiatorFluxMonitor:
		if !c.FeatureFluxMonitorV2() {
			return jb, errors.New("need to enable FEATURE_FLUX_MONITOR_V2=true to migrate FM jobs")
		}
		return migrateFluxMonitorJob(js)
	case models.InitiatorCron:
		return migrateCronJob(js)
	case models.InitiatorRunLog:
		return migrateRunLogJob(js)
	default:
		return jb, errors.Wrapf(ErrInvalidInitiatorType, "%v", v1JobType)
	}
}

func migrateFluxMonitorJob(js models.JobSpec) (job.Job, error) {
	var jb job.Job
	initr := js.Initiators[0]
	ca := ethkey.EIP55AddressFromAddress(initr.Address)
	jb = job.Job{
		Name: null.StringFrom(js.Name),
		FluxMonitorSpec: &job.FluxMonitorSpec{
			ContractAddress:   ca,
			Threshold:         initr.Threshold,
			AbsoluteThreshold: initr.AbsoluteThreshold,
			PollTimerPeriod:   initr.PollTimer.Period.Duration(),
			PollTimerDisabled: initr.PollTimer.Disabled,
			IdleTimerPeriod:   initr.IdleTimer.Duration.Duration(),
			IdleTimerDisabled: initr.IdleTimer.Disabled,
			MinPayment:        js.MinPayment,
			CreatedAt:         js.CreatedAt,
			UpdatedAt:         js.UpdatedAt,
		},
		Type:          job.FluxMonitor,
		SchemaVersion: 1,
		ExternalJobID: uuid.NewV4(),
	}
	ps, pd, err := BuildFMTaskDAG(js)
	if err != nil {
		return jb, err
	}
	jb.PipelineSpec = &pipeline.Spec{
		DotDagSource: ps,
	}
	jb.Pipeline = *pd
	return jb, nil
}

func BuildFMTaskDAG(js models.JobSpec) (string, *pipeline.Pipeline, error) {
	dg := pipeline.NewGraph()
	// First add the feeds as parallel HTTP tasks,
	// which all coalesce into a single median task.
	var medianTask = pipeline.NewGraphNode(dg.NewNode(), "median", map[string]string{
		"type": pipeline.TaskTypeMedian.String(),
	})
	dg.AddNode(medianTask)
	for i, feed := range js.Initiators[0].Feeds.Array() {
		// Apparently there are *no* urls directly used in production, its all bridges.
		// Support anyways just in case someone was using it without our knowledge.
		// ALL fm jobs are POSTs see
		// https://github.com/smartcontractkit/chainlink/blob/e5957895e3aa4947c2ddb5a4a8525041639962e9/core/services/fluxmonitor/fetchers.go#L67
		var httpTask *pipeline.GraphNode
		if feed.IsObject() && feed.Get("bridge").Exists() {
			httpTask = pipeline.NewGraphNode(dg.NewNode(), fmt.Sprintf("feed%d", i), map[string]string{
				"type":        pipeline.TaskTypeBridge.String(),
				"method":      "POST",
				"name":        feed.Get("bridge").String(),
				"requestData": js.Initiators[0].InitiatorParams.RequestData.String(),
			})
		} else {
			httpTask = pipeline.NewGraphNode(dg.NewNode(), fmt.Sprintf("feed%d", i), map[string]string{
				"type":        pipeline.TaskTypeHTTP.String(),
				"method":      "POST",
				"url":         feed.String(),
				"requestData": js.Initiators[0].InitiatorParams.RequestData.String(),
			})
		}
		dg.AddNode(httpTask)
		// We always implicity parse {"data": {"result": X}}
		parseTask := pipeline.NewGraphNode(dg.NewNode(), fmt.Sprintf("jsonparse%d", i), map[string]string{
			"type": pipeline.TaskTypeJSONParse.String(),
			"path": "data,result",
		})
		dg.AddNode(parseTask)
		dg.SetEdge(dg.NewEdge(httpTask, parseTask))
		dg.SetEdge(dg.NewEdge(parseTask, medianTask))
	}
	// Now add tasks linearly from the median task.
	var foundEthTx = false
	var last = medianTask
	for i, ts := range js.Tasks {
		switch ts.Type {
		case adapters.TaskTypeMultiply:
			attrs := map[string]string{
				"type": pipeline.TaskTypeMultiply.String(),
			}
			if ts.Params.Get("times").Exists() {
				attrs["times"] = ts.Params.Get("times").String()
			} else {
				return "", nil, errors.New("no times param on multiply task")
			}
			n := pipeline.NewGraphNode(dg.NewNode(), fmt.Sprintf("multiply%d", i), attrs)
			dg.AddNode(n)
			dg.SetEdge(dg.NewEdge(last, n))
			last = n
		case adapters.TaskTypeEthUint256, adapters.TaskTypeEthInt256:
			// Do nothing. This is implicit in FMv2.
		case adapters.TaskTypeEthTx:
			// Do nothing. This is implicit in FMV2.
			foundEthTx = true
		default:
			return "", nil, errors.Errorf("unsupported task type %v", ts.Type)
		}
	}
	if !foundEthTx {
		return "", nil, errors.New("expected ethtx in FM v1 job spec")
	}
	s, err := dot.Marshal(dg, "", "", "")
	if err != nil {
		return "", nil, err
	}

	// Double check we can unmarshal it
	generatedDotDagSource := string(s)
	generatedDotDagSource = strings.Replace(generatedDotDagSource, "strict digraph {", "", 1)
	generatedDotDagSource = generatedDotDagSource[:len(generatedDotDagSource)-1] // Remove final }
	p, err := pipeline.Parse(generatedDotDagSource)
	if err != nil {
		return "", nil, err
	}
	return generatedDotDagSource, p, err
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
			Requesters:               initr.Requesters,
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
			"abi":    "OracleRequest(bytes32 indexed specId, address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes data)",
			"data":   "$(jobRun.logData)",
			"topics": "$(jobRun.logTopics)",
		}
		n := pipeline.NewGraphNode(dg.NewNode(), "decode_log", attrs)
		dg.AddNode(n)
		last = n

		attrs2 := map[string]string{
			"type": "cborparse",
			"data": "$(decode_log.data)",
			"mode": "diet",
		}
		n2 := pipeline.NewGraphNode(dg.NewNode(), "decode_cbor", attrs2)
		dg.AddNode(n2)
		last = n2
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
			//TODO:----------------------------------
			//TODO: add merge

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
			//TODO:----------------------------------
			//TODO: add merge

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

				template := fmt.Sprintf("%%REQ_DATA_%v%%", i)
				i++
				attrs := map[string]string{
					"type": "ethabiencode",
					"abi":  "(uint256 value)",
					"data": template,
				}
				replacements["\""+template+"\""] = fmt.Sprintf(`{ "value": $(%v) }`, last.DOTID())

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
				replacements["\""+template+"\""] = fmt.Sprintf(`{
"requestId":          $(decode_log.requestId),
"payment":            $(decode_log.payment),
"callbackAddress":    $(decode_log.callbackAddr),
"callbackFunctionId": $(decode_log.callbackFunctionId),
"expiration":         $(decode_log.cancelExpiration),
"data":               $(%v)
}
`, last.DOTID())

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
