// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {OptimismL1Fees} from "./OptimismL1Fees.sol";
import {VRFV2PlusWrapper} from "./VRFV2PlusWrapper.sol";

contract VRFV2PlusWrapper_Optimism is VRFV2PlusWrapper, OptimismL1Fees {
  error UnsupportedL1FeeCalculationMode(uint8 mode);

  constructor(
    address _link,
    address _linkNativeFeed,
    address _coordinator,
    uint256 _subId
  ) VRFV2PlusWrapper(_link, _linkNativeFeed, _coordinator, _subId) {
    // default calculation mode from OptimismL1Fees is not supported on the wrapper
    // switch to the next available one
    s_l1FeeCalculationMode = L1_CALLDATA_GAS_COST_MODE;
  }

  /**
   * @notice Overriding the setL1FeeCalculation function in VRFV2PlusWrapper for Optimism
   * @notice ensures that L1_GAS_FEES_MODE can't be set for the wrapper contract.
   */
  function setL1FeeCalculation(uint8 mode, uint8 coefficient) external override onlyOwner {
    if (mode == L1_GAS_FEES_MODE) {
      revert UnsupportedL1FeeCalculationMode(mode);
    }
    OptimismL1Fees._setL1FeeCalculationInternal(mode, coefficient);
  }

  /**
   * @notice Returns estimated L1 gas fee cost for fulfillment calldata payload once
   * @notice the request has been made through VRFV2PlusWrapper (direct funding model).
   */
  function _getL1CostWei() internal view override returns (uint256) {
    return OptimismL1Fees._getL1CostWeiForCalldataSize(s_fulfillmentTxSizeBytes);
  }
}
