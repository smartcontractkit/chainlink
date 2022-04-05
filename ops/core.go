package ops

import (
	"errors"
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/smartcontractkit/chainlink-relay/ops/adapter"
	"github.com/smartcontractkit/chainlink-relay/ops/chainlink"
	"github.com/smartcontractkit/chainlink-relay/ops/database"
	"github.com/smartcontractkit/chainlink-relay/ops/utils"
	"github.com/smartcontractkit/integrations-framework/client"
	"gopkg.in/guregu/null.v4"
)

// Deployer interface for deploying contracts
type Deployer interface {
	Load() error                             // upload contracts (may not be necessary)
	DeployLINK() error                       // deploy LINK contract
	DeployOCR() error                        // deploy OCR contract
	TransferLINK() error                     // transfer LINK to OCR contract
	InitOCR(keys []chainlink.NodeKeys) error // initialize OCR contract with provided keys
	Fund(addresses []string) error           // fund the nodes
	OCR2Address() string                     // fetch deployed OCR contract address
	Addresses() map[int]string               // map of all deployed addresses (ocr2, validators, etc)
}

// ObservationSource creates the observation source for the CL node jobs
type ObservationSource func(priceAdapter string) string

// RelayConfig creates the stringified config for the job spec
type RelayConfig func(ctx *pulumi.Context, addresses map[int]string) (map[string]string, error)

func New(ctx *pulumi.Context, deployer Deployer, obsSource ObservationSource, juelsObsSource ObservationSource, relayConfigFunc RelayConfig) error {
	// check these parameters at the beginning to prevent getting to the end and erroring if they are not present
	chain := config.Require(ctx, "CL-RELAY_NAME")
	onlyContainers := config.GetBool(ctx, "ENV-ONLY_BOOT_CONTAINERS")
	mixedUpgrade := config.GetBool(ctx, "ENV-MIX_UPGRADE")
	buildLocal := config.GetBool(ctx, "CL-BUILD_LOCALLY")

	if mixedUpgrade && buildLocal {
		return errors.New("MIX_UPGRADE and BUILD_LOCALLY cannot be enabled together")
	}

	if mixedUpgrade {
		if !onlyContainers {
			fmt.Println("WARN: MIX_UPGRADE requires ONLY_BOOT_CONTAINERS=true changing config to true")
			onlyContainers = true
		}
	}

	img := map[string]*utils.Image{}

	// fetch postgres
	img["psql"] = &utils.Image{
		Name: "postgres-image",
		Tag:  "postgres:latest", // always use latest postgres
	}

	if !buildLocal {
		// fetch chainlink image
		img["chainlink"] = &utils.Image{
			Name: "chainlink-remote-image",
			Tag:  "public.ecr.aws/chainlink/chainlink:" + config.Require(ctx, "CL-NODE_VERSION"),
		}

		if mixedUpgrade {
			img["chainlink-upgrade"] = &utils.Image{
				Name: "chainlink-remote-image-new",
				Tag:  "public.ecr.aws/chainlink/chainlink:" + config.Require(ctx, "ENV-MIX_UPGRADE_IMAGE"),
			}
		}
	}

	// fetch list of EAs
	eas := []string{}
	if err := config.GetObject(ctx, "EA-NAMES", &eas); err != nil {
		return err
	}
	for _, n := range eas {
		img[n] = &utils.Image{
			Name: n + "-adapter-image",
			Tag:  fmt.Sprintf("public.ecr.aws/chainlink/adapters/%s-adapter:develop-latest", n),
		}
	}

	// pull remote images
	for i := range img {
		if err := img[i].Pull(ctx); err != nil {
			return err
		}
	}

	// build local chainlink node
	if buildLocal {
		img["chainlink"] = &utils.Image{
			Name: "chainlink-local-build",
			Tag:  "chainlink:local",
		}
		if err := img["chainlink"].Build(ctx, config.Require(ctx, "CL-BUILD_CONTEXT"), config.Require(ctx, "CL-BUILD_DOCKERFILE")); err != nil {
			return err
		}
	}

	// validate number of relays
	nodeNum := config.GetInt(ctx, "CL-COUNT")
	if nodeNum < 4 {
		return fmt.Errorf("Minimum number of chainlink nodes (4) not met (%d)", nodeNum)
	}

	// create network
	nwName := utils.GetDefaultNetworkName(ctx)
	_, err := utils.CreateNetwork(ctx, nwName)

	if err != nil {
		return err
	}

	// start pg + create DBs
	db, err := database.New(ctx, img["psql"].Img)
	if err != nil {
		return err
	}
	if !ctx.DryRun() {
		// wait for readiness check
		if err := db.Ready(); err != nil {
			return err
		}

		// create DB names
		dbNames := []string{"chainlink_bootstrap"}
		for i := 0; i < nodeNum; i++ {
			dbNames = append(dbNames, fmt.Sprintf("chainlink_%d", i))
		}

		// create DBs
		for _, n := range dbNames {
			if err := db.Create(n); err != nil {
				return err
			}
		}
	}

	// start EAs
	adapters := []client.BridgeTypeAttributes{}
	for i, ea := range eas {
		a, err := adapter.New(ctx, img[ea], i)
		if err != nil {
			return err
		}
		adapters = append(adapters, a)
	}

	// start chainlink nodes
	nodes := map[string]*chainlink.Node{}
	mixNodeArr := []bool{}
	err = config.GetObject(ctx, "ENV-MIX_UPGRADE_NODES", &mixNodeArr)
	if err != nil && mixedUpgrade {
		return err // only return error if mixedUpgrade is true
	}
	// only check array length if mixedUpgrade is true
	if mixedUpgrade && nodeNum+1 != len(mixNodeArr) {
		return fmt.Errorf("incorrect MIX_UPGRADE_NODES length (%d), expected %d", len(mixNodeArr), nodeNum+1)
	}

	for i := 0; i <= nodeNum; i++ {
		var cl chainlink.Node
		var err error

		// mixed upgrade containers
		if mixedUpgrade && mixNodeArr[i] {
			cl, err = chainlink.New(ctx, img["chainlink-upgrade"], db.Port, i)
			fmt.Printf("⚠️  Upgrading %s\n", cl.Name)
		} else {
			// start container
			cl, err = chainlink.New(ctx, img["chainlink"], db.Port, i)
		}
		if err != nil {
			return err
		}
		nodes[cl.Name] = &cl // store in map
	}

	if onlyContainers {
		fmt.Println("ONLY BOOTING CONTAINERS")
		return nil
	}

	if !ctx.DryRun() {
		for _, cl := range nodes {
			// wait for readiness check
			if err := cl.Ready(); err != nil {
				return err
			}

			// delete all jobs if any exist
			if err := cl.DeleteAllJobs(); err != nil {
				return err
			}

			// add adapters to CL node
			for _, a := range adapters {
				if err := cl.AddBridge(a.Name, a.URL); err != nil {
					return err
				}
			}
		}
	}

	if !ctx.DryRun() {
		// fetch keys from relays
		for k := range nodes {
			if err := nodes[k].GetKeys(chain); err != nil {
				return err
			}
		}

		// upload contracts
		if err = deployer.Load(); err != nil {
			return err
		}
		// deploy LINK
		if err = deployer.DeployLINK(); err != nil {
			return err
		}

		// deploy OCR2 contract (w/ dummy access controller addresses)
		if err = deployer.DeployOCR(); err != nil {
			return err
		}

		// transfer tokens to OCR2 contract
		if err = deployer.TransferLINK(); err != nil {
			return err
		}

		// set OCR2 config
		var keys []chainlink.NodeKeys
		for k := range nodes {
			// skip if bootstrap node
			if k == "chainlink-bootstrap" {
				continue
			}
			keys = append(keys, nodes[k].Keys)
		}
		if err = deployer.InitOCR(keys); err != nil {
			return err
		}

		// create relay config
		relayConfig, err := relayConfigFunc(ctx, deployer.Addresses())
		if err != nil {
			return err
		}

		// create job specs
		var addresses []string
		i := 0
		for k := range nodes {
			// add chain & node to CL node
			// TODO: refactor under chainlink folder, pass configs in as interfaces defined from each individual file
			// also check if nodes/chains exist to prevent recreating
			switch chain {
			case "terra":
				msg := utils.LogStatus(fmt.Sprintf("Adding terra chain to '%s'", k))
				chainAttrs := client.TerraChainAttributes{
					ChainID: relayConfig["chainID"],
					Config: client.TerraChainConfig{
						BlocksUntilTxTimeout:  null.IntFrom(1),
						ConfirmPollPeriod:     null.StringFrom("1m0s"),
						FallbackGasPriceULuna: null.StringFrom("9.999"),
						GasLimitMultiplier:    null.FloatFrom(1.55555),
						MaxMsgsPerBatch:       null.IntFrom(10),
					},
				}
				_, err = nodes[k].Call.CreateTerraChain(&chainAttrs)
				if msg.Check(err) != nil {
					return err
				}
				msg = utils.LogStatus(fmt.Sprintf("Adding terra node to '%s'", k))
				nodeAttrs := client.TerraNodeAttributes{
					Name:          "Terra Node Localhost",
					TerraChainID:  relayConfig["chainID"],
					TendermintURL: relayConfig["tendermintURL"],
					FCDURL:        relayConfig["fcdURL"],
				}
				_, err = nodes[k].Call.CreateTerraNode(&nodeAttrs)
				if msg.Check(err) != nil {
					return err
				}
			case "solana":
				msg := utils.LogStatus(fmt.Sprintf("Adding solana chain to '%s'", k))
				chainAttrs := client.SolanaChainAttributes{
					ChainID: relayConfig["chainID"],
				}
				_, err = nodes[k].Call.CreateSolanaChain(&chainAttrs)
				if msg.Check(err) != nil {
					return err
				}
				msg = utils.LogStatus(fmt.Sprintf("Adding solana node to '%s'", k))
				nodeAttrs := client.SolanaNodeAttributes{
					Name:          "Solana Node Localhost",
					SolanaChainID: relayConfig["chainID"],
					SolanaURL:     relayConfig["solanaURL"],
				}
				_, err = nodes[k].Call.CreateSolanaNode(&nodeAttrs)
				if msg.Check(err) != nil {
					return err
				}
			default:
				fmt.Printf("WARN: No chain config to add to '%s'\n", k)
			}

			// create specs + add to CL node
			ea := eas[i%len(eas)]
			msg := utils.LogStatus(fmt.Sprintf("Adding job spec to '%s' with '%s' EA", k, ea))

			jobType := "offchainreporting2"
			if k == "chainlink-bootstrap" {
				jobType = "bootstrap"
			}

			spec := &client.OCR2TaskJobSpec{
				Name:        "local testing job",
				JobType:     jobType,
				ContractID:  deployer.OCR2Address(),
				Relay:       chain,
				RelayConfig: relayConfig,
				PluginType:  "median",
				P2PPeerID:   nodes[k].Keys.P2PID,
				P2PBootstrapPeers: []client.P2PData{
					nodes["chainlink-bootstrap"].P2P,
				},
				OCRKeyBundleID:        nodes[k].Keys.OCR2KeyID,
				TransmitterID:         nodes[k].Keys.OCR2TransmitterID,
				ObservationSource:     obsSource(ea),
				JuelsPerFeeCoinSource: juelsObsSource(ea),
			}
			_, err = nodes[k].Call.CreateJob(spec)
			if msg.Check(err) != nil {
				return err
			}
			i++

			// retrieve transmitter address for funding
			addresses = append(addresses, nodes[k].Keys.OCR2Transmitter)
		}

		// fund nodes
		if err = deployer.Fund(addresses); err != nil {
			return err
		}
	}

	return nil
}
