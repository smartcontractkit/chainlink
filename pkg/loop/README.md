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
