// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice: IMPORTANT NOTICE for anyone who wants to use this contract
/// @notice Source: https://github.com/ethereum-optimism/optimism/blob/71b93116738ee98c9f8713b1a5dfe626ce06c1b2/packages/contracts-bedrock/src/libraries/L1BlockErrors.sol
/// @notice The original code was trimmed down to include only the necessary interface elements required to interact with GasPriceOracle
/// @notice We need this file so that Solidity compiler will not complain because some functions don't exist
/// @notice In reality, we don't embed this code into our own contracts, instead we make cross-contract calls on predeployed GasPriceOracle contract

/// @notice Error returns when a non-depositor account tries to set L1 block values.
error NotDepositor();

/// @notice Error when a chain ID is not in the interop dependency set.
error NotDependency();

/// @notice Error when the interop dependency set size is too large.
error DependencySetSizeTooLarge();

/// @notice Error when a chain ID already in the interop dependency set is attempted to be added.
error AlreadyDependency();

/// @notice Error when the chain's chain ID is attempted to be removed from the interop dependency set.
error CantRemovedDependency();
