// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {VRFCoordinatorV2_5} from "./VRFCoordinatorV2_5.sol";
import {OptimismL1Fees} from "./OptimismL1Fees.sol";

/// @dev VRFCoordinatorV2_5_Optimism combines VRFCoordinatorV2_5 base contract with
/// @dev Optimism specific opcodes and L1 gas fee calculations.
/// @dev This coordinator contract is used for all chains in the OP stack (e.g. Base).
contract VRFCoordinatorV2_5_Optimism is VRFCoordinatorV2_5, OptimismL1Fees {
  constructor(address blockhashStore) VRFCoordinatorV2_5(blockhashStore) {}

  /// @notice no need to override getBlockhash and getBlockNumber from VRFCoordinatorV2_5
  /// @notice on OP stack, they will work with the default implementation

  /// @notice Override getL1CostWei function from VRFCoordinatorV2_5 to activate Optimism getL1Fee computation
  function _getL1CostWei(bytes calldata data) internal view override returns (uint256) {
    return _getL1CostWeiForCalldata(data);
  }
}
