package database

import (
	"github.com/pulumi/pulumi-docker/sdk/v3/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/smartcontractkit/chainlink-relay/ops/utils"
)

// New spins up a postgres db docker image
func New(ctx *pulumi.Context, image *docker.RemoteImage) (Database, error) {
	port, err := config.TryInt(ctx, "PG-PORT")
	if err != nil {
		port = 5432 // default port
	}

	wait, err := config.TryInt(ctx, "PG-HEALTH-TIMEOUT")
	if err != nil {
		wait = 30 // default timeout period
	}

	db := Database{
		User:    "postgres",
		Host:    "localhost",
		Sslmode: "disable",
		Port:    port,
		Timeout: wait,
	}

	_, err = docker.NewContainer(ctx, "postgres", &docker.ContainerArgs{
		Image:       image.Name,
		Envs:        pulumi.StringArrayInput(pulumi.ToStringArray([]string{"POSTGRES_HOST_AUTH_METHOD=trust"})),
		NetworkMode: pulumi.String(utils.GetDefaultNetworkName(ctx)),
		Hostname:    pulumi.String("postgres"),
		Ports: docker.ContainerPortArray{
			docker.ContainerPortArgs{
				Internal: pulumi.Int(5432),
				External: pulumi.Int(db.Port),
			},
		},
	},
		pulumi.IgnoreChanges([]string{"image"}), // ignore changes to image to preserve db
	)

	return db, err
}
