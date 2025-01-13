// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";

import {Pool} from "../libraries/Pool.sol";
import {TokenPool} from "./TokenPool.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice A variation on Lock Release token pools where liquidity is shared among some chains, and stored independently
/// for others. Chains which do not share liquidity are known as siloed chains.
contract SiloedLockReleaseTokenPool is TokenPool, ITypeAndVersion {
  using SafeERC20 for IERC20;

  error InsufficientLiquidity();
  error LiquidityNotAccepted();
  error ChainNotSiloed(uint64 remoteChainSelector);

  event LiquidityTransferred(uint64 remoteChainSelector, address indexed from, uint256 amount);
  event LiquidityAdded(uint64 remoteChainSelector, address indexed provider, uint256 indexed amount);
  event LiquidityRemoved(uint64 remoteChainSelector, address indexed provider, uint256 indexed amount);
  event ChainSiloeDesignationUpdated(uint64 remoteChainSelector, bool isSiloed);

  event RebalancerSet(uint64 indexed remoteChainSelector, address oldRebalancer, address newRebalancer);

  string public constant override typeAndVersion = "SiloedLockReleaseTokenPool 1.6.0";

  /// @notice The amount of tokens available for remote chains which are not siloed as an additional security precaution.
  uint256 internal s_unsiloedTokenBalance;

  struct ChainSiloConfig {
    uint256 siloedBalance; // The amount of tokens available for incoming messages, either locked or as liquidity.
    address rebalancer; // ───────╮ The address allowed to add liquidity for the given siloed chain.
    bool isSiloed; // ────────────╯ Whether funds should be isolated from all other chains or shared amongst all non-siloed chains.
  }

  /// @notice The configuration for each chain that is siloed, or not. By default chains are not siloed.
  mapping(uint64 remoteChainSelector => ChainSiloConfig) internal s_siloedChainConfigs;

  constructor(
    IERC20 token,
    uint8 localTokenDecimals,
    address[] memory allowlist,
    address rmnProxy,
    address router
  ) TokenPool(token, localTokenDecimals, allowlist, rmnProxy, router) {}

  /// @notice Locks the token in the pool
  /// @dev The _validateLockOrBurn check is an essential security check
  function lockOrBurn(
    Pool.LockOrBurnInV1 calldata lockOrBurnIn
  ) external virtual override returns (Pool.LockOrBurnOutV1 memory) {
    _validateLockOrBurn(lockOrBurnIn);

    // If funds need to be siloed, update internal accounting;
    if (s_siloedChainConfigs[lockOrBurnIn.remoteChainSelector].isSiloed) {
      s_siloedChainConfigs[lockOrBurnIn.remoteChainSelector].siloedBalance += lockOrBurnIn.amount;
    }
    // If the messages is going to a chain without siloed funds, update state accounting accordingly.
    else {
      s_unsiloedTokenBalance += lockOrBurnIn.amount;
    }

    emit Locked(msg.sender, lockOrBurnIn.amount);

    return Pool.LockOrBurnOutV1({
      destTokenAddress: getRemoteToken(lockOrBurnIn.remoteChainSelector),
      destPoolData: _encodeLocalDecimals()
    });
  }

  /// @notice Release tokens from the pool to the recipient
  /// @dev The _validateReleaseOrMint check is an essential security check
  function releaseOrMint(
    Pool.ReleaseOrMintInV1 calldata releaseOrMintIn
  ) external virtual override returns (Pool.ReleaseOrMintOutV1 memory) {
    _validateReleaseOrMint(releaseOrMintIn);

    // Calculate the local amount
    uint256 localAmount =
      _calculateLocalAmount(releaseOrMintIn.amount, _parseRemoteDecimals(releaseOrMintIn.sourcePoolData));

    // Tracking balances independently by chain is a security measure to prevent liquidity for one chain from being
    // released by another chain. Since all tokens are stored in the same contract, and not isolated physically,
    // it must be ensured that the remoteChainSelector can be trusted.
    if (s_siloedChainConfigs[releaseOrMintIn.remoteChainSelector].isSiloed) {
      s_siloedChainConfigs[releaseOrMintIn.remoteChainSelector].siloedBalance -= localAmount;
    } else {
      s_unsiloedTokenBalance -= localAmount;
    }

    // Release to the recipient
    getToken().safeTransfer(releaseOrMintIn.receiver, localAmount);

    emit Released(msg.sender, releaseOrMintIn.receiver, localAmount);

    return Pool.ReleaseOrMintOutV1({destinationAmount: localAmount});
  }

  /// @notice Returns whether the tokens locked for a given remote chain should be siloed independently
  /// from all other remote chains.
  /// @param remoteChainSelector the CCIP specific selector for the remote chain being interacted with.
  /// @return isSiloed Whether the funds should be isolated from all the others.
  function chainFundsAreSiloed(
    uint64 remoteChainSelector
  ) external view returns (bool isSiloed) {
    return s_siloedChainConfigs[remoteChainSelector].isSiloed;
  }

  /// @notice Returns the amount of tokens in the token pool that were siloed for a specific remote chain selector.
  /// @param remoteChainSelector the CCIP specific selector for the remote chain being interacted with.
  /// @return lockedTokens The tokens locked into this token pool for the given selector. If the chain is not siloed,
  /// the amount will be the amount of liquidity shared among all unsiloed chains.
  function getSiloedTokensByChain(
    uint64 remoteChainSelector
  ) external view returns (uint256 lockedTokens) {
    if (s_siloedChainConfigs[remoteChainSelector].isSiloed) {
      return s_siloedChainConfigs[remoteChainSelector].siloedBalance;
    }

    return s_unsiloedTokenBalance;
  }

  /// @notice Returns the amount of tokens in the token pool that are shared among all unsiloed chains.
  /// @return unsiloedTokens amount of tokens available to all unsiloed chains.
  function getliquidityForUnsiloedChains() public view returns (uint256) {
    return s_unsiloedTokenBalance;
  }

  /// @notice Updates designations for chains on whether to mark funds as Siloed or not
  /// @param removes A list of chain selectors to disable Siloing. Their funds will be moved into the unsiloed pool.
  /// @param adds A list of chain selectors to enable Siloing. Adding a chain to siloing will not set the rebalancer.
  /// The rebalancer will need to be set separately.
  function updateSiloDesignationForChainSelectors(uint64[] calldata removes, uint64[] calldata adds) external onlyOwner {
    for (uint256 i = 0; i < removes.length; ++i) {
      // When a chain is removed from siloing, the funds are moved to the accounting pool shared by all unsiloed chain.
      s_unsiloedTokenBalance += s_siloedChainConfigs[removes[i]].siloedBalance;

      /// @notice Removing the entire configuration will also delete the rebalancer. If the chain ever becomes siloed
      /// again, the rebalancer will need to be set again.
      delete s_siloedChainConfigs[removes[i]];

      emit ChainSiloeDesignationUpdated(removes[i], false);
    }

    for (uint256 i = 0; i < adds.length; ++i) {
      /// @notice Adding a chain to siloing will not set the rebalancer. The rebalancer will need to be set
      s_siloedChainConfigs[adds[i]].isSiloed = true;
      emit ChainSiloeDesignationUpdated(adds[i], true);
    }
  }

  /// @notice Gets the rebalancer able to provide liquidity for a remote chain selector
  /// @param remoteChainSelector The CCIP specific selector for the remote chain being interacted with.
  /// @return The current liquidity manager, contract owner if the chain's funds are not siloed.
  function getRebalancerByChain(
    uint64 remoteChainSelector
  ) public view returns (address) {
    ChainSiloConfig memory remoteConfig = s_siloedChainConfigs[remoteChainSelector];
    if (remoteConfig.isSiloed) return remoteConfig.rebalancer;
    else return owner();
  }

  /// @notice Sets the Rebalancer address for a given remoteChainSelector.
  /// @dev Only callable by the owner.
  /// @param remoteChainSelector the remote chain to set.
  /// @param newRebalancer the address allowed to add liquidity for the given siloed chain.
  function setRebalancer(uint64 remoteChainSelector, address newRebalancer) external onlyOwner {
    // Storage is used instead of memory to save gas, as the state may need to be updated if the chain is siloed.
    ChainSiloConfig storage remoteConfig = s_siloedChainConfigs[remoteChainSelector];

    if (!remoteConfig.isSiloed) revert ChainNotSiloed(remoteChainSelector);

    address oldRebalancer = remoteConfig.rebalancer;

    remoteConfig.rebalancer = newRebalancer;

    emit RebalancerSet(remoteChainSelector, newRebalancer, oldRebalancer);
  }

  /// @notice Adds liquidity to the pool. The tokens should be approved first.
  /// @param remoteChainSelector the remote chain to set. If the chain is not siloed, the liquidity will be shared among all
  /// non-siloed chains.
  /// @param amount The amount of liquidity to provide.
  /// @dev If the chain is siloed, only the rebalancer can provide liquidity, otherwise the contract owner can.
  function provideLiquidity(uint64 remoteChainSelector, uint256 amount) external {
    if (msg.sender != getRebalancerByChain(remoteChainSelector)) revert Unauthorized(msg.sender);

    // Storage is used instead of memory to save gas, as the state may need to be updated if the chain is siloed.
    ChainSiloConfig storage remoteConfig = s_siloedChainConfigs[remoteChainSelector];

    if (remoteConfig.isSiloed) {
      remoteConfig.siloedBalance += amount;
    } else {
      s_unsiloedTokenBalance += amount;
    }

    i_token.safeTransferFrom(msg.sender, address(this), amount);
    emit LiquidityAdded(remoteChainSelector, msg.sender, amount);
  }

  /// @notice Removed liquidity to the pool. The tokens will be sent to msg.sender.
  /// @param remoteChainSelector the remote chain to set. If the chain is not siloed, then no accounting will be updated,
  /// which can be considered the liquidity for all non-siloed chains sharing liquidity.
  /// @param amount The amount of liquidity to remove.
  /// @dev Only the owner can remove liquidity from the contract, for both siloed and unsiloed chains.
  function withdrawLiquidity(uint64 remoteChainSelector, uint256 amount) external onlyOwner {
    // If funds are siloed by chain, prevent more than has been locked from being removed from the token pool.
    if (s_siloedChainConfigs[remoteChainSelector].isSiloed) {
      s_siloedChainConfigs[remoteChainSelector].siloedBalance -= amount;
    } else {
      s_unsiloedTokenBalance -= amount;
    }

    i_token.safeTransfer(msg.sender, amount);
    emit LiquidityRemoved(remoteChainSelector, msg.sender, amount);
  }
}
