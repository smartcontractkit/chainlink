package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	clipkg "github.com/urfave/cli"
)

// getPassword retrieves the password from the file specified on the CL, or errors
func getPassword(c *clipkg.Context) ([]byte, error) {
	if c.String("password") == "" {
		return nil, fmt.Errorf("must specify password file")
	}
	rawPassword, err := passwordFromFile(c.String("password"))
	if err != nil {
		return nil, errors.Wrapf(err, "could not read password from file %s",
			c.String("password"))
	}
	return []byte(rawPassword), nil
}
