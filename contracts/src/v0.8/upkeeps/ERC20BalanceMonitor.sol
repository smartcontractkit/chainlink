// SPDX-License-Identifier: MIT

pragma solidity ^0.8.4;

import "../ConfirmedOwner.sol";
import "../interfaces/KeeperCompatibleInterface.sol";
import "../vendor/openzeppelin-solidity/v4.7.0/contracts/security/Pausable.sol";
import "../vendor/openzeppelin-solidity/v4.7.0/contracts/token/ERC20/utils/SafeERC20.sol";

/**
 * @title The ERC20BalanceMonitor contract.
 * @notice A keeper-compatible contract that monitors and funds ERC20 tokens.
 */
contract ERC20BalanceMonitor is ConfirmedOwner, Pausable, KeeperCompatibleInterface {
  uint16 private constant MAX_WATCHLIST_SIZE = 300;
  uint256 private constant MIN_GAS_FOR_TRANSFER = 55_000;

  event FundsWithdrawn(uint256 amountWithdrawn, address payee);
  event TopUpSucceeded(address indexed topUpAddress);
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
    uint96 topUpLevel;
    uint56 lastTopUpTimestamp;
  }

  IERC20 private s_erc20Token;
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
   * @param topUpLevels the amount to top up to for each subscription
   */
  function setWatchList(
    address[] calldata addresses,
    uint96[] calldata minBalances,
    uint96[] calldata topUpLevels
  ) external onlyOwner {
    if (
      addresses.length != minBalances.length ||
      addresses.length != topUpLevels.length ||
      addresses.length > MAX_WATCHLIST_SIZE
    ) {
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
      if (topUpLevels[idx] <= minBalances[idx]) {
        revert InvalidWatchList();
      }
      s_targets[addresses[idx]] = Target({
        isActive: true,
        minBalance: minBalances[idx],
        topUpLevel: topUpLevels[idx],
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
    uint256 minWaitPeriodSeconds = s_minWaitPeriodSeconds;
    uint256 contractBalance = s_erc20Token.balanceOf(address(this));
    Target memory target;
    for (uint256 idx = 0; idx < watchList.length; idx++) {
      target = s_targets[watchList[idx]];
      uint256 targetTokenBalance = s_erc20Token.balanceOf(watchList[idx]);
      if (
        target.lastTopUpTimestamp + minWaitPeriodSeconds <= block.timestamp &&
        targetTokenBalance < target.minBalance &&
        contractBalance >= (target.topUpLevel - targetTokenBalance)
      ) {
        uint256 topUpAmount = target.topUpLevel - targetTokenBalance;
        needsFunding[count] = watchList[idx];
        count++;
        contractBalance -= topUpAmount;
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
    uint256 contractBalance = s_erc20Token.balanceOf(address(this));
    for (uint256 idx = 0; idx < needsFunding.length; idx++) {
      target = s_targets[needsFunding[idx]];
      uint256 targetTokenBalance = s_erc20Token.balanceOf(needsFunding[idx]);
      if (
        target.isActive &&
        target.lastTopUpTimestamp + minWaitPeriodSeconds <= block.timestamp &&
        targetTokenBalance < target.minBalance &&
        contractBalance >= (target.topUpLevel - targetTokenBalance)
      ) {
        uint256 topUpAmount = target.topUpLevel - targetTokenBalance;
        s_targets[needsFunding[idx]].lastTopUpTimestamp = uint56(block.timestamp);
        contractBalance -= topUpAmount;
        SafeERC20.safeTransfer(s_erc20Token, needsFunding[idx], topUpAmount);
        emit TopUpSucceeded(needsFunding[idx]);
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
    SafeERC20.safeTransfer(s_erc20Token, payee, amount);
    emit FundsWithdrawn(amount, payee);
  }

  /**
   * @notice Sets the ERC20 token address.
   */
  function setERC20TokenAddress(address erc20TokenAddress) public onlyOwner {
    require(erc20TokenAddress != address(0));
    emit ERC20TokenAddressUpdated(address(s_erc20Token), erc20TokenAddress);
    s_erc20Token = IERC20(erc20TokenAddress);
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
    return address(s_erc20Token);
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
      uint96 topUpLevel,
      uint56 lastTopUpTimestamp
    )
  {
    Target memory target = s_targets[targetAddress];
    return (target.isActive, target.minBalance, target.topUpLevel, target.lastTopUpTimestamp);
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
