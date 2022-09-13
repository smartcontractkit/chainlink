// SPDX-License-Identifier: MIT

pragma solidity ^0.8.4;

import "../ConfirmedOwner.sol";
import "../interfaces/KeeperCompatibleInterface.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/**
 * @title The ERC20BalanceMonitor contract.
 * @notice A keeper-compatible contract that monitors and funds ERC20 tokens.
 */
contract ERC20BalanceMonitor is ConfirmedOwner, Pausable, KeeperCompatibleInterface {
  IERC20 private ERC20Token;

  uint256 private constant MIN_GAS_FOR_TRANSFER = 55_000;

  event FundsAdded(uint256 amountAdded, uint256 newBalance, address sender);
  event FundsWithdrawn(uint256 amountWithdrawn, address payee);
  event TopUpSucceeded(address indexed topUpAddress);
  event TopUpFailed(address indexed topUpAddress);
  event WatchlistUpdated(address[] oldWatchlist, address[] newWatchlist);
  event KeeperRegistryAddressUpdated(address oldAddress, address newAddress);
  event ERC20TokenAddressUpdated(address oldAddress, address newAddress);
  event MinWaitPeriodUpdated(uint256 oldMinWaitPeriod, uint256 newMinWaitPeriod);

  error InvalidWatchList();
  error OnlyKeeperRegistry();
  error DuplicateAddress(address duplicate);

  struct Target {
    bool isActive;
    uint96 minBalance;
    uint96 topUpAmount;
    uint56 lastTopUpTimestamp;
  }

  address private s_keeperRegistryAddress;
  uint256 private s_minWaitPeriodSeconds;
  address[] private s_watchList;
  mapping(address => Target) internal s_targets;

  /**
   * @param erc20TokenAddress the ERC20 token address
   * @param keeperRegistryAddress the address of the keeper registry contract
   * @param minWaitPeriodSeconds the minimum wait period for addresses between funding
   */
  constructor(
    address erc20TokenAddress,
    address keeperRegistryAddress,
    uint256 minWaitPeriodSeconds
  ) ConfirmedOwner(msg.sender) {
    setERC20TokenAddress(erc20TokenAddress);
    setKeeperRegistryAddress(keeperRegistryAddress);
    setMinWaitPeriodSeconds(minWaitPeriodSeconds);
  }

  /**
   * @notice Sets the list of subscriptions to watch and their funding parameters.
   * @param addresses the list of subscription ids to watch
   * @param minBalances the minimum balances for each subscription
   * @param topUpAmounts the amount to top up each subscription
   */
  function setWatchList(
    address[] calldata addresses,
    uint96[] calldata minBalances,
    uint96[] calldata topUpAmounts
  ) external onlyOwner {
    if (addresses.length != minBalances.length || addresses.length != topUpAmounts.length) {
      revert InvalidWatchList();
    }
    address[] memory oldWatchList = s_watchList;
    for (uint256 idx = 0; idx < oldWatchList.length; idx++) {
      s_targets[oldWatchList[idx]].isActive = false;
    }
    for (uint256 idx = 0; idx < addresses.length; idx++) {
      if (s_targets[addresses[idx]].isActive) {
        revert DuplicateAddress(addresses[idx]);
      }
      if (addresses[idx] == address(0)) {
        revert InvalidWatchList();
      }
      if (topUpAmounts[idx] == 0) {
        revert InvalidWatchList();
      }
      s_targets[addresses[idx]] = Target({
        isActive: true,
        minBalance: minBalances[idx],
        topUpAmount: topUpAmounts[idx],
        lastTopUpTimestamp: 0
      });
    }
    s_watchList = addresses;
    emit WatchlistUpdated(oldWatchList, addresses);
  }

  /**
   * @notice Gets a list of subscriptions that are underfunded.
   * @return list of subscriptions that are underfunded
   */
  function getUnderfundedAddresses() public view returns (address[] memory) {
    address[] memory watchList = s_watchList;
    address[] memory needsFunding = new address[](watchList.length);
    uint256 count = 0;
    uint256 minWaitPeriod = s_minWaitPeriodSeconds;
    uint256 contractBalance = ERC20Token.balanceOf(address(this));
    Target memory target;
    for (uint256 idx = 0; idx < watchList.length; idx++) {
      target = s_targets[watchList[idx]];
      uint256 targetTokenBalance = ERC20Token.balanceOf(watchList[idx]);
      if (
        target.lastTopUpTimestamp + minWaitPeriod <= block.timestamp &&
        contractBalance >= target.topUpAmount &&
        targetTokenBalance < target.minBalance
      ) {
        needsFunding[count] = watchList[idx];
        count++;
        contractBalance -= target.topUpAmount;
      }
    }
    if (count != watchList.length) {
      assembly {
        mstore(needsFunding, count) // resize array to number of valid targets
      }
    }
    return needsFunding;
  }

  /**
   * @notice Send funds to the subscriptions provided.
   * @param needsFunding the list of subscriptions to fund
   */
  function topUp(address[] memory needsFunding) public whenNotPaused {
    uint256 minWaitPeriodSeconds = s_minWaitPeriodSeconds;
    Target memory target;
    for (uint256 idx = 0; idx < needsFunding.length; idx++) {
      target = s_targets[needsFunding[idx]];
      uint256 targetTokenBalance = ERC20Token.balanceOf(needsFunding[idx]);
      uint256 contractBalance = ERC20Token.balanceOf(address(this));
      if (
        target.isActive &&
        target.lastTopUpTimestamp + minWaitPeriodSeconds <= block.timestamp &&
        targetTokenBalance < target.minBalance &&
        contractBalance >= target.topUpAmount
      ) {
        bool success = ERC20Token.transfer(needsFunding[idx], target.topUpAmount);
        if (success) {
          s_targets[needsFunding[idx]].lastTopUpTimestamp = uint56(block.timestamp);
          emit TopUpSucceeded(needsFunding[idx]);
        } else {
          emit TopUpFailed(needsFunding[idx]);
        }
      }
      if (gasleft() < MIN_GAS_FOR_TRANSFER) {
        return;
      }
    }
  }

  /**
   * @notice Gets list of subscription ids that are underfunded and returns a keeper-compatible payload.
   * @return upkeepNeeded signals if upkeep is needed, performData is an abi encoded list of subscription ids that need funds
   */
  function checkUpkeep(bytes calldata)
    external
    view
    override
    whenNotPaused
    returns (bool upkeepNeeded, bytes memory performData)
  {
    address[] memory needsFunding = getUnderfundedAddresses();
    upkeepNeeded = needsFunding.length > 0;
    performData = abi.encode(needsFunding);
    return (upkeepNeeded, performData);
  }

  /**
   * @notice Called by the keeper to send funds to underfunded addresses.
   * @param performData the abi encoded list of addresses to fund
   */
  function performUpkeep(bytes calldata performData) external override onlyKeeperRegistry whenNotPaused {
    address[] memory needsFunding = abi.decode(performData, (address[]));
    topUp(needsFunding);
  }

  /**
   * @notice Withdraws the contract balance in the ERC20 token.
   * @param amount the amount of the ERC20 to withdraw
   * @param payee the address to pay
   */
  function withdraw(uint256 amount, address payable payee) external onlyOwner {
    require(payee != address(0));
    emit FundsWithdrawn(amount, payee);
    ERC20Token.transfer(payee, amount);
  }

  /**
   * @notice Sets the ERC20 token address.
   */
  function setERC20TokenAddress(address erc20TokenAddress) public onlyOwner {
    require(erc20TokenAddress != address(0));
    emit ERC20TokenAddressUpdated(address(ERC20Token), erc20TokenAddress);
    ERC20Token = IERC20(erc20TokenAddress);
  }

  /**
   * @notice Sets the keeper registry address.
   */
  function setKeeperRegistryAddress(address keeperRegistryAddress) public onlyOwner {
    require(keeperRegistryAddress != address(0));
    emit KeeperRegistryAddressUpdated(s_keeperRegistryAddress, keeperRegistryAddress);
    s_keeperRegistryAddress = keeperRegistryAddress;
  }

  /**
   * @notice Sets the minimum wait period (in seconds) for subscription ids between funding.
   */
  function setMinWaitPeriodSeconds(uint256 period) public onlyOwner {
    emit MinWaitPeriodUpdated(s_minWaitPeriodSeconds, period);
    s_minWaitPeriodSeconds = period;
  }

  /**
   * @notice Gets the ERC20 token address.
   */
  function getERC20TokenAddress() external view returns (address) {
    return address(ERC20Token);
  }

  /**
   * @notice Gets the keeper registry address.
   */
  function getKeeperRegistryAddress() external view returns (address) {
    return s_keeperRegistryAddress;
  }

  /**
   * @notice Gets the minimum wait period.
   */
  function getMinWaitPeriodSeconds() external view returns (uint256) {
    return s_minWaitPeriodSeconds;
  }

  /**
   * @notice Gets the list of subscription ids being watched.
   */
  function getWatchList() external view returns (address[] memory) {
    return s_watchList;
  }

  /**
   * @notice Gets configuration information for an address on the watchlist
   */
  function getAccountInfo(address targetAddress)
    external
    view
    returns (
      bool isActive,
      uint96 minBalance,
      uint96 topUpAmount,
      uint56 lastTopUpTimestamp
    )
  {
    Target memory target = s_targets[targetAddress];
    return (target.isActive, target.minBalance, target.topUpAmount, target.lastTopUpTimestamp);
  }

  /**
   * @notice Pause the contract, which prevents executing performUpkeep.
   */
  function pause() external onlyOwner {
    _pause();
  }

  /**
   * @notice Unpause the contract.
   */
  function unpause() external onlyOwner {
    _unpause();
  }

  modifier onlyKeeperRegistry() {
    if (msg.sender != s_keeperRegistryAddress) {
      revert OnlyKeeperRegistry();
    }
    _;
  }
}
