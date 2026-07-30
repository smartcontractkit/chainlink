# Blank Workflow Example

This template provides a blank workflow example. It aims to give a starting point for writing a workflow from scratch and to get started with local simulation.

Steps to run the example

## 1. Update .env file

You need to add a private key to env file. This is specifically required if you want to simulate chain writes. For that to work the key should be valid and funded.
If your workflow does not do any chain write then you can just put any dummy key as a private key. e.g.
```
CRE_ETH_PRIVATE_KEY=0000000000000000000000000000000000000000000000000000000000000001
```

## 2. Simulate the workflow
Run the command from <b>project root directory</b>

```bash
cre workflow simulate <path-to-workflow> --target=staging-settings
```

It is recommended to look into other existing examples to see how to write a workflow. You can generate then by running the `cre init` command.
