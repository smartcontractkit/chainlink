package gogauntlet

import (
	"fmt"
	"io"
)

type DefaultConfig struct {
	network string
	pk      string
	nodeurl string
	out     io.Writer
}

func (cd *DefaultConfig) Out() io.Writer {
	return cd.out
}

func (cd *DefaultConfig) SetOut(out io.Writer) {
	cd.out = out
}

func (cd *DefaultConfig) NodeURL() string {
	return cd.nodeurl
}

func (cd *DefaultConfig) SetNodeURL(nodeurl string) {
	cd.nodeurl = nodeurl
}

// SetNetwork Selects the 'network' on which to execute the specified command
func (cd *DefaultConfig) SetNetwork(network string) {
	cd.network = network
}

func (cd *DefaultConfig) Network() string {
	return cd.network
}

// SetPrivatekey Selects the privatekey of the wallet to use for the specified command
func (cd *DefaultConfig) SetPrivatekey(pk string) {
	cd.pk = pk
}

func (cd *DefaultConfig) Privatekey() string {
	return cd.pk
}

func NewDefaultConfig(nodeUrl, privatekey string, out io.Writer) *DefaultConfig {
	return &DefaultConfig{
		nodeurl: nodeUrl,
		network: "",
		pk:      privatekey,
		out:     out,
	}
}

// CreateDefaultFlags Generates the network, privatekey, and export flag if they are specified
func (cd *DefaultConfig) CreateDefaultFlags() []string {
	var output []string
	// We prioritize a nodeurl over a network default. Network requires user to create and set .env variables
	if cd.nodeurl != "" {
		output = append(output, fmt.Sprintf("--nodeurl=%s", cd.nodeurl))
	} else if cd.network != "" {
		output = append(output, fmt.Sprintf("--network=%s", cd.network))
	}
	if cd.pk != "" {
		output = append(output, fmt.Sprintf("--privatekey=%s", cd.pk))
	}
	return output
}
