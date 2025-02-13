## Overview

This directory contains **changesets**, which define specific operations that can be invoked by migrations to modify the on-chain system state.

### Naming Conventions

- **`deploy_`**: Deploys a new on-chain contract.
- **`call_`**: Calls a function on an existing on-chain contract.
- **`jd_`**: Interacts with the Job Distributor.

Files named exactly after these prefixes contain shared code or utility functions relevant to that category of changeset.
