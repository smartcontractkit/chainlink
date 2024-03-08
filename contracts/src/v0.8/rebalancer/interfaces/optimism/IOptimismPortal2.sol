// SPDX-License-Identifier: MIT
// Copied from https://github.com/ethereum-optimism/optimism/blob/v1.7.0/packages/contracts-bedrock/src/L1/OptimismPortal2.sol
pragma solidity ^0.8.0;

import {GameType} from "./DisputeTypes.sol";

interface IOptimismPortal2 {
  /// @notice The game type that the OptimismPortal consults for output proposals.
  function respectedGameType() external view returns (GameType);
}
