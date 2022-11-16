// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./OCR2DRRegistryInterface.sol";

/**
 * @title OCR2DR billable oracle interface.
 */
interface OCR2DRBillableInterface {
  /**
   * @notice Determine the fee charged by the DON that will be paid to Node Operators for servicing a request
   * @param data Encoded OCR2DR request data, use OCR2DRClient API to encode a request
   * @param billing The request's billing configuration
   * @return fee Cost in Juels (1e18) of LINK
   */
  function getRequiredFee(bytes calldata data, OCR2DRRegistryInterface.RequestBilling calldata billing)
    external
    view
    returns (uint96);

  /**
   * @notice Estimate the total cost that will be charged to a subscription to make a request: gas re-imbursement, plus DON fee, plus Registry fee
   * @param data Encoded OCR2DR request data, use OCR2DRClient API to encode a request
   * @param billing The request's billing configuration
   * @return billedCost Cost in Juels (1e18) of LINK
   */
  function estimateCost(bytes calldata data, OCR2DRRegistryInterface.RequestBilling calldata billing)
    external
    view
    returns (uint96);
}
