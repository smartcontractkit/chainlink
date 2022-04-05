package adapter

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-docker/sdk/v3/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/smartcontractkit/chainlink-relay/ops/utils"
	"github.com/smartcontractkit/integrations-framework/client"
)

// New spins up docker images for adapters
func New(ctx *pulumi.Context, img *utils.Image, i int) (client.BridgeTypeAttributes, error) {
	port, err := config.TryInt(ctx, "EA-PORT")
	if err != nil {
		port = 8080 // default port
	}

	// get env vars from YAML file
	envs, err := utils.GetEnvVars(ctx, "EA")
	if err != nil {
		return client.BridgeTypeAttributes{}, err
	}

	name := strings.Split(img.Name, "-")[0]

	_, err = docker.NewContainer(ctx, name+"-adapter", &docker.ContainerArgs{
		Image:       img.Img.Name,
		NetworkMode: pulumi.String(utils.GetDefaultNetworkName(ctx)),
		Envs:        pulumi.StringArrayInput(pulumi.ToStringArray(envs)),
		Hostname:    pulumi.String(name + "-adapter"),
		Ports: docker.ContainerPortArray{
			docker.ContainerPortArgs{
				Internal: pulumi.Int(port),
				External: pulumi.Int(port + i),
			},
		},
	})

	return client.BridgeTypeAttributes{Name: name, URL: fmt.Sprintf("http://%s-adapter:%d", name, port)}, err
}
