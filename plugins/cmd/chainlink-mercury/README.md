This directory house the Mercury LOOPP

# Running Integration Tests Locally

Running the tests is as simple as
- building this binary
- setting the CL_MERCURY_CMD env var to the *fully resolved* binary path
- running the test(s)


The interesting tests are `TestIntegration_MercuryV*` in ` github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury`

In detail:
```
sh

make install-mercury-loop # builds `mercury` binary in this dir
CL_MERCURY_CMD=<YOUR_REPO_ROOT_DIR>/plugins/cmd/mercury/mercury go test -v -timeout 120s -run ^TestIntegration_MercuryV github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury 2>&1 | tee /tmp/mercury_loop.log
```