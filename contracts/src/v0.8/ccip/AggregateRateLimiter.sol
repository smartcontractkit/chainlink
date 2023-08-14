// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IPriceRegistry} from "./interfaces/IPriceRegistry.sol";

import {OwnerIsCreator} from "./../shared/access/OwnerIsCreator.sol";
import {Client} from "./libraries/Client.sol";
import {RateLimiter} from "./libraries/RateLimiter.sol";
import {USDPriceWith18Decimals} from "./libraries/USDPriceWith18Decimals.sol";

contract AggregateRateLimiter is OwnerIsCreator {
  using RateLimiter for RateLimiter.TokenBucket;
  using USDPriceWith18Decimals for uint192;

  error PriceNotFoundForToken(address token);
  event AdminSet(address newAdmin);

  // The address of the token limit admin that has the same permissions as the owner.
  address internal s_admin;

  // The token bucket object that contains the bucket state.
  RateLimiter.TokenBucket private s_rateLimiter;

  /// @param config The RateLimiter.Config containing the capacity and refill rate
  /// of the bucket, plus the admin address.
  constructor(RateLimiter.Config memory config) {
    s_rateLimiter = RateLimiter.TokenBucket({
      rate: config.rate,
      capacity: config.capacity,
      tokens: config.capacity,
      lastUpdated: uint32(block.timestamp),
      isEnabled: config.isEnabled
    });
  }

  /// @notice Consumes value from the rate limiter bucket based on the
  /// token value given. First, calculate the prices
  function _rateLimitValue(Client.EVMTokenAmount[] memory tokenAmounts, IPriceRegistry priceRegistry) internal {
    uint256 numberOfTokens = tokenAmounts.length;

    uint256 value = 0;
    for (uint256 i = 0; i < numberOfTokens; ++i) {
      // not fetching validated price, as price staleness is not important for value-based rate limiting
      // we only need to verify price is not 0
      uint192 pricePerToken = priceRegistry.getTokenPrice(tokenAmounts[i].token).value;
      if (pricePerToken == 0) revert PriceNotFoundForToken(tokenAmounts[i].token);
      value += pricePerToken._calcUSDValueFromTokenAmount(tokenAmounts[i].amount);
    }

    s_rateLimiter._consume(value, address(0));
  }

  /// @notice Gets the token bucket with its values for the block it was requested at.
  /// @return The token bucket.
  function currentRateLimiterState() external view returns (RateLimiter.TokenBucket memory) {
    return s_rateLimiter._currentTokenBucketState();
  }

  /// @notice Sets the rate limited config.
  /// @param config The new rate limiter config.
  /// @dev should only be callable by the owner or token limit admin.
  function setRateLimiterConfig(RateLimiter.Config memory config) external onlyAdminOrOwner {
    s_rateLimiter._setTokenBucketConfig(config);
  }

  // ================================================================
  // |                           Access                             |
  // ================================================================

  /// @notice Gets the token limit admin address.
  /// @return the token limit admin address.
  function getTokenLimitAdmin() external view returns (address) {
    return s_admin;
  }

  /// @notice Sets the token limit admin address.
  /// @param newAdmin the address of the new admin.
  /// @dev setting this to address(0) indicates there is no active admin.
  function setAdmin(address newAdmin) external onlyAdminOrOwner {
    s_admin = newAdmin;
    emit AdminSet(newAdmin);
  }

  /// @notice a modifier that allows the owner or the s_tokenLimitAdmin call the functions
  /// it is applied to.
  modifier onlyAdminOrOwner() {
    if (msg.sender != owner() && msg.sender != s_admin) revert RateLimiter.OnlyCallableByAdminOrOwner();
    _;
  }
}
