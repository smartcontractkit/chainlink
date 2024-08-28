// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ILiquidityContainer} from "../../../liquiditymanager/interfaces/ILiquidityContainer.sol";
import {ITokenMessenger} from "../USDC/ITokenMessenger.sol";

import {Pool} from "../../libraries/Pool.sol";
import {TokenPool} from "../TokenPool.sol";
import {USDCTokenPool} from "../USDC/USDCTokenPool.sol";
import {USDCBridgeMigrator} from "./USDCBridgeMigrator.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {EnumerableSet} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";

/// @notice A token pool for USDC which uses CCTP for supported chains and Lock/Release for all others
/// @dev The functionality from LockReleaseTokenPool.sol has been duplicated due to lack of compiler support for shared
/// constructors between parents
/// @dev The primary token mechanism in this pool is Burn/Mint with CCTP, with Lock/Release as the
/// secondary, opt in mechanism for chains not currently supporting CCTP.
contract HybridLockReleaseUSDCTokenPool is USDCTokenPool, USDCBridgeMigrator {
  using SafeERC20 for IERC20;
  using EnumerableSet for EnumerableSet.UintSet;

  event LiquidityTransferred(address indexed from, uint64 indexed remoteChainSelector, uint256 amount);
  event LiquidityProviderSet(
    address indexed oldProvider, address indexed newProvider, uint64 indexed remoteChainSelector
  );

  event LockReleaseEnabled(uint64 indexed remoteChainSelector);
  event LockReleaseDisabled(uint64 indexed remoteChainSelector);

  error LanePausedForCCTPMigration(uint64 remoteChainSelector);
  error TokenLockingNotAllowedAfterMigration(uint64 remoteChainSelector);

  /// @notice The address of the liquidity provider for a specific chain.
  /// External liquidity is not required when there is one canonical token deployed to a chain,
  /// and CCIP is facilitating mint/burn on all the other chains, in which case the invariant
  /// balanceOf(pool) on home chain >= sum(totalSupply(mint/burn "wrapped" token) on all remote chains) should always hold
  mapping(uint64 remoteChainSelector => address liquidityProvider) internal s_liquidityProvider;

  constructor(
    ITokenMessenger tokenMessenger,
    IERC20 token,
    address[] memory allowlist,
    address rmnProxy,
    address router
  ) USDCTokenPool(tokenMessenger, token, allowlist, rmnProxy, router) USDCBridgeMigrator(address(token), router) {}

  // ================================================================
  // │                   Incoming/Outgoing Mechanisms               |
  // ================================================================

  /// @notice Locks the token in the pool
  /// @dev The _validateLockOrBurn check is an essential security check
  function lockOrBurn(Pool.LockOrBurnInV1 calldata lockOrBurnIn)
    public
    virtual
    override
    returns (Pool.LockOrBurnOutV1 memory)
  {
    // // If the alternative mechanism (L/R) for chains which have it enabled
    if (!shouldUseLockRelease(lockOrBurnIn.remoteChainSelector)) {
      return super.lockOrBurn(lockOrBurnIn);
    }

    // Circle requires a supply-lock to prevent outgoing messages once the migration process begins.
    // This prevents new outgoing messages once the migration has begun to ensure any the procedure runs as expected
    if (s_proposedUSDCMigrationChain == lockOrBurnIn.remoteChainSelector) {
      revert LanePausedForCCTPMigration(s_proposedUSDCMigrationChain);
    }

    return _lockReleaseOutgoingMessage(lockOrBurnIn);
  }

  /// @notice Release tokens from the pool to the recipient
  /// @dev The _validateReleaseOrMint check is an essential security check
  function releaseOrMint(Pool.ReleaseOrMintInV1 calldata releaseOrMintIn)
    public
    virtual
    override
    returns (Pool.ReleaseOrMintOutV1 memory)
  {
    if (!shouldUseLockRelease(releaseOrMintIn.remoteChainSelector)) {
      return super.releaseOrMint(releaseOrMintIn);
    }
    return _lockReleaseIncomingMessage(releaseOrMintIn);
  }

  /// @notice Contains the alternative mechanism for incoming tokens, in this implementation is "Release" incoming tokens
  function _lockReleaseIncomingMessage(Pool.ReleaseOrMintInV1 calldata releaseOrMintIn)
    internal
    virtual
    returns (Pool.ReleaseOrMintOutV1 memory)
  {
    _validateReleaseOrMint(releaseOrMintIn);

    // Decrease internal tracking of locked tokens to ensure accurate accounting for burnLockedUSDC() migration
    s_lockedTokensByChainSelector[releaseOrMintIn.remoteChainSelector] -= releaseOrMintIn.amount;

    // Release to the offRamp, which forwards it to the recipient
    getToken().safeTransfer(releaseOrMintIn.receiver, releaseOrMintIn.amount);

    emit Released(msg.sender, releaseOrMintIn.receiver, releaseOrMintIn.amount);

    return Pool.ReleaseOrMintOutV1({destinationAmount: releaseOrMintIn.amount});
  }

  /// @notice Contains the alternative mechanism, in this implementation is "Lock" on outgoing tokens
  function _lockReleaseOutgoingMessage(Pool.LockOrBurnInV1 calldata lockOrBurnIn)
    internal
    virtual
    returns (Pool.LockOrBurnOutV1 memory)
  {
    _validateLockOrBurn(lockOrBurnIn);

    // Increase internal accounting of locked tokens for burnLockedUSDC() migration
    s_lockedTokensByChainSelector[lockOrBurnIn.remoteChainSelector] += lockOrBurnIn.amount;

    emit Locked(msg.sender, lockOrBurnIn.amount);

    return Pool.LockOrBurnOutV1({destTokenAddress: getRemoteToken(lockOrBurnIn.remoteChainSelector), destPoolData: ""});
  }

  // ================================================================
  // │                   Liquidity Management                       |
  // ================================================================

  /// @notice Gets LiquidityManager, can be address(0) if none is configured.
  /// @return The current liquidity manager for the given chain selector
  function getLiquidityProvider(uint64 remoteChainSelector) external view returns (address) {
    return s_liquidityProvider[remoteChainSelector];
  }

  /// @notice Sets the LiquidityManager address.
  /// @dev Only callable by the owner.
  function setLiquidityProvider(uint64 remoteChainSelector, address liquidityProvider) external onlyOwner {
    address oldProvider = s_liquidityProvider[remoteChainSelector];

    s_liquidityProvider[remoteChainSelector] = liquidityProvider;

    emit LiquidityProviderSet(oldProvider, liquidityProvider, remoteChainSelector);
  }

  /// @notice Adds liquidity to the pool for a specific chain. The tokens should be approved first.
  /// @dev Liquidity is expected to be added on a per chain basis. Parties are expected to provide liquidity for their
  /// own chain which implements non canonical USDC and liquidity is not shared across lanes.
  /// @param amount The amount of liquidity to provide.
  /// @param remoteChainSelector The chain for which liquidity is provided to. Necessary to ensure there's accurate
  /// parity between locked USDC in this contract and the circulating supply on the remote chain
  function provideLiquidity(uint64 remoteChainSelector, uint256 amount) external {
    if (s_liquidityProvider[remoteChainSelector] != msg.sender) revert TokenPool.Unauthorized(msg.sender);

    s_lockedTokensByChainSelector[remoteChainSelector] += amount;

    i_token.safeTransferFrom(msg.sender, address(this), amount);

    emit ILiquidityContainer.LiquidityAdded(msg.sender, amount);
  }

  /// @notice Removed liquidity to the pool. The tokens will be sent to msg.sender.
  /// @param remoteChainSelector The chain where liquidity is being released.
  /// @param amount The amount of liquidity to remove.
  /// @dev The function should only be called if non canonical USDC on the remote chain has been burned and is not being
  /// withdrawn on this chain, otherwise a mismatch may occur between locked token balance and remote circulating supply
  /// which may block a potential future migration of the chain to CCTP.
  function withdrawLiquidity(uint64 remoteChainSelector, uint256 amount) external {
    if (s_liquidityProvider[remoteChainSelector] != msg.sender) revert TokenPool.Unauthorized(msg.sender);

    s_lockedTokensByChainSelector[remoteChainSelector] -= amount;

    i_token.safeTransfer(msg.sender, amount);
    emit ILiquidityContainer.LiquidityRemoved(msg.sender, amount);
  }

  /// @notice This function can be used to transfer liquidity from an older version of the pool to this pool. To do so
  /// this pool will have to be set as the liquidity provider in the older version of the pool. This allows it to transfer the
  /// funds in the old pool to the new pool.
  /// @dev When upgrading a LockRelease pool, this function can be called at the same time as the pool is changed in the
  /// TokenAdminRegistry. This allows for a smooth transition of both liquidity and transactions to the new pool.
  /// Alternatively, when no multicall is available, a portion of the funds can be transferred to the new pool before
  /// changing which pool CCIP uses, to ensure both pools can operate. Then the pool should be changed in the
  /// TokenAdminRegistry, which will activate the new pool. All new transactions will use the new pool and its
  /// liquidity. Finally, the remaining liquidity can be transferred to the new pool using this function one more time.
  /// @param from The address of the old pool.
  /// @param amount The amount of liquidity to transfer.
  function transferLiquidity(address from, uint64 remoteChainSelector, uint256 amount) external onlyOwner {
    HybridLockReleaseUSDCTokenPool(from).withdrawLiquidity(remoteChainSelector, amount);

    s_lockedTokensByChainSelector[remoteChainSelector] += amount;

    emit LiquidityTransferred(from, remoteChainSelector, amount);
  }

  // ================================================================
  // │                   Alt Mechanism Logic                        |
  // ================================================================

  /// @notice Return whether a lane should use the alternative L/R mechanism in the token pool.
  /// @param remoteChainSelector the remote chain the lane is interacting with
  /// @return bool Return true if the alternative L/R mechanism should be used
  function shouldUseLockRelease(uint64 remoteChainSelector) public view virtual returns (bool) {
    return s_shouldUseLockRelease[remoteChainSelector];
  }

  /// @notice Updates Updates designations for chains on whether to use primary or alt mechanism on CCIP messages
  /// @param removes A list of chain selectors to disable Lock-Release, and enforce BM
  /// @param adds A list of chain selectors to enable LR instead of BM
  function updateChainSelectorMechanisms(uint64[] calldata removes, uint64[] calldata adds) external onlyOwner {
    for (uint256 i = 0; i < removes.length; ++i) {
      delete s_shouldUseLockRelease[removes[i]];
      emit LockReleaseDisabled(removes[i]);
    }

    for (uint256 i = 0; i < adds.length; ++i) {
      s_shouldUseLockRelease[adds[i]] = true;
      emit LockReleaseEnabled(adds[i]);
    }
  }
}
