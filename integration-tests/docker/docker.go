package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
)

func CreateNetwork(name string) (string, error) {
	if name == "" {
		name = fmt.Sprintf("network-%s", uuid.NewString())
	}
	exists, err := CheckNetworkExists(name)
	if err != nil {
		return name, err
	}
	if exists {
		return name, nil
	}
	network, err := tc.GenericNetwork(context.Background(), tc.GenericNetworkRequest{
		NetworkRequest: tc.NetworkRequest{
			Name: name,
		},
	})
	if err != nil {
		return name, err
	}
	dockerNetwork, ok := network.(*tc.DockerNetwork)
	if !ok {
		return "", fmt.Errorf("failed to cast network to *dockertest.Network")
	}
	log.Info().Str("Network", dockerNetwork.Name).Msg("Created docker network")
	return name, nil
}

func CheckNetworkExists(name string) (bool, error) {
	c, err := tc.NewDockerClient()
	if err != nil {
		return false, err
	}
	nr, err := c.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return false, err
	}
	for _, netr := range nr {
		if netr.Name == name {
			return true, nil
		}
	}
	return false, nil
}
