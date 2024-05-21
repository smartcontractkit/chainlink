# Operator UI

NOTE: If you're looking for the source of operator UI, it has now been moved to https://github.com/smartcontractkit/operator-ui

This directory instead now as a collection of scripts for maintaining the version of operator UI to pull in when developing and building the chainlink node.

## About

This package is responsible for rendering the UI of the chainlink node, which allows interactions with node jobs, jobs runs, configuration and any other related tasks.

## Installation

### Requirements

The `install.sh` script handles installing the specified tag of operator UI within the [tag file](./TAG). When executed, it downloads then moves the static assets of operator UI into the `core/web/assets` path. Then, when the chainlink binary is built, these assets are included into the build that gets served.

## Updates

### Requirements

- gh cli ^2.15.0 https://github.com/cli/cli/releases/tag/v2.15.0
- jq ^1.6 https://stedolan.github.io/jq/

The `update.sh` script will check for the latest release from the `smartcontractkit/operator-ui` repository, if the latest release is newer than the current tag, it'll update the [tag file](./TAG) with the corresponding latest tag. Checking for updates is automatically [handled by CI](../.github/workflows/operator-ui.yml), where any new detected updates will be pushed to a branch and have a PR opened against `develop`.

See https://docs.github.com/en/rest/releases/releases#get-the-latest-release for how a latest release is determined.
