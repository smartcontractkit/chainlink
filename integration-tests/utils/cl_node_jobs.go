package utils

import (
	"bytes"
	"fmt"
	"net/url"
	"text/template"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func BuildBootstrapSpec(verifierAddr common.Address, chainID int64, feedId [32]byte) *nodeclient.OCR2TaskJobSpec {
	hash := common.BytesToHash(feedId[:])
	return &nodeclient.OCR2TaskJobSpec{
		Name:    fmt.Sprintf("bootstrap-%s", uuid.NewString()),
		JobType: "bootstrap",
		OCR2OracleSpec: job.OCR2OracleSpec{
			ContractID: verifierAddr.String(),
			Relay:      "evm",
			FeedID:     &hash,
			RelayConfig: map[string]interface{}{
				"chainID": int(chainID),
			},
			ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
		},
	}
}

func BuildOCRSpec(
	verifierAddr common.Address, chainID int64, fromBlock uint64,
	feedId [32]byte, bridges []nodeclient.BridgeTypeAttributes,
	csaPubKey string, msRemoteUrl string, msPubKey string,
	nodeOCRKey string, p2pV2Bootstrapper string, allowedFaults int) *nodeclient.OCR2TaskJobSpec {

	tmpl, err := template.New("os").Parse(`
{{range $i, $b := .Bridges}}
{{$b.Name}}_payload      [type=bridge name="{{$b.Name}}" timeout="50ms" requestData="{}"];
{{$b.Name}}_median       [type=jsonparse path="data,result"];
{{$b.Name}}_bid          [type=jsonparse path="data,result"];
{{$b.Name}}_ask          [type=jsonparse path="data,result"];

{{$b.Name}}_median_multiply          [type=multiply times=10];
{{$b.Name}}_bid_multiply          [type=multiply times=10];
{{$b.Name}}_ask_multiply          [type=multiply times=10];
{{end}}


{{range $i, $b := .Bridges}}
{{$b.Name}}_payload        -> {{$b.Name}}_median        -> {{$b.Name}}_median_multiply        -> benchmark_price;
{{end}}

benchmark_price [type=median allowedFaults={{.AllowedFaults}} index=0];

{{range $i, $b := .Bridges}}
{{$b.Name}}_payload        -> {{$b.Name}}_bid         -> {{$b.Name}}_bid_multiply        -> bid_price;
{{end}}

bid_price [type=median allowedFaults={{.AllowedFaults}} index=1];

{{range $i, $b := .Bridges}}
{{$b.Name}}_payload        -> {{$b.Name}}_ask         -> {{$b.Name}}_ask_multiply        -> ask_price;
{{end}}

ask_price [type=median allowedFaults={{.AllowedFaults}} index=2];
	`)
	if err != nil {
		panic(err)
	}
	data := struct {
		Bridges       []nodeclient.BridgeTypeAttributes
		AllowedFaults int
	}{
		Bridges:       bridges,
		AllowedFaults: allowedFaults,
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	observationSource := buf.String()

	hash := common.BytesToHash(feedId[:])
	return &nodeclient.OCR2TaskJobSpec{
		Name:              fmt.Sprintf("ocr2-%s", uuid.NewString()),
		JobType:           "offchainreporting2",
		MaxTaskDuration:   "1s",
		ForwardingAllowed: false,
		OCR2OracleSpec: job.OCR2OracleSpec{
			PluginType: "mercury",
			PluginConfig: map[string]interface{}{
				"serverURL":    fmt.Sprintf("\"%s\"", msRemoteUrl),
				"serverPubKey": fmt.Sprintf("\"%s\"", msPubKey),
			},
			Relay: "evm",
			RelayConfig: map[string]interface{}{
				"chainID":   int(chainID),
				"fromBlock": fromBlock,
			},
			ContractConfigTrackerPollInterval: *models.NewInterval(time.Second * 15),
			ContractID:                        verifierAddr.String(),
			FeedID:                            &hash,
			OCRKeyBundleID:                    null.StringFrom(nodeOCRKey),
			TransmitterID:                     null.StringFrom(csaPubKey),
			P2PV2Bootstrappers:                pq.StringArray{p2pV2Bootstrapper},
		},
		ObservationSource: observationSource,
	}
}

func BuildBridges(eaUrls []*url.URL) []nodeclient.BridgeTypeAttributes {
	var bridges []nodeclient.BridgeTypeAttributes
	for _, url := range eaUrls {
		bridges = append(bridges, nodeclient.BridgeTypeAttributes{
			Name:        fmt.Sprintf("bridge_%s", uuid.NewString()[0:6]),
			URL:         url.String(),
			RequestData: "{}",
		})
	}
	return bridges
}
