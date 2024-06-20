// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IMessageInterceptor} from "./interfaces/IMessageInterceptor.sol";
import {IPriceRegistry} from "./interfaces/IPriceRegistry.sol";

import {OwnerIsCreator} from "./../shared/access/OwnerIsCreator.sol";
import {EnumerableMapAddresses} from "./../shared/enumerable/EnumerableMapAddresses.sol";
import {Client} from "./libraries/Client.sol";
import {RateLimiter} from "./libraries/RateLimiter.sol";
import {USDPriceWith18Decimals} from "./libraries/USDPriceWith18Decimals.sol";

import {EnumerableSet} from "./../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";

/// @notice The aggregate rate limiter is a wrapper of the token bucket rate limiter
/// which permits rate limiting based on the aggregate value of a group of
/// token transfers, using a price registry to convert to a numeraire asset (e.g. USD).
/// The contract is a standalone multi-lane message validator contract, which can be called by authorized
/// ramp contracts to apply rate limit changes to lanes, and revert when the rate limits get breached.
contract MultiAggregateRateLimiter is IMessageInterceptor, OwnerIsCreator {
  using RateLimiter for RateLimiter.TokenBucket;
  using USDPriceWith18Decimals for uint224;
  using EnumerableMapAddresses for EnumerableMapAddresses.AddressToBytes32Map;
  using EnumerableSet for EnumerableSet.AddressSet;

  error UnauthorizedCaller(address caller);
  error PriceNotFoundForToken(address token);
  error ZeroAddressNotAllowed();
  error ZeroChainSelectorNotAllowed();

  event RateLimiterConfigUpdated(uint64 indexed remoteChainSelector, bool isOutgoingLane, RateLimiter.Config config);
  event PriceRegistrySet(address newPriceRegistry);
  event TokenAggregateRateLimitAdded(uint64 remoteChainSelector, bytes32 remoteToken, address localToken);
  event TokenAggregateRateLimitRemoved(uint64 remoteChainSelector, address localToken);
  event AuthorizedCallerAdded(address caller);
  event AuthorizedCallerRemoved(address caller);

  /// @notice RemoteRateLimitToken struct containing the local token address with the chain selector
  /// The struct is used for removals and updates, since the local -> remote token mappings are scoped per-chain
  struct LocalRateLimitToken {
    uint64 remoteChainSelector; // ────╮ Remote chain selector for which to update the rate limit token mapping
    address localToken; // ────────────╯ Token on the chain on which the multi-ARL is deployed
  }

  /// @notice RateLimitToken struct containing both the local and remote token addresses
  struct RateLimitTokenArgs {
    LocalRateLimitToken localTokenArgs; // Local token update args scoped to one remote chain
    bytes32 remoteToken; // Token on the remote chain (for OnRamp - dest, of OffRamp - source)
  }

  /// @notice Update args for changing the authorized callers
  struct AuthorizedCallerArgs {
    address[] addedCallers;
    address[] removedCallers;
  }

  /// @notice Update args for a single rate limiter config update
  struct RateLimiterConfigArgs {
    uint64 remoteChainSelector; // ────╮ Chain selector to set config for
    bool isOutgoingLane; // ───────────╯ If set to true, represents the outgoing message lane (OnRamp), and the incoming message lane otherwise (OffRamp)
    RateLimiter.Config rateLimiterConfig; // Rate limiter config to set
  }

  /// @notice Struct to store rate limit token buckets for both lane directions
  struct RateLimiterBuckets {
    RateLimiter.TokenBucket incomingLaneBucket; // Bucket for the incoming lane (remote -> local)
    RateLimiter.TokenBucket outgoingLaneBucket; // Bucket for the outgoing lane (local -> remote)
  }

  /// @dev Tokens that should be included in Aggregate Rate Limiting (from local chain (this chain) -> remote),
  /// grouped per-remote chain.
  mapping(uint64 remoteChainSelector => EnumerableMapAddresses.AddressToBytes32Map tokensLocalToRemote) internal
    s_rateLimitedTokensLocalToRemote;

  /// @dev Set of callers that can call the validation functions (this is required since the validations modify state)
  EnumerableSet.AddressSet internal s_authorizedCallers;

  /// @notice The address of the PriceRegistry used to query token values for ratelimiting
  address internal s_priceRegistry;

  /// @notice Rate limiter token bucket states per chain, with separate buckets for incoming and outgoing lanes.
  mapping(uint64 remoteChainSelector => RateLimiterBuckets buckets) internal s_rateLimitersByChainSelector;

  /// @param priceRegistry the price registry to set
  /// @param authorizedCallers the authorized callers to set
  constructor(address priceRegistry, address[] memory authorizedCallers) {
    _setPriceRegistry(priceRegistry);
    _applyAuthorizedCallerUpdates(
      AuthorizedCallerArgs({addedCallers: authorizedCallers, removedCallers: new address[](0)})
    );
  }

  /// @inheritdoc IMessageInterceptor
  function onIncomingMessage(Client.Any2EVMMessage memory message) external {
    if (!s_authorizedCallers.contains(msg.sender)) {
      revert UnauthorizedCaller(msg.sender);
    }

    uint64 remoteChainSelector = message.sourceChainSelector;
    RateLimiter.TokenBucket storage tokenBucket = _getTokenBucket(remoteChainSelector, false);

    // Skip rate limiting if it is disabled
    if (tokenBucket.isEnabled) {
      uint256 value;
      Client.EVMTokenAmount[] memory destTokenAmounts = message.destTokenAmounts;
      for (uint256 i = 0; i < destTokenAmounts.length; ++i) {
        if (s_rateLimitedTokensLocalToRemote[remoteChainSelector].contains(destTokenAmounts[i].token)) {
          value += _getTokenValue(destTokenAmounts[i]);
        }
      }

      if (value > 0) tokenBucket._consume(value, address(0));
    }
  }

  /// @inheritdoc IMessageInterceptor
  function onOutgoingMessage(Client.EVM2AnyMessage memory message, uint64 destChainSelector) external {
    // TODO: to be implemented (assuming the same rate limiter states are shared for incoming and outgoing messages)
  }

  /// @param remoteChainSelector chain selector to retrieve token bucket for
  /// @param isOutgoingLane if set to true, fetches the bucket for the outgoing message lane (OnRamp).
  /// Otherwise fetches for the incoming message lane (OffRamp).
  /// @return bucket Storage pointer to the token bucket representing a specific lane
  function _getTokenBucket(
    uint64 remoteChainSelector,
    bool isOutgoingLane
  ) internal view returns (RateLimiter.TokenBucket storage) {
    RateLimiterBuckets storage rateLimiterBuckets = s_rateLimitersByChainSelector[remoteChainSelector];
    if (isOutgoingLane) {
      return rateLimiterBuckets.outgoingLaneBucket;
    } else {
      return rateLimiterBuckets.incomingLaneBucket;
    }
  }

  /// @notice Retrieves the token value for a token using the PriceRegistry
  /// @return tokenValue USD value in 18 decimals
  function _getTokenValue(Client.EVMTokenAmount memory tokenAmount) internal view returns (uint256) {
    // not fetching validated price, as price staleness is not important for value-based rate limiting
    // we only need to verify the price is not 0
    uint224 pricePerToken = IPriceRegistry(s_priceRegistry).getTokenPrice(tokenAmount.token).value;
    if (pricePerToken == 0) revert PriceNotFoundForToken(tokenAmount.token);
    return pricePerToken._calcUSDValueFromTokenAmount(tokenAmount.amount);
  }

  /// @notice Gets the token bucket with its values for the block it was requested at.
  /// @param remoteChainSelector chain selector to retrieve state for
  /// @param isOutgoingLane if set to true, fetches the rate limit state for the outgoing message lane (OnRamp).
  /// Otherwise fetches for the incoming message lane (OffRamp).
  /// The outgoing and incoming message rate limit state is completely separated.
  /// @return The token bucket.
  function currentRateLimiterState(
    uint64 remoteChainSelector,
    bool isOutgoingLane
  ) external view returns (RateLimiter.TokenBucket memory) {
    return _getTokenBucket(remoteChainSelector, isOutgoingLane)._currentTokenBucketState();
  }

  /// @notice Applies the provided rate limiter config updates.
  /// @param rateLimiterUpdates Rate limiter updates
  /// @dev should only be callable by the owner
  function applyRateLimiterConfigUpdates(RateLimiterConfigArgs[] memory rateLimiterUpdates) external onlyOwner {
    for (uint256 i = 0; i < rateLimiterUpdates.length; ++i) {
      RateLimiterConfigArgs memory updateArgs = rateLimiterUpdates[i];
      RateLimiter.Config memory configUpdate = updateArgs.rateLimiterConfig;
      uint64 remoteChainSelector = updateArgs.remoteChainSelector;

      if (remoteChainSelector == 0) {
        revert ZeroChainSelectorNotAllowed();
      }

      bool isOutgoingLane = updateArgs.isOutgoingLane;

      RateLimiter.TokenBucket storage tokenBucket = _getTokenBucket(remoteChainSelector, isOutgoingLane);

      if (tokenBucket.lastUpdated == 0) {
        // Token bucket needs to be newly added
        RateLimiter.TokenBucket memory newTokenBucket = RateLimiter.TokenBucket({
          rate: configUpdate.rate,
          capacity: configUpdate.capacity,
          tokens: configUpdate.capacity,
          lastUpdated: uint32(block.timestamp),
          isEnabled: configUpdate.isEnabled
        });

        if (isOutgoingLane) {
          s_rateLimitersByChainSelector[remoteChainSelector].outgoingLaneBucket = newTokenBucket;
        } else {
          s_rateLimitersByChainSelector[remoteChainSelector].incomingLaneBucket = newTokenBucket;
        }
      } else {
        tokenBucket._setTokenBucketConfig(configUpdate);
      }
      emit RateLimiterConfigUpdated(remoteChainSelector, isOutgoingLane, configUpdate);
    }
  }

  /// @notice Get all tokens which are included in Aggregate Rate Limiting.
  /// @param remoteChainSelector chain selector to get rate limit tokens for
  /// @return localTokens The local chain representation of the tokens that are rate limited.
  /// @return remoteTokens The remote representation of the tokens that are rate limited.
  /// @dev the order of IDs in the list is **not guaranteed**, therefore, if ordering matters when
  /// making successive calls, one should keep the block height constant to ensure a consistent result.
  function getAllRateLimitTokens(uint64 remoteChainSelector)
    external
    view
    returns (address[] memory localTokens, bytes32[] memory remoteTokens)
  {
    uint256 tokenCount = s_rateLimitedTokensLocalToRemote[remoteChainSelector].length();

    localTokens = new address[](tokenCount);
    remoteTokens = new bytes32[](tokenCount);

    for (uint256 i = 0; i < tokenCount; ++i) {
      (address localToken, bytes32 remoteToken) = s_rateLimitedTokensLocalToRemote[remoteChainSelector].at(i);
      localTokens[i] = localToken;
      remoteTokens[i] = remoteToken;
    }
    return (localTokens, remoteTokens);
  }

  /// @notice Adds or removes tokens from being used in Aggregate Rate Limiting.
  /// @param removes - A list of one or more tokens to be removed.
  /// @param adds - A list of one or more tokens to be added.
  function updateRateLimitTokens(
    LocalRateLimitToken[] memory removes,
    RateLimitTokenArgs[] memory adds
  ) external onlyOwner {
    for (uint256 i = 0; i < removes.length; ++i) {
      address localToken = removes[i].localToken;
      uint64 remoteChainSelector = removes[i].remoteChainSelector;

      if (s_rateLimitedTokensLocalToRemote[remoteChainSelector].remove(localToken)) {
        emit TokenAggregateRateLimitRemoved(remoteChainSelector, localToken);
      }
    }

    for (uint256 i = 0; i < adds.length; ++i) {
      LocalRateLimitToken memory localTokenArgs = adds[i].localTokenArgs;
      bytes32 remoteToken = adds[i].remoteToken;
      address localToken = localTokenArgs.localToken;

      if (localToken == address(0) || remoteToken == bytes32("")) {
        revert ZeroAddressNotAllowed();
      }

      uint64 remoteChainSelector = localTokenArgs.remoteChainSelector;

      if (s_rateLimitedTokensLocalToRemote[remoteChainSelector].set(localToken, remoteToken)) {
        emit TokenAggregateRateLimitAdded(remoteChainSelector, remoteToken, localToken);
      }
    }
  }

  /// @return priceRegistry The configured PriceRegistry address
  function getPriceRegistry() external view returns (address) {
    return s_priceRegistry;
  }

  /// @notice Sets the Price Registry address
  /// @param newPriceRegistry the address of the new PriceRegistry
  /// @dev precondition The address must be a non-zero address
  function setPriceRegistry(address newPriceRegistry) external onlyOwner {
    _setPriceRegistry(newPriceRegistry);
  }

  /// @notice Sets the Price Registry address
  /// @param newPriceRegistry the address of the new PriceRegistry
  /// @dev precondition The address must be a non-zero address
  function _setPriceRegistry(address newPriceRegistry) internal {
    if (newPriceRegistry == address(0)) {
      revert ZeroAddressNotAllowed();
    }

    s_priceRegistry = newPriceRegistry;
    emit PriceRegistrySet(newPriceRegistry);
  }

  // ================================================================
  // │                           Access                             │
  // ================================================================

  /// @return authorizedCallers Returns all callers that are authorized to call the validation functions
  function getAllAuthorizedCallers() external view returns (address[] memory) {
    return s_authorizedCallers.values();
  }

  /// @notice Updates the callers that are authorized to call the message validation functions
  /// @param authorizedCallerArgs Callers to add and remove
  function applyAuthorizedCallerUpdates(AuthorizedCallerArgs memory authorizedCallerArgs) external onlyOwner {
    _applyAuthorizedCallerUpdates(authorizedCallerArgs);
  }

  /// @notice Updates the callers that are authorized to call the message validation functions
  /// @param authorizedCallerArgs Callers to add and remove
  function _applyAuthorizedCallerUpdates(AuthorizedCallerArgs memory authorizedCallerArgs) internal {
    address[] memory removedCallers = authorizedCallerArgs.removedCallers;
    for (uint256 i = 0; i < removedCallers.length; ++i) {
      address caller = removedCallers[i];

      if (s_authorizedCallers.remove(caller)) {
        emit AuthorizedCallerRemoved(caller);
      }
    }

    address[] memory addedCallers = authorizedCallerArgs.addedCallers;
    for (uint256 i = 0; i < addedCallers.length; ++i) {
      address caller = addedCallers[i];

      if (caller == address(0)) {
        revert ZeroAddressNotAllowed();
      }

      s_authorizedCallers.add(caller);
      emit AuthorizedCallerAdded(caller);
    }
  }
}
