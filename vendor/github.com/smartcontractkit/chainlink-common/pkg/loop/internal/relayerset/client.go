package relayerset

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/relayerset"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

type connProvider interface {
	ClientConn() grpc.ClientConnInterface
}

type Client struct {
	*net.BrokerExt
	*goplugin.ServiceClient

	log logger.Logger

	relayerSetClient relayerset.RelayerSetClient
}

func NewRelayerSetClient(log logger.Logger, b *net.BrokerExt, conn grpc.ClientConnInterface) *Client {
	b = b.WithName("ChainRelayerClient")
	return &Client{log: log, BrokerExt: b, ServiceClient: goplugin.NewServiceClient(b, conn), relayerSetClient: relayerset.NewRelayerSetClient(conn)}
}

func (k *Client) Get(ctx context.Context, relayID types.RelayID) (core.Relayer, error) {
	_, err := k.relayerSetClient.Get(ctx, &relayerset.GetRelayerRequest{Id: &relayerset.RelayerId{ChainId: relayID.ChainID, Network: relayID.Network}})
	if err != nil {
		return nil, fmt.Errorf("error getting relayer: %w", err)
	}

	return newRelayerClient(k.log, k, relayID), nil
}

func (k *Client) List(ctx context.Context, relayIDs ...types.RelayID) (map[types.RelayID]core.Relayer, error) {
	var ids []*relayerset.RelayerId
	for _, id := range relayIDs {
		ids = append(ids, &relayerset.RelayerId{ChainId: id.ChainID, Network: id.Network})
	}

	resp, err := k.relayerSetClient.List(ctx, &relayerset.ListAllRelayersRequest{Ids: ids})
	if err != nil {
		return nil, fmt.Errorf("error getting all relayers: %w", err)
	}

	result := map[types.RelayID]core.Relayer{}
	for _, id := range resp.Ids {
		relayID := types.RelayID{ChainID: id.ChainId, Network: id.Network}
		result[relayID] = newRelayerClient(k.log, k, relayID)
	}

	return result, nil
}

func (k *Client) StartRelayer(ctx context.Context, relayID types.RelayID) error {
	_, err := k.relayerSetClient.StartRelayer(ctx, &relayerset.RelayerId{ChainId: relayID.ChainID, Network: relayID.Network})
	return err
}

func (k *Client) CloseRelayer(ctx context.Context, relayID types.RelayID) error {
	_, err := k.relayerSetClient.CloseRelayer(ctx, &relayerset.RelayerId{ChainId: relayID.ChainID, Network: relayID.Network})
	return err
}

func (k *Client) RelayerReady(ctx context.Context, relayID types.RelayID) error {
	_, err := k.relayerSetClient.RelayerReady(ctx, &relayerset.RelayerId{ChainId: relayID.ChainID, Network: relayID.Network})
	return err
}

func (k *Client) RelayerHealthReport(ctx context.Context, relayID types.RelayID) (map[string]error, error) {
	report, err := k.relayerSetClient.RelayerHealthReport(ctx, &relayerset.RelayerId{ChainId: relayID.ChainID, Network: relayID.Network})
	if err != nil {
		return nil, fmt.Errorf("error getting health report: %w", err)
	}

	result := map[string]error{}
	for k, v := range report.Report {
		result[k] = errors.New(v)
	}

	return result, nil
}

func (k *Client) RelayerName(ctx context.Context, relayID types.RelayID) (string, error) {
	resp, err := k.relayerSetClient.RelayerName(ctx, &relayerset.RelayerId{ChainId: relayID.ChainID, Network: relayID.Network})
	if err != nil {
		return "", fmt.Errorf("error getting name: %w", err)
	}

	return resp.Name, nil
}

func (k *Client) NewPluginProvider(ctx context.Context, relayID types.RelayID, relayArgs core.RelayArgs, pluginArgs core.PluginArgs) (uint32, error) {
	// TODO at a later phase these credentials should be set as part of the relay config and not as a separate field
	var mercuryCredentials *relayerset.MercuryCredentials
	if relayArgs.MercuryCredentials != nil {
		mercuryCredentials = &relayerset.MercuryCredentials{
			LegacyUrl: relayArgs.MercuryCredentials.LegacyURL,
			Url:       relayArgs.MercuryCredentials.URL,
			Username:  relayArgs.MercuryCredentials.Username,
			Password:  relayArgs.MercuryCredentials.Password,
		}
	}

	req := &relayerset.NewPluginProviderRequest{
		RelayerId:  &relayerset.RelayerId{ChainId: relayID.ChainID, Network: relayID.Network},
		RelayArgs:  &relayerset.RelayArgs{ContractID: relayArgs.ContractID, RelayConfig: relayArgs.RelayConfig, ProviderType: relayArgs.ProviderType, MercuryCredentials: mercuryCredentials},
		PluginArgs: &relayerset.PluginArgs{TransmitterID: pluginArgs.TransmitterID, PluginConfig: pluginArgs.PluginConfig},
	}

	resp, err := k.relayerSetClient.NewPluginProvider(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("error getting new plugin provider: %w", err)
	}
	return resp.PluginProviderId, nil
}
