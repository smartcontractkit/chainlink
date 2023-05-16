package loop_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"
)

func TestPluginRelayer(t *testing.T) {
	t.Parallel()

	stopCh := make(chan struct{})
	if d, ok := t.Deadline(); ok {
		time.AfterFunc(time.Until(d), func() { close(stopCh) })
	}
	testPlugin(t, loop.PluginRelayerName, &loop.GRPCPluginRelayer{Logger: logger.Test(t), PluginServer: staticPluginRelayer{}, StopCh: stopCh}, testPluginRelayer)
}

func TestPluginRelayerExec(t *testing.T) {
	t.Parallel()
	relayer := loop.GRPCPluginRelayer{Logger: logger.Test(t)}
	cc := relayer.ClientConfig()
	cc.Cmd = helperProcess(loop.PluginRelayerName)
	c := plugin.NewClient(cc)
	client, err := c.Client()
	require.NoError(t, err)
	defer client.Close()
	require.NoError(t, client.Ping())
	i, err := client.Dispense(loop.PluginRelayerName)
	require.NoError(t, err)

	testPluginRelayer(t, i.(loop.PluginRelayer))
}

func testPluginRelayer(t *testing.T, p loop.PluginRelayer) {
	ctx := utils.Context(t)

	t.Run("Relayer", func(t *testing.T) {
		relayer, err := p.NewRelayer(ctx, configTOML, staticKeystore{})
		require.NoError(t, err)
		require.NoError(t, relayer.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, relayer.Close()) })
		testRelayer(t, relayer)
	})
}

func testPlugin[I any](t *testing.T, name string, p plugin.Plugin, testFn func(*testing.T, I)) {
	ctx, cancel := context.WithCancel(utils.Context(t))
	defer cancel()

	ch := make(chan *plugin.ReattachConfig, 1)
	closeCh := make(chan struct{})
	go plugin.Serve(&plugin.ServeConfig{
		Test: &plugin.ServeTestConfig{
			Context:          ctx,
			ReattachConfigCh: ch,
			CloseCh:          closeCh,
		},
		GRPCServer: plugin.DefaultGRPCServer,
		Plugins:    map[string]plugin.Plugin{name: p},
	})

	// We should get a config
	var config *plugin.ReattachConfig
	select {
	case config = <-ch:
	case <-time.After(5 * time.Second):
		t.Fatal("should've received reattach")
	}
	require.NotNil(t, config)

	c := plugin.NewClient(&plugin.ClientConfig{
		Reattach: config,
		Plugins:  map[string]plugin.Plugin{name: p},
	})
	clientProtocol, err := c.Client()
	require.NoError(t, err)
	defer clientProtocol.Close()
	i, err := clientProtocol.Dispense(name)
	require.NoError(t, err)

	testFn(t, i.(I))

	// stop plugin
	cancel()
	select {
	case <-closeCh:
	case <-time.After(5 * time.Second):
		t.Fatal("should've stopped")
	}
	require.Error(t, clientProtocol.Ping())
}

func testRelayer(t *testing.T, relayer loop.Relayer) {
	ctx := utils.Context(t)

	t.Run("ConfigProvider", func(t *testing.T) {
		t.Parallel()
		configProvider, err := relayer.NewConfigProvider(ctx, rargs)
		require.NoError(t, err)
		require.NoError(t, configProvider.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, configProvider.Close()) })

		t.Run("OffchainConfigDigester", func(t *testing.T) {
			t.Parallel()
			ocd := configProvider.OffchainConfigDigester()
			gotConfigDigestPrefix := ocd.ConfigDigestPrefix()
			assert.Equal(t, configDigestPrefix, gotConfigDigestPrefix)
			gotConfigDigest, err := ocd.ConfigDigest(contractConfig)
			require.NoError(t, err)
			assert.Equal(t, configDigest, gotConfigDigest)
		})
		t.Run("ContractConfigTracker", func(t *testing.T) {
			t.Parallel()
			cct := configProvider.ContractConfigTracker()
			gotBlockHeight, err := cct.LatestBlockHeight(ctx)
			require.NoError(t, err)
			assert.Equal(t, blockHeight, gotBlockHeight)
			gotChangedInBlock, gotConfigDigest, err := cct.LatestConfigDetails(ctx)
			require.NoError(t, err)
			assert.Equal(t, changedInBlock, gotChangedInBlock)
			assert.Equal(t, configDigest, gotConfigDigest)
			gotContractConfig, err := cct.LatestConfig(ctx, changedInBlock)
			require.NoError(t, err)
			assert.Equal(t, contractConfig, gotContractConfig)
		})
	})

	t.Run("MedianProvider", func(t *testing.T) {
		t.Parallel()
		provider, err := relayer.NewMedianProvider(ctx, rargs, pargs)
		require.NoError(t, err)
		require.NoError(t, provider.Start(ctx))
		t.Cleanup(func() { assert.NoError(t, provider.Close()) })

		t.Run("ReportingPluginProvider", func(t *testing.T) {
			t.Parallel()

			t.Run("OffchainConfigDigester", func(t *testing.T) {
				t.Parallel()
				ocd := provider.OffchainConfigDigester()
				gotConfigDigestPrefix := ocd.ConfigDigestPrefix()
				assert.Equal(t, configDigestPrefix, gotConfigDigestPrefix)
				gotConfigDigest, err := ocd.ConfigDigest(contractConfig)
				require.NoError(t, err)
				assert.Equal(t, configDigest, gotConfigDigest)
			})
			t.Run("ContractConfigTracker", func(t *testing.T) {
				t.Parallel()
				cct := provider.ContractConfigTracker()
				gotBlockHeight, err := cct.LatestBlockHeight(ctx)
				require.NoError(t, err)
				assert.Equal(t, blockHeight, gotBlockHeight)
				gotChangedInBlock, gotConfigDigest, err := cct.LatestConfigDetails(ctx)
				require.NoError(t, err)
				assert.Equal(t, changedInBlock, gotChangedInBlock)
				assert.Equal(t, configDigest, gotConfigDigest)
				gotContractConfig, err := cct.LatestConfig(ctx, changedInBlock)
				require.NoError(t, err)
				assert.Equal(t, contractConfig, gotContractConfig)
			})
			t.Run("ContractTransmitter", func(t *testing.T) {
				t.Parallel()
				ct := provider.ContractTransmitter()
				gotAccount := ct.FromAccount()
				assert.Equal(t, account, gotAccount)
				gotConfigDigest, gotEpoch, err := ct.LatestConfigDigestAndEpoch(ctx)
				require.NoError(t, err)
				assert.Equal(t, configDigest, gotConfigDigest)
				assert.Equal(t, epoch, gotEpoch)
				err = ct.Transmit(ctx, reportContext, report, sigs)
				require.NoError(t, err)
			})
			t.Run("ReportCodec", func(t *testing.T) {
				t.Parallel()
				rc := provider.ReportCodec()
				gotReport, err := rc.BuildReport(pobs)
				require.NoError(t, err)
				assert.Equal(t, report, gotReport)
				gotMedianValue, err := rc.MedianFromReport(report)
				require.NoError(t, err)
				assert.Equal(t, medianValue, gotMedianValue)
				gotMax := rc.MaxReportLength(n)
				assert.Equal(t, max, gotMax)
			})
			t.Run("MedianContract", func(t *testing.T) {
				t.Parallel()
				mc := provider.MedianContract()
				gotConfigDigest, gotEpoch, gotRound, err := mc.LatestRoundRequested(ctx, lookbackDuration)
				require.NoError(t, err)
				assert.Equal(t, configDigest, gotConfigDigest)
				assert.Equal(t, epoch, gotEpoch)
				assert.Equal(t, round, gotRound)
				gotConfigDigest, gotEpoch, gotRound, gotLatestAnswer, gotLatestTimestamp, err := mc.LatestTransmissionDetails(ctx)
				require.NoError(t, err)
				assert.Equal(t, configDigest, gotConfigDigest)
				assert.Equal(t, epoch, gotEpoch)
				assert.Equal(t, round, gotRound)
				assert.Equal(t, latestAnswer, gotLatestAnswer)
				assert.WithinDuration(t, latestTimestamp, gotLatestTimestamp, time.Second)
			})
			t.Run("OnchainConfigCodec", func(t *testing.T) {
				t.Parallel()
				occ := provider.OnchainConfigCodec()
				gotEncoded, err := occ.Encode(onchainConfig)
				require.NoError(t, err)
				assert.Equal(t, encoded, gotEncoded)
				gotDecoded, err := occ.Decode(encoded)
				require.NoError(t, err)
				assert.Equal(t, onchainConfig, gotDecoded)
			})
		})
	})

	t.Run("ChainStatus", func(t *testing.T) {
		t.Parallel()
		gotChain, err := relayer.ChainStatus(ctx, chainID)
		require.NoError(t, err)
		assert.Equal(t, chain, gotChain)
	})

	t.Run("ChainStatuses", func(t *testing.T) {
		t.Parallel()
		gotChains, gotCount, err := relayer.ChainStatuses(ctx, offset, limit)
		require.NoError(t, err)
		assert.Equal(t, chains, gotChains)
		assert.Equal(t, count, gotCount)
	})

	t.Run("NodeStatuses", func(t *testing.T) {
		t.Parallel()
		gotNodes, gotCount, err := relayer.NodeStatuses(ctx, offset, limit, chainIDs...)
		require.NoError(t, err)
		assert.Equal(t, nodes, gotNodes)
		assert.Equal(t, count, gotCount)
	})

	t.Run("SendTx", func(t *testing.T) {
		t.Parallel()
		err := relayer.SendTx(ctx, chainID, from, to, amount, balanceCheck)
		require.NoError(t, err)
	})
}

type staticKeystore struct{}

func (s staticKeystore) Accounts(ctx context.Context) (accounts []string, err error) {
	return []string{string(account)}, nil
}

func (s staticKeystore) Sign(ctx context.Context, id string, data []byte) ([]byte, error) {
	if string(account) != id {
		return nil, fmt.Errorf("expected id %q but got %q", account, id)
	}
	if !bytes.Equal(encoded, data) {
		return nil, fmt.Errorf("expected encoded data %x but got %x", encoded, data)
	}
	return signed, nil
}

type staticPluginRelayer struct{}

func (s staticPluginRelayer) NewRelayer(ctx context.Context, config string, keystore loop.Keystore) (loop.Relayer, error) {
	if config != configTOML {
		return nil, fmt.Errorf("expected config %q but got %q", configTOML, config)
	}
	keys, err := keystore.Accounts(ctx)
	if err != nil {
		return nil, err
	}
	if !reflect.DeepEqual([]string{string(account)}, keys) {
		return nil, fmt.Errorf("expected keys %v but got %v", []string{string(account)}, keys)
	}
	gotSigned, err := keystore.Sign(ctx, string(account), encoded)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(signed, gotSigned) {
		return nil, fmt.Errorf("expected signed bytes %x but got %x", signed, gotSigned)
	}
	return staticRelayer{}, nil
}

func equalRelayArgs(a, b types.RelayArgs) bool {
	return a.ExternalJobID == b.ExternalJobID &&
		a.JobID == b.JobID &&
		a.ContractID == b.ContractID &&
		a.New == b.New &&
		bytes.Equal(a.RelayConfig, b.RelayConfig)
}

func helperProcess(s ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, s...)
	env := []string{
		"GO_WANT_HELPER_PROCESS=1",
	}

	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(env, os.Environ()...)
	return cmd
}

// This is not a real test. This is just a helper process kicked off by
// tests.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}

		args = args[1:]
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]

	limit := -1
	if len(args) > 0 {
		var err error
		limit, err = strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse integer limit: %s\n", err)
			os.Exit(2)
		}
	}

	grpcServer := func(opts []grpc.ServerOption) *grpc.Server { return grpc.NewServer(opts...) }
	if limit > -1 {
		unary, stream := limitInterceptors(limit)
		grpcServer = func(opts []grpc.ServerOption) *grpc.Server {
			opts = append(opts, grpc.UnaryInterceptor(unary), grpc.StreamInterceptor(stream))
			return grpc.NewServer(opts...)
		}
	}

	stopCh := make(chan struct{})
	defer close(stopCh)
	switch cmd {
	case loop.PluginRelayerName:
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: loop.PluginRelayerHandshakeConfig(),
			Plugins: map[string]plugin.Plugin{
				loop.PluginRelayerName: &loop.GRPCPluginRelayer{Logger: logger.Test(t), PluginServer: staticPluginRelayer{}, StopCh: stopCh},
			},
			GRPCServer: grpcServer,
		})
		os.Exit(0)

	case loop.PluginMedianName:
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: loop.PluginMedianHandshakeConfig(),
			Plugins: map[string]plugin.Plugin{
				loop.PluginRelayerName: &loop.GRPCPluginMedian{Logger: logger.Test(t), PluginServer: staticPluginMedian{}, StopCh: stopCh},
			},
			GRPCServer: grpcServer,
		})
		os.Exit(0)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %q\n", cmd)
		os.Exit(2)
	}
}

// limitInterceptors returns a pair of interceptors which increment a shared count for each call and exit the program
// when limit is reached.
func limitInterceptors(limit int) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor) {
	count := make(chan struct{})
	go func() {
		for i := 0; i < limit; i++ {
			<-count
		}
		os.Exit(3)
	}()
	return limitUnaryInterceptor(count), limitStreamInterceptor(count)
}

func limitUnaryInterceptor(count chan<- struct{}) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		count <- struct{}{}
		return handler(ctx, req)
	}
}

func limitStreamInterceptor(count chan<- struct{}) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		count <- struct{}{}
		return handler(srv, ss)
	}
}
