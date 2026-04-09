CI/CD for the Chainlink Go monorepo.

<resources>
Prefer runs-on runners: https://runs-on.com/docs/
Resolve smartcontractkit/.github actions and workflows from the local checkout first. Ask the user for a path if missing. Fetch the web only for a specific version, commit, or when local behavior disagrees with docs.
Use [octometrics-action](https://github.com/kalverra/octometrics-action) for debugging resource usage.
```yaml
example-job:
  name: Example Job
  runs-on: ubuntu-latest
  steps:
    - name: Monitor
      uses: kalverra/octometrics-action
      with:
        job_name: Example Job
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }} # Optional, but highly recommended to prevent rate limiting
    - name: Checkout code
      uses: actions/checkout@v4
    - name: Run rest of workflow
      run: |
        echo "Hello World!"
```
</resources>

<priorities>
When changing workflows or actions, rank goals in this order:
1. Maintainability: Prefer composable Actions and small scripts over large inline YAML or bash.
2. Security: Apply least privilege and sound secret handling.
3. Reliability: Fail clearly; handle transient errors where appropriate.
4. Speed: Reduce wall-clock time.
5. Cost: Reduce runner spend.
</priorities>

<rules>
Pin 3rd party action versions to commit hashes, not version tags.
```yaml
- name: Enable S3 Cache for Self-Hosted Runners
  uses: runs-on/action@742bf56072eb4845a0f94b3394673e4903c90ff0 # v2.1.0
  with:
    metrics: cpu,network,memory,disk
```
</rules>
