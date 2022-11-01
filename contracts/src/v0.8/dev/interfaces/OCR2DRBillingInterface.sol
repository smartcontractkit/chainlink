// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

/**
 * @title OCR2DR billing interface.
 */
interface OCR2DRBillingInterface {
  struct BillingConfig {
    // There is some variance of gas costs for execution, so use a
    // fixed value representing the average cost so that NOPs are fairly compensated
    uint32 gasOverhead;
    // Price Oracles to look up currency conversion rates to allow paying in any ERC20 token
    mapping(address => address) feeTokenPriceOracles; /* feeToken => oracle */
    // An additional fee charged to use a fee token
    mapping(address => uint32) requiredFeeByToken; /* feeToken => fee */
  }

  struct RequestBilling {
    // required fee + execution fee
    uint32 totalFee;
    // ERC20 compatible token that the user wants to pay in
    address feeToken;
    // customer specified gas limit for the fulfillment callback
    uint32 gasLimit;
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
