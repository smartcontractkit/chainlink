// SPDX-License-Identifier: MIT
// Copied from https://github.com/ethereum-optimism/optimism/blob/v1.7.0/packages/contracts-bedrock/src/L1/OptimismPortal2.sol
pragma solidity ^0.8.0;
import {GameType} from "./DisputeTypes.sol";

interface IOptimismPortal2 {
  /// @notice The dispute game factory address.
  /// @dev See https://github.com/ethereum-optimism/optimism/blob/f707883038d527cbf1e9f8ea513fe33255deadbc/packages/contracts-bedrock/src/L1/OptimismPortal2.sol#L79.
  function disputeGameFactory() external view returns (address);
  /// @notice The game type that the OptimismPortal consults for output proposals.
  function respectedGameType() external view returns (GameType);
}
