// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ILiquidityContainer} from "../../liquiditymanager/interfaces/ILiquidityContainer.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";

import {Pool} from "../libraries/Pool.sol";
import {TokenPool} from "./TokenPool.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice Token pool used for tokens on their native chain. This uses a lock and release mechanism.
/// Because of lock/unlock requiring liquidity, this pool contract also has function to add and remove
/// liquidity. This allows for proper bookkeeping for both user and liquidity provider balances.
/// @dev One token per LockReleaseTokenPool.
contract LockReleaseTokenPool is TokenPool, ILiquidityContainer, ITypeAndVersion {
  using SafeERC20 for IERC20;

  error InsufficientLiquidity();
  error LiquidityNotAccepted();

  event LiquidityTransferred(address indexed from, uint256 amount);

  string public constant override typeAndVersion = "LockReleaseTokenPool 1.5.0";

  /// @dev Whether or not the pool accepts liquidity.
  /// External liquidity is not required when there is one canonical token deployed to a chain,
  /// and CCIP is facilitating mint/burn on all the other chains, in which case the invariant
  /// balanceOf(pool) on home chain >= sum(totalSupply(mint/burn "wrapped" token) on all remote chains) should always hold
  bool internal immutable i_acceptLiquidity;
  /// @notice The address of the rebalancer.
  address internal s_rebalancer;

  constructor(
    IERC20 token,
    address[] memory allowlist,
    address rmnProxy,
    bool acceptLiquidity,
    address router
  ) TokenPool(token, allowlist, rmnProxy, router) {
    i_acceptLiquidity = acceptLiquidity;
  }

  /// @notice Locks the token in the pool
  /// @dev The _validateLockOrBurn check is an essential security check
  function lockOrBurn(
    Pool.LockOrBurnInV1 calldata lockOrBurnIn
  ) external virtual override returns (Pool.LockOrBurnOutV1 memory) {
    _validateLockOrBurn(lockOrBurnIn);

    emit Locked(msg.sender, lockOrBurnIn.amount);

    return Pool.LockOrBurnOutV1({destTokenAddress: getRemoteToken(lockOrBurnIn.remoteChainSelector), destPoolData: ""});
  }

  /// @notice Release tokens from the pool to the recipient
  /// @dev The _validateReleaseOrMint check is an essential security check
  function releaseOrMint(
    Pool.ReleaseOrMintInV1 calldata releaseOrMintIn
  ) external virtual override returns (Pool.ReleaseOrMintOutV1 memory) {
    _validateReleaseOrMint(releaseOrMintIn);

    // Release to the recipient
    getToken().safeTransfer(releaseOrMintIn.receiver, releaseOrMintIn.amount);

    emit Released(msg.sender, releaseOrMintIn.receiver, releaseOrMintIn.amount);

    return Pool.ReleaseOrMintOutV1({destinationAmount: releaseOrMintIn.amount});
  }

  // @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) public pure virtual override returns (bool) {
    return interfaceId == type(ILiquidityContainer).interfaceId || super.supportsInterface(interfaceId);
  }

  /// @notice Gets LiquidityManager, can be address(0) if none is configured.
  /// @return The current liquidity manager.
  function getRebalancer() external view returns (address) {
    return s_rebalancer;
  }

  /// @notice Sets the LiquidityManager address.
  /// @dev Only callable by the owner.
  function setRebalancer(address rebalancer) external onlyOwner {
    s_rebalancer = rebalancer;
  }

  /// @notice Checks if the pool can accept liquidity.
  /// @return true if the pool can accept liquidity, false otherwise.
  function canAcceptLiquidity() external view returns (bool) {
    return i_acceptLiquidity;
  }

  /// @notice Adds liquidity to the pool. The tokens should be approved first.
  /// @param amount The amount of liquidity to provide.
  function provideLiquidity(uint256 amount) external {
    if (!i_acceptLiquidity) revert LiquidityNotAccepted();
    if (s_rebalancer != msg.sender) revert Unauthorized(msg.sender);

    i_token.safeTransferFrom(msg.sender, address(this), amount);
    emit LiquidityAdded(msg.sender, amount);
  }

  /// @notice Removed liquidity to the pool. The tokens will be sent to msg.sender.
  /// @param amount The amount of liquidity to remove.
  function withdrawLiquidity(uint256 amount) external {
    if (s_rebalancer != msg.sender) revert Unauthorized(msg.sender);

    if (i_token.balanceOf(address(this)) < amount) revert InsufficientLiquidity();
    i_token.safeTransfer(msg.sender, amount);
    emit LiquidityRemoved(msg.sender, amount);
  }

  /// @notice This function can be used to transfer liquidity from an older version of the pool to this pool. To do so
  /// this pool will have to be set as the rebalancer in the older version of the pool. This allows it to transfer the
  /// funds in the old pool to the new pool.
  /// @dev When upgrading a LockRelease pool, this function can be called at the same time as the pool is changed in the
  /// TokenAdminRegistry. This allows for a smooth transition of both liquidity and transactions to the new pool.
  /// Alternatively, when no multicall is available, a portion of the funds can be transferred to the new pool before
  /// changing which pool CCIP uses, to ensure both pools can operate. Then the pool should be changed in the
  /// TokenAdminRegistry, which will activate the new pool. All new transactions will use the new pool and its
  /// liquidity. Finally, the remaining liquidity can be transferred to the new pool using this function one more time.
  /// @param from The address of the old pool.
  /// @param amount The amount of liquidity to transfer.
  function transferLiquidity(address from, uint256 amount) external onlyOwner {
    LockReleaseTokenPool(from).withdrawLiquidity(amount);

    emit LiquidityTransferred(from, amount);
  }
}
