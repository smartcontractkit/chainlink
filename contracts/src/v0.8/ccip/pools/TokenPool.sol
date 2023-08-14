// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IPool} from "../interfaces/pools/IPool.sol";
import {IARM} from "../interfaces/IARM.sol";

import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";
import {RateLimiter} from "../libraries/RateLimiter.sol";

import {Pausable} from "../../vendor/Pausable.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.0/token/ERC20/IERC20.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.0/utils/introspection/IERC165.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.8.0/utils/structs/EnumerableSet.sol";

/// @notice Base abstract class with common functions for all token pools.
abstract contract TokenPool is IPool, OwnerIsCreator, IERC165 {
  using EnumerableSet for EnumerableSet.AddressSet;
  using RateLimiter for RateLimiter.TokenBucket;

  error PermissionsError();
  error ZeroAddressNotAllowed();
  error SenderNotAllowed(address sender);
  error AllowListNotEnabled();
  error NonExistentRamp(address ramp);
  error BadARMSignal();
  error RampAlreadyExists(address ramp);

  event Locked(address indexed sender, uint256 amount);
  event Burned(address indexed sender, uint256 amount);
  event Released(address indexed sender, address indexed recipient, uint256 amount);
  event Minted(address indexed sender, address indexed recipient, uint256 amount);
  event OnRampAdded(address onRamp, RateLimiter.Config rateLimiterConfig);
  event OnRampConfigured(address onRamp, RateLimiter.Config rateLimiterConfig);
  event OnRampRemoved(address onRamp);
  event OffRampAdded(address offRamp, RateLimiter.Config rateLimiterConfig);
  event OffRampConfigured(address offRamp, RateLimiter.Config rateLimiterConfig);
  event OffRampRemoved(address offRamp);
  event AllowListAdd(address sender);
  event AllowListRemove(address sender);

  struct RampUpdate {
    address ramp;
    bool allowed;
    RateLimiter.Config rateLimiterConfig;
  }

  // The immutable token that belongs to this pool.
  IERC20 internal immutable i_token;
  address internal immutable i_armProxy;
  // The immutable flag that indicates if the pool is access-controlled.
  bool internal immutable i_allowlistEnabled;
  // A set of addresses allowed to trigger lockOrBurn as original senders.
  EnumerableSet.AddressSet internal s_allowList;

  // A set of allowed onRamps. We want the whitelist to be enumerable to
  // be able to quickly determine (without parsing logs) who can access the pool.
  EnumerableSet.AddressSet internal s_onRamps;
  mapping(address => RateLimiter.TokenBucket) internal s_onRampRateLimits;
  // A set of allowed offRamps.
  EnumerableSet.AddressSet internal s_offRamps;
  mapping(address => RateLimiter.TokenBucket) internal s_offRampRateLimits;

  constructor(IERC20 token, address[] memory allowlist, address armProxy) {
    if (address(token) == address(0)) revert ZeroAddressNotAllowed();
    i_token = token;
    i_armProxy = armProxy;

    // Pool can be set as permissioned or permissionless at deployment time only.
    i_allowlistEnabled = allowlist.length > 0;
    if (i_allowlistEnabled) {
      _applyAllowListUpdates(new address[](0), allowlist);
    }
  }

  /// @notice Get ARM proxy address
  /// @return armProxy Address of arm proxy
  function getArmProxy() public view returns (address armProxy) {
    return i_armProxy;
  }

  /// @inheritdoc IPool
  function getToken() public view override returns (IERC20 token) {
    return i_token;
  }

  /// @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) public pure virtual override returns (bool) {
    return interfaceId == type(IPool).interfaceId || interfaceId == type(IERC165).interfaceId;
  }

  // ================================================================
  // |                      Ramp permissions                        |
  // ================================================================

  /// @notice Checks whether something is a permissioned onRamp on this contract.
  /// @return true if the given address is a permissioned onRamp.
  function isOnRamp(address onRamp) public view returns (bool) {
    return s_onRamps.contains(onRamp);
  }

  /// @notice Checks whether something is a permissioned offRamp on this contract.
  /// @return true if the given address is a permissioned offRamp.
  function isOffRamp(address offRamp) public view returns (bool) {
    return s_offRamps.contains(offRamp);
  }

  /// @notice Get onRamp whitelist
  /// @return list of onramps.
  function getOnRamps() public view returns (address[] memory) {
    return s_onRamps.values();
  }

  /// @notice Get offRamp whitelist
  /// @return list of offramps
  function getOffRamps() public view returns (address[] memory) {
    return s_offRamps.values();
  }

  /// @notice Sets permissions for all on and offRamps.
  /// @dev Only callable by the owner
  /// @param onRamps A list of onRamps and their new permission status/rate limits
  /// @param offRamps A list of offRamps and their new permission status/rate limits
  function applyRampUpdates(RampUpdate[] calldata onRamps, RampUpdate[] calldata offRamps) external virtual onlyOwner {
    _applyRampUpdates(onRamps, offRamps);
  }

  function _applyRampUpdates(RampUpdate[] calldata onRamps, RampUpdate[] calldata offRamps) internal onlyOwner {
    for (uint256 i = 0; i < onRamps.length; ++i) {
      RampUpdate memory update = onRamps[i];
      if (update.allowed) {
        if (s_onRamps.add(update.ramp)) {
          s_onRampRateLimits[update.ramp] = RateLimiter.TokenBucket({
            rate: update.rateLimiterConfig.rate,
            capacity: update.rateLimiterConfig.capacity,
            tokens: update.rateLimiterConfig.capacity,
            lastUpdated: uint32(block.timestamp),
            isEnabled: update.rateLimiterConfig.isEnabled
          });
          emit OnRampAdded(update.ramp, update.rateLimiterConfig);
        } else {
          revert RampAlreadyExists(update.ramp);
        }
      } else {
        if (s_onRamps.remove(update.ramp)) {
          delete s_onRampRateLimits[update.ramp];
          emit OnRampRemoved(update.ramp);
        } else {
          // Cannot remove a non-existent onRamp.
          revert NonExistentRamp(update.ramp);
        }
      }
    }

    for (uint256 i = 0; i < offRamps.length; ++i) {
      RampUpdate memory update = offRamps[i];
      if (update.allowed) {
        if (s_offRamps.add(update.ramp)) {
          s_offRampRateLimits[update.ramp] = RateLimiter.TokenBucket({
            rate: update.rateLimiterConfig.rate,
            capacity: update.rateLimiterConfig.capacity,
            tokens: update.rateLimiterConfig.capacity,
            lastUpdated: uint32(block.timestamp),
            isEnabled: update.rateLimiterConfig.isEnabled
          });
          emit OffRampAdded(update.ramp, update.rateLimiterConfig);
        } else {
          revert RampAlreadyExists(update.ramp);
        }
      } else {
        if (s_offRamps.remove(update.ramp)) {
          delete s_offRampRateLimits[update.ramp];
          emit OffRampRemoved(update.ramp);
        } else {
          // Cannot remove a non-existent offRamp.
          revert NonExistentRamp(update.ramp);
        }
      }
    }
  }

  // ================================================================
  // |                        Rate limiting                         |
  // ================================================================

  /// @notice Consumes outbound rate limiting capacity in this pool
  function _consumeOnRampRateLimit(uint256 amount) internal {
    s_onRampRateLimits[msg.sender]._consume(amount, address(i_token));
  }

  /// @notice Consumes inbound rate limiting capacity in this pool
  function _consumeOffRampRateLimit(uint256 amount) internal {
    s_offRampRateLimits[msg.sender]._consume(amount, address(i_token));
  }

  /// @notice Gets the token bucket with its values for the block it was requested at.
  /// @return The token bucket.
  function currentOnRampRateLimiterState(address onRamp) external view returns (RateLimiter.TokenBucket memory) {
    return s_onRampRateLimits[onRamp]._currentTokenBucketState();
  }

  /// @notice Gets the token bucket with its values for the block it was requested at.
  /// @return The token bucket.
  function currentOffRampRateLimiterState(address offRamp) external view returns (RateLimiter.TokenBucket memory) {
    return s_offRampRateLimits[offRamp]._currentTokenBucketState();
  }

  /// @notice Sets the onramp rate limited config.
  /// @param config The new rate limiter config.
  function setOnRampRateLimiterConfig(address onRamp, RateLimiter.Config memory config) external onlyOwner {
    if (!isOnRamp(onRamp)) revert NonExistentRamp(onRamp);
    s_onRampRateLimits[onRamp]._setTokenBucketConfig(config);
    emit OnRampConfigured(onRamp, config);
  }

  /// @notice Sets the offramp rate limited config.
  /// @param config The new rate limiter config.
  function setOffRampRateLimiterConfig(address offRamp, RateLimiter.Config memory config) external onlyOwner {
    if (!isOffRamp(offRamp)) revert NonExistentRamp(offRamp);
    s_offRampRateLimits[offRamp]._setTokenBucketConfig(config);
    emit OffRampConfigured(offRamp, config);
  }

  // ================================================================
  // |                           Access                             |
  // ================================================================

  /// @notice Checks whether the msg.sender is a permissioned onRamp on this contract
  /// @dev Reverts with a PermissionsError if check fails
  modifier onlyOnRamp() {
    if (!isOnRamp(msg.sender)) revert PermissionsError();
    _;
  }

  /// @notice Checks whether the msg.sender is a permissioned offRamp on this contract
  /// @dev Reverts with a PermissionsError if check fails
  modifier onlyOffRamp() {
    if (!isOffRamp(msg.sender)) revert PermissionsError();
    _;
  }

  // ================================================================
  // |                          Allowlist                           |
  // ================================================================

  modifier checkAllowList(address sender) {
    if (i_allowlistEnabled && !s_allowList.contains(sender)) revert SenderNotAllowed(sender);
    _;
  }

  /// @notice Gets whether the allowList functionality is enabled.
  /// @return true is enabled, false if not.
  function getAllowListEnabled() external view returns (bool) {
    return i_allowlistEnabled;
  }

  /// @notice Gets the allowed addresses.
  /// @return The allowed addresses.
  function getAllowList() external view returns (address[] memory) {
    return s_allowList.values();
  }

  /// @notice Apply updates to the allow list.
  /// @param removes The addresses to be removed.
  /// @param adds The addresses to be added.
  /// @dev allowListing will be removed before public launch
  function applyAllowListUpdates(address[] calldata removes, address[] calldata adds) external onlyOwner {
    _applyAllowListUpdates(removes, adds);
  }

  /// @notice Internal version of applyAllowListUpdates to allow for reuse in the constructor.
  function _applyAllowListUpdates(address[] memory removes, address[] memory adds) internal {
    if (!i_allowlistEnabled) revert AllowListNotEnabled();

    for (uint256 i = 0; i < removes.length; ++i) {
      address toRemove = removes[i];
      if (s_allowList.remove(toRemove)) {
        emit AllowListRemove(toRemove);
      }
    }
    for (uint256 i = 0; i < adds.length; ++i) {
      address toAdd = adds[i];
      if (toAdd == address(0)) {
        continue;
      }
      if (s_allowList.add(toAdd)) {
        emit AllowListAdd(toAdd);
      }
    }
  }

  /// @notice Ensure that the ARM has not emitted a bad signal, and that the latest heartbeat is not stale.
  modifier whenHealthy() {
    if (IARM(i_armProxy).isCursed()) revert BadARMSignal();
    _;
  }
}
