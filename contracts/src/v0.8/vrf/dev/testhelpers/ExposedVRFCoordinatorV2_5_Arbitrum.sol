// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {VRFCoordinatorV2_5_Arbitrum} from "../VRFCoordinatorV2_5_Arbitrum.sol";

contract ExposedVRFCoordinatorV2_5_Arbitrum is VRFCoordinatorV2_5_Arbitrum {
  constructor(address blockhashStore) VRFCoordinatorV2_5_Arbitrum(blockhashStore) {}

  function getBlockNumberExternal() external view returns (uint256) {
    return _getBlockNumber();
  }

  function getBlockhashExternal(uint64 blockNumber) external view returns (bytes32) {
    return _getBlockhash(blockNumber);
  }

  function calculatePaymentAmountNativeExternal(
    uint256 startGas,
    uint256 weiPerUnitGas,
    bool onlyPremium
  ) external returns (uint96) {
    return _calculatePaymentAmountNative(startGas, weiPerUnitGas, onlyPremium);
  }

  function calculatePaymentAmountLinkExternal(
    uint256 startGas,
    uint256 weiPerUnitGas,
    bool onlyPremium
  ) external returns (uint96, bool) {
    return _calculatePaymentAmountLink(startGas, weiPerUnitGas, onlyPremium);
  }
}
