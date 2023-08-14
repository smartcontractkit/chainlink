// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IPriceRegistry} from "./interfaces/IPriceRegistry.sol";

import {OwnerIsCreator} from "./../shared/access/OwnerIsCreator.sol";
import {Internal} from "./libraries/Internal.sol";
import {USDPriceWith18Decimals} from "./libraries/USDPriceWith18Decimals.sol";

import {EnumerableSet} from "../vendor/openzeppelin-solidity/v4.8.0/utils/structs/EnumerableSet.sol";

/// @notice The PriceRegistry contract responsibility is to store the current gas price in USD for a given destination chain,
/// and the price of a token in USD allowing the owner or priceUpdater to update this value.
contract PriceRegistry is IPriceRegistry, OwnerIsCreator {
  using EnumerableSet for EnumerableSet.AddressSet;
  using USDPriceWith18Decimals for uint192;

  error TokenNotSupported(address token);
  error ChainNotSupported(uint64 chain);
  error OnlyCallableByUpdaterOrOwner();
  error StaleGasPrice(uint64 destChainSelector, uint256 threshold, uint256 timePassed);
  error StaleTokenPrice(address token, uint256 threshold, uint256 timePassed);
  error InvalidStalenessThreshold();

  event PriceUpdaterSet(address indexed priceUpdater);
  event PriceUpdaterRemoved(address indexed priceUpdater);
  event FeeTokenAdded(address indexed feeToken);
  event FeeTokenRemoved(address indexed feeToken);
  event UsdPerUnitGasUpdated(uint64 indexed destChain, uint256 value, uint256 timestamp);
  event UsdPerTokenUpdated(address indexed token, uint256 value, uint256 timestamp);

  /// @dev The price, in USD with 18 decimals, of 1 unit of gas for a given destination chain.
  /// @dev Price of 1e18 is 1 USD. Examples:
  ///     Very Expensive:   1 unit of gas costs 1 USD                  -> 1e18
  ///     Expensive:        1 unit of gas costs 0.1 USD                -> 1e17
  ///     Cheap:            1 unit of gas costs 0.000001 USD           -> 1e12
  mapping(uint64 destChainSelector => Internal.TimestampedUint192Value price)
    private s_usdPerUnitGasByDestChainSelector;

  /// @dev The price, in USD with 18 decimals, per 1e18 of the smallest token denomination.
  /// @dev Price of 1e18 represents 1 USD per 1e18 token amount.
  ///     1 USDC = 1.00 USD per full token, each full token is 1e6 units -> 1 * 1e18 * 1e18 / 1e6 = 1e30
  ///     1 ETH = 2,000 USD per full token, each full token is 1e18 units -> 2000 * 1e18 * 1e18 / 1e18 = 2_000e18
  ///     1 LINK = 5.00 USD per full token, each full token is 1e18 units -> 5 * 1e18 * 1e18 / 1e18 = 5e18
  mapping(address token => Internal.TimestampedUint192Value price) private s_usdPerToken;

  // Price updaters are allowed to update the prices.
  EnumerableSet.AddressSet private s_priceUpdaters;
  // Subset of tokens which prices tracked by this registry which are fee tokens.
  EnumerableSet.AddressSet private s_feeTokens;
  // The amount of time a price can be stale before it is considered invalid.
  uint32 private immutable i_stalenessThreshold;

  constructor(address[] memory priceUpdaters, address[] memory feeTokens, uint32 stalenessThreshold) {
    _applyPriceUpdatersUpdates(priceUpdaters, new address[](0));
    _applyFeeTokensUpdates(feeTokens, new address[](0));
    if (stalenessThreshold == 0) revert InvalidStalenessThreshold();
    i_stalenessThreshold = stalenessThreshold;
  }

  // ================================================================
  // |                     Price calculations                       |
  // ================================================================

  // @inheritdoc IPriceRegistry
  function getTokenPrice(address token) public view override returns (Internal.TimestampedUint192Value memory) {
    return s_usdPerToken[token];
  }

  // @inheritdoc IPriceRegistry
  function getValidatedTokenPrice(address token) external view override returns (uint192) {
    return _getValidatedTokenPrice(token);
  }

  // @inheritdoc IPriceRegistry
  function getTokenPrices(
    address[] calldata tokens
  ) external view override returns (Internal.TimestampedUint192Value[] memory) {
    uint256 length = tokens.length;
    Internal.TimestampedUint192Value[] memory tokenPrices = new Internal.TimestampedUint192Value[](length);
    for (uint256 i = 0; i < length; ++i) {
      tokenPrices[i] = getTokenPrice(tokens[i]);
    }
    return tokenPrices;
  }

  /// @notice Get the staleness threshold.
  /// @return stalenessThreshold The staleness threshold.
  function getStalenessThreshold() external view returns (uint128) {
    return i_stalenessThreshold;
  }

  // @inheritdoc IPriceRegistry
  function getDestinationChainGasPrice(
    uint64 destChainSelector
  ) external view override returns (Internal.TimestampedUint192Value memory) {
    return s_usdPerUnitGasByDestChainSelector[destChainSelector];
  }

  function getTokenAndGasPrices(
    address token,
    uint64 destChainSelector
  ) external view override returns (uint192 tokenPrice, uint192 gasPriceValue) {
    Internal.TimestampedUint192Value memory gasPrice = s_usdPerUnitGasByDestChainSelector[destChainSelector];
    // We do allow a gas price of 0, but no stale or unset gas prices
    if (gasPrice.timestamp == 0) revert ChainNotSupported(destChainSelector);
    uint256 timePassed = block.timestamp - gasPrice.timestamp;
    if (timePassed > i_stalenessThreshold) revert StaleGasPrice(destChainSelector, i_stalenessThreshold, timePassed);

    return (_getValidatedTokenPrice(token), gasPrice.value);
  }

  /// @inheritdoc IPriceRegistry
  /// @dev this function assumes that no more than 1e59 dollars are sent as payment.
  /// If more is sent, the multiplication of feeTokenAmount and feeTokenValue will overflow.
  /// Since there isn't even close to 1e59 dollars in the world economy this is safe.
  function convertTokenAmount(
    address fromToken,
    uint256 fromTokenAmount,
    address toToken
  ) external view override returns (uint256) {
    /// Example:
    /// fromTokenAmount:   1e18      // 1 ETH
    /// ETH:               2_000e18
    /// LINK:              5e18
    /// return:            1e18 * 2_000e18 / 5e18 = 400e18 (400 LINK)
    return (fromTokenAmount * _getValidatedTokenPrice(fromToken)) / _getValidatedTokenPrice(toToken);
  }

  /// @notice Gets the token price for a given token and revert if the token is either
  /// not supported or the price is stale.
  /// @param token The address of the token to get the price for
  /// @return the token price
  function _getValidatedTokenPrice(address token) internal view returns (uint192) {
    Internal.TimestampedUint192Value memory tokenPrice = s_usdPerToken[token];
    if (tokenPrice.timestamp == 0 || tokenPrice.value == 0) revert TokenNotSupported(token);
    uint256 timePassed = block.timestamp - tokenPrice.timestamp;
    if (timePassed > i_stalenessThreshold) revert StaleTokenPrice(token, i_stalenessThreshold, timePassed);
    return tokenPrice.value;
  }

  // ================================================================
  // |                         Fee tokens                           |
  // ================================================================

  /// @notice Get the list of fee tokens.
  /// @return The tokens set as fee tokens.
  function getFeeTokens() external view returns (address[] memory) {
    return s_feeTokens.values();
  }

  /// @notice Add and remove tokens from feeTokens set.
  /// @param feeTokensToAdd The addresses of the tokens which are now considered fee tokens
  /// and can be used to calculate fees.
  /// @param feeTokensToRemove The addresses of the tokens which are no longer considered feeTokens.
  function applyFeeTokensUpdates(
    address[] memory feeTokensToAdd,
    address[] memory feeTokensToRemove
  ) external onlyOwner {
    _applyFeeTokensUpdates(feeTokensToAdd, feeTokensToRemove);
  }

  /// @notice Add and remove tokens from feeTokens set.
  /// @param feeTokensToAdd The addresses of the tokens which are now considered fee tokens
  /// and can be used to calculate fees.
  /// @param feeTokensToRemove The addresses of the tokens which are no longer considered feeTokens.
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
  // |                       Price updates                          |
  // ================================================================

  // @inheritdoc IPriceRegistry
  function updatePrices(Internal.PriceUpdates calldata priceUpdates) external override requireUpdaterOrOwner {
    uint256 priceUpdatesLength = priceUpdates.tokenPriceUpdates.length;

    for (uint256 i = 0; i < priceUpdatesLength; ++i) {
      Internal.TokenPriceUpdate memory update = priceUpdates.tokenPriceUpdates[i];
      s_usdPerToken[update.sourceToken] = Internal.TimestampedUint192Value({
        value: update.usdPerToken,
        timestamp: uint64(block.timestamp)
      });
      emit UsdPerTokenUpdated(update.sourceToken, update.usdPerToken, block.timestamp);
    }

    if (priceUpdates.destChainSelector != 0) {
      s_usdPerUnitGasByDestChainSelector[priceUpdates.destChainSelector] = Internal.TimestampedUint192Value({
        value: priceUpdates.usdPerUnitGas,
        timestamp: uint64(block.timestamp)
      });
      emit UsdPerUnitGasUpdated(priceUpdates.destChainSelector, priceUpdates.usdPerUnitGas, block.timestamp);
    }
  }

  // ================================================================
  // |                           Access                             |
  // ================================================================

  /// @notice Get the list of price updaters.
  /// @return The price updaters.
  function getPriceUpdaters() external view returns (address[] memory) {
    return s_priceUpdaters.values();
  }

  /// @notice Adds new priceUpdaters and remove existing ones.
  /// @param priceUpdatersToAdd The addresses of the priceUpdaters that are now allowed
  /// to send fee updates.
  /// @param priceUpdatersToRemove The addresses of the priceUpdaters that are no longer allowed
  /// to send fee updates.
  function applyPriceUpdatersUpdates(
    address[] memory priceUpdatersToAdd,
    address[] memory priceUpdatersToRemove
  ) external onlyOwner {
    _applyPriceUpdatersUpdates(priceUpdatersToAdd, priceUpdatersToRemove);
  }

  /// @notice Adds new priceUpdaters and remove existing ones.
  /// @param priceUpdatersToAdd The addresses of the priceUpdaters that are now allowed
  /// to send fee updates.
  /// @param priceUpdatersToRemove The addresses of the priceUpdaters that are no longer allowed
  /// to send fee updates.
  function _applyPriceUpdatersUpdates(
    address[] memory priceUpdatersToAdd,
    address[] memory priceUpdatersToRemove
  ) private {
    for (uint256 i = 0; i < priceUpdatersToAdd.length; ++i) {
      if (s_priceUpdaters.add(priceUpdatersToAdd[i])) {
        emit PriceUpdaterSet(priceUpdatersToAdd[i]);
      }
    }
    for (uint256 i = 0; i < priceUpdatersToRemove.length; ++i) {
      if (s_priceUpdaters.remove(priceUpdatersToRemove[i])) {
        emit PriceUpdaterRemoved(priceUpdatersToRemove[i]);
      }
    }
  }

  /// @notice Require that the caller is the owner or a fee updater.
  modifier requireUpdaterOrOwner() {
    if (msg.sender != owner() && !s_priceUpdaters.contains(msg.sender)) revert OnlyCallableByUpdaterOrOwner();
    _;
  }
}
