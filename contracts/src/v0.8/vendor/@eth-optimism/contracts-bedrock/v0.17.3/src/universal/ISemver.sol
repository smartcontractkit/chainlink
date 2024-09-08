// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice: IMPORTANT NOTICE for anyone who wants to use this contract
/// @notice Source: https://github.com/ethereum-optimism/optimism/blob/71b93116738ee98c9f8713b1a5dfe626ce06c1b2/packages/contracts-bedrock/src/universal/ISemver.sol
/// @notice The original code was trimmed down to include only the necessary interface elements required to interact with GasPriceOracle
/// @notice We need this file so that Solidity compiler will not complain because some functions don't exist
/// @notice In reality, we don't embed this code into our own contracts, instead we make cross-contract calls on predeployed GasPriceOracle contract

/// @title ISemver
/// @notice ISemver is a simple contract for ensuring that contracts are
///         versioned using semantic versioning.
interface ISemver {
  /// @notice Getter for the semantic version of the contract. This is not
  ///         meant to be used onchain but instead meant to be used by offchain
  ///         tooling.
  /// @return Semver contract version as a string.
  function version() external view returns (string memory);
}
