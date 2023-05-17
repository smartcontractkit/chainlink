# LOOP Plugins

Local out of process (LOOP) plugins using [github.com/hashicorp/go-plugin](https://github.com/hashicorp/go-plugin).

## Packages

```mermaid
flowchart
    loop
    internal[loop/internal]
    pb[loop/internal/pb]
    test[loop/internal/test]

    internal --> pb
    test --> internal
    loop --> internal
    loop --> test

```

### `package loop`

Public API and `hashicorp/go-plugin` integration.

### `package test`

Testing utilities.

### `package internal`

GRPC client & server implementations.

### `package pb`

Protocol buffer definitions & generated code.
