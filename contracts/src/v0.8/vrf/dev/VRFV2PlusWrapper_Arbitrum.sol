// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ArbitrumL1Fees} from "./ArbitrumL1Fees.sol";
import {VRFV2PlusWrapper} from "./VRFV2PlusWrapper.sol";

contract VRFV2PlusWrapper_Arbitrum is VRFV2PlusWrapper, ArbitrumL1Fees {
  constructor(
    address _link,
    address _linkNativeFeed,
    address _coordinator,
    uint256 _subId
  ) VRFV2PlusWrapper(_link, _linkNativeFeed, _coordinator, _subId) {}

  /**
   * @notice Returns estimated L1 gas fee cost for fulfillment calldata payload once
   * @notice the request has been made through VRFV2PlusWrapper (direct funding model).
   */
  function _getL1CostWei() internal view override returns (uint256) {
    return ArbitrumL1Fees._getL1CostWeiForCalldataSize(s_fulfillmentTxSizeBytes);
  }
}
