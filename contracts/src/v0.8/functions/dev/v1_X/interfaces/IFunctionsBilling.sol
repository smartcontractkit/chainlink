// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @title Chainlink Functions DON billing interface.
interface IFunctionsBilling {
  /// @notice Return the current conversion from WEI of ETH to LINK from the configured Chainlink data feed
  /// @return weiPerUnitLink - The amount of WEI in one LINK
  function getWeiPerUnitLink() external view returns (uint256);

  /// @notice Return the current conversion from LINK to USD from the configured Chainlink data feed
  /// @return weiPerUnitLink - The amount of USD that one LINK is worth
  /// @return decimals - The number of decimals that should be represented in the price feed's response
  function getUsdPerUnitLink() external view returns (uint256, uint8);

  /// @notice Determine the fee that will be split between Node Operators for servicing a request
  /// @param requestCBOR - CBOR encoded Chainlink Functions request data, use FunctionsRequest library to encode a request
  /// @return fee - Cost in Juels (1e18) of LINK
  function getDONFeeJuels(bytes memory requestCBOR) external view returns (uint72);

  /// @notice Determine the fee that will be paid to the Coordinator owner for operating the network
  /// @return fee - Cost in Juels (1e18) of LINK
  function getOperationFeeJuels() external view returns (uint72);

  /// @notice Determine the fee that will be paid to the Router owner for operating the network
  /// @return fee - Cost in Juels (1e18) of LINK
  function getAdminFeeJuels() external view returns (uint72);

  /// @notice Estimate the total cost that will be charged to a subscription to make a request: transmitter gas re-reimbursement, plus DON fee, plus Registry fee
  /// @param - subscriptionId An identifier of the billing account
  /// @param - data Encoded Chainlink Functions request data, use FunctionsClient API to encode a request
  /// @param - callbackGasLimit Gas limit for the fulfillment callback
  /// @param - gasPriceWei The blockchain's gas price to estimate with
  /// @return - billedCost Cost in Juels (1e18) of LINK
  function estimateCost(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 callbackGasLimit,
    uint256 gasPriceWei
  ) external view returns (uint96);

  /// @notice Remove a request commitment that the Router has determined to be stale
  /// @param requestId - The request ID to remove
  function deleteCommitment(bytes32 requestId) external;

  /// @notice Oracle withdraw LINK earned through fulfilling requests
  /// @notice If amount is 0 the full balance will be withdrawn
  /// @param recipient where to send the funds
  /// @param amount amount to withdraw
  function oracleWithdraw(address recipient, uint96 amount) external;

  /// @notice Withdraw all LINK earned by Oracles through fulfilling requests
  /// @dev transmitter addresses must support LINK tokens to avoid tokens from getting stuck as oracleWithdrawAll() calls will forward tokens directly to transmitters
  function oracleWithdrawAll() external;
}

// ================================================================
// |                     Configuration state                      |
// ================================================================

struct FunctionsBillingConfig {
  uint32 fulfillmentGasPriceOverEstimationBP; // ══╗ Percentage of gas price overestimation to account for changes in gas price between request and response. Held as basis points (one hundredth of 1 percentage point)
  uint32 feedStalenessSeconds; //                  ║ How long before we consider the feed price to be stale and fallback to fallbackNativePerUnitLink.
  uint32 gasOverheadBeforeCallback; //             ║ Represents the average gas execution cost before the fulfillment callback. This amount is always billed for every request.
  uint32 gasOverheadAfterCallback; //              ║ Represents the average gas execution cost after the fulfillment callback. This amount is always billed for every request.
  uint40 minimumEstimateGasPriceWei; //            ║ The lowest amount of wei that will be used as the tx.gasprice when estimating the cost to fulfill the request
  uint16 maxSupportedRequestDataVersion; //        ║ The highest support request data version supported by the node. All lower versions should also be supported.
  uint64 fallbackUsdPerUnitLink; //                ║ Fallback LINK / USD conversion rate if the data feed is stale
  uint8 fallbackUsdPerUnitLinkDecimals; // ════════╝ Fallback LINK / USD conversion rate decimal places if the data feed is stale
  uint224 fallbackNativePerUnitLink; // ═══════════╗ Fallback NATIVE CURRENCY / LINK conversion rate if the data feed is stale
  uint32 requestTimeoutSeconds; // ════════════════╝ How many seconds it takes before we consider a request to be timed out
  uint16 donFeeCentsUsd; // ═══════════════════════════════╗ Additional flat fee (denominated in cents of USD, paid as LINK) that will be split between Node Operators.
  uint16 operationFeeCentsUsd; // ═════════════════════════╝ Additional flat fee (denominated in cents of USD, paid as LINK) that will be paid to the owner of the Coordinator contract.
}
