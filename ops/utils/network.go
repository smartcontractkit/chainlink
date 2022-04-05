package utils

import (
	"github.com/pulumi/pulumi-docker/sdk/v3/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// Create network
func CreateNetwork(ctx *pulumi.Context, nwName string) (*docker.Network, error) {
	network, err := docker.GetNetwork(ctx, nwName, nil, nil, nil)
	if err != nil {
		network, err := docker.NewNetwork(ctx, nwName, &docker.NetworkArgs{Name: pulumi.String(nwName)}, nil)
		return network, err
	}
	return network, err
}

func GetDefaultNetworkName(ctx *pulumi.Context) string {
	name := config.Get(ctx, "NETWORK_NAME")
	if name != "" {
		return name
	}
	return "pulumi-local"
}
