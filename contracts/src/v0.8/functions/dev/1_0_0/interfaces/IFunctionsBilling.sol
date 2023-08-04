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
    uint256 expectedGasPriceGwei;
    // Flat fee (in Juels of LINK) that will be paid to the Router owner for operation of the network
    uint96 adminFee;
  }

  struct Config {
    // Maximum amount of gas that can be given to a request's client callback
    uint32 maxCallbackGasLimit;
    // How long before we consider the feed price to be stale
    // and fallback to fallbackNativePerUnitLink.
    uint32 feedStalenessSeconds;
    // Represents the average gas execution cost before the fulfillment callback.
    // This amount is always billed for every request
    uint32 gasOverheadBeforeCallback;
    // Represents the average gas execution cost after the fulfillment callback.
    // This amount is always billed for every request
    uint32 gasOverheadAfterCallback;
    // How many seconds it takes before we consider a request to be timed out
    uint32 requestTimeoutSeconds;
    // Additional flat fee (in Juels of LINK) that will be split between Node Operators
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

  // @notice Gets the Chainlink Coordinator's billing configuration
  // @return config
  function getConfig() external view returns (Config memory);

  // @notice Sets the Chainlink Coordinator's billing configuration
  // @param config - See the contents of the Config struct in IFunctionsBilling.Config for more information
  function updateConfig(Config memory config) external;

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

  // @notice Estimate the total cost that will be charged to a subscription to make a request: transmitter gas re-reimbursement, plus DON fee, plus Registry fee
  // @param subscriptionId An identifier of the billing account
  // @param data Encoded Chainlink Functions request data, use FunctionsClient API to encode a request
  // @param callbackGasLimit Gas limit for the fulfillment callback
  // @param gasPriceGwei The blockchain's gas price to estimate with
  // @return billedCost Cost in Juels (1e18) of LINK
  function estimateCost(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 callbackGasLimit,
    uint256 gasPriceGwei
  ) external view returns (uint96);

  // @notice Remove a request commitment that the Router has determined to be stale
  // @param requestId - The request ID to remove
  function deleteCommitment(bytes32 requestId) external returns (bool);

  // @notice Oracle withdraw LINK earned through fulfilling requests
  // @notice If amount is 0 the full balance will be withdrawn
  // @param recipient where to send the funds
  // @param amount amount to withdraw
  function oracleWithdraw(address recipient, uint96 amount) external;

  // @notice Withdraw all LINK earned by Oracles through fulfilling requests
  function oracleWithdrawAll() external;
}
