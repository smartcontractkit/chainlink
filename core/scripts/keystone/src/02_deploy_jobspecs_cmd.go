package src

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
)

type deployJobSpecs struct {
}

func NewDeployJobSpecsCommand() *deployJobSpecs {
	return &deployJobSpecs{}
}

func (g *deployJobSpecs) Name() string {
	return "deploy-jobspecs"
}

func (g *deployJobSpecs) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	chainID := fs.Int64("chainid", 11155111, "chain id")
	p2pPort := fs.Int64("p2pport", 6690, "p2p port")
	onlyReplay := fs.Bool("onlyreplay", false, "only replay the block from the OCR3 contract setConfig transaction")
	err := fs.Parse(args)
	if err != nil || chainID == nil || *chainID == 0 || p2pPort == nil || *p2pPort == 0 || onlyReplay == nil {
		fs.Usage()
		os.Exit(1)
	}
	if *onlyReplay {
		fmt.Println("Only replaying OCR3 contract setConfig transaction")
	} else {
		fmt.Println("Deploying OCR3 job specs")
	}

	nodes := downloadNodeAPICredentialsDefault()
	deployedContracts, err := LoadDeployedContracts()
	PanicErr(err)

	jobspecs := genSpecs(
		".cache/PublicKeys.json", ".cache/NodeList.txt", "templates",
		*chainID, *p2pPort, deployedContracts.OCRContract.Hex(),
	)
	flattenedSpecs := []hostSpec{jobspecs.bootstrap}
	flattenedSpecs = append(flattenedSpecs, jobspecs.oracles...)

	// sanity check arr lengths
	if len(nodes) != len(flattenedSpecs) {
		PanicErr(errors.New("Mismatched node and job spec lengths"))
	}

	for i, n := range nodes {
		output := &bytes.Buffer{}
		client, app := newApp(n, output)
		fmt.Println("Logging in:", n.url)
		loginFs := flag.NewFlagSet("test", flag.ContinueOnError)
		loginFs.Bool("bypass-version-check", true, "")
		loginCtx := cli.NewContext(app, loginFs, nil)
		err := client.RemoteLogin(loginCtx)
		helpers.PanicErr(err)
		output.Reset()

		if !*onlyReplay {
			specToDeploy := flattenedSpecs[i].spec.ToString()
			specFragment := flattenedSpecs[i].spec[0:1]
			fmt.Printf("Deploying jobspec: %s\n... \n", specFragment)
			fs := flag.NewFlagSet("test", flag.ExitOnError)
			err = fs.Parse([]string{specToDeploy})

			helpers.PanicErr(err)
			err = client.CreateJob(cli.NewContext(app, fs, nil))
			if err != nil {
				fmt.Println("Failed to deploy job spec:", specFragment, "Error:", err)
			}
			output.Reset()
		}

		replayFs := flag.NewFlagSet("test", flag.ExitOnError)
		flagSetApplyFromAction(client.ReplayFromBlock, replayFs, "")
		err = replayFs.Set("block-number", fmt.Sprint(deployedContracts.SetConfigTxBlock))
		helpers.PanicErr(err)
		err = replayFs.Set("evm-chain-id", fmt.Sprint(*chainID))
		helpers.PanicErr(err)

		fmt.Printf("Replaying from block: %d\n", deployedContracts.SetConfigTxBlock)
		fmt.Printf("EVM Chain ID: %d\n\n", *chainID)
		replayCtx := cli.NewContext(app, replayFs, nil)
		err = client.ReplayFromBlock(replayCtx)
		helpers.PanicErr(err)
	}
}

// flagSetApplyFromAction applies the flags from action to the flagSet.
//
// `parentCommand` will filter the app commands and only applies the flags if the command/subcommand has a parent with that name, if left empty no filtering is done
//
// Taken from: https://github.com/smartcontractkit/chainlink/blob/develop/core/cmd/shell_test.go#L590
func flagSetApplyFromAction(action interface{}, flagSet *flag.FlagSet, parentCommand string) {
	cliApp := cmd.Shell{}
	app := cmd.NewApp(&cliApp)

	foundName := parentCommand == ""
	actionFuncName := getFuncName(action)

	for _, command := range app.Commands {
		flags := recursiveFindFlagsWithName(actionFuncName, command, parentCommand, foundName)

		for _, flag := range flags {
			flag.Apply(flagSet)
		}
	}
}

func recursiveFindFlagsWithName(actionFuncName string, command cli.Command, parent string, foundName bool) []cli.Flag {
	if command.Action != nil {
		if actionFuncName == getFuncName(command.Action) && foundName {
			return command.Flags
		}
	}

	for _, subcommand := range command.Subcommands {
		if !foundName {
			foundName = strings.EqualFold(subcommand.Name, parent)
		}

		found := recursiveFindFlagsWithName(actionFuncName, subcommand, parent, foundName)
		if found != nil {
			return found
		}
	}
	return nil
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
