// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

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
