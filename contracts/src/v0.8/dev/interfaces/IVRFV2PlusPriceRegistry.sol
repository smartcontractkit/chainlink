// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface IVRFV2PlusPriceRegistry {
  /**
   * @notice Calculate the payment amount of a VRF request
   * @param startGas starting gas, retreived at the beginning of the fulfillment transaction using gasleft()
   * @param weiPerUnitGas the gas price at which to calculate the payment, typically tx.gasprice
   * @param nativePayment whether the payment is to be calculated in link or in native currency (i.e ether)
  */
  function calculatePaymentAmount(
    uint256 startGas,
    uint256 weiPerUnitGas,
    bool nativePayment
  ) external view returns (uint96);
}
