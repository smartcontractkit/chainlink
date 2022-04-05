package utils

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// GetEnvVars fetches all of the specified environment variables for a given prefix
func GetEnvVars(ctx *pulumi.Context, p string) (out []string, err error) {
	vars, err := GetEnvList(ctx, p)
	if err != nil {
		return
	}

	out = GetVars(ctx, p, vars)
	return
}

// GetEnvList fetches the list of environment variables for a given prefix
func GetEnvList(ctx *pulumi.Context, p string) ([]string, error) {
	var list []string
	return list, config.GetObject(ctx, p+"-ENV_VARS", &list)
}

// GetVars fetches the environment variables as the expected string format in an array
// given the prefix and the envVar name
func GetVars(ctx *pulumi.Context, p string, envVars []string) (out []string) {
	for _, env := range envVars {
		out = append(out, fmt.Sprintf("%s=%s", env, config.Get(ctx, p+"-"+env)))
	}
	return
}
