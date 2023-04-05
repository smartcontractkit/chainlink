package src

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
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
	nodesFile := fs.String("nodes", "", "a file containing nodes urls, logins and passwords")
	err := fs.Parse(args)
	if err != nil || nodesFile == nil || *nodesFile == "" {
		fs.Usage()
		os.Exit(1)
	}

	nodes := mustReadNodesList(*nodesFile)
	for _, node := range nodes {
		output := &bytes.Buffer{}
		client, app := newApp(node, output)

		fmt.Println("Logging in:", node.url)
		loginFs := flag.NewFlagSet("test", flag.ContinueOnError)
		loginFs.Bool("bypass-version-check", true, "")
		loginCtx := cli.NewContext(app, loginFs, nil)
		err := client.RemoteLogin(loginCtx)
		helpers.PanicErr(err)
		output.Reset()

		tomlPath := filepath.Join(artefactsDir, node.url.Host+".toml")
		tomlPath, err = filepath.Abs(tomlPath)
		if err != nil {
			helpers.PanicErr(err)
		}
		fmt.Println("Deploying jobspec:", tomlPath)
		if _, err := os.Stat(tomlPath); err != nil {
			helpers.PanicErr(errors.New("toml file does not exist"))
		}

		fileFs := flag.NewFlagSet("test", flag.ExitOnError)
		fileFs.Parse([]string{tomlPath})
		err = client.CreateJob(cli.NewContext(app, fileFs, nil))
		helpers.PanicErr(err)
		output.Reset()
	}
}
