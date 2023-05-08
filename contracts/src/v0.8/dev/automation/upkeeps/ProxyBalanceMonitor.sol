// SPDX-License-Identifier: MIT

pragma solidity ^0.8.4;

import "../../../ConfirmedOwner.sol";
import "../../../interfaces/automation/KeeperCompatibleInterface.sol";
import "../../../vendor/openzeppelin-solidity/v4.7.0/contracts/security/Pausable.sol";
import "../../../vendor/openzeppelin-solidity/v4.7.0/contracts/token/ERC20/utils/SafeERC20.sol";

interface IAggregatorProxy {
  function aggregator() external view returns (address);
}

interface IAggregator {
  function linkAvailableForPayment() external view returns (int256 availableBalance);

  function transmitters() external view returns (address[] memory);
}

// TODO - add max number of proxies that can be funced in a single tx

/**
 * @title The ProxyBalanceMonitor contract.
 * @notice A keeper-compatible contract that aggregator proxies and funds them with LINK
 */
contract ProxyBalanceMonitor is ConfirmedOwner, Pausable, KeeperCompatibleInterface {
  uint256 private constant MIN_GAS_FOR_TRANSFER = 55_000;

  event FundsWithdrawn(uint256 amountWithdrawn, address payee);
  event TopUpSucceeded(address indexed topUpAddress);
  event WatchlistUpdated();

  error InvalidWatchList();
  error DuplicateAddress(address duplicate);

  struct Target {
    bool isActive;
    uint96 minBalance;
    uint96 topUpLevel;
    uint56 lastTopUpTimestamp;
  }

  IERC20 private immutable s_linkToken;
  address[] private s_watchList;
  mapping(address => Target) internal s_targets;

  /**
   * @param linkTokenAddress the ERC20 token address
   */
  constructor(
    address linkTokenAddress,
  ) ConfirmedOwner(msg.sender) {
    require(linkTokenAddress != address(0));
    s_linkToken = IERC20(linkTokenAddress);
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
    if (addresses.length != minBalances.length || addresses.length != topUpLevels.length) {
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
    emit WatchlistUpdated();
  }

  /**
   * @notice Adds addresses to the watchlist without overwriting existing members
   * @param addresses the list of subscription ids to watch
   * @param minBalances the minimum balances for each subscription
   * @param topUpLevels the amount to top up to for each subscription
   */
  function addToWatchList(
    address[] calldata addresses,
    uint96[] calldata minBalances,
    uint96[] calldata topUpLevels
  ) external onlyOwner {
    if (addresses.length != minBalances.length || addresses.length != topUpLevels.length) {
      revert InvalidWatchList();
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
      s_watchList.push(addresses[idx]);
    }
    emit WatchlistUpdated();
  }

  /**
   * @notice Gets a list of subscriptions that are underfunded.
   * @return list of subscriptions that are underfunded
   */
  function getUnderfundedAddresses() public view returns (address[] memory) {
    // TODO - shuffle to that size doesn't limit which addresses can be checked
    address[] memory watchList = s_watchList;
    address[] memory needsFunding = new address[](watchList.length);
    uint256 count = 0;
    uint256 upkeepBalance = s_linkToken.balanceOf(address(this));
    Target memory target;
    for (uint256 idx = 0; idx < watchList.length; idx++) {
      target = s_targets[watchList[idx]];
      try IAggregatorProxy(watchList[idx]).aggregator() returns (address aggregator) {
        try IAggregator(aggregator).transmitters() returns (address[] memory transmitters) {
          if (transmitters.length == 0) {
            continue; // aggregator is "dead" and does not require funding
          }
        } catch {
          continue; // not an aggregator, skip
        }
        try IAggregator(aggregator).linkAvailableForPayment() returns (int256 balance) {
          uint256 s_minBalance = 0; // TODO
          uint256 s_topUpAmount = 0; // TODO
          if (balance < s_minBalance && upkeepBalance > s_topUpAmount) {
            count++;
            upkeepBalance -= topUpAmount;
          }
        } catch {
          continue; // not an aggregator, skip
        }
      } catch {
        continue; // not an aggregator proxy, skip
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
    Target memory target;
    uint256 upkeepBalance = s_linkToken.balanceOf(address(this));
    for (uint256 idx = 0; idx < needsFunding.length; idx++) {
      target = s_targets[needsFunding[idx]];
      uint256 targetTokenBalance = s_linkToken.balanceOf(needsFunding[idx]);
      if (
        target.isActive &&
        targetTokenBalance < target.minBalance &&
        upkeepBalance >= (target.topUpLevel - targetTokenBalance)
      ) {
        uint256 topUpAmount = target.topUpLevel - targetTokenBalance;
        s_targets[needsFunding[idx]].lastTopUpTimestamp = uint56(block.timestamp);
        upkeepBalance -= topUpAmount;
        SafeERC20.safeTransfer(s_linkToken, needsFunding[idx], topUpAmount);
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
  function performUpkeep(bytes calldata performData) external override whenNotPaused {
    address[] memory needsFunding = abi.decode(performData, (address[]));
    topUp(needsFunding);
  }

  /**
   * @notice Withdraws the contract balance in the LINK token.
   * @param amount the amount of the LINK to withdraw
   * @param payee the address to pay
   */
  function withdraw(uint256 amount, address payable payee) external onlyOwner {
    require(payee != address(0));
    SafeERC20.safeTransfer(s_linkToken, payee, amount);
    emit FundsWithdrawn(amount, payee);
  }

  /**
   * @notice Gets the LINK token address.
   */
  function getLINKTokenAddress() external view returns (address) {
    return address(s_linkToken);
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
}
