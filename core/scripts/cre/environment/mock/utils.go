package mock

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	mockcapability "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock"
	pb2 "github.com/smartcontractkit/chainlink/system-tests/lib/cre/mock/pb"
)

var MockCommand = &cobra.Command{
	Use:   "mock",
	Short: "Mock CRE capability tools",
}

// newMockCapabilityController creates a new MockCapabilityController with a standard logger
func newMockCapabilityController() *mockcapability.Controller {
	lggr := zerolog.New(os.Stdout)
	c := mockcapability.NewMockCapabilityController(lggr)
	err := c.ConnectAll([]string{"172.28.0.13:7777", "172.28.0.14:7777", "172.28.0.15:7777"}, true, false) // TODO: @george-dorin get addresses
	if err != nil {
		panic(err)
	}
	return c
}

func Create(ctx context.Context, capabilityWithVersion string, capabilityType string, description string) error {
	mocks := newMockCapabilityController()
	exists, err := mocks.HasCapability(ctx, capabilityWithVersion)
	if err != nil {
		return err
	}
	if !exists {
		err = mocks.CreateCapability(ctx, &pb2.CapabilityInfo{
			ID:             capabilityWithVersion,
			CapabilityType: mockcapability.StringToCapabilityType(capabilityType),
			Description:    description,
			DON:            nil,
			IsLocal:        true,
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func Delete(ctx context.Context, capabilityWithVersion string) error {
	mocks := newMockCapabilityController()
	return mocks.DeleteCapability(ctx, capabilityWithVersion)
}

func WatchExecutables(ctx context.Context) error {
	mocks := newMockCapabilityController()
	responseCh := make(chan capabilities.CapabilityRequest)
	err := mocks.HookExecutables(ctx, responseCh)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	fmt.Printf("Watching for executable requests\n")
	for {
		// Print all requests from chan
		select {
		case response := <-responseCh:
			fmt.Printf("Received response: %s\n", response.CapabilityId)
			spew.Dump(response)
		case <-ctx.Done():
			fmt.Printf("Context done\n")
			return nil
		}
	}
}

func SendTrigger(ctx context.Context, id string, mockDataType string, frequency time.Duration, duration time.Duration) error {
	mocks := newMockCapabilityController()

	// Check if trigger has subscribers
	err := mocks.WaitForTriggerSubscribers(context.Background(), id, time.Second*30)
	if err != nil {
		return err
	}

	var runCtx context.Context
	var cancel context.CancelFunc

	if duration > 0 {
		// Create a timeout context if duration is specified
		runCtx, cancel = context.WithTimeout(ctx, duration)
	} else {
		// Use the parent context if no duration specified (run forever)
		runCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	// Create ticker for the specified frequency
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	// Statistics
	triggersSent := 0
	errorCount := 0
	startTime := time.Now()

	// Create a ticker for periodic statistics reporting (every 5 seconds)
	statsTicker := time.NewTicker(5 * time.Second)
	defer statsTicker.Stop()

	logger := zerolog.New(os.Stdout)

	for {
		select {
		case <-runCtx.Done():
			// When finished, log final statistics
			elapsed := time.Since(startTime)
			logger.Info().
				Int("triggers_sent", triggersSent).
				Int("errors", errorCount).
				Dur("elapsed_time", elapsed).
				Float64("triggers_per_second", float64(triggersSent)/elapsed.Seconds()).
				Msg("SendTrigger completed")
			return nil

		case <-statsTicker.C:
			// Periodically log statistics while running
			elapsed := time.Since(startTime)
			if elapsed.Seconds() > 0 {
				logger.Info().
					Int("triggers_sent", triggersSent).
					Int("errors", errorCount).
					Dur("elapsed_time", elapsed).
					Float64("triggers_per_second", float64(triggersSent)/elapsed.Seconds()).
					Msg("SendTrigger progress")
			}

		case <-ticker.C:
			// Send trigger at each tick
			data, err := getTriggerRequest(MockTriggerType(mockDataType))
			if err != nil {
				return err
			}
			data.TriggerID = id
			data.TriggerType = id
			spew.Dump(data)
			err = mocks.SendTrigger(runCtx, data)

			triggersSent++

			if err != nil {
				errorCount++
				logger.Error().Err(err).Msg("Error sending trigger")
			}
		}
	}
	return nil
}

func runWatchExecutables(cmd *cobra.Command, args []string) error {
	return WatchExecutables(cmd.Context())
}

func runCreate(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("id")
	capType, _ := cmd.Flags().GetString("type")
	desc, _ := cmd.Flags().GetString("description")

	return Create(cmd.Context(), id, capType, desc)
}

func runDelete(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("id")

	return Delete(cmd.Context(), id)
}

func runSendTrigger(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("id")
	dataType, _ := cmd.Flags().GetString("type")
	frequency, _ := cmd.Flags().GetDuration("frequency")
	duration, _ := cmd.Flags().GetDuration("duration")

	return SendTrigger(cmd.Context(), id, dataType, frequency, duration)
}

func init() {
	// Create command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a capability",
		RunE:  runCreate,
	}
	createCmd.Flags().String("id", "", "Capability ID with version")
	createCmd.Flags().String("type", "", "Capability type")
	createCmd.Flags().String("description", "", "Capability description")
	createCmd.MarkFlagRequired("id")
	createCmd.MarkFlagRequired("type")

	// Delete command
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a capability",
		RunE:  runDelete,
	}
	deleteCmd.Flags().String("id", "", "Capability ID with version")
	deleteCmd.MarkFlagRequired("id")

	// SendTrigger command
	triggerCmd := &cobra.Command{
		Use:   "trigger",
		Short: "Send trigger events at specified intervals",
		RunE:  runSendTrigger,
	}
	triggerCmd.Flags().String("id", "", "Trigger ID")
	triggerCmd.Flags().String("type", string(TriggerTypeCron), "Mock data type (empty, cron)")
	triggerCmd.Flags().Duration("frequency", 1*time.Second, "Frequency to send triggers")
	triggerCmd.Flags().Duration("duration", 0*time.Second, "Duration to send triggers in seconds (0 for unlimited)")
	triggerCmd.MarkFlagRequired("id")

	// Watch command
	watchCmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch for executable requests",
		RunE:  runWatchExecutables,
	}

	MockCommand.AddCommand(createCmd, deleteCmd, triggerCmd, watchCmd)
}
