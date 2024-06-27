package test_env

import (
	tc "github.com/testcontainers/testcontainers-go"
)

type EnvComponent struct {
	ContainerName string
	Container     tc.Container
	Networks      []string
}

type EnvComponentOption = func(c *EnvComponent)

func WithContainerName(name string) EnvComponentOption {
	return func(c *EnvComponent) {
		if name != "" {
			c.ContainerName = name
		}
	}
}
