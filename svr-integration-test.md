# Implement SVR Secondary Transmission Test Verification

## Current State
The test `TestIntegration_secondary_feed_transmission` in `core/internal/features/svr/txmv2_test.go`:
- Sets up a DON with bootstrap + 4 oracle nodes
- Configures dual transmission with primary (keys[0]) and secondary (keys[1]) transmitters
- Creates OCR2 jobs with `enableDualTransmission: true`
- Currently just waits for blocks indefinitely without verification

## How to Run the Test

The test can be run with the following command (postgres container must be running):

```bash
CL_DATABASE_URL=postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_development_test?sslmode=disable go test -run ^TestIntegration_secondary_feed_transmission$ github.com/smartcontractkit/chainlink/v2/core/internal/features/svr -v
```

## Implementation Plan

### 1. Check Node Logs for Transaction Evidence
   - Primary transaction: Verify in node logs that a transaction was sent to chain
     - Look for log messages indicating primary transmission was sent
     - Can also verify on-chain if needed (transaction receipt exists)
   - Secondary transaction: Verify in node logs that a transaction was sent to Flashbots
     - Look for log messages indicating secondary transmission was sent to Flashbots
     - Note: Secondary transaction won't appear on-chain in this test setup

### 2. Replace Infinite Loop with Eventually Pattern
   - Remove the current infinite block-waiting loop (lines 341-355)
   - Use `gomega.Eventually` to wait for both primary and secondary transaction log evidence
   - Set a reasonable timeout (e.g., 5 minutes)
   - Continue committing blocks during the wait period using a ticker
   - When timeout occurs, fail the test with clear error message indicating which transaction type(s) were missing

### 3. Add Verification Logic
   - Access node logs from the test setup (nodes have `ObservedLogs` field)
   - Search logs for evidence of:
     - Primary transmission: Look for log messages about sending primary transaction to chain
     - Secondary transmission: Look for log messages about sending secondary transaction to Flashbots
   - Assert that at least one of each type occurred in the logs
   - Only verify that transactions were attempted/sent, not their on-chain success

### 4. Implementation Details
   - Use `node.ObservedLogs` (type `*observer.ObservedLogs`) to access logs
   - Search for log messages containing keywords like:
     - Primary: "transmit", "primary", "transaction", etc.
     - Secondary: "transmitSecondary", "secondary", "flashbots", etc.
   - Use `observer.ObservedLogs.FilterMessage()` or similar to search logs
   - Check logs from all oracle nodes (not bootstrap)
   - Track whether both transaction types were found in logs

### 5. Test Structure
   ```go
   // Start block ticker
   tick := time.NewTicker(1 * time.Second)
   defer tick.Stop()
   go func() {
       for range tick.C {
           backend.Commit()
       }
   }()
   
   // Use Eventually to wait for both transactions in logs
   gomega.NewGomegaWithT(t).Eventually(func() bool {
       primaryFound := false
       secondaryFound := false
       
       // Check logs from all oracle nodes
       for _, node := range nodes {
           logs := node.ObservedLogs.All()
           for _, log := range logs {
               msg := log.Message
               // Check for primary transmission
               if strings.Contains(msg, "transmit") && 
                  (strings.Contains(msg, "primary") || !strings.Contains(msg, "secondary")) {
                   primaryFound = true
               }
               // Check for secondary transmission to Flashbots
               if strings.Contains(msg, "transmitSecondary") || 
                  (strings.Contains(msg, "secondary") && strings.Contains(msg, "flashbots")) {
                   secondaryFound = true
               }
           }
       }
       return primaryFound && secondaryFound
   }, 5*time.Minute, 1*time.Second).Should(gomega.BeTrue(), 
       "Expected both primary (to chain) and secondary (to Flashbots) transmissions in logs")
   ```

## Files to Modify
- `core/internal/features/svr/txmv2_test.go` - Main test file

## Verification Approach
1. Access logs from all oracle nodes using `node.ObservedLogs`
2. Search log messages for evidence of primary transmission (to chain)
3. Search log messages for evidence of secondary transmission (to Flashbots)
4. Track whether both types were found
5. Use `gomega.Eventually` with timeout to wait for both
6. Fail test if timeout is reached without both transaction types appearing in logs

## Notes
- Secondary transaction goes to Flashbots, not directly to chain
- We verify transaction attempts in logs, not on-chain success
- Need to identify the exact log message patterns used by TXM for dual transmission

## Todos
1. Identify log message patterns for primary and secondary transmissions in TXM code
2. Add log searching logic to check node.ObservedLogs for evidence of both transaction types
3. Replace infinite block-waiting loop with gomega.Eventually pattern and block ticker
4. Add timeout handling and clear error messages when verification fails
