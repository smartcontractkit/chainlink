# LOOP Plugins

Local out of process (LOOP) plugins using [github.com/hashicorp/go-plugin](https://github.com/hashicorp/go-plugin).

## Packages

```mermaid
flowchart
    subgraph chainlink-relay/pkg
        loop
        internal[loop/internal]
        pb[loop/internal/pb]
        test[loop/internal/test]

        internal --> pb
        test --> internal
        loop --> internal
        loop --> test
    end
    
    grpc[google.golang.org/grpc]
    hashicorp[hashicorp/go-plugin]

    loop ---> hashicorp
    loop ---> grpc
    test ---> grpc
    internal ---> grpc
    pb ---> grpc
    hashicorp --> grpc

```

### `package loop`

Public API and `hashicorp/go-plugin` integration.

### `package test`

Testing utilities.

### `package internal`

GRPC client & server implementations.

### `package pb`

Protocol buffer definitions & generated code.

## Communication

GRPC client/server pairs are used to communicated between the host and each plugin.
Plugins cannot communicate directly with one another, but the host can proxy a connection between them.

Here are the main components for the case of Median:  
```mermaid
sequenceDiagram
    autonumber
    participant relayer as Relayer (plugin)
    participant core as Chainlink (host)
    participant median as Median (plugin)

    Note over core: KeystoreServer
    core->>+relayer: NewRelayer
    Note over relayer: KeystoreClient
    Note over relayer: RelayerServer
    relayer->>-core: Relayer ID 
    Note over core: RelayerClient

    core->>+relayer: NewMedianProvider
    Note over relayer: MedianProviderServer
    relayer->>-core: MedianProvider ID
    Note over core: GRPC Proxy
    core->>+median: NewMedianFactory
    Note over median: MedianProviderClient
    Note over median: MedianFactoryServer
    median->>-core: MedianFactory ID
    Note over core: MedianFactoryClient


    core->>+median: NewReportingPlugin
    Note over median: ReportingPluginServer
    median->>-core: ReportingPlugin ID
    Note over core: ReportingPluginClient
```