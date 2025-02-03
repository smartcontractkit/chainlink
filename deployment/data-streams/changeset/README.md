## Overview

The files in this directory define changesets or operations to be performed. These changesets can be called by
migrations, which are the actual actions which change the state of the system.

We use several prefixes:

- `deploy_` refers to changesets that deploy a new onchain contract
- `call_` refers to changesets that call a function on an existing onchain contract

Aside from the prefixed files, we also have files with names that match the prefixes. Those contain shared code.
