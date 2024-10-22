## Developing Components

To build a scalable framework that enables the reuse of our product deployments (contracts or services in Docker), we need to establish a clear component structure.
```
package mycomponent

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

type Input struct {
    // inputs fields that component exposes for configuration
    ...
    // outputs are embedded into inputs so framework can automatically save them
	Out                      *Output  `toml:"out"`
}

type Output struct {
    // outputs that will be dumped to config and cached
}


func NewComponent(input *Input) (*Output, error) {
	if input.Out != nil && framework.UseCache() {
		return input.Out, nil
	}
	
	// component logic here
	// deploy a docker container(s)
	// or deploy a set of smart contracts
	
	input.Out = &Output{...}
	return out, nil
}
```
Each component can define inputs and outputs, following these rules:

- Inputs and outputs must avoid complex structures (e.g., loggers, clients) to ensure they remain cacheable.
- Outputs should be included within inputs.
- Every input field should have `validate: "required"` to maintain consistent configuration and ensure default values are always present. Test won't start if there are fields that are missing in configuration but declared in your component structures.

Specific docker components good practices for [testcontainers-go](https://golang.testcontainers.org/):
- `ContainerRequest` must contain labels, network and alias required for local observability stack and deployment isolation
```
		Labels:   framework.DefaultTCLabels(),
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
```
- In order to copy files into container use `framework.WriteTmpFile(data string, fileName string)`
```
	userSecretsOverridesFile, err := WriteTmpFile(in.Node.UserSecretsOverrides, "user-secrets-overrides.toml")
	if err != nil {
		return nil, err
	}
```
- Output of docker component must contain all the URLs component exposes for access, both for internal docker usage and external test (host) usage
```
	host, err := framework.GetHost(c)
	if err != nil {
		return nil, err
	}
	mp, err := c.MappedPort(ctx, nat.Port(bindPort))
	if err != nil {
		return nil, err
	}

	return &NodeOut{
		DockerURL: fmt.Sprintf("http://%s:%s", containerName, in.Node.Port),
		HostURL:   fmt.Sprintf("http://%s:%s", host, mp.Port()),
	}, nil
```

## Combining components in tests
We explicitly combine components in tests. To add your component, include its inputs in the test's Config struct and connect other outputs or inputs to it within the test.
```
type Config struct {
	ComponentA        *blockchain.Input `toml:"component_a" validate:"required"`
	ComponentB        *blockchain.Input `toml:"component_b" validate:"required"`
}

func TestDON(t *testing.T) {
    // load configuration
	in, err := framework.Load[Config](t)
	require.NoError(t, err)

	// deploy docker components
	bc, err := myComponent.NewComponent(in.ComponentA)
	require.NoError(t, err)
}
```
