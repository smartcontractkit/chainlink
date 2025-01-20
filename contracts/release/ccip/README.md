# CCIP Contracts Releases

This directory contains the changelogs, version (via `package.json`), and the changesets for the CCIP contracts.

## Overview

The actual CCIP contracts code currently lives in the `contracts/src/*/ccip` directory in order to share code with other Chainlink contracts. Even though this CCIP code directory is under the `@chainlink/contracts`'s `package.json` file, it's not part of the `@chainlink/contracts` NPM package and should be versioned, released, and published separately which is why this directory exists.

## Directory Structure

```
ccip/
├── .changeset/     # Contains changesets for versioning
├── CHANGELOG.md    # Auto-generated changelog from changesets
└── package.json    # @chainlink/contracts-ccip package configuration for versioning
```

## Create a Changeset

To be ran from the (`./contracts`) directory.

1. Create a changeset for your changes:

   ```shell
   pnpm changeset:ccip
   ```

2. Follow the prompts to describe your changes
3. Commit the generated changeset file
