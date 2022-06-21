package chainlink_test

import (
	"context"
	"fmt"
	"os"

	"github.com/smartcontractkit/chainlink/core/services/chainlink"
)

func ExampleConfigDump_clear() {
	os.Clearenv()
	var app chainlink.ChainlinkApplication
	s, err := app.ConfigDump(context.Background())
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(s)
	}
	// Output:
	//
}

func ExampleConfigDump_Dev() {
	os.Clearenv()
	if err := os.Setenv("CHAINLINK_DEV", "true"); err != nil {
		fmt.Println(err)
		return
	}
	var app chainlink.ChainlinkApplication
	s, err := app.ConfigDump(context.Background())
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(s)
	}
	// Output:
	// Dev = true
}
