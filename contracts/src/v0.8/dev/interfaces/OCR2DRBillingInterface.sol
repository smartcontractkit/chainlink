// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

/**
 * @title OCR2DR billing interface.
 */
interface OCR2DRBillingInterface {
  struct BillingConfig {
    uint32 gasOverhead;
    mapping(address => address) feeTokenPriceOracles; /* feeToken => oracle */
    mapping(address => uint32) requiredFeeByToken; /* feeToken => fee */
  }

  struct RequestBilling {
    uint32 totalFee; // required fee + execution fee
    address feeToken; // ERC20 compatible token that the user wants to pay in
    uint32 gasLimit; // customer specified gas limit for the fulfillment callback
  }

  /**
   * @notice Determine the fee that will be paid to Node Operators of the DON for servicing a request
   * @param data Encoded OCR2DR request data, use OCR2DRClient API to encode a request
   * @param billing The request's billing configuration
   */
  function getRequiredFee(bytes calldata data, RequestBilling calldata billing) external returns (uint32);

  /**
   * @notice Determine the execution cost that will be reimbursed to the Node Operator who transmits the data on-chain
   * @param billing The request's billing configuration
   */
  function estimateExecutionFee(RequestBilling calldata billing) external returns (uint256);
}
