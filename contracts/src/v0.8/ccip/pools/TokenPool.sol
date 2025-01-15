// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IPoolV1} from "../interfaces/IPool.sol";
import {IRMN} from "../interfaces/IRMN.sol";
import {IRouter} from "../interfaces/IRouter.sol";

import {Ownable2StepMsgSender} from "../../shared/access/Ownable2StepMsgSender.sol";
import {Pool} from "../libraries/Pool.sol";
import {RateLimiter} from "../libraries/RateLimiter.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {IERC20Metadata} from
  "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";

/// @notice Base abstract class with common functions for all token pools.
/// A token pool serves as isolated place for holding tokens and token specific logic
/// that may execute as tokens move across the bridge.
/// @dev This pool supports different decimals on different chains but using this feature could impact the total number
/// of tokens in circulation. Since all of the tokens are locked/burned on the source, and a rounded amount is
/// minted/released on the destination, the number of tokens minted/released could be less than the number of tokens
/// burned/locked. This is because the source chain does not know about the destination token decimals. This is not a
/// problem if the decimals are the same on both chains.
///
/// Example:
/// Assume there is a token with 6 decimals on chain A and 3 decimals on chain B.
/// - 1.234567 tokens are burned on chain A.
/// - 1.234    tokens are minted on chain B.
/// When sending the 1.234 tokens back to chain A, you will receive 1.234000 tokens on chain A, effectively losing
/// 0.000567 tokens.
/// In the case of a burnMint pool on chain A, these funds are burned in the pool on chain A.
/// In the case of a lockRelease pool on chain A, these funds accumulate in the pool on chain A.
abstract contract TokenPool is IPoolV1, Ownable2StepMsgSender {
  using EnumerableSet for EnumerableSet.Bytes32Set;
  using EnumerableSet for EnumerableSet.AddressSet;
  using EnumerableSet for EnumerableSet.UintSet;
  using RateLimiter for RateLimiter.TokenBucket;

  error CallerIsNotARampOnRouter(address caller);
  error ZeroAddressNotAllowed();
  error SenderNotAllowed(address sender);
  error AllowListNotEnabled();
  error NonExistentChain(uint64 remoteChainSelector);
  error ChainNotAllowed(uint64 remoteChainSelector);
  error CursedByRMN();
  error ChainAlreadyExists(uint64 chainSelector);
  error InvalidSourcePoolAddress(bytes sourcePoolAddress);
  error InvalidToken(address token);
  error Unauthorized(address caller);
  error PoolAlreadyAdded(uint64 remoteChainSelector, bytes remotePoolAddress);
  error InvalidRemotePoolForChain(uint64 remoteChainSelector, bytes remotePoolAddress);
  error InvalidRemoteChainDecimals(bytes sourcePoolData);
  error MismatchedArrayLengths();
  error OverflowDetected(uint8 remoteDecimals, uint8 localDecimals, uint256 remoteAmount);
  error InvalidDecimalArgs(uint8 expected, uint8 actual);

  event Locked(address indexed sender, uint256 amount);
  event Burned(address indexed sender, uint256 amount);
  event Released(address indexed sender, address indexed recipient, uint256 amount);
  event Minted(address indexed sender, address indexed recipient, uint256 amount);
  event ChainAdded(
    uint64 remoteChainSelector,
    bytes remoteToken,
    RateLimiter.Config outboundRateLimiterConfig,
    RateLimiter.Config inboundRateLimiterConfig
  );
  event ChainConfigured(
    uint64 remoteChainSelector,
    RateLimiter.Config outboundRateLimiterConfig,
    RateLimiter.Config inboundRateLimiterConfig
  );
  event ChainRemoved(uint64 remoteChainSelector);
  event RemotePoolAdded(uint64 indexed remoteChainSelector, bytes remotePoolAddress);
  event RemotePoolRemoved(uint64 indexed remoteChainSelector, bytes remotePoolAddress);
  event AllowListAdd(address sender);
  event AllowListRemove(address sender);
  event RouterUpdated(address oldRouter, address newRouter);
  event RateLimitAdminSet(address rateLimitAdmin);

  struct ChainUpdate {
    uint64 remoteChainSelector; // Remote chain selector
    bytes[] remotePoolAddresses; // Address of the remote pool, ABI encoded in the case of a remote EVM chain.
    bytes remoteTokenAddress; // Address of the remote token, ABI encoded in the case of a remote EVM chain.
    RateLimiter.Config outboundRateLimiterConfig; // Outbound rate limited config, meaning the rate limits for all of the onRamps for the given chain
    RateLimiter.Config inboundRateLimiterConfig; // Inbound rate limited config, meaning the rate limits for all of the offRamps for the given chain
  }

  struct RemoteChainConfig {
    RateLimiter.TokenBucket outboundRateLimiterConfig; // Outbound rate limited config, meaning the rate limits for all of the onRamps for the given chain
    RateLimiter.TokenBucket inboundRateLimiterConfig; // Inbound rate limited config, meaning the rate limits for all of the offRamps for the given chain
    bytes remoteTokenAddress; // Address of the remote token, ABI encoded in the case of a remote EVM chain.
    EnumerableSet.Bytes32Set remotePools; // Set of remote pool hashes, ABI encoded in the case of a remote EVM chain.
  }

  /// @dev The bridgeable token that is managed by this pool. Pools could support multiple tokens at the same time if
  /// required, but this implementation only supports one token.
  IERC20 internal immutable i_token;
  /// @dev The number of decimals of the token managed by this pool.
  uint8 internal immutable i_tokenDecimals;
  /// @dev The address of the RMN proxy
  address internal immutable i_rmnProxy;
  /// @dev The immutable flag that indicates if the pool is access-controlled.
  bool internal immutable i_allowlistEnabled;
  /// @dev A set of addresses allowed to trigger lockOrBurn as original senders.
  /// Only takes effect if i_allowlistEnabled is true.
  /// This can be used to ensure only token-issuer specified addresses can move tokens.
  EnumerableSet.AddressSet internal s_allowlist;
  /// @dev The address of the router
  IRouter internal s_router;
  /// @dev A set of allowed chain selectors. We want the allowlist to be enumerable to
  /// be able to quickly determine (without parsing logs) who can access the pool.
  /// @dev The chain selectors are in uint256 format because of the EnumerableSet implementation.
  EnumerableSet.UintSet internal s_remoteChainSelectors;
  mapping(uint64 remoteChainSelector => RemoteChainConfig) internal s_remoteChainConfigs;
  /// @notice A mapping of hashed pool addresses to their unhashed form. This is used to be able to find the actually
  /// configured pools and not just their hashed versions.
  mapping(bytes32 poolAddressHash => bytes poolAddress) internal s_remotePoolAddresses;
  /// @notice The address of the rate limiter admin.
  /// @dev Can be address(0) if none is configured.
  address internal s_rateLimitAdmin;

  constructor(IERC20 token, uint8 localTokenDecimals, address[] memory allowlist, address rmnProxy, address router) {
    if (address(token) == address(0) || router == address(0) || rmnProxy == address(0)) revert ZeroAddressNotAllowed();
    i_token = token;
    i_rmnProxy = rmnProxy;

    try IERC20Metadata(address(token)).decimals() returns (uint8 actualTokenDecimals) {
      if (localTokenDecimals != actualTokenDecimals) {
        revert InvalidDecimalArgs(localTokenDecimals, actualTokenDecimals);
      }
    } catch {
      // The decimals function doesn't exist, which is possible since it's optional in the ERC20 spec. We skip the check and
      // assume the supplied token decimals are correct.
    }
    i_tokenDecimals = localTokenDecimals;

    s_router = IRouter(router);

    // Pool can be set as permissioned or permissionless at deployment time only to save hot-path gas.
    i_allowlistEnabled = allowlist.length > 0;
    if (i_allowlistEnabled) {
      _applyAllowListUpdates(new address[](0), allowlist);
    }
  }

  /// @inheritdoc IPoolV1
  function isSupportedToken(
    address token
  ) public view virtual returns (bool) {
    return token == address(i_token);
  }

  /// @notice Gets the IERC20 token that this pool can lock or burn.
  /// @return token The IERC20 token representation.
  function getToken() public view returns (IERC20 token) {
    return i_token;
  }

  /// @notice Get RMN proxy address
  /// @return rmnProxy Address of RMN proxy
  function getRmnProxy() public view returns (address rmnProxy) {
    return i_rmnProxy;
  }

  /// @notice Gets the pool's Router
  /// @return router The pool's Router
  function getRouter() public view returns (address router) {
    return address(s_router);
  }

  /// @notice Sets the pool's Router
  /// @param newRouter The new Router
  function setRouter(
    address newRouter
  ) public onlyOwner {
    if (newRouter == address(0)) revert ZeroAddressNotAllowed();
    address oldRouter = address(s_router);
    s_router = IRouter(newRouter);

    emit RouterUpdated(oldRouter, newRouter);
  }

  /// @notice Signals which version of the pool interface is supported
  function supportsInterface(
    bytes4 interfaceId
  ) public pure virtual override returns (bool) {
    return interfaceId == Pool.CCIP_POOL_V1 || interfaceId == type(IPoolV1).interfaceId
      || interfaceId == type(IERC165).interfaceId;
  }

  // ================================================================
  // │                         Validation                           │
  // ================================================================

  /// @notice Validates the lock or burn input for correctness on
  /// - token to be locked or burned
  /// - RMN curse status
  /// - allowlist status
  /// - if the sender is a valid onRamp
  /// - rate limit status
  /// @param lockOrBurnIn The input to validate.
  /// @dev This function should always be called before executing a lock or burn. Not doing so would allow
  /// for various exploits.
  function _validateLockOrBurn(
    Pool.LockOrBurnInV1 calldata lockOrBurnIn
  ) internal {
    if (!isSupportedToken(lockOrBurnIn.localToken)) revert InvalidToken(lockOrBurnIn.localToken);
    if (IRMN(i_rmnProxy).isCursed(bytes16(uint128(lockOrBurnIn.remoteChainSelector)))) revert CursedByRMN();
    _checkAllowList(lockOrBurnIn.originalSender);

    _onlyOnRamp(lockOrBurnIn.remoteChainSelector);
    _consumeOutboundRateLimit(lockOrBurnIn.remoteChainSelector, lockOrBurnIn.amount);
  }

  /// @notice Validates the release or mint input for correctness on
  /// - token to be released or minted
  /// - RMN curse status
  /// - if the sender is a valid offRamp
  /// - if the source pool is valid
  /// - rate limit status
  /// @param releaseOrMintIn The input to validate.
  /// @dev This function should always be called before executing a release or mint. Not doing so would allow
  /// for various exploits.
  function _validateReleaseOrMint(
    Pool.ReleaseOrMintInV1 calldata releaseOrMintIn
  ) internal {
    if (!isSupportedToken(releaseOrMintIn.localToken)) revert InvalidToken(releaseOrMintIn.localToken);
    if (IRMN(i_rmnProxy).isCursed(bytes16(uint128(releaseOrMintIn.remoteChainSelector)))) revert CursedByRMN();
    _onlyOffRamp(releaseOrMintIn.remoteChainSelector);

    // Validates that the source pool address is configured on this pool.
    if (!isRemotePool(releaseOrMintIn.remoteChainSelector, releaseOrMintIn.sourcePoolAddress)) {
      revert InvalidSourcePoolAddress(releaseOrMintIn.sourcePoolAddress);
    }

    _consumeInboundRateLimit(releaseOrMintIn.remoteChainSelector, releaseOrMintIn.amount);
  }

  // ================================================================
  // │                      Token decimals                          │
  // ================================================================

  /// @notice Gets the IERC20 token decimals on the local chain.
  function getTokenDecimals() public view virtual returns (uint8 decimals) {
    return i_tokenDecimals;
  }

  function _encodeLocalDecimals() internal view virtual returns (bytes memory) {
    return abi.encode(i_tokenDecimals);
  }

  function _parseRemoteDecimals(
    bytes memory sourcePoolData
  ) internal view virtual returns (uint8) {
    // Fallback to the local token decimals if the source pool data is empty. This allows for backwards compatibility.
    if (sourcePoolData.length == 0) {
      return i_tokenDecimals;
    }
    if (sourcePoolData.length != 32) {
      revert InvalidRemoteChainDecimals(sourcePoolData);
    }
    uint256 remoteDecimals = abi.decode(sourcePoolData, (uint256));
    if (remoteDecimals > type(uint8).max) {
      revert InvalidRemoteChainDecimals(sourcePoolData);
    }
    return uint8(remoteDecimals);
  }

  /// @notice Calculates the local amount based on the remote amount and decimals.
  /// @param remoteAmount The amount on the remote chain.
  /// @param remoteDecimals The decimals of the token on the remote chain.
  /// @return The local amount.
  /// @dev This function protects against overflows. If there is a transaction that hits the overflow check, it is
  /// probably incorrect as that means the amount cannot be represented on this chain. If the local decimals have been
  /// wrongly configured, the token issuer could redeploy the pool with the correct decimals and manually re-execute the
  /// CCIP tx to fix the issue.
  function _calculateLocalAmount(uint256 remoteAmount, uint8 remoteDecimals) internal view virtual returns (uint256) {
    if (remoteDecimals == i_tokenDecimals) {
      return remoteAmount;
    }
    if (remoteDecimals > i_tokenDecimals) {
      uint8 decimalsDiff = remoteDecimals - i_tokenDecimals;
      if (decimalsDiff > 77) {
        // This is a safety check to prevent overflow in the next calculation.
        revert OverflowDetected(remoteDecimals, i_tokenDecimals, remoteAmount);
      }
      // Solidity rounds down so there is no risk of minting more tokens than the remote chain sent.
      return remoteAmount / (10 ** decimalsDiff);
    }

    // This is a safety check to prevent overflow in the next calculation.
    // More than 77 would never fit in a uint256 and would cause an overflow. We also check if the resulting amount
    // would overflow.
    uint8 diffDecimals = i_tokenDecimals - remoteDecimals;
    if (diffDecimals > 77 || remoteAmount > type(uint256).max / (10 ** diffDecimals)) {
      revert OverflowDetected(remoteDecimals, i_tokenDecimals, remoteAmount);
    }

    return remoteAmount * (10 ** diffDecimals);
  }

  // ================================================================
  // │                     Chain permissions                        │
  // ================================================================

  /// @notice Gets the pool address on the remote chain.
  /// @param remoteChainSelector Remote chain selector.
  /// @dev To support non-evm chains, this value is encoded into bytes
  function getRemotePools(
    uint64 remoteChainSelector
  ) public view returns (bytes[] memory) {
    bytes32[] memory remotePoolHashes = s_remoteChainConfigs[remoteChainSelector].remotePools.values();

    bytes[] memory remotePools = new bytes[](remotePoolHashes.length);
    for (uint256 i = 0; i < remotePoolHashes.length; ++i) {
      remotePools[i] = s_remotePoolAddresses[remotePoolHashes[i]];
    }

    return remotePools;
  }

  /// @notice Checks if the pool address is configured on the remote chain.
  /// @param remoteChainSelector Remote chain selector.
  /// @param remotePoolAddress The address of the remote pool.
  function isRemotePool(uint64 remoteChainSelector, bytes calldata remotePoolAddress) public view returns (bool) {
    return s_remoteChainConfigs[remoteChainSelector].remotePools.contains(keccak256(remotePoolAddress));
  }

  /// @notice Gets the token address on the remote chain.
  /// @param remoteChainSelector Remote chain selector.
  /// @dev To support non-evm chains, this value is encoded into bytes
  function getRemoteToken(
    uint64 remoteChainSelector
  ) public view returns (bytes memory) {
    return s_remoteChainConfigs[remoteChainSelector].remoteTokenAddress;
  }

  /// @notice Adds a remote pool for a given chain selector. This could be due to a pool being upgraded on the remote
  /// chain. We don't simply want to replace the old pool as there could still be valid inflight messages from the old
  /// pool. This function allows for multiple pools to be added for a single chain selector.
  /// @param remoteChainSelector The remote chain selector for which the remote pool address is being added.
  /// @param remotePoolAddress The address of the new remote pool.
  function addRemotePool(uint64 remoteChainSelector, bytes calldata remotePoolAddress) external onlyOwner {
    if (!isSupportedChain(remoteChainSelector)) revert NonExistentChain(remoteChainSelector);

    _setRemotePool(remoteChainSelector, remotePoolAddress);
  }

  /// @notice Removes the remote pool address for a given chain selector.
  /// @dev All inflight txs from the remote pool will be rejected after it is removed. To ensure no loss of funds, there
  /// should be no inflight txs from the given pool.
  function removeRemotePool(uint64 remoteChainSelector, bytes calldata remotePoolAddress) external onlyOwner {
    if (!isSupportedChain(remoteChainSelector)) revert NonExistentChain(remoteChainSelector);

    if (!s_remoteChainConfigs[remoteChainSelector].remotePools.remove(keccak256(remotePoolAddress))) {
      revert InvalidRemotePoolForChain(remoteChainSelector, remotePoolAddress);
    }

    emit RemotePoolRemoved(remoteChainSelector, remotePoolAddress);
  }

  /// @inheritdoc IPoolV1
  function isSupportedChain(
    uint64 remoteChainSelector
  ) public view returns (bool) {
    return s_remoteChainSelectors.contains(remoteChainSelector);
  }

  /// @notice Get list of allowed chains
  /// @return list of chains.
  function getSupportedChains() public view returns (uint64[] memory) {
    uint256[] memory uint256ChainSelectors = s_remoteChainSelectors.values();
    uint64[] memory chainSelectors = new uint64[](uint256ChainSelectors.length);
    for (uint256 i = 0; i < uint256ChainSelectors.length; ++i) {
      chainSelectors[i] = uint64(uint256ChainSelectors[i]);
    }

    return chainSelectors;
  }

  /// @notice Sets the permissions for a list of chains selectors. Actual senders for these chains
  /// need to be allowed on the Router to interact with this pool.
  /// @param remoteChainSelectorsToRemove A list of chain selectors to remove.
  /// @param chainsToAdd A list of chains and their new permission status & rate limits. Rate limits
  /// are only used when the chain is being added through `allowed` being true.
  /// @dev Only callable by the owner
  function applyChainUpdates(
    uint64[] calldata remoteChainSelectorsToRemove,
    ChainUpdate[] calldata chainsToAdd
  ) external virtual onlyOwner {
    for (uint256 i = 0; i < remoteChainSelectorsToRemove.length; ++i) {
      uint64 remoteChainSelectorToRemove = remoteChainSelectorsToRemove[i];
      // If the chain doesn't exist, revert
      if (!s_remoteChainSelectors.remove(remoteChainSelectorToRemove)) {
        revert NonExistentChain(remoteChainSelectorToRemove);
      }

      // Remove all remote pool hashes for the chain
      bytes32[] memory remotePools = s_remoteChainConfigs[remoteChainSelectorToRemove].remotePools.values();
      for (uint256 j = 0; j < remotePools.length; ++j) {
        s_remoteChainConfigs[remoteChainSelectorToRemove].remotePools.remove(remotePools[j]);
      }

      delete s_remoteChainConfigs[remoteChainSelectorToRemove];

      emit ChainRemoved(remoteChainSelectorToRemove);
    }

    for (uint256 i = 0; i < chainsToAdd.length; ++i) {
      ChainUpdate memory newChain = chainsToAdd[i];
      RateLimiter._validateTokenBucketConfig(newChain.outboundRateLimiterConfig, false);
      RateLimiter._validateTokenBucketConfig(newChain.inboundRateLimiterConfig, false);

      if (newChain.remoteTokenAddress.length == 0) {
        revert ZeroAddressNotAllowed();
      }

      // If the chain already exists, revert
      if (!s_remoteChainSelectors.add(newChain.remoteChainSelector)) {
        revert ChainAlreadyExists(newChain.remoteChainSelector);
      }

      RemoteChainConfig storage remoteChainConfig = s_remoteChainConfigs[newChain.remoteChainSelector];

      remoteChainConfig.outboundRateLimiterConfig = RateLimiter.TokenBucket({
        rate: newChain.outboundRateLimiterConfig.rate,
        capacity: newChain.outboundRateLimiterConfig.capacity,
        tokens: newChain.outboundRateLimiterConfig.capacity,
        lastUpdated: uint32(block.timestamp),
        isEnabled: newChain.outboundRateLimiterConfig.isEnabled
      });
      remoteChainConfig.inboundRateLimiterConfig = RateLimiter.TokenBucket({
        rate: newChain.inboundRateLimiterConfig.rate,
        capacity: newChain.inboundRateLimiterConfig.capacity,
        tokens: newChain.inboundRateLimiterConfig.capacity,
        lastUpdated: uint32(block.timestamp),
        isEnabled: newChain.inboundRateLimiterConfig.isEnabled
      });
      remoteChainConfig.remoteTokenAddress = newChain.remoteTokenAddress;

      for (uint256 j = 0; j < newChain.remotePoolAddresses.length; ++j) {
        _setRemotePool(newChain.remoteChainSelector, newChain.remotePoolAddresses[j]);
      }

      emit ChainAdded(
        newChain.remoteChainSelector,
        newChain.remoteTokenAddress,
        newChain.outboundRateLimiterConfig,
        newChain.inboundRateLimiterConfig
      );
    }
  }

  /// @notice Adds a pool address to the allowed remote token pools for a particular chain.
  /// @param remoteChainSelector The remote chain selector for which the remote pool address is being added.
  /// @param remotePoolAddress The address of the new remote pool.
  function _setRemotePool(uint64 remoteChainSelector, bytes memory remotePoolAddress) internal {
    if (remotePoolAddress.length == 0) {
      revert ZeroAddressNotAllowed();
    }

    bytes32 poolHash = keccak256(remotePoolAddress);

    // Check if the pool already exists.
    if (!s_remoteChainConfigs[remoteChainSelector].remotePools.add(poolHash)) {
      revert PoolAlreadyAdded(remoteChainSelector, remotePoolAddress);
    }

    // Add the pool to the mapping to be able to un-hash it later.
    s_remotePoolAddresses[poolHash] = remotePoolAddress;

    emit RemotePoolAdded(remoteChainSelector, remotePoolAddress);
  }

  // ================================================================
  // │                        Rate limiting                         │
  // ================================================================

  /// @dev The inbound rate limits should be slightly higher than the outbound rate limits. This is because many chains
  /// finalize blocks in batches. CCIP also commits messages in batches: the commit plugin bundles multiple messages in
  /// a single merkle root.
  /// Imagine the following scenario.
  /// - Chain A has an inbound and outbound rate limit of 100 tokens capacity and 1 token per second refill rate.
  /// - Chain B has an inbound and outbound rate limit of 100 tokens capacity and 1 token per second refill rate.
  ///
  /// At time 0:
  /// - Chain A sends 100 tokens to Chain B.
  /// At time 5:
  /// - Chain A sends 5 tokens to Chain B.
  /// At time 6:
  /// The epoch that contains blocks [0-5] is finalized.
  /// Both transactions will be included in the same merkle root and become executable at the same time. This means
  /// the token pool on chain B requires a capacity of 105 to successfully execute both messages at the same time.
  /// The exact additional capacity required depends on the refill rate and the size of the source chain epochs and the
  /// CCIP round time. For simplicity, a 5-10% buffer should be sufficient in most cases.

  /// @notice Sets the rate limiter admin address.
  /// @dev Only callable by the owner.
  /// @param rateLimitAdmin The new rate limiter admin address.
  function setRateLimitAdmin(
    address rateLimitAdmin
  ) external onlyOwner {
    s_rateLimitAdmin = rateLimitAdmin;
    emit RateLimitAdminSet(rateLimitAdmin);
  }

  /// @notice Gets the rate limiter admin address.
  function getRateLimitAdmin() external view returns (address) {
    return s_rateLimitAdmin;
  }

  /// @notice Consumes outbound rate limiting capacity in this pool
  function _consumeOutboundRateLimit(uint64 remoteChainSelector, uint256 amount) internal {
    s_remoteChainConfigs[remoteChainSelector].outboundRateLimiterConfig._consume(amount, address(i_token));
  }

  /// @notice Consumes inbound rate limiting capacity in this pool
  function _consumeInboundRateLimit(uint64 remoteChainSelector, uint256 amount) internal {
    s_remoteChainConfigs[remoteChainSelector].inboundRateLimiterConfig._consume(amount, address(i_token));
  }

  /// @notice Gets the token bucket with its values for the block it was requested at.
  /// @return The token bucket.
  function getCurrentOutboundRateLimiterState(
    uint64 remoteChainSelector
  ) external view returns (RateLimiter.TokenBucket memory) {
    return s_remoteChainConfigs[remoteChainSelector].outboundRateLimiterConfig._currentTokenBucketState();
  }

  /// @notice Gets the token bucket with its values for the block it was requested at.
  /// @return The token bucket.
  function getCurrentInboundRateLimiterState(
    uint64 remoteChainSelector
  ) external view returns (RateLimiter.TokenBucket memory) {
    return s_remoteChainConfigs[remoteChainSelector].inboundRateLimiterConfig._currentTokenBucketState();
  }

  /// @notice Sets multiple chain rate limiter configs.
  /// @param remoteChainSelectors The remote chain selector for which the rate limits apply.
  /// @param outboundConfigs The new outbound rate limiter config, meaning the onRamp rate limits for the given chain.
  /// @param inboundConfigs The new inbound rate limiter config, meaning the offRamp rate limits for the given chain.
  function setChainRateLimiterConfigs(
    uint64[] calldata remoteChainSelectors,
    RateLimiter.Config[] calldata outboundConfigs,
    RateLimiter.Config[] calldata inboundConfigs
  ) external {
    if (msg.sender != s_rateLimitAdmin && msg.sender != owner()) revert Unauthorized(msg.sender);
    if (remoteChainSelectors.length != outboundConfigs.length || remoteChainSelectors.length != inboundConfigs.length) {
      revert MismatchedArrayLengths();
    }

    for (uint256 i = 0; i < remoteChainSelectors.length; ++i) {
      _setRateLimitConfig(remoteChainSelectors[i], outboundConfigs[i], inboundConfigs[i]);
    }
  }

  /// @notice Sets the chain rate limiter config.
  /// @param remoteChainSelector The remote chain selector for which the rate limits apply.
  /// @param outboundConfig The new outbound rate limiter config, meaning the onRamp rate limits for the given chain.
  /// @param inboundConfig The new inbound rate limiter config, meaning the offRamp rate limits for the given chain.
  function setChainRateLimiterConfig(
    uint64 remoteChainSelector,
    RateLimiter.Config memory outboundConfig,
    RateLimiter.Config memory inboundConfig
  ) external {
    if (msg.sender != s_rateLimitAdmin && msg.sender != owner()) revert Unauthorized(msg.sender);

    _setRateLimitConfig(remoteChainSelector, outboundConfig, inboundConfig);
  }

  function _setRateLimitConfig(
    uint64 remoteChainSelector,
    RateLimiter.Config memory outboundConfig,
    RateLimiter.Config memory inboundConfig
  ) internal {
    if (!isSupportedChain(remoteChainSelector)) revert NonExistentChain(remoteChainSelector);
    RateLimiter._validateTokenBucketConfig(outboundConfig, false);
    s_remoteChainConfigs[remoteChainSelector].outboundRateLimiterConfig._setTokenBucketConfig(outboundConfig);
    RateLimiter._validateTokenBucketConfig(inboundConfig, false);
    s_remoteChainConfigs[remoteChainSelector].inboundRateLimiterConfig._setTokenBucketConfig(inboundConfig);
    emit ChainConfigured(remoteChainSelector, outboundConfig, inboundConfig);
  }

  // ================================================================
  // │                           Access                             │
  // ================================================================

  /// @notice Checks whether remote chain selector is configured on this contract, and if the msg.sender
  /// is a permissioned onRamp for the given chain on the Router.
  function _onlyOnRamp(
    uint64 remoteChainSelector
  ) internal view {
    if (!isSupportedChain(remoteChainSelector)) revert ChainNotAllowed(remoteChainSelector);
    if (!(msg.sender == s_router.getOnRamp(remoteChainSelector))) revert CallerIsNotARampOnRouter(msg.sender);
  }

  /// @notice Checks whether remote chain selector is configured on this contract, and if the msg.sender
  /// is a permissioned offRamp for the given chain on the Router.
  function _onlyOffRamp(
    uint64 remoteChainSelector
  ) internal view {
    if (!isSupportedChain(remoteChainSelector)) revert ChainNotAllowed(remoteChainSelector);
    if (!s_router.isOffRamp(remoteChainSelector, msg.sender)) revert CallerIsNotARampOnRouter(msg.sender);
  }

  // ================================================================
  // │                          Allowlist                           │
  // ================================================================

  function _checkAllowList(
    address sender
  ) internal view {
    if (i_allowlistEnabled) {
      if (!s_allowlist.contains(sender)) {
        revert SenderNotAllowed(sender);
      }
    }
  }

  /// @notice Gets whether the allowlist functionality is enabled.
  /// @return true is enabled, false if not.
  function getAllowListEnabled() external view returns (bool) {
    return i_allowlistEnabled;
  }

  /// @notice Gets the allowed addresses.
  /// @return The allowed addresses.
  function getAllowList() external view returns (address[] memory) {
    return s_allowlist.values();
  }

  /// @notice Apply updates to the allow list.
  /// @param removes The addresses to be removed.
  /// @param adds The addresses to be added.
  function applyAllowListUpdates(address[] calldata removes, address[] calldata adds) external onlyOwner {
    _applyAllowListUpdates(removes, adds);
  }

  /// @notice Internal version of applyAllowListUpdates to allow for reuse in the constructor.
  function _applyAllowListUpdates(address[] memory removes, address[] memory adds) internal {
    if (!i_allowlistEnabled) revert AllowListNotEnabled();

    for (uint256 i = 0; i < removes.length; ++i) {
      address toRemove = removes[i];
      if (s_allowlist.remove(toRemove)) {
        emit AllowListRemove(toRemove);
      }
    }
    for (uint256 i = 0; i < adds.length; ++i) {
      address toAdd = adds[i];
      if (toAdd == address(0)) {
        continue;
      }
      if (s_allowlist.add(toAdd)) {
        emit AllowListAdd(toAdd);
      }
    }
  }
}
