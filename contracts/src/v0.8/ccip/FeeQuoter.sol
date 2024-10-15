// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../shared/interfaces/ITypeAndVersion.sol";
import {IFeeQuoter} from "./interfaces/IFeeQuoter.sol";
import {IPriceRegistry} from "./interfaces/IPriceRegistry.sol";

import {AuthorizedCallers} from "../shared/access/AuthorizedCallers.sol";
import {AggregatorV3Interface} from "./../shared/interfaces/AggregatorV3Interface.sol";
import {Client} from "./libraries/Client.sol";
import {Internal} from "./libraries/Internal.sol";
import {Pool} from "./libraries/Pool.sol";
import {USDPriceWith18Decimals} from "./libraries/USDPriceWith18Decimals.sol";

import {KeystoneFeedsPermissionHandler} from "../keystone/KeystoneFeedsPermissionHandler.sol";
import {IReceiver} from "../keystone/interfaces/IReceiver.sol";
import {KeystoneFeedDefaultMetadataLib} from "../keystone/lib/KeystoneFeedDefaultMetadataLib.sol";
import {EnumerableSet} from "../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";

/// @notice The FeeQuoter contract responsibility is to:
///   - Store the current gas price in USD for a given destination chain,
///   - Store the price of a token in USD allowing the owner or priceUpdater to update this value.
///   - Manage chain specific fee calculations.
/// The authorized callers in the contract represent the fee price updaters.
contract FeeQuoter is AuthorizedCallers, IFeeQuoter, ITypeAndVersion, IReceiver, KeystoneFeedsPermissionHandler {
  using EnumerableSet for EnumerableSet.AddressSet;
  using USDPriceWith18Decimals for uint224;
  using KeystoneFeedDefaultMetadataLib for bytes;

  error TokenNotSupported(address token);
  error FeeTokenNotSupported(address token);
  error ChainNotSupported(uint64 chain);
  error StaleGasPrice(uint64 destChainSelector, uint256 threshold, uint256 timePassed);
  error StaleKeystoneUpdate(address token, uint256 feedTimestamp, uint256 storedTimeStamp);
  error DataFeedValueOutOfUint224Range();
  error InvalidDestBytesOverhead(address token, uint32 destBytesOverhead);
  error MessageGasLimitTooHigh();
  error DestinationChainNotEnabled(uint64 destChainSelector);
  error ExtraArgOutOfOrderExecutionMustBeTrue();
  error InvalidExtraArgsTag();
  error SourceTokenDataTooLarge(address token);
  error InvalidDestChainConfig(uint64 destChainSelector);
  error MessageFeeTooHigh(uint256 msgFeeJuels, uint256 maxFeeJuelsPerMsg);
  error InvalidStaticConfig();
  error MessageTooLarge(uint256 maxSize, uint256 actualSize);
  error UnsupportedNumberOfTokens();

  event FeeTokenAdded(address indexed feeToken);
  event FeeTokenRemoved(address indexed feeToken);
  event UsdPerUnitGasUpdated(uint64 indexed destChain, uint256 value, uint256 timestamp);
  event UsdPerTokenUpdated(address indexed token, uint256 value, uint256 timestamp);
  event PriceFeedPerTokenUpdated(address indexed token, TokenPriceFeedConfig priceFeedConfig);
  event TokenTransferFeeConfigUpdated(
    uint64 indexed destChainSelector, address indexed token, TokenTransferFeeConfig tokenTransferFeeConfig
  );
  event TokenTransferFeeConfigDeleted(uint64 indexed destChainSelector, address indexed token);
  event PremiumMultiplierWeiPerEthUpdated(address indexed token, uint64 premiumMultiplierWeiPerEth);
  event DestChainConfigUpdated(uint64 indexed destChainSelector, DestChainConfig destChainConfig);
  event DestChainAdded(uint64 indexed destChainSelector, DestChainConfig destChainConfig);

  /// @dev Token price data feed configuration
  struct TokenPriceFeedConfig {
    address dataFeedAddress; // ──╮ AggregatorV3Interface contract (0 - feed is unset)
    uint8 tokenDecimals; // ──────╯ Decimals of the token that the feed represents
  }

  /// @dev Token price data feed update
  struct TokenPriceFeedUpdate {
    address sourceToken; // Source token to update feed for
    TokenPriceFeedConfig feedConfig; // Feed config update data
  }

  /// @dev Struct that contains the static configuration
  /// RMN depends on this struct, if changing, please notify the RMN maintainers.
  // solhint-disable-next-line gas-struct-packing
  struct StaticConfig {
    uint96 maxFeeJuelsPerMsg; // ─╮ Maximum fee that can be charged for a message
    address linkToken; // ────────╯ LINK token address
    uint32 stalenessThreshold; // The amount of time a gas price can be stale before it is considered invalid.
  }

  /// @dev The struct representing the received CCIP feed report from keystone IReceiver.onReport()
  struct ReceivedCCIPFeedReport {
    address token; // Token address
    uint224 price; // ─────────╮ Price of the token in USD with 18 decimals
    uint32 timestamp; // ──────╯ Timestamp of the price update
  }

  /// @dev Struct to hold the fee & validation configs for a destination chain
  struct DestChainConfig {
    bool isEnabled; // ──────────────────────────╮ Whether this destination chain is enabled
    uint16 maxNumberOfTokensPerMsg; //           │ Maximum number of distinct ERC20 token transferred per message
    uint32 maxDataBytes; //                      │ Maximum payload data size in bytes
    uint32 maxPerMsgGasLimit; //                 │ Maximum gas limit for messages targeting EVMs
    uint32 destGasOverhead; //                   │ Gas charged on top of the gasLimit to cover destination chain costs
    uint16 destGasPerPayloadByte; //             │ Destination chain gas charged for passing each byte of `data` payload to receiver
    uint32 destDataAvailabilityOverheadGas; //   │ Extra data availability gas charged on top of the message, e.g. for OCR
    uint16 destGasPerDataAvailabilityByte; //    │ Amount of gas to charge per byte of message data that needs availability
    uint16 destDataAvailabilityMultiplierBps; // │ Multiplier for data availability gas, multiples of bps, or 0.0001
    // The following three properties are defaults, they can be overridden by setting the TokenTransferFeeConfig for a token
    uint16 defaultTokenFeeUSDCents; //           │ Default token fee charged per token transfer
    uint32 defaultTokenDestGasOverhead; // ──────╯ Default gas charged to execute the token transfer on the destination chain
    uint32 defaultTxGasLimit; //─────────────────╮ Default gas limit for a tx
    uint64 gasMultiplierWeiPerEth; //            │ Multiplier for gas costs, 1e18 based so 11e17 = 10% extra cost.
    uint32 networkFeeUSDCents; //                │ Flat network fee to charge for messages, multiples of 0.01 USD
    bool enforceOutOfOrder; //                   │ Whether to enforce the allowOutOfOrderExecution extraArg value to be true.
    bytes4 chainFamilySelector; // ──────────────╯ Selector that identifies the destination chain's family. Used to determine the correct validations to perform for the dest chain.
  }

  /// @dev Struct to hold the configs and its destination chain selector
  /// Same as DestChainConfig but with the destChainSelector so that an array of these
  /// can be passed in the constructor and the applyDestChainConfigUpdates function
  //solhint-disable gas-struct-packing
  struct DestChainConfigArgs {
    uint64 destChainSelector; // Destination chain selector
    DestChainConfig destChainConfig; // Config to update for the chain selector
  }

  /// @dev Struct to hold the transfer fee configuration for token transfers
  struct TokenTransferFeeConfig {
    uint32 minFeeUSDCents; // ──────────╮ Minimum fee to charge per token transfer, multiples of 0.01 USD
    uint32 maxFeeUSDCents; //           │ Maximum fee to charge per token transfer, multiples of 0.01 USD
    uint16 deciBps; //                  │ Basis points charged on token transfers, multiples of 0.1bps, or 1e-5
    uint32 destGasOverhead; //          │ Gas charged to execute the token transfer on the destination chain
    //                                  │ Extra data availability bytes that are returned from the source pool and sent
    uint32 destBytesOverhead; //        │ to the destination pool. Must be >= Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES
    bool isEnabled; // ─────────────────╯ Whether this token has custom transfer fees
  }

  /// @dev Struct to hold the token transfer fee configurations for a token, same as TokenTransferFeeConfig but with the token address included so
  /// that an array of these can be passed in the TokenTransferFeeConfigArgs struct to set the mapping
  struct TokenTransferFeeConfigSingleTokenArgs {
    address token; // Token address
    TokenTransferFeeConfig tokenTransferFeeConfig; // Struct to hold the transfer fee configuration for token transfers
  }

  /// @dev Struct to hold the token transfer fee configurations for a destination chain and a set of tokens. Same as TokenTransferFeeConfigSingleTokenArgs
  /// but with the destChainSelector and an array of TokenTransferFeeConfigSingleTokenArgs included so that an array of these can be passed in the constructor
  /// and the applyTokenTransferFeeConfigUpdates function
  struct TokenTransferFeeConfigArgs {
    uint64 destChainSelector; // Destination chain selector
    TokenTransferFeeConfigSingleTokenArgs[] tokenTransferFeeConfigs; // Array of token transfer fee configurations
  }

  /// @dev Struct to hold a pair of destination chain selector and token address so that an array of these can be passed in the
  /// applyTokenTransferFeeConfigUpdates function to remove the token transfer fee configuration for a token
  struct TokenTransferFeeConfigRemoveArgs {
    uint64 destChainSelector; // ─╮ Destination chain selector
    address token; // ────────────╯ Token address
  }

  /// @dev Struct to hold the fee token configuration for a token, same as the s_premiumMultiplierWeiPerEth but with
  /// the token address included so that an array of these can be passed in the constructor and
  /// applyPremiumMultiplierWeiPerEthUpdates to set the mapping
  struct PremiumMultiplierWeiPerEthArgs {
    address token; // // ───────────────────╮ Token address
    uint64 premiumMultiplierWeiPerEth; // ──╯ Multiplier for destination chain specific premiums.
  }

  /// @dev The base decimals for cost calculations
  uint256 public constant FEE_BASE_DECIMALS = 36;
  /// @dev The decimals that Keystone reports prices in
  uint256 public constant KEYSTONE_PRICE_DECIMALS = 18;

  string public constant override typeAndVersion = "FeeQuoter 1.6.0-dev";

  /// @dev The gas price per unit of gas for a given destination chain, in USD with 18 decimals.
  /// Multiple gas prices can be encoded into the same value. Each price takes {Internal.GAS_PRICE_BITS} bits.
  /// For example, if Optimism is the destination chain, gas price can include L1 base fee and L2 gas price.
  /// Logic to parse the price components is chain-specific, and should live in OnRamp.
  /// @dev Price of 1e18 is 1 USD. Examples:
  ///     Very Expensive:   1 unit of gas costs 1 USD                  -> 1e18
  ///     Expensive:        1 unit of gas costs 0.1 USD                -> 1e17
  ///     Cheap:            1 unit of gas costs 0.000001 USD           -> 1e12
  mapping(uint64 destChainSelector => Internal.TimestampedPackedUint224 price) private
    s_usdPerUnitGasByDestChainSelector;

  /// @dev The price, in USD with 18 decimals, per 1e18 of the smallest token denomination.
  /// @dev Price of 1e18 represents 1 USD per 1e18 token amount.
  ///     1 USDC = 1.00 USD per full token, each full token is 1e6 units -> 1 * 1e18 * 1e18 / 1e6 = 1e30
  ///     1 ETH = 2,000 USD per full token, each full token is 1e18 units -> 2000 * 1e18 * 1e18 / 1e18 = 2_000e18
  ///     1 LINK = 5.00 USD per full token, each full token is 1e18 units -> 5 * 1e18 * 1e18 / 1e18 = 5e18
  mapping(address token => Internal.TimestampedPackedUint224 price) private s_usdPerToken;

  /// @dev Stores the price data feed configurations per token.
  mapping(address token => TokenPriceFeedConfig dataFeedAddress) private s_usdPriceFeedsPerToken;

  /// @dev The multiplier for destination chain specific premiums that can be set by the owner or fee admin
  mapping(address token => uint64 premiumMultiplierWeiPerEth) private s_premiumMultiplierWeiPerEth;

  /// @dev The destination chain specific fee configs
  mapping(uint64 destChainSelector => DestChainConfig destChainConfig) internal s_destChainConfigs;

  /// @dev The token transfer fee config that can be set by the owner or fee admin
  mapping(uint64 destChainSelector => mapping(address token => TokenTransferFeeConfig tranferFeeConfig)) private
    s_tokenTransferFeeConfig;

  /// @dev Maximum fee that can be charged for a message. This is a guard to prevent massively overcharging due to misconfiguation.
  uint96 internal immutable i_maxFeeJuelsPerMsg;
  /// @dev The link token address
  address internal immutable i_linkToken;

  /// @dev Subset of tokens which prices tracked by this registry which are fee tokens.
  EnumerableSet.AddressSet private s_feeTokens;
  /// @dev The amount of time a gas price can be stale before it is considered invalid.
  uint32 private immutable i_stalenessThreshold;

  constructor(
    StaticConfig memory staticConfig,
    address[] memory priceUpdaters,
    address[] memory feeTokens,
    TokenPriceFeedUpdate[] memory tokenPriceFeeds,
    TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs,
    PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs,
    DestChainConfigArgs[] memory destChainConfigArgs
  ) AuthorizedCallers(priceUpdaters) {
    if (
      staticConfig.linkToken == address(0) || staticConfig.maxFeeJuelsPerMsg == 0
        || staticConfig.stalenessThreshold == 0
    ) {
      revert InvalidStaticConfig();
    }

    i_linkToken = staticConfig.linkToken;
    i_maxFeeJuelsPerMsg = staticConfig.maxFeeJuelsPerMsg;
    i_stalenessThreshold = staticConfig.stalenessThreshold;

    _applyFeeTokensUpdates(feeTokens, new address[](0));
    _updateTokenPriceFeeds(tokenPriceFeeds);
    _applyDestChainConfigUpdates(destChainConfigArgs);
    _applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);
    _applyTokenTransferFeeConfigUpdates(tokenTransferFeeConfigArgs, new TokenTransferFeeConfigRemoveArgs[](0));
  }

  // ================================================================
  // │                     Price calculations                       │
  // ================================================================

  /// @inheritdoc IPriceRegistry
  function getTokenPrice(
    address token
  ) public view override returns (Internal.TimestampedPackedUint224 memory) {
    Internal.TimestampedPackedUint224 memory tokenPrice = s_usdPerToken[token];

    // If the token price is not stale, return it
    if (block.timestamp - tokenPrice.timestamp < i_stalenessThreshold) {
      return tokenPrice;
    }

    TokenPriceFeedConfig memory priceFeedConfig = s_usdPriceFeedsPerToken[token];

    // If the token price feed is not set, return the stale price
    if (priceFeedConfig.dataFeedAddress == address(0)) {
      return tokenPrice;
    }

    // If the token price feed is set, return the price from the feed
    // The price feed is the fallback because we do not expect it to be the default source due to the gas cost of reading from it
    return _getTokenPriceFromDataFeed(priceFeedConfig);
  }

  /// @notice Get the `tokenPrice` for a given token, checks if the price is valid.
  /// @param token The token to get the price for.
  /// @return tokenPrice The tokenPrice for the given token if it exists and is valid.
  function getValidatedTokenPrice(
    address token
  ) external view returns (uint224) {
    return _getValidatedTokenPrice(token);
  }

  /// @notice Get the `tokenPrice` for an array of tokens.
  /// @param tokens The tokens to get prices for.
  /// @return tokenPrices The tokenPrices for the given tokens.
  function getTokenPrices(
    address[] calldata tokens
  ) external view returns (Internal.TimestampedPackedUint224[] memory) {
    uint256 length = tokens.length;
    Internal.TimestampedPackedUint224[] memory tokenPrices = new Internal.TimestampedPackedUint224[](length);
    for (uint256 i = 0; i < length; ++i) {
      tokenPrices[i] = getTokenPrice(tokens[i]);
    }
    return tokenPrices;
  }

  /// @notice Returns the token price data feed configuration
  /// @param token The token to retrieve the feed config for
  /// @return tokenPriceFeedConfig The token price data feed config (if feed address is 0, the feed config is disabled)
  function getTokenPriceFeedConfig(
    address token
  ) external view returns (TokenPriceFeedConfig memory) {
    return s_usdPriceFeedsPerToken[token];
  }

  /// @notice Get an encoded `gasPrice` for a given destination chain ID.
  /// The 224-bit result encodes necessary gas price components.
  /// On L1 chains like Ethereum or Avax, the only component is the gas price.
  /// On Optimistic Rollups, there are two components - the L2 gas price, and L1 base fee for data availability.
  /// On future chains, there could be more or differing price components.
  /// PriceRegistry does not contain chain-specific logic to parse destination chain price components.
  /// @param destChainSelector The destination chain to get the price for.
  /// @return gasPrice The encoded gasPrice for the given destination chain ID.
  function getDestinationChainGasPrice(
    uint64 destChainSelector
  ) external view returns (Internal.TimestampedPackedUint224 memory) {
    return s_usdPerUnitGasByDestChainSelector[destChainSelector];
  }

  /// @notice Gets the fee token price and the gas price, both denominated in dollars.
  /// @param token The source token to get the price for.
  /// @param destChainSelector The destination chain to get the gas price for.
  /// @return tokenPrice The price of the feeToken in 1e18 dollars per base unit.
  /// @return gasPriceValue The price of gas in 1e18 dollars per base unit.
  function getTokenAndGasPrices(
    address token,
    uint64 destChainSelector
  ) public view returns (uint224 tokenPrice, uint224 gasPriceValue) {
    Internal.TimestampedPackedUint224 memory gasPrice = s_usdPerUnitGasByDestChainSelector[destChainSelector];
    // We do allow a gas price of 0, but no stale or unset gas prices
    if (gasPrice.timestamp == 0) revert ChainNotSupported(destChainSelector);
    uint256 timePassed = block.timestamp - gasPrice.timestamp;
    if (timePassed > i_stalenessThreshold) revert StaleGasPrice(destChainSelector, i_stalenessThreshold, timePassed);

    return (_getValidatedTokenPrice(token), gasPrice.value);
  }

  /// @notice Convert a given token amount to target token amount.
  /// @dev this function assumes that no more than 1e59 dollars are sent as payment.
  /// If more is sent, the multiplication of feeTokenAmount and feeTokenValue will overflow.
  /// Since there isn't even close to 1e59 dollars in the world economy this is safe.
  /// @param fromToken The given token address.
  /// @param fromTokenAmount The given token amount.
  /// @param toToken The target token address.
  /// @return toTokenAmount The target token amount.
  function convertTokenAmount(
    address fromToken,
    uint256 fromTokenAmount,
    address toToken
  ) public view returns (uint256) {
    /// Example:
    /// fromTokenAmount:   1e18      // 1 ETH
    /// ETH:               2_000e18
    /// LINK:              5e18
    /// return:            1e18 * 2_000e18 / 5e18 = 400e18 (400 LINK)
    return (fromTokenAmount * _getValidatedTokenPrice(fromToken)) / _getValidatedTokenPrice(toToken);
  }

  /// @notice Gets the token price for a given token and reverts if the token is not supported
  /// @param token The address of the token to get the price for
  /// @return tokenPriceValue The token price
  function _getValidatedTokenPrice(
    address token
  ) internal view returns (uint224) {
    Internal.TimestampedPackedUint224 memory tokenPrice = getTokenPrice(token);
    // Token price must be set at least once
    if (tokenPrice.timestamp == 0 || tokenPrice.value == 0) revert TokenNotSupported(token);
    return tokenPrice.value;
  }

  /// @notice Gets the token price from a data feed address, rebased to the same units as s_usdPerToken
  /// @param priceFeedConfig token data feed configuration with valid data feed address (used to retrieve price & timestamp)
  /// @return tokenPrice data feed price answer rebased to s_usdPerToken units, with latest block timestamp
  function _getTokenPriceFromDataFeed(
    TokenPriceFeedConfig memory priceFeedConfig
  ) internal view returns (Internal.TimestampedPackedUint224 memory tokenPrice) {
    AggregatorV3Interface dataFeedContract = AggregatorV3Interface(priceFeedConfig.dataFeedAddress);
    (
      /* uint80 roundID */
      ,
      int256 dataFeedAnswer,
      /* uint startedAt */
      ,
      /* uint256 updatedAt */
      ,
      /* uint80 answeredInRound */
    ) = dataFeedContract.latestRoundData();

    if (dataFeedAnswer < 0) {
      revert DataFeedValueOutOfUint224Range();
    }
    uint224 rebasedValue =
      _calculateRebasedValue(dataFeedContract.decimals(), priceFeedConfig.tokenDecimals, uint256(dataFeedAnswer));

    // Data feed staleness is unchecked to decouple the FeeQuoter from data feed delay issues
    return Internal.TimestampedPackedUint224({value: rebasedValue, timestamp: uint32(block.timestamp)});
  }

  // ================================================================
  // │                         Fee tokens                           │
  // ================================================================

  /// @inheritdoc IPriceRegistry
  function getFeeTokens() external view returns (address[] memory) {
    return s_feeTokens.values();
  }

  /// @notice Add and remove tokens from feeTokens set.
  /// @param feeTokensToRemove The addresses of the tokens which are no longer considered feeTokens.
  /// @param feeTokensToAdd The addresses of the tokens which are now considered fee tokens
  /// and can be used to calculate fees.
  function applyFeeTokensUpdates(
    address[] memory feeTokensToAdd,
    address[] memory feeTokensToRemove
  ) external onlyOwner {
    _applyFeeTokensUpdates(feeTokensToAdd, feeTokensToRemove);
  }

  /// @notice Add and remove tokens from feeTokens set.
  /// @param feeTokensToRemove The addresses of the tokens which are no longer considered feeTokens.
  /// @param feeTokensToAdd The addresses of the tokens which are now considered fee tokens
  /// and can be used to calculate fees.
  function _applyFeeTokensUpdates(address[] memory feeTokensToAdd, address[] memory feeTokensToRemove) private {
    for (uint256 i = 0; i < feeTokensToAdd.length; ++i) {
      if (s_feeTokens.add(feeTokensToAdd[i])) {
        emit FeeTokenAdded(feeTokensToAdd[i]);
      }
    }
    for (uint256 i = 0; i < feeTokensToRemove.length; ++i) {
      if (s_feeTokens.remove(feeTokensToRemove[i])) {
        emit FeeTokenRemoved(feeTokensToRemove[i]);
      }
    }
  }

  // ================================================================
  // │                       Price updates                          │
  // ================================================================

  /// @inheritdoc IPriceRegistry
  function updatePrices(
    Internal.PriceUpdates calldata priceUpdates
  ) external override {
    // The caller must be the fee updater
    _validateCaller();

    uint256 tokenUpdatesLength = priceUpdates.tokenPriceUpdates.length;

    for (uint256 i = 0; i < tokenUpdatesLength; ++i) {
      Internal.TokenPriceUpdate memory update = priceUpdates.tokenPriceUpdates[i];
      s_usdPerToken[update.sourceToken] =
        Internal.TimestampedPackedUint224({value: update.usdPerToken, timestamp: uint32(block.timestamp)});
      emit UsdPerTokenUpdated(update.sourceToken, update.usdPerToken, block.timestamp);
    }

    uint256 gasUpdatesLength = priceUpdates.gasPriceUpdates.length;

    for (uint256 i = 0; i < gasUpdatesLength; ++i) {
      Internal.GasPriceUpdate memory update = priceUpdates.gasPriceUpdates[i];
      s_usdPerUnitGasByDestChainSelector[update.destChainSelector] =
        Internal.TimestampedPackedUint224({value: update.usdPerUnitGas, timestamp: uint32(block.timestamp)});
      emit UsdPerUnitGasUpdated(update.destChainSelector, update.usdPerUnitGas, block.timestamp);
    }
  }

  /// @notice Updates the USD token price feeds for given tokens
  /// @param tokenPriceFeedUpdates Token price feed updates to apply
  function updateTokenPriceFeeds(
    TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates
  ) external onlyOwner {
    _updateTokenPriceFeeds(tokenPriceFeedUpdates);
  }

  /// @notice Updates the USD token price feeds for given tokens
  /// @param tokenPriceFeedUpdates Token price feed updates to apply
  function _updateTokenPriceFeeds(
    TokenPriceFeedUpdate[] memory tokenPriceFeedUpdates
  ) private {
    for (uint256 i; i < tokenPriceFeedUpdates.length; ++i) {
      TokenPriceFeedUpdate memory update = tokenPriceFeedUpdates[i];
      address sourceToken = update.sourceToken;
      TokenPriceFeedConfig memory tokenPriceFeedConfig = update.feedConfig;

      s_usdPriceFeedsPerToken[sourceToken] = tokenPriceFeedConfig;
      emit PriceFeedPerTokenUpdated(sourceToken, tokenPriceFeedConfig);
    }
  }

  /// @notice Handles the report containing price feeds and updates the internal price storage
  /// @inheritdoc IReceiver
  /// @dev This function is called to process incoming price feed data.
  /// @param metadata Arbitrary metadata associated with the report (not used in this implementation).
  /// @param report Encoded report containing an array of `ReceivedCCIPFeedReport` structs.
  function onReport(bytes calldata metadata, bytes calldata report) external {
    (bytes10 workflowName, address workflowOwner, bytes2 reportName) = metadata._extractMetadataInfo();

    _validateReportPermission(msg.sender, workflowOwner, workflowName, reportName);

    ReceivedCCIPFeedReport[] memory feeds = abi.decode(report, (ReceivedCCIPFeedReport[]));

    for (uint256 i = 0; i < feeds.length; ++i) {
      uint8 tokenDecimals = s_usdPriceFeedsPerToken[feeds[i].token].tokenDecimals;
      if (tokenDecimals == 0) {
        revert TokenNotSupported(feeds[i].token);
      }
      // Keystone reports prices in USD with 18 decimals, so we passing it as 18 in the _calculateRebasedValue function
      uint224 rebasedValue = _calculateRebasedValue(18, tokenDecimals, feeds[i].price);

      //if stale update then revert
      if (feeds[i].timestamp < s_usdPerToken[feeds[i].token].timestamp) {
        revert StaleKeystoneUpdate(feeds[i].token, feeds[i].timestamp, s_usdPerToken[feeds[i].token].timestamp);
      }

      s_usdPerToken[feeds[i].token] =
        Internal.TimestampedPackedUint224({value: rebasedValue, timestamp: feeds[i].timestamp});
      emit UsdPerTokenUpdated(feeds[i].token, rebasedValue, feeds[i].timestamp);
    }
  }

  // ================================================================
  // │                       Fee quoting                            │
  // ================================================================

  /// @inheritdoc IFeeQuoter
  /// @dev The function should always validate message.extraArgs, message.receiver and family-specific configs
  function getValidatedFee(
    uint64 destChainSelector,
    Client.EVM2AnyMessage calldata message
  ) external view returns (uint256 feeTokenAmount) {
    DestChainConfig memory destChainConfig = s_destChainConfigs[destChainSelector];
    if (!destChainConfig.isEnabled) revert DestinationChainNotEnabled(destChainSelector);
    if (!s_feeTokens.contains(message.feeToken)) revert FeeTokenNotSupported(message.feeToken);

    uint256 numberOfTokens = message.tokenAmounts.length;
    _validateMessage(destChainConfig, message.data.length, numberOfTokens, message.receiver);

    // The below call asserts that feeToken is a supported token
    (uint224 feeTokenPrice, uint224 packedGasPrice) = getTokenAndGasPrices(message.feeToken, destChainSelector);

    // Calculate premiumFee in USD with 18 decimals precision first.
    // If message-only and no token transfers, a flat network fee is charged.
    // If there are token transfers, premiumFee is calculated from token transfer fee.
    // If there are both token transfers and message, premiumFee is only calculated from token transfer fee.
    uint256 premiumFee = 0;
    uint32 tokenTransferGas = 0;
    uint32 tokenTransferBytesOverhead = 0;
    if (numberOfTokens > 0) {
      (premiumFee, tokenTransferGas, tokenTransferBytesOverhead) =
        _getTokenTransferCost(destChainConfig, destChainSelector, message.feeToken, feeTokenPrice, message.tokenAmounts);
    } else {
      // Convert USD cents with 2 decimals to 18 decimals.
      premiumFee = uint256(destChainConfig.networkFeeUSDCents) * 1e16;
    }

    // Calculate data availability cost in USD with 36 decimals. Data availability cost exists on rollups that need to post
    // transaction calldata onto another storage layer, e.g. Eth mainnet, incurring additional storage gas costs.
    uint256 dataAvailabilityCost = 0;

    // Only calculate data availability cost if data availability multiplier is non-zero.
    // The multiplier should be set to 0 if destination chain does not charge data availability cost.
    if (destChainConfig.destDataAvailabilityMultiplierBps > 0) {
      dataAvailabilityCost = _getDataAvailabilityCost(
        destChainConfig,
        // Parse the data availability gas price stored in the higher-order 112 bits of the encoded gas price.
        uint112(packedGasPrice >> Internal.GAS_PRICE_BITS),
        message.data.length,
        numberOfTokens,
        tokenTransferBytesOverhead
      );
    }

    // Calculate execution gas fee on destination chain in USD with 36 decimals.
    // We add the message gas limit, the overhead gas, the gas of passing message data to receiver, and token transfer gas together.
    // We then multiply this gas total with the gas multiplier and gas price, converting it into USD with 36 decimals.
    // uint112(packedGasPrice) = executionGasPrice

    // NOTE: when supporting non-EVM chains, revisit how generic this fee logic can be
    // NOTE: revisit parsing non-EVM args
    uint256 executionCost = uint112(packedGasPrice)
      * (
        destChainConfig.destGasOverhead + (message.data.length * destChainConfig.destGasPerPayloadByte) + tokenTransferGas
          + _parseEVMExtraArgsFromBytes(message.extraArgs, destChainConfig).gasLimit
      ) * destChainConfig.gasMultiplierWeiPerEth;

    // Calculate number of fee tokens to charge.
    // Total USD fee is in 36 decimals, feeTokenPrice is in 18 decimals USD for 1e18 smallest token denominations.
    // Result of the division is the number of smallest token denominations.
    return ((premiumFee * s_premiumMultiplierWeiPerEth[message.feeToken]) + executionCost + dataAvailabilityCost)
      / feeTokenPrice;
  }

  /// @notice Sets the fee configuration for a token.
  /// @param premiumMultiplierWeiPerEthArgs Array of PremiumMultiplierWeiPerEthArgs structs.
  function applyPremiumMultiplierWeiPerEthUpdates(
    PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs
  ) external onlyOwner {
    _applyPremiumMultiplierWeiPerEthUpdates(premiumMultiplierWeiPerEthArgs);
  }

  /// @dev Sets the fee config.
  /// @param premiumMultiplierWeiPerEthArgs The multiplier for destination chain specific premiums.
  function _applyPremiumMultiplierWeiPerEthUpdates(
    PremiumMultiplierWeiPerEthArgs[] memory premiumMultiplierWeiPerEthArgs
  ) internal {
    for (uint256 i = 0; i < premiumMultiplierWeiPerEthArgs.length; ++i) {
      address token = premiumMultiplierWeiPerEthArgs[i].token;
      uint64 premiumMultiplierWeiPerEth = premiumMultiplierWeiPerEthArgs[i].premiumMultiplierWeiPerEth;
      s_premiumMultiplierWeiPerEth[token] = premiumMultiplierWeiPerEth;

      emit PremiumMultiplierWeiPerEthUpdated(token, premiumMultiplierWeiPerEth);
    }
  }

  /// @notice Gets the fee configuration for a token.
  /// @param token The token to get the fee configuration for.
  /// @return premiumMultiplierWeiPerEth The multiplier for destination chain specific premiums.
  function getPremiumMultiplierWeiPerEth(
    address token
  ) external view returns (uint64 premiumMultiplierWeiPerEth) {
    return s_premiumMultiplierWeiPerEth[token];
  }

  /// @notice Returns the token transfer cost parameters.
  /// A basis point fee is calculated from the USD value of each token transfer.
  /// For each individual transfer, this fee is between [minFeeUSD, maxFeeUSD].
  /// Total transfer fee is the sum of each individual token transfer fee.
  /// @dev Assumes that tokenAmounts are validated to be listed tokens elsewhere.
  /// @dev Splitting one token transfer into multiple transfers is discouraged,
  /// as it will result in a transferFee equal or greater than the same amount aggregated/de-duped.
  /// @param destChainConfig the config configured for the destination chain selector.
  /// @param destChainSelector the destination chain selector.
  /// @param feeToken address of the feeToken.
  /// @param feeTokenPrice price of feeToken in USD with 18 decimals.
  /// @param tokenAmounts token transfers in the message.
  /// @return tokenTransferFeeUSDWei total token transfer bps fee in USD with 18 decimals.
  /// @return tokenTransferGas total execution gas of the token transfers.
  /// @return tokenTransferBytesOverhead additional token transfer data passed to destination, e.g. USDC attestation.
  function _getTokenTransferCost(
    DestChainConfig memory destChainConfig,
    uint64 destChainSelector,
    address feeToken,
    uint224 feeTokenPrice,
    Client.EVMTokenAmount[] calldata tokenAmounts
  ) internal view returns (uint256 tokenTransferFeeUSDWei, uint32 tokenTransferGas, uint32 tokenTransferBytesOverhead) {
    uint256 numberOfTokens = tokenAmounts.length;

    for (uint256 i = 0; i < numberOfTokens; ++i) {
      Client.EVMTokenAmount memory tokenAmount = tokenAmounts[i];
      TokenTransferFeeConfig memory transferFeeConfig = s_tokenTransferFeeConfig[destChainSelector][tokenAmount.token];

      // If the token has no specific overrides configured, we use the global defaults.
      if (!transferFeeConfig.isEnabled) {
        tokenTransferFeeUSDWei += uint256(destChainConfig.defaultTokenFeeUSDCents) * 1e16;
        tokenTransferGas += destChainConfig.defaultTokenDestGasOverhead;
        tokenTransferBytesOverhead += Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES;
        continue;
      }

      uint256 bpsFeeUSDWei = 0;
      // Only calculate bps fee if ratio is greater than 0. Ratio of 0 means no bps fee for a token.
      // Useful for when the FeeQuoter cannot return a valid price for the token.
      if (transferFeeConfig.deciBps > 0) {
        uint224 tokenPrice = 0;
        if (tokenAmount.token != feeToken) {
          tokenPrice = _getValidatedTokenPrice(tokenAmount.token);
        } else {
          tokenPrice = feeTokenPrice;
        }

        // Calculate token transfer value, then apply fee ratio
        // ratio represents multiples of 0.1bps, or 1e-5
        bpsFeeUSDWei = (tokenPrice._calcUSDValueFromTokenAmount(tokenAmount.amount) * transferFeeConfig.deciBps) / 1e5;
      }

      tokenTransferGas += transferFeeConfig.destGasOverhead;
      tokenTransferBytesOverhead += transferFeeConfig.destBytesOverhead;

      // Bps fees should be kept within range of [minFeeUSD, maxFeeUSD].
      // Convert USD values with 2 decimals to 18 decimals.
      uint256 minFeeUSDWei = uint256(transferFeeConfig.minFeeUSDCents) * 1e16;
      if (bpsFeeUSDWei < minFeeUSDWei) {
        tokenTransferFeeUSDWei += minFeeUSDWei;
        continue;
      }

      uint256 maxFeeUSDWei = uint256(transferFeeConfig.maxFeeUSDCents) * 1e16;
      if (bpsFeeUSDWei > maxFeeUSDWei) {
        tokenTransferFeeUSDWei += maxFeeUSDWei;
        continue;
      }

      tokenTransferFeeUSDWei += bpsFeeUSDWei;
    }

    return (tokenTransferFeeUSDWei, tokenTransferGas, tokenTransferBytesOverhead);
  }

  /// @notice calculates the rebased value for 1e18 smallest token denomination
  /// @param dataFeedDecimal decimal of the data feed
  /// @param tokenDecimal decimal of the token
  /// @param feedValue value of the data feed
  /// @return rebasedValue rebased value
  function _calculateRebasedValue(
    uint8 dataFeedDecimal,
    uint8 tokenDecimal,
    uint256 feedValue
  ) internal pure returns (uint224 rebasedValue) {
    // Rebase formula for units in smallest token denomination: usdValue * (1e18 * 1e18) / 1eTokenDecimals
    // feedValue * (10 ** (18 - feedDecimals)) * (10 ** (18 - erc20Decimals))
    // feedValue * (10 ** ((18 - feedDecimals) + (18 - erc20Decimals)))
    // feedValue * (10 ** (36 - feedDecimals - erc20Decimals))
    // feedValue * (10 ** (36 - (feedDecimals + erc20Decimals)))
    // feedValue * (10 ** (36 - excessDecimals))
    // If excessDecimals > 36 => flip it to feedValue / (10 ** (excessDecimals - 36))
    uint8 excessDecimals = dataFeedDecimal + tokenDecimal;
    uint256 rebasedVal;

    if (excessDecimals > FEE_BASE_DECIMALS) {
      rebasedVal = feedValue / (10 ** (excessDecimals - FEE_BASE_DECIMALS));
    } else {
      rebasedVal = feedValue * (10 ** (FEE_BASE_DECIMALS - excessDecimals));
    }

    if (rebasedVal > type(uint224).max) {
      revert DataFeedValueOutOfUint224Range();
    }

    return uint224(rebasedVal);
  }

  /// @notice Returns the estimated data availability cost of the message.
  /// @dev To save on gas, we use a single destGasPerDataAvailabilityByte value for both zero and non-zero bytes.
  /// @param destChainConfig the config configured for the destination chain selector.
  /// @param dataAvailabilityGasPrice USD per data availability gas in 18 decimals.
  /// @param messageDataLength length of the data field in the message.
  /// @param numberOfTokens number of distinct token transfers in the message.
  /// @param tokenTransferBytesOverhead additional token transfer data passed to destination, e.g. USDC attestation.
  /// @return dataAvailabilityCostUSD36Decimal total data availability cost in USD with 36 decimals.
  function _getDataAvailabilityCost(
    DestChainConfig memory destChainConfig,
    uint112 dataAvailabilityGasPrice,
    uint256 messageDataLength,
    uint256 numberOfTokens,
    uint32 tokenTransferBytesOverhead
  ) internal pure returns (uint256 dataAvailabilityCostUSD36Decimal) {
    // dataAvailabilityLengthBytes sums up byte lengths of fixed message fields and dynamic message fields.
    // Fixed message fields do account for the offset and length slot of the dynamic fields.
    uint256 dataAvailabilityLengthBytes = Internal.MESSAGE_FIXED_BYTES + messageDataLength
      + (numberOfTokens * Internal.MESSAGE_FIXED_BYTES_PER_TOKEN) + tokenTransferBytesOverhead;

    // destDataAvailabilityOverheadGas is a separate config value for flexibility to be updated independently of message cost.
    // Its value is determined by CCIP lane implementation, e.g. the overhead data posted for OCR.
    uint256 dataAvailabilityGas = (dataAvailabilityLengthBytes * destChainConfig.destGasPerDataAvailabilityByte)
      + destChainConfig.destDataAvailabilityOverheadGas;

    // dataAvailabilityGasPrice is in 18 decimals, destDataAvailabilityMultiplierBps is in 4 decimals
    // We pad 14 decimals to bring the result to 36 decimals, in line with token bps and execution fee.
    return ((dataAvailabilityGas * dataAvailabilityGasPrice) * destChainConfig.destDataAvailabilityMultiplierBps) * 1e14;
  }

  /// @notice Gets the transfer fee config for a given token.
  /// @param destChainSelector The destination chain selector.
  /// @param token The token address.
  /// @return tokenTransferFeeConfig The transfer fee config for the token.
  function getTokenTransferFeeConfig(
    uint64 destChainSelector,
    address token
  ) external view returns (TokenTransferFeeConfig memory tokenTransferFeeConfig) {
    return s_tokenTransferFeeConfig[destChainSelector][token];
  }

  /// @notice Sets the transfer fee config.
  /// @dev only callable by the owner or admin.
  function applyTokenTransferFeeConfigUpdates(
    TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs,
    TokenTransferFeeConfigRemoveArgs[] memory tokensToUseDefaultFeeConfigs
  ) external onlyOwner {
    _applyTokenTransferFeeConfigUpdates(tokenTransferFeeConfigArgs, tokensToUseDefaultFeeConfigs);
  }

  /// @notice internal helper to set the token transfer fee config.
  function _applyTokenTransferFeeConfigUpdates(
    TokenTransferFeeConfigArgs[] memory tokenTransferFeeConfigArgs,
    TokenTransferFeeConfigRemoveArgs[] memory tokensToUseDefaultFeeConfigs
  ) internal {
    for (uint256 i = 0; i < tokenTransferFeeConfigArgs.length; ++i) {
      TokenTransferFeeConfigArgs memory tokenTransferFeeConfigArg = tokenTransferFeeConfigArgs[i];
      uint64 destChainSelector = tokenTransferFeeConfigArg.destChainSelector;

      for (uint256 j = 0; j < tokenTransferFeeConfigArg.tokenTransferFeeConfigs.length; ++j) {
        TokenTransferFeeConfig memory tokenTransferFeeConfig =
          tokenTransferFeeConfigArg.tokenTransferFeeConfigs[j].tokenTransferFeeConfig;
        address token = tokenTransferFeeConfigArg.tokenTransferFeeConfigs[j].token;

        if (tokenTransferFeeConfig.destBytesOverhead < Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) {
          revert InvalidDestBytesOverhead(token, tokenTransferFeeConfig.destBytesOverhead);
        }

        s_tokenTransferFeeConfig[destChainSelector][token] = tokenTransferFeeConfig;

        emit TokenTransferFeeConfigUpdated(destChainSelector, token, tokenTransferFeeConfig);
      }
    }

    // Remove the custom fee configs for the tokens that are in the tokensToUseDefaultFeeConfigs array
    for (uint256 i = 0; i < tokensToUseDefaultFeeConfigs.length; ++i) {
      uint64 destChainSelector = tokensToUseDefaultFeeConfigs[i].destChainSelector;
      address token = tokensToUseDefaultFeeConfigs[i].token;
      delete s_tokenTransferFeeConfig[destChainSelector][token];
      emit TokenTransferFeeConfigDeleted(destChainSelector, token);
    }
  }

  // ================================================================
  // │             Validations & message processing                 │
  // ================================================================

  /// @notice Validates that the destAddress matches the expected format of the family.
  /// @param chainFamilySelector Tag to identify the target family.
  /// @param destAddress Dest address to validate.
  /// @dev precondition - assumes the family tag is correct and validated.
  function _validateDestFamilyAddress(bytes4 chainFamilySelector, bytes memory destAddress) internal pure {
    if (chainFamilySelector == Internal.CHAIN_FAMILY_SELECTOR_EVM) {
      Internal._validateEVMAddress(destAddress);
    }
  }

  /// @dev Convert the extra args bytes into a struct with validations against the dest chain config.
  /// @param extraArgs The extra args bytes.
  /// @param destChainConfig Dest chain config to validate against.
  /// @return evmExtraArgs The EVMExtraArgs struct (latest version).
  function _parseEVMExtraArgsFromBytes(
    bytes calldata extraArgs,
    DestChainConfig memory destChainConfig
  ) internal pure returns (Client.EVMExtraArgsV2 memory) {
    Client.EVMExtraArgsV2 memory evmExtraArgs =
      _parseUnvalidatedEVMExtraArgsFromBytes(extraArgs, destChainConfig.defaultTxGasLimit);

    if (evmExtraArgs.gasLimit > uint256(destChainConfig.maxPerMsgGasLimit)) revert MessageGasLimitTooHigh();
    if (destChainConfig.enforceOutOfOrder && !evmExtraArgs.allowOutOfOrderExecution) {
      revert ExtraArgOutOfOrderExecutionMustBeTrue();
    }

    return evmExtraArgs;
  }

  /// @dev Convert the extra args bytes into a struct.
  /// @param extraArgs The extra args bytes.
  /// @param defaultTxGasLimit default tx gas limit to use in the absence of extra args.
  /// @return EVMExtraArgs the extra args struct (latest version)
  function _parseUnvalidatedEVMExtraArgsFromBytes(
    bytes calldata extraArgs,
    uint64 defaultTxGasLimit
  ) private pure returns (Client.EVMExtraArgsV2 memory) {
    if (extraArgs.length == 0) {
      // If extra args are empty, generate default values
      return Client.EVMExtraArgsV2({gasLimit: defaultTxGasLimit, allowOutOfOrderExecution: false});
    }

    bytes4 extraArgsTag = bytes4(extraArgs);
    bytes memory argsData = extraArgs[4:];

    if (extraArgsTag == Client.EVM_EXTRA_ARGS_V2_TAG) {
      return abi.decode(argsData, (Client.EVMExtraArgsV2));
    } else if (extraArgsTag == Client.EVM_EXTRA_ARGS_V1_TAG) {
      // EVMExtraArgsV1 originally included a second boolean (strict) field which has been deprecated.
      // Clients may still include it but it will be ignored.
      return Client.EVMExtraArgsV2({gasLimit: abi.decode(argsData, (uint256)), allowOutOfOrderExecution: false});
    }

    revert InvalidExtraArgsTag();
  }

  /// @notice Validate the forwarded message to ensure it matches the configuration limits (message length, number of tokens)
  /// and family-specific expectations (address format).
  /// @param destChainConfig The destination chain config.
  /// @param dataLength The length of the data field of the message.
  /// @param numberOfTokens The number of tokens to be sent.
  /// @param receiver Message receiver on the dest chain.
  function _validateMessage(
    DestChainConfig memory destChainConfig,
    uint256 dataLength,
    uint256 numberOfTokens,
    bytes memory receiver
  ) internal pure {
    // Check that payload is formed correctly
    if (dataLength > uint256(destChainConfig.maxDataBytes)) {
      revert MessageTooLarge(uint256(destChainConfig.maxDataBytes), dataLength);
    }
    if (numberOfTokens > uint256(destChainConfig.maxNumberOfTokensPerMsg)) revert UnsupportedNumberOfTokens();
    _validateDestFamilyAddress(destChainConfig.chainFamilySelector, receiver);
  }

  /// @inheritdoc IFeeQuoter
  /// @dev precondition - onRampTokenTransfers and sourceTokenAmounts lengths must be equal
  function processMessageArgs(
    uint64 destChainSelector,
    address feeToken,
    uint256 feeTokenAmount,
    bytes calldata extraArgs,
    Internal.EVM2AnyTokenTransfer[] calldata onRampTokenTransfers,
    Client.EVMTokenAmount[] calldata sourceTokenAmounts
  )
    external
    view
    returns (
      uint256 msgFeeJuels,
      bool isOutOfOrderExecution,
      bytes memory convertedExtraArgs,
      bytes[] memory destExecDataPerToken
    )
  {
    // Convert feeToken to link if not already in link
    if (feeToken == i_linkToken) {
      msgFeeJuels = feeTokenAmount;
    } else {
      msgFeeJuels = convertTokenAmount(feeToken, feeTokenAmount, i_linkToken);
    }

    if (msgFeeJuels > i_maxFeeJuelsPerMsg) revert MessageFeeTooHigh(msgFeeJuels, i_maxFeeJuelsPerMsg);

    uint64 defaultTxGasLimit = s_destChainConfigs[destChainSelector].defaultTxGasLimit;
    // NOTE: when supporting non-EVM chains, revisit this and parse non-EVM args.
    // We can parse unvalidated args since this message is called after getFee (which will already validate the params)
    Client.EVMExtraArgsV2 memory parsedExtraArgs = _parseUnvalidatedEVMExtraArgsFromBytes(extraArgs, defaultTxGasLimit);
    isOutOfOrderExecution = parsedExtraArgs.allowOutOfOrderExecution;
    destExecDataPerToken = _processPoolReturnData(destChainSelector, onRampTokenTransfers, sourceTokenAmounts);

    return (msgFeeJuels, isOutOfOrderExecution, Client._argsToBytes(parsedExtraArgs), destExecDataPerToken);
  }

  /// @notice Validates pool return data
  /// @param destChainSelector Destination chain selector to which the token amounts are sent to
  /// @param onRampTokenTransfers Token amounts with populated pool return data
  /// @param sourceTokenAmounts Token amounts originally sent in a Client.EVM2AnyMessage message
  /// @return destExecDataPerToken Destination chain execution data
  function _processPoolReturnData(
    uint64 destChainSelector,
    Internal.EVM2AnyTokenTransfer[] calldata onRampTokenTransfers,
    Client.EVMTokenAmount[] calldata sourceTokenAmounts
  ) internal view returns (bytes[] memory destExecDataPerToken) {
    bytes4 chainFamilySelector = s_destChainConfigs[destChainSelector].chainFamilySelector;
    destExecDataPerToken = new bytes[](onRampTokenTransfers.length);
    for (uint256 i = 0; i < onRampTokenTransfers.length; ++i) {
      address sourceToken = sourceTokenAmounts[i].token;

      // Since the DON has to pay for the extraData to be included on the destination chain, we cap the length of the
      // extraData. This prevents gas bomb attacks on the NOPs. As destBytesOverhead accounts for both
      // extraData and offchainData, this caps the worst case abuse to the number of bytes reserved for offchainData.
      uint256 destPoolDataLength = onRampTokenTransfers[i].extraData.length;
      if (destPoolDataLength > Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES) {
        if (destPoolDataLength > s_tokenTransferFeeConfig[destChainSelector][sourceToken].destBytesOverhead) {
          revert SourceTokenDataTooLarge(sourceToken);
        }
      }

      _validateDestFamilyAddress(chainFamilySelector, onRampTokenTransfers[i].destTokenAddress);
      FeeQuoter.TokenTransferFeeConfig memory tokenTransferFeeConfig =
        s_tokenTransferFeeConfig[destChainSelector][sourceToken];
      uint32 defaultGasOverhead = s_destChainConfigs[destChainSelector].defaultTokenDestGasOverhead;
      // NOTE: Revisit this when adding new non-EVM chain family selector support
      uint32 destGasAmount =
        tokenTransferFeeConfig.isEnabled ? tokenTransferFeeConfig.destGasOverhead : defaultGasOverhead;

      // The user will be billed either the default or the override, so we send the exact amount that we billed for
      // to the destination chain to be used for the token releaseOrMint and transfer.
      destExecDataPerToken[i] = abi.encode(destGasAmount);
    }
    return destExecDataPerToken;
  }

  // ================================================================
  // │                           Configs                            │
  // ================================================================

  /// @notice Returns the configured config for the dest chain selector.
  /// @param destChainSelector Destination chain selector to fetch config for.
  /// @return destChainConfig Config for the destination chain.
  function getDestChainConfig(
    uint64 destChainSelector
  ) external view returns (DestChainConfig memory) {
    return s_destChainConfigs[destChainSelector];
  }

  /// @notice Updates the destination chain specific config.
  /// @param destChainConfigArgs Array of source chain specific configs.
  function applyDestChainConfigUpdates(
    DestChainConfigArgs[] memory destChainConfigArgs
  ) external onlyOwner {
    _applyDestChainConfigUpdates(destChainConfigArgs);
  }

  /// @notice Internal version of applyDestChainConfigUpdates.
  function _applyDestChainConfigUpdates(
    DestChainConfigArgs[] memory destChainConfigArgs
  ) internal {
    for (uint256 i = 0; i < destChainConfigArgs.length; ++i) {
      DestChainConfigArgs memory destChainConfigArg = destChainConfigArgs[i];
      uint64 destChainSelector = destChainConfigArgs[i].destChainSelector;
      DestChainConfig memory destChainConfig = destChainConfigArg.destChainConfig;

      // NOTE: when supporting non-EVM chains, update chainFamilySelector validations
      if (
        destChainSelector == 0 || destChainConfig.defaultTxGasLimit == 0
          || destChainConfig.chainFamilySelector != Internal.CHAIN_FAMILY_SELECTOR_EVM
          || destChainConfig.defaultTxGasLimit > destChainConfig.maxPerMsgGasLimit
      ) {
        revert InvalidDestChainConfig(destChainSelector);
      }

      // The chain family selector cannot be zero - indicates that it is a new chain
      if (s_destChainConfigs[destChainSelector].chainFamilySelector == 0) {
        emit DestChainAdded(destChainSelector, destChainConfig);
      } else {
        emit DestChainConfigUpdated(destChainSelector, destChainConfig);
      }

      s_destChainConfigs[destChainSelector] = destChainConfig;
    }
  }

  /// @notice Returns the static FeeQuoter config.
  /// @dev RMN depends on this function, if updated, please notify the RMN maintainers.
  /// @return staticConfig The static configuration.
  function getStaticConfig() external view returns (StaticConfig memory) {
    return StaticConfig({
      maxFeeJuelsPerMsg: i_maxFeeJuelsPerMsg,
      linkToken: i_linkToken,
      stalenessThreshold: i_stalenessThreshold
    });
  }
}
