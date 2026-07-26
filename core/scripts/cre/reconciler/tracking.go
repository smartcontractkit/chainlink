package reconciler

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/tracking"
)

const (
	// GetDX configuration
	getDXAPITokenVariableName = "API_TOKEN_CRE_RECONCILER"
	getDXProductName          = "cre_reconciler"

	metricCommandUsed = "cre.reconciler.command.used"
)

var (
	dxTracker      tracking.Tracker
	invokedCommand string
)

func initDxTracker() {
	if dxTracker != nil {
		return
	}
	var err error
	dxTracker, err = tracking.NewDxTracker(getDXAPITokenVariableName, getDXProductName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create DX tracker: %s\n", err)
		dxTracker = &tracking.NoOpTracker{}
	}
}

func trackCommandPreRun(cmd *cobra.Command, _ []string) {
	initDxTracker()
	invokedCommand = cmd.Name()
}

func trackCommandResult(result string) {
	if dxTracker == nil {
		return
	}
	metadata := map[string]any{"command": invokedCommand, "result": result}
	if trackErr := dxTracker.Track(metricCommandUsed, metadata); trackErr != nil {
		fmt.Fprintf(os.Stderr, "failed to track command usage: %s\n", trackErr)
	}
}
