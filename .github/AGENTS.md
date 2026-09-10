GitHub Actions in the Chainlink Go monorepo.

## Rules

- Prefer runs-on runners when ubuntu-latest is insufficient.
- Minimize YAML and shell in workflows. Utilize the [Go CI CLI](tools/ci/) for any bash requiring more than basic commands.
- Resolve smartcontractkit/.github from a local clone. Ask the user for the path if you cannot find it.

## New Features

- Can [reference local actions without checkout step](https://github.blog/changelog/2026-07-30-reference-same-repository-actions-with-self-repository-syntax/) with `uses: $/path/to/action`

## Docs & Tools

- [runs-on](https://runs-on.com/docs/): Docs for our runs-on runners
- [octometrics](https://github.com/kalverra/octometrics): Debugging and analyzing workflows
- [octometrics-action](https://github.com/kalverra/octometrics-action): Runner resource monitoring
- [Go CI CLI](tools/ci/): Custom logic for /chainlink workflows
