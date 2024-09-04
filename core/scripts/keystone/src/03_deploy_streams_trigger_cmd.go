package src

// This package deploys "offchainreporting2" job specs, which setup the streams trigger
// for the targetted node set
// See https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/services/ocr2/plugins/mercury/integration_test.go#L92
// for how to setup the mercury portion of the streams trigger
//  You can see how all fields are being used here: https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/services/ocr2/plugins/mercury/helpers_test.go#L314
//  https://github.com/smartcontractkit/infra-k8s/blob/be47098adfb605d79b5bab6aa601bcf443a6c48b/projects/chainlink/files/chainlink-clusters/cl-keystone-cap-one/config.yaml#L1
//  Trigger gets added to the registry here: https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/services/relay/evm/evm.go#L360
//  See integration workflow here: https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/capabilities/integration_tests/workflow.go#L15
//  ^ setup.go provides good insight too
import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"net/url"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/web/presenters"
)

type deployStreamsTrigger struct{}

func NewDeployStreamsTriggerCommand() *deployStreamsTrigger {
	return &deployStreamsTrigger{}
}

func (g *deployStreamsTrigger) Name() string {
	return "deploy-streams-trigger"
}

func (g *deployStreamsTrigger) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	chainID := fs.Int64("chainid", 11155111, "chain id")
	templatesLocation := fs.String("templates", "", "Custom templates location")
	feedID := fs.String("feedid", "", "Feed ID")
	linkFeedID := fs.String("linkfeedid", "", "Link Feed ID")
	nativeFeedID := fs.String("nativefeedid", "", "Native Feed ID")
	fromBlock := fs.Int64("fromblock", 0, "From block")
	nodeList := fs.String("nodes", "", "Custom node list location")
	publicKeys := fs.String("publickeys", "", "Custom public keys json location")
	verifierContractAddress := fs.String("verifiercontractaddress", "", "Verifier contract address")
	verifierProxyContractAddress := fs.String("verifierproxycontractaddress", "", "Verifier proxy contract address")
	dryrun := fs.Bool("dryrun", false, "Dry run")

	err := fs.Parse(args)
	if err != nil || chainID == nil || *chainID == 0 ||
		feedID == nil || *feedID == "" ||
		linkFeedID == nil || *linkFeedID == "" ||
		nativeFeedID == nil || *nativeFeedID == "" ||
		fromBlock == nil || *fromBlock == 0 ||
		verifierContractAddress == nil || *verifierContractAddress == "" ||
		verifierProxyContractAddress == nil || *verifierProxyContractAddress == "" {
		fs.Usage()
		os.Exit(1)
	}

	if *publicKeys == "" {
		*publicKeys = defaultPublicKeys
	}
	if *nodeList == "" {
		*nodeList = defaultNodeList
	}
	if *templatesLocation == "" {
		*templatesLocation = "templates"
	}

	nodes := downloadNodeAPICredentials(*nodeList)

	jobspecs := genStreamsTriggerJobSpecs(
		*publicKeys,
		*nodeList,
		*templatesLocation,

		*feedID,
		*linkFeedID,
		*nativeFeedID,

		*chainID,
		*fromBlock,

		*verifierContractAddress,
		*verifierProxyContractAddress,
	)

	// sanity check arr lengths
	if len(nodes) != len(jobspecs) {
		PanicErr(errors.New("Mismatched node and job spec lengths"))
	}

	for i, n := range nodes {
		api := newNodeAPI(n)

		specToDeploy := strings.Join(jobspecs[i], "\n")
		specFragment := jobspecs[i][0:2]
		if *dryrun {
			fmt.Println("Dry run, skipping job deployment and bridge setup")
			fmt.Printf("Deploying jobspec: %s\n... \n", specToDeploy)
			continue
		} else {
			fmt.Printf("Deploying jobspec: %s\n... \n", specFragment)
		}

		_, err := api.withArg(specToDeploy).exec(api.methods.CreateJob)
		if err != nil {
			fmt.Println("Failed to deploy job spec:", specFragment, "Error:", err)
		}

		// hard coded bridges for now
		createBridgeIfDoesNotExist(api, "bridge-coinmetrics", "http://localhost:4001")
		createBridgeIfDoesNotExist(api, "bridge-tiingo", "http://localhost:4002")
		createBridgeIfDoesNotExist(api, "bridge-ncfx", "http://localhost:4003")

	}
}

func createBridgeIfDoesNotExist(api *nodeAPI, name string, eaURL string) {
	if doesBridgeExist(api, name) {
		fmt.Println("Bridge", name, "already exists, skipping creation")
		return
	}

	u, err := url.Parse(eaURL)
	url := models.WebURL(*u)
	// Confirmations and MinimumContractPayment are not used, so we can leave them as 0
	b := bridges.BridgeTypeRequest{
		Name: bridges.MustParseBridgeName(name),
		URL:  url,
	}
	payload, err := json.Marshal(b)
	helpers.PanicErr(err)

	resp := api.withArg(string(payload)).mustExec(api.methods.CreateBridge)
	resource := mustJSON[presenters.BridgeResource](resp)
	fmt.Printf("Created bridge: %s %s\n", resource.Name, resource.URL)
}

func doesBridgeExist(api *nodeAPI, name string) bool {
	resp, err := api.withArg(name).exec(api.methods.ShowBridge)

	if err != nil {
		return false
	}

	b := mustJSON[presenters.BridgeResource](resp)
	fmt.Printf("Found bridge: %s with URL: %s\n", b.Name, b.URL)
	return true
}

func genStreamsTriggerJobSpecs(
	pubkeysPath string,
	nodeListPath string,
	templatesDir string,

	feedID string,
	linkFeedID string,
	nativeFeedID string,

	chainID int64,
	fromBlock int64,

	verifierContractAddress string,
	verifierProxyContractAddress string,
) (output [][]string) {
	nodes := downloadNodeAPICredentials(nodeListPath)
	nca := downloadNodePubKeys(nodeListPath, chainID, pubkeysPath)
	lines, err := readLines(filepath.Join(templatesDir, streamsTriggerSpecTemplate))
	if err != nil {
		PanicErr(err)
	}

	for i := 0; i < len(nodes); i++ {
		n := nca[i]
		specLines := renderStreamsTriggerJobSpec(
			lines,

			feedID,
			linkFeedID,
			nativeFeedID,

			chainID,
			fromBlock,

			verifierContractAddress,
			verifierProxyContractAddress,

			n,
		)
		output = append(output, specLines)
	}

	return output
}

func renderStreamsTriggerJobSpec(
	lines []string,

	feedID string,
	linkFeedID string,
	nativeFeedID string,

	chainID int64,
	fromBlock int64,

	verifierContractAddress string,
	verifierProxyContractAddress string,

	node NodeKeys,
) (output []string) {
	chainIDStr := strconv.FormatInt(chainID, 10)
	fromBlockStr := strconv.FormatInt(fromBlock, 10)
	for _, l := range lines {
		l = strings.Replace(l, "{{ feed_id }}", feedID, 1)
		l = strings.Replace(l, "{{ link_feed_id }}", linkFeedID, 1)
		l = strings.Replace(l, "{{ native_feed_id }}", nativeFeedID, 1)

		l = strings.Replace(l, "{{ chain_id }}", chainIDStr, 1)
		l = strings.Replace(l, "{{ from_block }}", fromBlockStr, 1)

		// Verifier contract https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/services/ocr2/plugins/mercury/integration_test.go#L111
		l = strings.Replace(l, "{{ contract_id }}", verifierContractAddress, 1)
		// Ends up just being part of the name as documentation, it's the proxy to the verifier contract
		l = strings.Replace(l, "{{ verifier_proxy_id }}", verifierProxyContractAddress, 1)

		//  TransmitterID is the CSA key of the node since it's offchain
		//  https://github.com/smartcontractkit/chainlink/blob/4d5fc1943bd6a60b49cbc3d263c0aa47dc3cecb7/core/services/ocr2/plugins/mercury/helpers_test.go#L219
		// https://github.com/smartcontractkit/chainlink-common/blob/9ee1e8cc8b9774c8f3eb92a722af5269469f46f4/pkg/types/mercury/types.go#L39
		l = strings.Replace(l, "{{ transmitter_id }}", node.CSAPublicKey, 1)
		output = append(output, l)
	}
	return
}
