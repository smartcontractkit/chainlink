## GQL SDK

This package exports a `Client` for interacting with the `feeds-manager` service of core. The implementation is based on code generated via [genqlient](https://github.com/Khan/genqlient).

### Structure

```
client/          # Client interface, implementation, and types
internal/        # GraphQL queries/mutations and generated code
  genqlient.graphql
  generated/
```

### Prerequisites

- [go-task](https://taskfile.dev/) (`task` CLI)
- Go 1.21+

### Extending the Client

If your feature requires **new GraphQL operations**, add queries or mutations to `genqlient.graphql` and then regenerate the implementation via the Taskfile:

```bash
$ task generate
```

Next, extend the `Client` interface and the `client` implementation.

If your feature **composes existing operations** (e.g. filtering or combining results from already-generated queries), you can extend the `Client` interface and implementation directly without modifying `genqlient.graphql` or regenerating.

### Usage

See `devenv/don.go` and `keystone/scripts/main.go` for real-world usage examples.
