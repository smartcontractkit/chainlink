// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title Chainlink Functions billing subscription registry interface.
 */
interface IFunctionsBilling {
  struct RequestBilling {
    // a unique subscription ID allocated by billing system,
    uint64 subscriptionId;
    // the client contract that initiated the request to the DON
    // to use the subscription it must be added as a consumer on the subscription
    address client;
    // customer specified gas limit for the fulfillment callback
    uint32 callbackGasLimit;
    // the expected gas price used to execute the transaction
    uint256 expectedGasPrice;
  }

  /**
   * @notice Gets the configuration of the Chainlink Functions billing registry
   * @return maxGasLimit global max for request gas limit
   * @return stalenessSeconds if the eth/link feed is more stale then this, use the fallback price
   * @return gasAfterPaymentCalculation gas used in doing accounting after completing the gas measurement
   * @return fallbackWeiPerUnitLink fallback eth/link price in the case of a stale feed
   * @return gasOverhead average gas execution cost used in estimating total cost
   * @return linkPriceFeed address of contract for a conversion price between LINK token and native token
   * @return maxSupportedRequestDataVersion The highest support request data version supported by the node
   * @return fulfillmentGasPriceOverEstimationBP Percentage of gas price overestimation to account for changes in gas price between request and response. Held as basis points (one hundredth of 1 percentage point)
   */
  function getConfig()
    external
    view
    returns (
      uint32 maxGasLimit,
      uint32 stalenessSeconds,
      uint256 gasAfterPaymentCalculation,
      int256 fallbackWeiPerUnitLink,
      uint32 gasOverhead,
      address linkPriceFeed,
      uint16 maxSupportedRequestDataVersion,
      uint256 fulfillmentGasPriceOverEstimationBP
    );

  /**
   * @notice Determine the fee that will be split between Node Operators for servicing a request
   * @param requestData Encoded Chainlink Functions request data, use FunctionsClient API to encode a request
   * @param billing The request's billing configuration
   * @return fee Cost in Juels (1e18) of LINK
   */
  function getDONFee(bytes memory requestData, RequestBilling memory billing) external view returns (uint80);

  /**
   * @notice Determine the fee that will be paid to the Router owner for operating the network
   * @param requestData Encoded Chainlink Functions request data, use FunctionsClient API to encode a request
   * @param billing The request's billing configuration
   * @return fee Cost in Juels (1e18) of LINK
   */
  function getAdminFee(bytes memory requestData, RequestBilling memory billing) external view returns (uint96);

  /**
   * @notice Estimate the total cost that will be charged to a subscription to make a request: gas re-reimbursement, plus DON fee, plus Registry fee
   * @param subscriptionId An identifier of the billing account
   * @param data Encoded Chainlink Functions request data, use FunctionsClient API to encode a request
   * @param callbackGasLimit Gas limit for the fulfillment callback
   * @param gasPrice The blockchain's gas price to estimate with
   * @return billedCost Cost in Juels (1e18) of LINK
   */
  function estimateCost(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 callbackGasLimit,
    uint256 gasPrice
  ) external view returns (uint96);

  /**
   * @notice Remove a request commitment that the Router has determined to be stale
   * @dev Only callable by the Router
   * @param requestId - The request ID to remove
   */
  function deleteCommitment(bytes32 requestId) external returns (bool);
}
