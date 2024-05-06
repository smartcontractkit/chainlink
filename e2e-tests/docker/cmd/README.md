# QA Command-Line Tool

`coreqa` is a QA tool for Chainlink, designed to streamline the testing of Chainlink nodes and other components within a Docker-based environment.

## Available Commands

To view all available commands, run the following command:

```bash
go run main.go --help
```

## `start-nodes` Subcommand

The `start-nodes` subcommand initializes a Docker environment with multiple Chainlink nodes, setting up a private Ethereum network and necessary mock services for integration testing.

**Usage:**

```bash
go run main.go start-nodes --node-count=10
```

**Flag:**

- `--node-count` (default 6): Determines the number of Chainlink nodes to deploy

