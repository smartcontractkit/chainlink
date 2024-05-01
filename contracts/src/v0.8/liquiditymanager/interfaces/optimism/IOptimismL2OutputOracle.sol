// SPDX-License-Identifier: MIT
// Copied from https://github.com/ethereum-optimism/optimism/blob/v1.7.0/packages/contracts-bedrock/src/L1/L2OutputOracle.sol
pragma solidity ^0.8.0;

import {Types} from "./Types.sol";

interface IOptimismL2OutputOracle {
  /// @notice Returns the index of the L2 output that checkpoints a given L2 block number.
  ///         Uses a binary search to find the first output greater than or equal to the given
  ///         block.
  /// @param _l2BlockNumber L2 block number to find a checkpoint for.
  /// @return Index of the first checkpoint that commits to the given L2 block number.
  function getL2OutputIndexAfter(uint256 _l2BlockNumber) external view returns (uint256);

  /// @notice Returns an output by index. Needed to return a struct instead of a tuple.
  /// @param _l2OutputIndex Index of the output to return.
  /// @return The output at the given index.
  function getL2Output(uint256 _l2OutputIndex) external view returns (Types.OutputProposal memory);
}
