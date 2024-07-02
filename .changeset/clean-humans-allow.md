---
"chainlink": patch
---

# Added limit max gas cost for Optimism

A new configuration `CostMax` has been introduced. This configuration sets a maximum limit for the total gas fee that is allowed for any transaction on a chain.

## New Configuration: `CostMax`

- **Description**: Defines the maximum total gas fee allowed for any transaction.
- **Usage**: The `CostMax` is used to calculate the highest permissible gas price or fee to ensure that the total gas cost of a transaction does not exceed this limit.
