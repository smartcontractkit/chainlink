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
    core->>+relayer: NewRelayer(Config, Keystore)
    Note over relayer: KeystoreClient
    Note over relayer: RelayerServer
    relayer->>-core: Relayer ID 
    Note over core: RelayerClient

    core->>+relayer: NewMedianProvider(RelayArgs, PluginArgs)
    Note over relayer: MedianProviderServer
    relayer->>-core: MedianProvider ID
    Note over core: MedianProvider (Proxy)

    Note over core:  DataSourceServer
    Note over core:  ErrorLogServer

    core->>+median: NewMedianFactory(MedianProvider, DataSource, ErrorLog)
    Note over median: MedianProviderClient
    Note over median: DataSourceClient
    Note over median: ErrorLogClient
    Note over median: MedianFactoryServer
    median->>-core: MedianFactory ID
    Note over core: MedianFactoryClient

    core->>+median: NewReportingPlugin(ReportingPluginConfig)
    Note over median: ReportingPluginServer
    median->>-core: ReportingPlugin ID
    Note over core: ReportingPluginClient
```
Note: MedianProvider includes multiple component services on the same connection.
```mermaid
sequenceDiagram
    autonumber
    participant relayer as Relayer (plugin)
    participant core as Chainlink (host)
    participant median as Median (plugin)

    core->>+relayer: NewMedianProvider(RelayArgs, PluginArgs)
    Note over relayer: OffchainConfigDigesterServer
    Note over relayer: ContractConfigTrackerServer
    Note over relayer: ContractTransmitterServer
    Note over relayer: ReportCodecServer
    Note over relayer: MedianContractServer
    Note over relayer: OnchainConfigCodecServer
    
    relayer->>-core: MedianProvider ID
    Note over core: MedianProvider (Proxy)
    
    Note over core: OffchainConfigDigesterClient
    Note over core: ContractConfigTrackerClient
    Note over core: ContractTransmitterClient
    
    core->>+median: NewMedianFactory(MedianProvider, DataSource, ErrorLog)
    Note over median: ReportCodecClient
    Note over median: MedianContractClient
    Note over median: OnchainConfigCodecClient
```