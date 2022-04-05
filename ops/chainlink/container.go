package chainlink

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pulumi/pulumi-docker/sdk/v3/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/smartcontractkit/chainlink-relay/ops/utils"
	"github.com/smartcontractkit/integrations-framework/client"
)

// New spins up image for a chainlink node
func New(ctx *pulumi.Context, image *utils.Image, dbPort int, index int) (Node, error) {
	// treat index 0 as bootstrap
	indexStr := ""
	if index == 0 {
		indexStr = "bootstrap"
	} else {
		indexStr = strconv.Itoa(index - 1)
	}

	// TODO: Default ports?

	portHTTP := fmt.Sprintf("%d", config.RequireInt(ctx, "CL-PORT-START")+index)
	portP2P := fmt.Sprintf("%d", config.RequireInt(ctx, "CL-P2P_PORT-START")+index)

	node := Node{
		Name: "chainlink-" + indexStr,
		P2P: client.P2PData{
			RemoteIP:   "chainlink-" + indexStr,
			RemotePort: portP2P,
			// PeerID this is set later when GetKeys() is called
		},
		Config: client.ChainlinkConfig{
			URL:      "http://localhost:" + portHTTP,
			Email:    "admin@chain.link",
			Password: "twoChains",
			RemoteIP: "localhost",
		},
		Keys: NodeKeys{},
	}

	// get env vars from YAML file
	envs, err := utils.GetEnvVars(ctx, "CL")
	if err != nil {
		return Node{}, err
	}

	// add additional configs (collected or calculated from environment configs)
	envs = append(envs,
		fmt.Sprintf("DATABASE_URL=postgresql://postgres@postgres:%d/chainlink_%s?sslmode=disable", 5432, indexStr),
		fmt.Sprintf("CHAINLINK_PORT=%s", portHTTP),
		fmt.Sprintf("P2PV2_LISTEN_ADDRESSES=0.0.0.0:%s", portP2P),
		fmt.Sprintf("P2PV2_ANNOUNCE_ADDRESSES=0.0.0.0:%s", portP2P),
		fmt.Sprintf("CLIENT_NODE_URL=http://0.0.0.0:%s", portHTTP), // needs to point to the correct node for container CLI
	)

	// fetch additional env vars (specific to each chainlink node)
	envListR, err := utils.GetEnvList(ctx, "CL_X")
	if err != nil {
		return Node{}, err
	}
	envsR := utils.GetVars(ctx, "CL_"+strings.ToUpper(indexStr), envListR)
	envs = append(envs, envsR...)

	entrypoints := pulumi.ToStringArray([]string{"chainlink", "node", "start", "-d", "-p", "/run/secrets/node_password", "-a", "/run/secrets/apicredentials"})
	uploads := docker.ContainerUploadArray{docker.ContainerUploadArgs{File: pulumi.String("/run/secrets/node_password"), Content: pulumi.String("abcd1234ABCD!@#$")}, docker.ContainerUploadArgs{File: pulumi.String("/run/secrets/apicredentials"), Content: pulumi.String(node.CredentialsString())}}

	var imageName pulumi.StringInput
	if config.GetBool(ctx, "CL-BUILD_LOCALLY") {
		imageName = image.Local.BaseImageName
	} else {
		imageName = image.Img.Name
	}

	_, err = docker.NewContainer(ctx, node.Name, &docker.ContainerArgs{
		Image:       imageName,
		Logs:        pulumi.BoolPtr(true),
		NetworkMode: pulumi.String(utils.GetDefaultNetworkName(ctx)),
		Hostname:    pulumi.String("chainlink-" + indexStr),
		Ports: docker.ContainerPortArray{
			docker.ContainerPortArgs{
				Internal: pulumi.Int(config.RequireInt(ctx, "CL-PORT-START") + index),
				External: pulumi.Int(config.RequireInt(ctx, "CL-PORT-START") + index),
			},
			docker.ContainerPortArgs{
				Internal: pulumi.Int(config.RequireInt(ctx, "CL-P2P_PORT-START") + index),
				External: pulumi.Int(config.RequireInt(ctx, "CL-P2P_PORT-START") + index),
			},
		},
		Envs:          pulumi.StringArrayInput(pulumi.ToStringArray(envs)),
		Uploads:       uploads.ToContainerUploadArrayOutput(),
		Entrypoints:   entrypoints.ToStringArrayOutput(),
		Restart:       pulumi.String("on-failure"),
		MaxRetryCount: pulumi.Int(3),
		Hosts: docker.ContainerHostArray{
			docker.ContainerHostArgs{
				Host: pulumi.String("host.docker.internal"),
				Ip:   pulumi.String("host-gateway"),
			},
		},
		// Attach:        pulumi.BoolPtr(true),
	})
	return node, err
}
