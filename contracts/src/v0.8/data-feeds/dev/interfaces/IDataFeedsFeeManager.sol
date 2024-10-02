// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @notice IDataFeedsFeeManager
/// Responsible for all finance operations associated with the Data Feeds
/// product. Includes methods to set fee values (denominated in USD) and
/// discounts per token or user. Finance admins may be permissioned to
/// perform operations.
interface IDataFeedsFeeManager {
  enum Service {
    GetBenchmarks,
    GetReports,
    RequestUpkeep
  }

  /// @notice Process the fee for a set of requests of feeds on a service ID.
  /// @param sender Original msg.sender
  /// @param service Service requested by msg.sender
  /// @param feedIds List of feed IDs
  /// @param billingData Encoded data for additional flexibility
  function processFee(
    address sender,
    Service service,
    bytes32[] calldata feedIds,
    bytes calldata billingData
  ) external payable;

  /// @notice Fetch fee for a set of requests for feeds on a service ID.
  /// @param sender Original msg.sender
  /// @param service Service requested by msg.sender
  /// @param feedIds List of feed IDs
  /// @param billingData Encoded data for additional flexibility
  /// @return fee
  function getFee(
    address sender,
    Service service,
    bytes32[] calldata feedIds,
    bytes calldata billingData
  ) external view returns (uint256 fee);

  /// @notice Allow users to pull fees from the calling address.
  /// @param spendingAddresses User addresses making reqests to Data Feeds
  function enableSpendingAddresses(address[] calldata spendingAddresses) external;

  /// @notice Remove user permissions to pull fees from the calling address.
  /// @param spendingAddresses User addresses making reqests to Data Feeds
  function disableSpendingAddresses(address[] calldata spendingAddresses) external;

  /// @notice Add finance admin for functionality to call finance administrative
  /// functions like withdrawals and surcharges/discounts.
  /// @param financeAdmins Admin addresses
  function addFinanceAdmins(address[] calldata financeAdmins) external;

  /// @notice Remove finance admin.
  /// @param financeAdmins Admin addresses
  function removeFinanceAdmins(address[] calldata financeAdmins) external;

  /// @notice Add permitted ERC-20 token for payment. Feed ID represents the
  /// feed that gives the USD-quoted price for the token to convert the USD fee
  /// to the quantity of the given token.
  /// @param tokenAddresses Fee token address
  /// @param priceFeedIds Feed ID
  /// @param feeTokenDiscountConfigId Fee Token Discount Config ID (see below)
  function addFeeTokens(
    address[] calldata tokenAddresses,
    bytes32[] calldata priceFeedIds,
    bytes32 feeTokenDiscountConfigId
  ) external;

  /// @notice Removes permissions to pay fees in the given ERC-20 token.
  /// @param tokenAddresses Fee token address
  function removeFeeTokens(address[] calldata tokenAddresses) external;

  /// @notice Set the USD Fee amount to charge for a given service. May only
  /// be called by finance admins. All values use 10^18 = $1.00 multiplier.
  /// @param configIds IDs of the configs
  /// @param getBenchmarkUsdFees USD-denominated fee to request a single benchmark
  /// @param getReportUsdFees USD-denominated fee to request a single report
  /// @param requestUpkeepUsdFees USD-denominated fee to request an upkeep
  function setServiceFeeConfigs(
    bytes32[] calldata configIds,
    uint256[] calldata getBenchmarkUsdFees,
    uint256[] calldata getReportUsdFees,
    uint256[] calldata requestUpkeepUsdFees
  ) external;

  /// @notice Set the ServiceFeeConfig used for the set of feeds
  /// @param configId ID of the config
  /// @param feedIds IDs of the feeds to be associated with the service
  function setFeedServiceFees(bytes32 configId, bytes32[] calldata feedIds) external;

  /// @notice Set discount configurations for fee tokens
  /// @param configIds list of IDs for configs
  /// @param discounts discount percentage to set (10^18 = 100%)
  function setFeeTokenDiscountConfigs(bytes32[] calldata configIds, uint256[] calldata discounts) external;

  /// @notice Set the fee discount for a list of tokens. May only be called by
  /// finance admins.
  /// @param configId FeeTokenDiscount config ID
  /// @param tokenAddresses List of fee token addresses to apply discount to
  function setFeeTokenDiscounts(bytes32 configId, address[] calldata tokenAddresses) external;

  /// @notice Set discount configurations for senders
  /// @param configIds list of IDs for configs
  /// @param discounts discount percentage to set (10^18 = 100%)
  function setSenderDiscountConfigs(bytes32[] calldata configIds, uint256[] calldata discounts) external;

  /// @notice Set the fee discount for a list of senders. May only be called by
  /// finance admins.
  /// @param configId FeeTokenDiscount config ID
  /// @param senders List of sender addresses to apply discount to
  function setSenderDiscounts(bytes32 configId, address[] calldata senders) external;

  /// @notice Withdraw ERC-20 balance to the recipient address. May only
  /// be called by finance admins.
  /// @param tokenAddresses Fee token addresses
  /// @param quantities Amount of each token to withdraw
  /// @param recipientAddress Balance address to send the tokens to
  function withdraw(
    address[] calldata tokenAddresses,
    uint256[] calldata quantities,
    address recipientAddress
  ) external;
}
