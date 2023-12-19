# Notify Slack Jobs Result

Sends a Slack message to a specified channel detailing the results of one to many GHA job results using a regex.

```yaml
inputs:
  github_token:
    description: "The GitHub token to use for authentication (usually ${{ github.token }})"
    required: true
  github_repository:
    description: "The GitHub owner/repository to use for authentication (usually ${{ github.repository }}))"
    required: true
  workflow_run_id:
    description: "The workflow run ID to get the results from (usually ${{ github.run_id }})"
    required: true
  github_job_name_regex:
    description: "The regex to use to match 1..many job name(s) to collect results from. Should include a capture group named 'cap' for the part of the job name you want to display in the Slack message (e.g. ^Client Compatability Test (?<cap>.*?)$)"
    required: true
  message_title:
    description: "The title of the Slack message"
    required: true
  slack_channel_id:
    description: "The Slack channel ID to post the message to"
    required: true
  slack_thread_ts:
    description: "The Slack thread timestamp to post the message to, handy for keeping multiple related results in a single thread"
    required: false
```
