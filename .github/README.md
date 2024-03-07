# Chainlink Workflows

Here are all the workflows used for the main Chainlink repository.

Refer to the [chainlink-github-actions repo](https://github.com/smartcontractkit/chainlink-github-actions) where we prefer to keep most of the reusable actions. If you find yourself repeating actions that would be helpful not only here, but throughout the organization, please add them there.

If you are debugging an action, especially when it comes to runtime or resource usage, consider utilizing [catchpoint/workflow-telemetry-action](https://github.com/catchpoint/workflow-telemetry-action), and make sure to include the necessary permissions if you're using an `environment`.

```yaml
permissions:
  actions: read
  pull-requests: write # To post comments on PR
```
