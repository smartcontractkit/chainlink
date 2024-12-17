// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IPoolV1} from "../../interfaces/IPool.sol";
import {IRMN} from "../../interfaces/IRMN.sol";
import {IRouter} from "../../interfaces/IRouter.sol";

import {Ownable2StepMsgSender} from "../../../shared/access/Ownable2StepMsgSender.sol";
import {Pool} from "../../libraries/Pool.sol";
import {RateLimiter} from "../../libraries/RateLimiter.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {IERC165} from "../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/introspection/IERC165.sol";
import {EnumerableSet} from "../../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";

/// @notice This contract is a proof of concept and should NOT be used in production.
abstract contract MultiTokenPool is IPoolV1, Ownable2StepMsgSender {
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
  event RemotePoolSet(uint64 indexed remoteChainSelector, bytes previousPoolAddress, bytes remotePoolAddress);
  event AllowListAdd(address sender);
  event AllowListRemove(address sender);
  event RouterUpdated(address oldRouter, address newRouter);

  struct ChainUpdate {
    uint64 remoteChainSelector; // ──╮ Remote chain selector
    bool allowed; // ────────────────╯ Whether the chain is allowed
    bytes remotePoolAddress; //        Address of the remote pool, ABI encoded in the case of a remote EVM chain.
    bytes remoteTokenAddress; //       Address of the remote token, ABI encoded in the case of a remote EVM chain.
    RateLimiter.Config outboundRateLimiterConfig; // Outbound rate limited config, meaning the rate limits for all of the onRamps for the given chain
    RateLimiter.Config inboundRateLimiterConfig; // Inbound rate limited config, meaning the rate limits for all of the offRamps for the given chain
  }

  struct RemoteChainConfig {
    RateLimiter.TokenBucket outboundRateLimiterConfig; // Outbound rate limited config, meaning the rate limits for all of the onRamps for the given chain
    RateLimiter.TokenBucket inboundRateLimiterConfig; // Inbound rate limited config, meaning the rate limits for all of the offRamps for the given chain
    bytes remotePoolAddress; // Address of the remote pool, ABI encoded in the case of a remote EVM chain.
    bytes remoteTokenAddress; // Address of the remote token, ABI encoded in the case of a remote EVM chain.
  }

  /// @dev The IERC20 token that this pool supports
  EnumerableSet.AddressSet internal s_tokens;
  /// @dev The address of the RMN proxy
  address internal immutable i_rmnProxy;
  /// @dev The immutable flag that indicates if the pool is access-controlled.
  bool internal immutable i_allowlistEnabled;
  /// @dev A set of addresses allowed to trigger lockOrBurn as original senders.
  /// Only takes effect if i_allowlistEnabled is true.
  /// This can be used to ensure only token-issuer specified addresses can
  /// move tokens.
  EnumerableSet.AddressSet internal s_allowlist;
  /// @dev The address of the router
  IRouter internal s_router;
  /// @dev A set of allowed chain selectors. We want the allowlist to be enumerable to
  /// be able to quickly determine (without parsing logs) who can access the pool.
  /// @dev The chain selectors are in uin256 format because of the EnumerableSet implementation.
  EnumerableSet.UintSet internal s_remoteChainSelectors;
  mapping(address token => mapping(uint64 remoteChainSelector => RemoteChainConfig)) internal s_remoteChainConfigs;

  constructor(IERC20[] memory token, address[] memory allowlist, address rmnProxy, address router) {
    if (router == address(0) || rmnProxy == address(0)) revert ZeroAddressNotAllowed();
    for (uint256 i = 0; i < token.length; ++i) {
      s_tokens.add(address(token[i]));
    }
    i_rmnProxy = rmnProxy;
    s_router = IRouter(router);

    // Pool can be set as permissioned or permissionless at deployment time only to save hot-path gas.
    i_allowlistEnabled = allowlist.length > 0;
    if (i_allowlistEnabled) {
      _applyAllowListUpdates(new address[](0), allowlist);
    }
  }

  /// @notice Get RMN proxy address
  /// @return rmnProxy Address of RMN proxy
  function getRmnProxy() public view returns (address rmnProxy) {
    return i_rmnProxy;
  }

  /// @inheritdoc IPoolV1
  function isSupportedToken(
    address token
  ) public view virtual returns (bool) {
    return s_tokens.contains(token);
  }

  /// @notice Gets the IERC20 token that this pool can lock or burn.
  /// @return tokens The IERC20 token representation.
  function getTokens() public view returns (IERC20[] memory tokens) {
    tokens = new IERC20[](s_tokens.length());
    for (uint256 i = 0; i < s_tokens.length(); ++i) {
      tokens[i] = IERC20(s_tokens.at(i));
    }
    return tokens;
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
    Pool.LockOrBurnInV1 memory lockOrBurnIn
  ) internal {
    if (!isSupportedToken(lockOrBurnIn.localToken)) revert InvalidToken(lockOrBurnIn.localToken);
    if (IRMN(i_rmnProxy).isCursed(bytes16(uint128(lockOrBurnIn.remoteChainSelector)))) revert CursedByRMN();
    _checkAllowList(lockOrBurnIn.originalSender);

    _onlyOnRamp(lockOrBurnIn.remoteChainSelector);
    _consumeOutboundRateLimit(lockOrBurnIn.localToken, lockOrBurnIn.remoteChainSelector, lockOrBurnIn.amount);
  }

  /// @notice Validates the release or mint input for correctness on
  /// - token to be released or minted
  /// - RMN curse status
  /// - if the sender is a valid offRamp
  /// - if the source pool is valid
  /// - rate limit status
  /// @param releaseOrMintIn The input to validate.
  /// @dev This function should always be called before executing a lock or burn. Not doing so would allow
  /// for various exploits.
  function _validateReleaseOrMint(
    Pool.ReleaseOrMintInV1 memory releaseOrMintIn
  ) internal {
    if (!isSupportedToken(releaseOrMintIn.localToken)) revert InvalidToken(releaseOrMintIn.localToken);
    if (IRMN(i_rmnProxy).isCursed(bytes16(uint128(releaseOrMintIn.remoteChainSelector)))) revert CursedByRMN();
    _onlyOffRamp(releaseOrMintIn.remoteChainSelector);

    // Validates that the source pool address is configured on this pool.
    bytes memory configuredRemotePool = getRemotePool(releaseOrMintIn.localToken, releaseOrMintIn.remoteChainSelector);
    if (
      configuredRemotePool.length == 0
        || keccak256(releaseOrMintIn.sourcePoolAddress) != keccak256(configuredRemotePool)
    ) {
      revert InvalidSourcePoolAddress(releaseOrMintIn.sourcePoolAddress);
    }
    _consumeInboundRateLimit(releaseOrMintIn.localToken, releaseOrMintIn.remoteChainSelector, releaseOrMintIn.amount);
  }

  // ================================================================
  // │                     Chain permissions                        │
  // ================================================================

  /// @notice Gets the pool address on the remote chain.
  /// @param remoteChainSelector Remote chain selector.
  /// @dev To support non-evm chains, this value is encoded into bytes
  function getRemotePool(address token, uint64 remoteChainSelector) public view returns (bytes memory) {
    return s_remoteChainConfigs[token][remoteChainSelector].remotePoolAddress;
  }

  /// @notice Gets the token address on the remote chain.
  /// @param remoteChainSelector Remote chain selector.
  /// @dev To support non-evm chains, this value is encoded into bytes
  function getRemoteToken(address token, uint64 remoteChainSelector) public view returns (bytes memory) {
    return s_remoteChainConfigs[token][remoteChainSelector].remoteTokenAddress;
  }

  /// @notice Sets the remote pool address for a given chain selector.
  /// @param remoteChainSelector The remote chain selector for which the remote pool address is being set.
  /// @param remotePoolAddress The address of the remote pool.
  function setRemotePool(
    address token,
    uint64 remoteChainSelector,
    bytes calldata remotePoolAddress
  ) external onlyOwner {
    if (!isSupportedChain(remoteChainSelector)) revert NonExistentChain(remoteChainSelector);

    bytes memory prevAddress = s_remoteChainConfigs[token][remoteChainSelector].remotePoolAddress;
    s_remoteChainConfigs[token][remoteChainSelector].remotePoolAddress = remotePoolAddress;

    emit RemotePoolSet(remoteChainSelector, prevAddress, remotePoolAddress);
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
  /// @dev Only callable by the owner
  /// @param chains A list of chains and their new permission status & rate limits. Rate limits
  /// are only used when the chain is being added through `allowed` being true.
  function applyChainUpdates(address token, ChainUpdate[] calldata chains) external virtual onlyOwner {
    for (uint256 i = 0; i < chains.length; ++i) {
      ChainUpdate memory update = chains[i];
      RateLimiter._validateTokenBucketConfig(update.outboundRateLimiterConfig, !update.allowed);
      RateLimiter._validateTokenBucketConfig(update.inboundRateLimiterConfig, !update.allowed);

      if (update.allowed) {
        // If the chain already exists, revert
        if (!s_remoteChainSelectors.add(update.remoteChainSelector)) {
          revert ChainAlreadyExists(update.remoteChainSelector);
        }

        if (update.remotePoolAddress.length == 0 || update.remoteTokenAddress.length == 0) {
          revert ZeroAddressNotAllowed();
        }

        s_remoteChainConfigs[token][update.remoteChainSelector] = RemoteChainConfig({
          outboundRateLimiterConfig: RateLimiter.TokenBucket({
            rate: update.outboundRateLimiterConfig.rate,
            capacity: update.outboundRateLimiterConfig.capacity,
            tokens: update.outboundRateLimiterConfig.capacity,
            lastUpdated: uint32(block.timestamp),
            isEnabled: update.outboundRateLimiterConfig.isEnabled
          }),
          inboundRateLimiterConfig: RateLimiter.TokenBucket({
            rate: update.inboundRateLimiterConfig.rate,
            capacity: update.inboundRateLimiterConfig.capacity,
            tokens: update.inboundRateLimiterConfig.capacity,
            lastUpdated: uint32(block.timestamp),
            isEnabled: update.inboundRateLimiterConfig.isEnabled
          }),
          remotePoolAddress: update.remotePoolAddress,
          remoteTokenAddress: update.remoteTokenAddress
        });

        emit ChainAdded(
          update.remoteChainSelector,
          update.remoteTokenAddress,
          update.outboundRateLimiterConfig,
          update.inboundRateLimiterConfig
        );
      } else {
        // If the chain doesn't exist, revert
        if (!s_remoteChainSelectors.remove(update.remoteChainSelector)) {
          revert NonExistentChain(update.remoteChainSelector);
        }

        delete s_remoteChainConfigs[token][update.remoteChainSelector];

        emit ChainRemoved(update.remoteChainSelector);
      }
    }
  }

  // ================================================================
  // │                        Rate limiting                         │
  // ================================================================

  /// @notice Consumes outbound rate limiting capacity in this pool
  function _consumeOutboundRateLimit(address token, uint64 remoteChainSelector, uint256 amount) internal {
    s_remoteChainConfigs[token][remoteChainSelector].outboundRateLimiterConfig._consume(amount, token);
  }

  /// @notice Consumes inbound rate limiting capacity in this pool
  function _consumeInboundRateLimit(address token, uint64 remoteChainSelector, uint256 amount) internal {
    s_remoteChainConfigs[token][remoteChainSelector].inboundRateLimiterConfig._consume(amount, token);
  }

  /// @notice Gets the token bucket with its values for the block it was requested at.
  /// @return The token bucket.
  function getCurrentOutboundRateLimiterState(
    address token,
    uint64 remoteChainSelector
  ) external view returns (RateLimiter.TokenBucket memory) {
    return s_remoteChainConfigs[token][remoteChainSelector].outboundRateLimiterConfig._currentTokenBucketState();
  }

  /// @notice Gets the token bucket with its values for the block it was requested at.
  /// @return The token bucket.
  function getCurrentInboundRateLimiterState(
    address token,
    uint64 remoteChainSelector
  ) external view returns (RateLimiter.TokenBucket memory) {
    return s_remoteChainConfigs[token][remoteChainSelector].inboundRateLimiterConfig._currentTokenBucketState();
  }

  /// @notice Sets the chain rate limiter config.
  /// @param remoteChainSelector The remote chain selector for which the rate limits apply.
  /// @param outboundConfig The new outbound rate limiter config, meaning the onRamp rate limits for the given chain.
  /// @param inboundConfig The new inbound rate limiter config, meaning the offRamp rate limits for the given chain.
  function _setChainRateLimiterConfig(
    address token,
    uint64 remoteChainSelector,
    RateLimiter.Config memory outboundConfig,
    RateLimiter.Config memory inboundConfig
  ) internal {
    if (!isSupportedChain(remoteChainSelector)) revert NonExistentChain(remoteChainSelector);
    RateLimiter._validateTokenBucketConfig(outboundConfig, false);
    s_remoteChainConfigs[token][remoteChainSelector].outboundRateLimiterConfig._setTokenBucketConfig(outboundConfig);
    RateLimiter._validateTokenBucketConfig(inboundConfig, false);
    s_remoteChainConfigs[token][remoteChainSelector].inboundRateLimiterConfig._setTokenBucketConfig(inboundConfig);
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
    if (i_allowlistEnabled && !s_allowlist.contains(sender)) revert SenderNotAllowed(sender);
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
  /// @dev allowlisting will be removed before public launch
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
