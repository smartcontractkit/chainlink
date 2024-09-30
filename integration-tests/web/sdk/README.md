## GQL SDK

This package exports a `Client` for interacting with the `feeds-manager` service of core. The implementation is based on code generated via `[genqlient](https://github.com/Khan/genqlient)`.

### Extending the Client

Add additional queries or mutations to `genqlient.graphql` and then regenerate the implementation via the Taskfile:

```bash
$ task generate
```

Next, extend the `Client` interface and the `client` implementation.
