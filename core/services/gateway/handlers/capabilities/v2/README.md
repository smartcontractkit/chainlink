# HTTP Handlers V2

> Implements [Gateway Handler interface](../../handler.go). Enables Chainlink Runtime Environment (CRE) workflows to interact with external systems through HTTP actions and triggers.

---

## 1. System Overview

### 1.1 Purpose

The HTTP Handlers V2 is responsible for:
- Dispatching outbound HTTP requests originating from HTTP Action capabilities
- Receiving inbound HTTP trigger requests and forwarding to HTTP Trigger capabilities to initiate workflows
- Receiving workflow metadata from HTTP Trigger capabilities and aggregating the metadata for authentication and workflow lookup 

### 1.2 Core Functionality

- **HTTP Actions**: Outbound HTTP requests with caching, rate limiting
- **HTTP Triggers**: Inbound requests that trigger workflow executions with JWT-based authentication
- **Auth Metadata Management**: Collection and aggregation of workflow authorization data
- **Response Caching**: Caching of HTTP responses to reduce redundant requests
- **Rate Limiting**: Multi-dimensional rate limiting (global, per workflow owner, per capability node)
- **Response Aggregation**: Byzantine fault-tolerant aggregation of node responses

---

## 2. Architecture

### 2.1 Components

#### 2.1.1 Gateway Handler (`gatewayHandler`)
- **Purpose**: Main orchestrator that coordinates all HTTP handling functionality
- **Functions**: Request routing, lifecycle management, cache management
- **Key Features**: Supports both HTTP actions and trigger requests, manages sub-handlers

#### 2.1.2 HTTP Trigger Handler (`httpTriggerHandler`)
- **Purpose**: Processes inbound HTTP trigger requests that initiate workflow executions
- **Functions**: Request validation, authorization, rate limiting, response aggregation
- **Key Features**: JWT authentication, workflow resolution, dispatching requests to HTTP Trigger nodes, BFT aggregation

#### 2.1.3 Workflow Metadata Handler (`WorkflowMetadataHandler`)
- **Purpose**: Manages workflow authorization metadata and keys
- **Functions**: Workflow metadata collection, aggregation, JWT verification, workflow selector mapping
- **Key Features**: BFT aggregation of metadata, periodic synchronization

#### 2.1.4 Response Cache (`responseCache`)
- **Purpose**: Caches HTTP responses to avoid redundant outbound requests
- **Functions**: TTL-based caching that optionally returns cached values based on max age parameter
- **Key Features**: Workflow-scoped caching

---

## 3. HTTP Action Message Handling

### 3.1 Process Flow

1. **Request Reception**: Gateway receives HTTP action request from a workflow node
2. **Rate Limiting**: Validates node rate limits
3. **Request Parsing**: Extracts `OutboundHTTPRequest` from the JSON-RPC message
4. **Cache Check**: Determines if request should use cached response or fetch fresh data
5. **HTTP Execution**: Makes actual HTTP request to external endpoint
6. **Response Caching**: Stores cacheable responses (2xx, 4xx status codes) only if `CacheSettings.Store` is `true`
7. **Node Response**: Sends HTTP response back to requesting node

### 3.2 Caching Behavior

- **Cacheable Responses**: 2xx (success) and 4xx (client error) status codes.
- **Cache TTL**: Configurable, default 10 minutes
- **Cache Key**: Generated from workflow ID and request hash
- **Cache Invalidation**: Time-based expiration with periodic cleanup
- **Cache Strategy**: All cacheable responses are cached; Non-zero `CacheSettings.MaxAgeMs` determines whether to return a cached value or make a fresh request
- **Workflow Isolation**: Cache entries are scoped by workflow ID to prevent cross-workflow data leakage
---

## 4. HTTP Trigger Message Handling

### 4.1 Process Flow

1. **Request Validation**: Validates JSON-RPC format, method, and parameters
2. **Workflow Resolution**: Resolves workflow ID from selector (ID, owner, name, tag)
3. **Authentication**: Verifies JWT token (ECDSA signature) and checks authorized keys
4. **Rate Limiting**: Enforces per-workflow-owner rate limits
5. **Node Distribution**: Sends request to all DON members with retry logic
6. **Response Aggregation**: Collects and aggregates responses from nodes (2f + 1 identical responses required, where f is max faulty nodes)
7. **User Response**: Returns aggregated result to the original requester

---

## 5. Auth Metadata Messages and Aggregation Logic

### 5.1 Metadata Collection Process

The system implements a workflow metadata collection and aggregation system to sync workflow metadata from workflow nodes to gateway nodes.
There are 2 flows:

#### 5.1.1 Metadata Push (Registration Events)
- **Trigger**: Workflow registration event
- **Process**: HTTP capability nodes push workflow metadata to gateway

#### 5.1.2 Metadata Pull (Periodic Sync)
- **Trigger**: Periodic timer (default 1 minute intervals)
- **Process**: Gateway requests metadata from all HTTP capability nodes, which respond with batches of workflow metadata.

### 5.2 Aggregation Logic

The aggregation system (located in `/core/services/gateway/common/aggregation/`) implements Byzantine fault-tolerant metadata collection:

1. **Observation Collection**: Each node's metadata is hashed and stored by digest
2. **Threshold Consensus**: Requires f+1 identical observations (where f is max faulty nodes)
3. **Duplicate Prevention**: Ensures unique workflow ID and reference mappings
4. **Periodic Cleanup**: Removes expired observations to prevent memory leaks

### 5.3 Synchronization Flow

1. **Collection**: Nodes submit metadata observations
2. **Aggregation**: Aggregator identifies consensus metadata (f+1 agreements)
3. **Sync**: Gateway updates local cache with aggregated metadata
4. **Cleanup**: Expired observations are removed periodically
---

## 6. Configuration Specification

### 6.1 Service Configuration Schema

```go
type ServiceConfig struct {
    NodeRateLimiter               ratelimit.RateLimiterConfig `json:"nodeRateLimiter"`
    UserRateLimiter               ratelimit.RateLimiterConfig `json:"userRateLimiter"`
    MaxTriggerRequestDurationMs   int                         `json:"maxTriggerRequestDurationMs"`
    RetryConfig                   RetryConfig                 `json:"retryConfig"`
    CleanUpPeriodMs               int                         `json:"cleanUpPeriodMs"`
    MetadataPullIntervalMs        int                         `json:"metadataPullIntervalMs"`
    MetadataAggregationIntervalMs int                         `json:"metadataAggregationIntervalMs"`
    OutboundRequestCacheTTLMs     int                         `json:"outboundRequestCacheTTLMs"`
}
```

### 6.2 Rate Limiter Configuration

```go
type RateLimiterConfig struct {
    GlobalRPS      float64 `json:"globalRPS"`      // Global requests per second
    GlobalBurst    int     `json:"globalBurst"`    // Global burst capacity
    PerSenderRPS   float64 `json:"perSenderRPS"`   // Per-sender requests per second
    PerSenderBurst int     `json:"perSenderBurst"` // Per-sender burst capacity
}
```

### 6.3 Retry Configuration

```go
type RetryConfig struct {
    InitialIntervalMs int     `json:"initialIntervalMs"` // Initial retry interval
    MaxIntervalTimeMs int     `json:"maxIntervalTimeMs"` // Maximum retry interval
    Multiplier        float64 `json:"multiplier"`        // Backoff multiplier
}
```

### 6.4 Default Values

| Configuration | Default Value | Description |
|---------------|---------------|-------------|
| `CleanUpPeriodMs` | 600000 (10 min) | Cache and callback cleanup interval |
| `MaxTriggerRequestDurationMs` | 60000 (1 min) | Maximum time for trigger request processing |
| `MetadataPullIntervalMs` | 60000 (1 min) | Interval for pulling metadata from nodes |
| `MetadataAggregationIntervalMs` | 60000 (1 min) | Interval for aggregating collected metadata |
| `InitialIntervalMs` | 100 | Initial retry interval |
| `MaxIntervalTimeMs` | 30000 (30 sec) | Maximum retry interval |
| `Multiplier` | 2.0 | Exponential backoff multiplier |
| `OutboundRequestCacheTTLMs` | 600000 (10 min) | HTTP response cache TTL |

### 6.5 Configuration Example

```json
{
  "nodeRateLimiter": {
    "globalRPS": 100.0,
    "globalBurst": 100,
    "perSenderRPS": 10.0,
    "perSenderBurst": 20
  },
  "userRateLimiter": {
    "globalRPS": 50.0,
    "globalBurst": 50,
    "perSenderRPS": 5.0,
    "perSenderBurst": 10
  },
  "maxTriggerRequestDurationMs": 60000,
  "retryConfig": {
    "initialIntervalMs": 100,
    "maxIntervalTimeMs": 30000,
    "multiplier": 2.0
  },
  "cleanUpPeriodMs": 600000,
  "metadataPullIntervalMs": 60000,
  "metadataAggregationIntervalMs": 60000,
  "outboundRequestCacheTTLMs": 600000
}
```

---

## 7. Security Features

### 7.1 Authentication & Authorization

- **JWT Verification**: All trigger requests must include valid JWT tokens
- **Address Validation**: All addresses must be 0x-prefixed and lowercase
- **Workflow-Scoped Auth**: Each workflow maintains its own authorized key set

### 7.2 Rate Limiting

- **Dual Rate Limiting**: Separate limits for node and user requests
- **Per-Sender Limits**: Individual rate limits per sending entity
- **Global Limits**: System-wide rate limiting for overall protection

### 7.3 Input Validation

- **Request ID Validation**: Prevents malicious request ID injection
- **JSON Validation**: Ensures valid JSON input for workflow parameters
- **Workflow Field Validation**: Validates workflow selector format
- **Public Key Validation**: Ensures proper ECDSA key format

---

## 8. Error Handling

### 8.1 Error Types

- **Authentication Errors**: Invalid JWT or unauthorized keys
- **Validation Errors**: Malformed requests or invalid parameters
- **Rate Limit Errors**: Exceeded rate limits
- **Timeout Errors**: Request processing timeouts
- **Network Errors**: HTTP request failures
- **Internal Errors**: System or aggregation failures

### 8.2 Error Response Format

```json
{
  "jsonrpc": "2.0",
  "id": "request-id",
  "error": {
    "code": -32602,
    "message": "Invalid request: Auth failure"
  }
}
```

### 8.3 Error Codes

| Code | Description |
|------|-------------|
| `-32700` | Parse error |
| `-32600` | Invalid request |
| `-32601` | Method not found |
| `-32602` | Invalid params |
| `-32603` | Internal error |
| `-32000` | Rate limit exceeded |
| `-32001` | Conflict (duplicate request) |

---

## 9. Implementation Details

### 9.1 Request ID Format

- **User Requests**: Plain string identifiers (cannot contain "/")
- **Node Messages**: Format `<methodName>/<workflowID>/<uuid>` or `<methodName>/<workflowID>/<workflowExecutionID>/<uuid>`
- **Method Routing**: Gateway routes messages based on method name in request ID

### 9.2 Workflow ID Extraction

For HTTP action requests, workflow ID is extracted from the request path using the pattern:
```
<methodName>/<workflowID>/...
```
The workflow ID is the second segment after splitting by "/".

### 9.3 Response Aggregation Details

- **Aggregator Type**: `IdenticalNodeResponseAggregator`
- **Consensus Mechanism**: Requires exact response matching (by digest)
- **Node Tracking**: Tracks which response each node provided
- **Response Updates**: If a node provides a different response, it's moved to the new response group
- **Threshold**: 2f+1 identical responses needed for consensus

### 9.4 Broadcast Retry Logic

- **Trigger Requests**: Exponential backoff with jitter for failed node deliveries
- **Timeout**: Maximum trigger request duration of 1 minute (configurable)
- **Partial Success**: Continues retrying until all nodes receive the request or timeout

---

## 10. Development

### 10.1 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -coverpkg=./... -coverprofile=coverage.txt

# Run race condition tests
GORACE="log_path=$PWD/race"  go test -race ./...
```

---