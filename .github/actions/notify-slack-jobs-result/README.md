# Notify Slack Jobs Result

Composite action that posts a short Slack message for **one** job’s conclusion: a header (job name + status emoji), the raw status, and a link to the workflow run. Optionally replies in an existing thread (`thread_ts`).

Use it from a **follow-up job** that `needs` the job you care about and runs with `if: always()` so failures still notify.

## Environment

Set on the **job** or **step** that invokes this action (not as action inputs):

| Variable | Required | Purpose |
|----------|----------|---------|
| `SLACK_BOT_TOKEN` | Yes | Bot token with `chat:write` (and channel access). |
| `SLACK_CHANNEL_ID` | Yes | Channel ID for `chat.postMessage`. |

## Inputs

| Input | Required | Description |
|-------|----------|-------------|
| `status` | Yes | Conclusion string, e.g. `success`, `failure`, `cancelled`, `skipped`. Often `needs.<job_id>.result`. |
| `job_name` | Yes | Label shown in the Slack header. |
| `run_url` | Yes | URL to the workflow run. E.g. `format('{0}/{1}/actions/runs/{2}', github.server_url, github.repository, github.run_id)`. |
| `slack_thread_ts` | No | If set, the message is posted in that thread; if empty, it posts to the channel. |

Status is mapped to emoji: success ✅, failure ❌, cancelled ⚠️, skipped ⏭️, unknown ❔.

## Example

```yaml
jobs:
  tests:
    runs-on: ubuntu-latest
    steps:
      - run: ./run-tests.sh

  notify:
    runs-on: ubuntu-latest
    needs: [tests]
    if: always()
    steps:
      - uses: ./.github/actions/notify-slack-jobs-result
        env:
          SLACK_BOT_TOKEN: ${{ secrets.SLACK_BOT_TOKEN }}
          SLACK_CHANNEL_ID: ${{ secrets.SLACK_CHANNEL_ID }}
        with:
          status: ${{ needs.tests.result }}
          job_name: "Smoke tests"
          run_url: ${{ format('{0}/{1}/actions/runs/{2}', github.server_url, github.repository, github.run_id) }}
          slack_thread_ts: ${{ inputs.slack_thread_ts }} # optional
```

Implementation detail: the action builds a Block Kit payload with `jq` and sends it via [`slackapi/slack-github-action`](https://github.com/slackapi/slack-github-action) (`chat.postMessage`).
