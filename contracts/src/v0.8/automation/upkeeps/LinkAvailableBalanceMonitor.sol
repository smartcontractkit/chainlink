// SPDX-License-Identifier: MIT

pragma solidity 0.8.19;

import {AutomationCompatibleInterface} from "../interfaces/AutomationCompatibleInterface.sol";
import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {Pausable} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/security/Pausable.sol";


interface IAggregatorProxy {
  function aggregator() external view returns (address);
}

interface ILinkAvailable {
  function linkAvailableForPayment() external view returns (int256 availableBalance);
}

/// @title The LinkAvailableBalanceMonitor contract.
/// @notice A keeper-compatible contract that monitors target contracts for balance from a custom
/// function linkAvailableForPayment() and funds them with LINK if it falls below a defined
/// threshold. Also supports aggregator proxy contracts monitoring which require fetching the actual
/// target contract through a predefined interface.
/// @dev with 30 addresses as the s_maxPerform, the measured max gas usage of performUpkeep is around 2M
/// therefore, we recommend an upkeep gas limit of 3M (this has a 33% margin of safety). Although, nothing
/// prevents us from using 5M gas and increasing s_maxPerform, 30 seems like a reasonable batch size that
/// is probably plenty for most needs.
/// @dev with 130 addresses as the s_maxCheck, the measured max gas usage of checkUpkeep is around 3.5M,
/// which is 30% below the 5M limit.
/// Note that testing conditions DO NOT match live chain gas usage, hence the margins. Change
/// at your own risk!!!
/// @dev some areas for improvement / acknowledgement of limitations:
///  validate that all addresses conform to interface when adding them to the watchlist
///  this is a "trusless" upkeep, meaning it does not trust the caller of performUpkeep;
/// we could save a fair amount of gas and re-write this upkeep for use with Automation v2.0+,
/// which has significantly different trust assumptions
contract LinkAvailableBalanceMonitor is ConfirmedOwner, Pausable, AutomationCompatibleInterface {
  event BalanceUpdated(address indexed addr, uint256 oldBalance, uint256 newBalance);
  event FundsWithdrawn(uint256 amountWithdrawn, address payee);
  event MaxCheckUpdated(uint256 oldMaxCheck, uint256 newMaxCheck);
  event MaxPerformUpdated(uint256 oldMaxPerform, uint256 newMaxPerform);
  event MinWaitPeriodUpdated(uint256 s_minWaitPeriodSeconds, uint256 minWaitPeriodSeconds);
  event TopUpBlocked(address indexed topUpAddress);
  event TopUpFailed(address indexed recipient);
  event TopUpSucceeded(address indexed topUpAddress);
  event TopUpUpdated(address indexed addr, uint256 oldTopUpAmount, uint256 newTopUpAmount);
  event WatchlistUpdated();

  error InvalidWatchList();
  error DuplicateAddress(address duplicate);

  struct MonitoredAddress {
    uint96 minBalance;
    uint96 topUpAmount;
    uint56 lastTopUpTimestamp;
    bool isActive;
  }

  IERC20 private immutable i_linkToken;
  uint256 private s_minWaitPeriodSeconds;
  uint16 private s_maxPerform;
  uint16 private s_maxCheck;
  address[] private s_watchList;
  mapping(address targetAddress => MonitoredAddress targetProperties) internal s_targets;

  /// @param linkTokenAddress the LINK token address
  constructor(
    address linkTokenAddress,
    uint16 maxPerform,
    uint16 maxCheck,
    uint256 minWaitPeriodSeconds
  ) ConfirmedOwner(msg.sender) {
    require(linkTokenAddress != address(0), "LinkAvailableBalanceMonitor: invalid linkTokenAddress");
    i_linkToken = IERC20(linkTokenAddress);
    setMaxPerform(maxPerform);
    setMaxCheck(maxCheck);
    setMinWaitPeriodSeconds(minWaitPeriodSeconds);
  }

  /// @notice Sets the list of subscriptions to watch and their funding parameters
  /// @param addresses the list of target addresses to watch (could be direct target or IAggregatorProxy)
  /// @param minBalances the list of corresponding minBalance for the target address
  /// @param topUpAmounts the list of corresponding minTopUp for the target address
  function setWatchList(
    address[] calldata addresses,
    uint96[] calldata minBalances,
    uint96[] calldata topUpAmounts
  ) external onlyOwner {
    if (addresses.length != minBalances.length || addresses.length != topUpAmounts.length) {
      revert InvalidWatchList();
    }
    for (uint256 idx = 0; idx < s_watchList.length; idx++) {
      delete s_targets[s_watchList[idx]];
    }
    for (uint256 idx = 0; idx < addresses.length; idx++) {
      address targetAddress = addresses[idx];
      if (s_targets[targetAddress].isActive) {
        revert DuplicateAddress(addresses[idx]);
      }
      if (addresses[idx] == address(0)) {
        revert InvalidWatchList();
      }
      s_targets[targetAddress] = MonitoredAddress({
        isActive: true,
        minBalance: minBalances[idx],
        topUpAmount: topUpAmounts[idx],
        lastTopUpTimestamp: 0
      });
    }
    s_watchList = addresses;
    emit WatchlistUpdated();
  }

  /// @notice Gets a list of proxies that are underfunded, up to the s_maxPerform size
  /// @dev the function starts at a random index in the list to avoid biasing the first
  /// addresses in the list over latter ones.
  /// @dev the function will check at most s_maxCheck proxies in a single call
  /// @dev the function returns a list with a max length of s_maxPerform
  /// @return list of target addresses which are underfunded
  function sampleUnderfundedAddresses() public view returns (address[] memory) {
    uint16 maxPerform = s_maxPerform;
    uint16 maxCheck = s_maxCheck;
    uint256 numTargets = s_watchList.length;
    uint256 idx = uint256(blockhash(block.number - 1)) % numTargets;
    uint256 numToCheck = numTargets < maxCheck ? numTargets : maxCheck;
    uint256 numFound = 0;
    uint256 minWaitPeriod = s_minWaitPeriodSeconds;
    address[] memory targetsToFund = new address[](maxPerform);
    MonitoredAddress memory target;
    for (
      uint256 numChecked = 0;
      numChecked < numToCheck;
      (idx, numChecked) = ((idx + 1) % numTargets, numChecked + 1)
    ) {
      address targetAddress = s_watchList[idx];
      target = s_targets[targetAddress];
      (bool needsFunding, ) = _needsFunding(targetAddress, target.minBalance);
      if (needsFunding && target.lastTopUpTimestamp + minWaitPeriod <= block.timestamp) {
        targetsToFund[numFound] = targetAddress;
        numFound++;
        if (numFound == maxPerform) {
          break; // max number of addresses in batch reached
        }
      }
    }
    if (numFound != maxPerform) {
      assembly {
        mstore(targetsToFund, numFound) // resize array to number of valid targets
      }
    }
    return targetsToFund;
  }

  function topUp(address[] memory targetAddresses) public whenNotPaused {
    MonitoredAddress memory target;
    uint256 localBalance = i_linkToken.balanceOf(address(this));
    address[] memory targetsTopUp = new address[](targetAddresses.length);
    for (uint256 idx = 0; idx < targetAddresses.length; idx++) {
      address targetAddress = targetAddresses[idx];
      target = s_targets[targetAddress];
      (bool needsFunding, ) = _needsFunding(targetAddress, target.minBalance);
      if (target.isActive && localBalance >= target.topUpAmount && !(targetAddress == address(0)) && needsFunding) {
        targetsTopUp[idx] = targetAddress;
        localBalance -= target.topUpAmount;
      } else {
        emit TopUpBlocked(targetAddress);
      }
    }
    for (uint256 idx = 0; idx < targetsTopUp.length; idx++) {
      address targetAddress = targetsTopUp[idx];
      if (targetAddress == address(0)) {
        continue;
      }
      target = s_targets[targetAddress];
      bool success = i_linkToken.transfer(targetAddress, target.topUpAmount);
      if (success) {
        emit TopUpSucceeded(targetAddress);
      } else {
        emit TopUpFailed(targetAddress);
      }
    }
  }

  /// @notice checks the target (could be direct target or IAggregatorProxy), and determines
  /// if it is elligible for funding
  /// @param targetAddress the target to check
  /// @param minBalance minimum balance required for the target
  /// @return bool whether the target needs funding or not
  /// @return address the address of the contract needing funding
  function _needsFunding(address targetAddress, uint256 minBalance) private view returns (bool, address) {
    ILinkAvailable target;
    IAggregatorProxy proxy = IAggregatorProxy(targetAddress);
    // Explicitly check if the targetAddress is the zero address
    // or if it's not a contract. In both cases return with false,
    // to prevent target.linkAvailableForPayment from running,
    // which would revert the operation.
    if (targetAddress == address(0) || !(targetAddress.code.length > 0)) {
      return (false, address(0));
    }
    try proxy.aggregator() returns (address aggregatorAddress) {
      if (aggregatorAddress == address(0)) {
        return (false, address(0));
      }
      target = ILinkAvailable(aggregatorAddress);
    } catch {
      target = ILinkAvailable(targetAddress);
    }
    try target.linkAvailableForPayment() returns (int256 balance) {
      if (balance < 0 || uint256(balance) < minBalance) {
        return (true, address(target));
      }
    } catch {}
    return (false, address(0));
  }

  /// @notice Gets list of subscription ids that are underfunded and returns a keeper-compatible payload.
  /// @return upkeepNeeded signals if upkeep is needed
  /// @return performData is an abi encoded list of subscription ids that need funds
  function checkUpkeep(
    bytes calldata
  ) external view override whenNotPaused returns (bool upkeepNeeded, bytes memory performData) {
    address[] memory needsFunding = sampleUnderfundedAddresses();
    upkeepNeeded = needsFunding.length > 0;
    performData = abi.encode(needsFunding);
    return (upkeepNeeded, performData);
  }

  /// @notice Called by the keeper to send funds to underfunded addresses.
  /// @param performData the abi encoded list of addresses to fund
  function performUpkeep(bytes calldata performData) external override whenNotPaused {
    address[] memory needsFunding = abi.decode(performData, (address[]));
    topUp(needsFunding);
  }

  /// @notice Withdraws the contract balance in the LINK token.
  /// @param amount the amount of the LINK to withdraw
  /// @param payee the address to pay
  function withdraw(uint256 amount, address payable payee) external onlyOwner {
    require(payee != address(0), "LinkAvailableBalanceMonitor: invalid payee address");
    i_linkToken.transfer(payee, amount);
    emit FundsWithdrawn(amount, payee);
  }

  /// @notice Sets the minimum balance for the given target address
  function setMinBalance(address target, uint96 minBalance) external onlyOwner {
    require(target != address(0), "LinkAvailableBalanceMonitor: invalid target address");
    require(minBalance > 0, "LinkAvailableBalanceMonitor: invalid minBalance");
    if (!s_targets[target].isActive) {
      revert InvalidWatchList();
    }
    uint256 oldBalance = s_targets[target].minBalance;
    s_targets[target].minBalance = minBalance;
    emit BalanceUpdated(target, oldBalance, minBalance);
  }

  /// @notice Sets the minimum balance for the given target address
  function setTopUpAmount(address target, uint96 topUpAmount) external onlyOwner {
    require(target != address(0), "LinkAvailableBalanceMonitor: invalid target address");
    require(topUpAmount > 0, "LinkAvailableBalanceMonitor: invalid topUpAmount");
    if (!s_targets[target].isActive) {
      revert InvalidWatchList();
    }
    uint256 oldTopUpAmount = s_targets[target].topUpAmount;
    s_targets[target].topUpAmount = topUpAmount;
    emit BalanceUpdated(target, oldTopUpAmount, topUpAmount);
  }

  /// @notice Update s_maxPerform
  function setMaxPerform(uint16 maxPerform) public onlyOwner {
    s_maxPerform = maxPerform;
    emit MaxPerformUpdated(s_maxPerform, maxPerform);
  }

  /// @notice Update s_maxCheck
  function setMaxCheck(uint16 maxCheck) public onlyOwner {
    s_maxCheck = maxCheck;
    emit MaxCheckUpdated(s_maxCheck, maxCheck);
  }

  /// @notice Sets the minimum wait period (in seconds) for addresses between funding
  function setMinWaitPeriodSeconds(uint256 minWaitPeriodSeconds) public onlyOwner {
    s_minWaitPeriodSeconds = minWaitPeriodSeconds;
    emit MinWaitPeriodUpdated(s_minWaitPeriodSeconds, minWaitPeriodSeconds);
  }

  /// @notice Gets maxPerform
  function getMaxPerform() external view returns (uint16) {
    return s_maxPerform;
  }

  /// @notice Gets maxCheck
  function getMaxCheck() external view returns (uint16) {
    return s_maxCheck;
  }

  /// @notice Gets the minimum wait period
  function getMinWaitPeriodSeconds() external view returns (uint256) {
    return s_minWaitPeriodSeconds;
  }

  /// @notice Gets the list of subscription ids being watched
  function getWatchList() external view returns (address[] memory) {
    return s_watchList;
  }

  /// @notice Gets configuration information for an address on the watchlist
  function getAccountInfo(
    address targetAddress
  ) external view returns (bool isActive, uint256 minBalance, uint256 topUpAmount) {
    MonitoredAddress memory target = s_targets[targetAddress];
    return (target.isActive, target.minBalance, target.topUpAmount);
  }

  /// @notice Pause the contract, which prevents executing performUpkeep
  function pause() external onlyOwner {
    _pause();
  }

  /// @notice Unpause the contract
  function unpause() external onlyOwner {
    _unpause();
  }
}
