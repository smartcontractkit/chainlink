package plugins

import (
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	relay "github.com/smartcontractkit/chainlink-relay/pkg/plugin"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const SolanaName = "Solana"

type solana struct {
	utils.StartStopOnce

	close func()
	relay.Solana
}

func (s *solana) Start() error {
	return s.StartOnce("Solana", func() error {
		cmdPath := os.Getenv("PLUGIN_SOLANA")
		if cmdPath == "" {
			//TODO log or silent? debug?
			return nil
		}
		return s.launch(cmdPath)
	})
}

func (s *solana) Close() error {
	return s.StopOnce("Solana", func() error {
		if s.close != nil {
			s.close()
		}
		return nil
	})
	return nil
}

func (s *solana) launch(cmdPath string) (err error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: relay.SolanaHandshake,
		Plugins: map[string]plugin.Plugin{
			SolanaName: &relay.SolanaPlugin{},
		},
		Cmd: exec.Command(cmdPath),
	})
	var rpcClient plugin.ClientProtocol
	rpcClient, err = client.Client()
	if err != nil {
		client.Kill()
		return
	}
	var raw interface{}
	raw, err = rpcClient.Dispense(SolanaName)
	if err != nil {
		client.Kill()
		return
	}
	s.close = client.Kill
	s.Solana = raw.(relay.Solana)
	return
}
