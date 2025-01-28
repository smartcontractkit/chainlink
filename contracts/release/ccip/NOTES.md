# Chainlink CCIP Smart Contracts

This directory contains the changelogs, version (via `package.json`), and the changesets for the CCIP contracts.

## Overview

The actual CCIP contracts code currently lives in the `contracts/src/*/ccip` directory in order to share code with other Chainlink contracts. Even though this CCIP code directory is under the `@chainlink/contracts`'s `package.json` file, it's not part of the `@chainlink/contracts` NPM package and should be versioned, released, and published separately which is why this directory exists.

## Create a Changeset

To be ran from the (`./contracts`) directory.

1. Create a changeset for your changes:

   ```shell
   pnpm changeset:ccip
   ```

2. Follow the prompts to describe your changes
3. Commit the generated changeset file

## CCIP Contracts Release

To be ran from the (`./contracts`) directory. Copy files over from `./contracts/release/ccip`/ to `./contracts/`.

```shell
# To undo the copy, run `git checkout -- package.json README.md`
pnpm copy:ccip-files
```
