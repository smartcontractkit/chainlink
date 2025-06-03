// Package por provides integration example for the POR plugin with Chainlink jobs
package main

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	por_delegate "github.com/smartcontractkit/chainlink/v2/core/services/por"
)

// Example of how the POR delegate would be integrated into the Chainlink node
func integrateWithChainlinkNode() error {
	// 1. Create logger (this would come from the main chainlink application)
	lggr := logger.NewLogger()
	
	// 2. Create the POR delegate 
	porDelegate := por_delegate.NewDelegate(lggr)
	
	// 3. Register the delegate with the job spawner (this would be done in the main app)
	// jobSpawner.RegisterDelegate(job.POR, porDelegate)
	
	// 4. When a POR job is created, the spawner would call ServicesForSpec
	porJob := job.Job{
		ID:   12345,
		Name: "production-por-job",
		Type: job.POR,
		PORSpec: &job.PORSpec{
			// Real POR configuration would go here
			// ContractAddress: "0x...",
			// ChainID: 1,
			// UpdateInterval: time.Hour,
			// etc.
		},
	}
	
	// 5. Create services for the job
	ctx := context.Background()
	services, err := porDelegate.ServicesForSpec(ctx, porJob)
	if err != nil {
		return fmt.Errorf("failed to create POR services: %w", err)
	}
	
	// 6. Start the services (this would be done by the job spawner)
	for _, service := range services {
		if err := service.Start(ctx); err != nil {
			return fmt.Errorf("failed to start POR service: %w", err)
		}
	}
	
	fmt.Println("✅ POR plugin successfully integrated and running")
	
	// 7. Services would continue running until job is deleted or node shuts down
	// When stopping, the job spawner would call Close() on each service
	
	return nil
}

func main() {
	fmt.Println("=== Chainlink POR Plugin Integration Example ===")
	
	if err := integrateWithChainlinkNode(); err != nil {
		fmt.Printf("❌ Integration failed: %v\n", err)
		return
	}
	
	fmt.Println("=== Integration Complete ===")
}

/*
In a real Chainlink node, the integration would look like this in the main application:

// In cmd/chainlink/main.go or similar:
func setupJobTypes(db *sql.DB, lggr logger.Logger, ...) {
    // ... other job type setup ...
    
    // Register POR delegate
    porDelegate := por.NewDelegate(lggr)
    jobSpawner.RegisterDelegate(job.POR, porDelegate)
    
    // ... continue with other setup ...
}

// Job definitions would be stored in the database and look like:
{
    "type": "por",
    "schemaVersion": 1,
    "name": "ethereum-usdc-por", 
    "porSpec": {
        "contractAddress": "0x1234567890123456789012345678901234567890",
        "chainId": 1,
        "updateInterval": "1h",
        "observationSource": "usdc_reserves",
        "juelsPerFeeCoin": "1000000000000000000"
    }
}

// The job spawner would automatically:
// 1. Load POR jobs from the database
// 2. Call porDelegate.ServicesForSpec() for each job
// 3. Start the returned services
// 4. Monitor service health
// 5. Restart services if they fail
// 6. Stop services when jobs are deleted
*/
