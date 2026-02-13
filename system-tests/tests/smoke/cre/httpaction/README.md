# HTTP Action Workflow (Smoke & Regression)

This package contains the **HTTP Action** CRE workflow used by the smoke and regression system tests. The workflow is a single WASM binary whose behavior is driven by configuration; the test runner chooses which suite (smoke vs regression) deploys it and which test case to run.

## Theory and Design

### Purpose

The workflow exercises the **HTTP Action capability** (outbound HTTP from a CRE node) in a DON. It is used to verify:

1. **Happy path (smoke)**: CRUD-style requests (GET, POST, PUT, DELETE), correct status handling, and response handling.
2. **Multi-value headers (smoke)**: Backwards compatibility of response `Headers` (e.g. comma-joined `Set-Cookie`) and the newer response `MultiHeaders` with distinct values.
3. **Validation regression**: When both `Headers` and `MultiHeaders` are set on a request, the capability must reject it with a user error (regression test).

### Flow

1. **Trigger**: A **cron** trigger runs the workflow on a schedule (e.g. every 30 seconds).
2. **Orchestration**: The workflow runs in **node mode** so that the HTTP client can perform outbound requests from a node. The same workflow binary is used for all test cases; the **config field `testCase`** selects the behavior.
3. **Dispatch**: A map of test case names to handler functions selects either a special case (e.g. `multi-headers`, `mh-regression-both`) or the default **generic CRUD** path.
4. **Consensus**: Results are aggregated with **consensus identical aggregation** (all nodes must return the same string result).
5. **Result**: The aggregated result string is returned and the test runner asserts on it (e.g. "HTTP Action CRUD success test completed: …" or "HTTP Action multi-headers regression completed").

### Configuration

The workflow expects a YAML config (see `config/config.go`) with:

- **URL**: Target base URL for HTTP requests (e.g. fake server from the test).
- **TestCase**: Identifies which test to run (`crud-post-success`, `multi-headers`, `mh-regression-both`, etc.).
- **Method**, **Body**: HTTP method and body for the request (used by the default CRUD path and by regression).

Config is passed at workflow registration time and is read by the workflow when it runs.

## Alignment with CRE Smoke Test Spec

The parent [CRE README](../README.md) defines:

- **Smoke vs regression**: *"Everything that is not a happy path functional system-tests (i.e. edge cases, negative conditions) should go to a regression package."*
- **Test architecture**: Environment is created once per topology; multiple tests can run against the same environment; tests follow the `Test_CRE_` naming convention and standard structure.

### How this workflow fits

| Aspect | Implementation |
|--------|----------------|
| **Smoke (happy path)** | `Test_CRE_V2_HTTP_Action_Suite` and `Test_CRE_V2_HTTP_Action_Regression_Suite` (see [cre_suite_test.go](../cre_suite_test.go)) use the same workflow binary. The **smoke suite** runs success cases only: CRUD operations and multi-headers response test. |
| **Regression (edge/negative)** | The **regression suite** (`Test_CRE_V2_HTTP_Action_Regression_Suite`) runs the same workflow with `testCase: mh-regression-both`. That case sends a request with both `Headers` and `MultiHeaders` set; the capability must return a user error. The workflow asserts the error message and returns a fixed success string so the test can verify the regression. |
| **Single binary, many cases** | One compiled workflow (`main.go`) handles all cases via `testCaseHandlers` and config. No separate regression binary. |
| **Workflow compilation** | The test runner compiles this Go package to WASM (e.g. `creworkflow.CompileWorkflow`), deploys it, and registers it with the contract as described in the parent README (§11). |
| **Naming** | Workflow names are kept short (e.g. `mh-regression-both`) to stay under the workflow name length limit (64 chars) used at deploy time. |

### Test cases (config `testCase`)

- **Default (CRUD)**: Any unrecognized `testCase` runs a single HTTP request with `Content-Type: application/json`, checks 2xx status, and returns a success message including the case name.
- **`multi-headers`**: Sends two requests (Headers-only, then MultiHeaders-only), asserts response headers and Set-Cookie handling (comma-joined in `Headers`, multiple values in `MultiHeaders`).
- **`mh-regression-both`**: Sends one request with both `Headers` and `MultiHeaders` set; expects the capability to reject it with a user error containing specific substrings; returns a fixed success message for the regression suite.

## Files

- **`main.go`**: Workflow entry, cron handler, test dispatch, and all test-case logic (CRUD, multi-headers, regression). Built with `GOOS=wasip1`, `GOARCH=wasm`.
- **`config/config.go`**: Config struct (URL, TestCase, Method, Body) used for YAML unmarshalling and passed into the workflow at runtime.

## Running the tests

From the repo root (see parent README for CRE environment setup):

```bash
# Smoke: HTTP Action success cases only
go test -timeout 15m -run "^Test_CRE_V2_HTTP_Action_Suite" ./system-tests/tests/smoke/cre/...

# Regression: HTTP Action regression cases (e.g. both Headers and MultiHeaders rejected)
go test -timeout 15m -run "^Test_CRE_V2_HTTP_Action_Regression_Suite" ./system-tests/tests/smoke/cre/...
```
