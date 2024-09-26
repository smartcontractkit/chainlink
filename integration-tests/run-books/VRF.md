# VRF Tests Run Book 
* All test configs should be placed in the [integration-tests/testconfig](../testconfig) folder  
* All test configs for running tests in live testnets should be under [integration-tests/testconfig/vrfv2plus/overrides](../testconfig/vrfv2plus/overrides) folder 


## Functional Tests
### In CI - using On Demand Workflows
```bash
gh workflow run "on-demand-vrfv2plus-smoke-tests.yml" \
--ref develop \
-f=test_secrets_override_key=<your testsecrets id> \
-f test_config_override_path=<path to test toml config which should be in `integration-tests/testconfig/vrfv2plus/overrides` folder> \
-f test_suite="Selected Tests" \ # Optional, Options - "All Tests", "Selected Tests". Default is "All Tests". If "Selected Tests" is selected, then `test_list_regex` should be provided 
-f test_list_regex="<regex for tests to run>" \ # Optional, default is "TestVRFv2Plus$/(Link_Billing|Native_Billing|Direct_Funding)|TestVRFV2PlusWithBHS" which are P0 tests
-f chainlink_version="<>" # Optional, default is image created from develop branch. Not needed if you run tests against existing environment
-f notify_user_id_on_failure=<your slack user id> # Optional, default is empty. If provided, will notify the user on slack if the tests fail
```

#### Examples:
Run P0 tests against existing environment (Staging) on Arbitrum Sepolia
```bash
gh workflow run "on-demand-vrfv2plus-smoke-tests.yml" \
--ref develop \
-f=test_secrets_override_key=<your testsecrets id> \
-f test_config_override_path=integration-tests/testconfig/vrfv2plus/overrides/staging/arbitrum_sepolia_staging_test_config.toml \
-f test_suite="Selected Tests" 
```

Run all tests deploying all contracts, CL nodes with `2.15.0` version on Base Sepolia
```bash
gh workflow run "on-demand-vrfv2plus-smoke-tests.yml" \
--ref develop \
-f=test_secrets_override_key=<your testsecrets id> \
-f test_config_override_path=integration-tests/testconfig/vrfv2plus/overrides/new_env/base_sepolia_new_env_test_config.toml \
-f test_suite="All Tests" \
-f chainlink_version="2.15.0" 
```

### Locally
```bash
cd integration-tests
TEST_LOG_LEVEL=debug \
BASE64_CONFIG_OVERRIDE=$(cat <full path to chain's test toml config> | base64) \
go test -v -timeout 15m -run "<regex for tests to run>" ./smoke
```

## Performance Tests
```bash
gh workflow run "on-demand-vrfv2plus-performance-test.yml" \
--ref develop \
-f=test_secrets_override_key=<your testsecrets id> \
-f test_config_override_path=<path to test toml config which should be in `integration-tests/testconfig/vrfv2plus/overrides` folder> \
-f performanceTestType=“Smoke” # Options - "Smoke", "Soak", "Stress", "Load".
-f test_list_regex="<regex for tests to run>" # Optional, default is "TestVRFV2PlusPerformance"
```

#### Examples:
Run SOAK tests against existing environment (Staging) on Base Sepolia
```bash
gh workflow run "on-demand-vrfv2plus-performance-test.yml" \
--ref develop \
-f=test_secrets_override_key=<your testsecrets id> \
-f test_config_override_path=integration-tests/testconfig/vrfv2plus/overrides/staging/base_sepolia_staging_test_config.toml \
-f performanceTestType=“Soak”
```
