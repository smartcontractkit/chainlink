---
"chainlink": patch
---

Mercury jobs can now broadcast to multiple mercury servers.

Previously, a single mercury server would be specified in a job spec as so:

```toml
[pluginConfig]
serverURL = "example.com/foo"
serverPubKey = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93"
```

You may now specify multiple mercury servers, as so:

```toml
[pluginConfig]
servers = { "example.com/foo" = "724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93", "mercury2.example:1234/bar" = "524ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93" }
```

