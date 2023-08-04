// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// @title Chainlink Functions billing subscription registry interface.
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
    // Flat fee (in Juels of LINK) that will be paid to the Router owner for operation of the network
    uint96 adminFee;
  }

  struct Config {
    // Maximum amount of gas that can be given to a request's client callback
    uint32 maxCallbackGasLimit;
    // feedStalenessSeconds is how long before we consider the feed price to be stale
    // and fallback to fallbackNativePerUnitLink.
    uint32 feedStalenessSeconds;
    // Represents the average gas execution cost. Used in estimating cost beforehand.
    uint32 gasOverheadBeforeCallback;
    // Gas to cover transmitter oracle payment after we calculate the payment.
    // We make it configurable in case those operations are repriced.
    uint32 gasOverheadAfterCallback;
    // how many seconds it takes before we consider a request to be timed out
    uint32 requestTimeoutSeconds;
    // additional flat fee (in Juels of LINK) that will be split between Node Operators
    // Max value is 2^80 - 1 == 1.2m LINK.
    uint80 donFee;
    // The highest support request data version supported by the node
    // All lower versions should also be supported
    uint16 maxSupportedRequestDataVersion;
    // Percentage of gas price overestimation to account for changes in gas price between request and response
    // Held as basis points (one hundredth of 1 percentage point)
    uint256 fulfillmentGasPriceOverEstimationBP;
    // fallback NATIVE CURRENCY / LINK conversion rate if the data feed is stale
    int256 fallbackNativePerUnitLink;
  }

  // @notice Gets the configuration of the Chainlink Functions billing registry
  // @return config
  function getConfig() external view returns (Config memory);

  // @notice Return the current conversion from WEI of ETH to LINK from the configured Chainlink data feed
  // @return weiPerUnitLink - The amount of WEI in one LINK
  function getWeiPerUnitLink() external view returns (uint256);

  // @notice Determine the fee that will be split between Node Operators for servicing a request
  // @param requestData Encoded Chainlink Functions request data, use FunctionsClient API to encode a request
  // @param billing The request's billing configuration
  // @return fee Cost in Juels (1e18) of LINK
  function getDONFee(bytes memory requestData, RequestBilling memory billing) external view returns (uint80);

  // @notice Determine the fee that will be paid to the Router owner for operating the network
  // @return fee Cost in Juels (1e18) of LINK
  function getAdminFee() external view returns (uint96);

  // @notice Estimate the total cost that will be charged to a subscription to make a request: gas re-reimbursement, plus DON fee, plus Registry fee
  // @param subscriptionId An identifier of the billing account
  // @param data Encoded Chainlink Functions request data, use FunctionsClient API to encode a request
  // @param callbackGasLimit Gas limit for the fulfillment callback
  // @param gasPrice The blockchain's gas price to estimate with
  // @return billedCost Cost in Juels (1e18) of LINK
  function estimateCost(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 callbackGasLimit,
    uint256 gasPrice
  ) external view returns (uint96);

  // @notice Remove a request commitment that the Router has determined to be stale
  // @dev Only callable by the Router
  // @param requestId - The request ID to remove
  function deleteCommitment(bytes32 requestId) external returns (bool);

  // @notice Oracle withdraw LINK earned through fulfilling requests
  // @notice If amount is 0 the full balance will be withdrawn
  // @notice Both signing and transmitting wallets will have a balance to withdraw
  // @param recipient where to send the funds
  // @param amount amount to withdraw
  function oracleWithdraw(address recipient, uint96 amount) external;

  // @notice Oracle withdraw all LINK earned through fulfilling requests to all Node Operators
  function oracleWithdrawAll() external;
}
