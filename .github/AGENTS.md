<scope>
GitHub Actions in the Chainlink Go monorepo.
</scope>

<rules>
- Prefer runs-on runners when ubuntu-latest is insufficient.
- Minimize YAML and shell in workflows.
- Resolve smartcontractkit/.github from a local clone. Ask the user for the path if you cannot find it.
- Use `echo "key=value" | tee -a "$GITHUB_OUTPUT"` instead of `echo "key=value" >> "${GITHUB_OUTPUT}"`
</rules>

<docs>
- runs-on: https://runs-on.com/docs/
</docs>

<tools>
- https://github.com/kalverra/octometrics — per-workflow debugging
- https://github.com/kalverra/octometrics-action — runner resource monitoring
</tools>
