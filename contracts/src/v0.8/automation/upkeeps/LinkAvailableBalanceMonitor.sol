// SPDX-License-Identifier: MIT

pragma solidity 0.8.19;

import {AutomationCompatibleInterface} from "../interfaces/AutomationCompatibleInterface.sol";
import {AccessControl} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/access/AccessControl.sol";
import {EnumerableMap} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableMap.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";
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
///  this is a "trustless" upkeep, meaning it does not trust the caller of performUpkeep;
/// we could save a fair amount of gas and re-write this upkeep for use with Automation v2.0+,
/// which has significantly different trust assumptions
contract LinkAvailableBalanceMonitor is AccessControl, AutomationCompatibleInterface, Pausable {
  using EnumerableMap for EnumerableMap.UintToAddressMap;
  using EnumerableSet for EnumerableSet.AddressSet;

  event BalanceUpdated(address indexed addr, uint256 oldBalance, uint256 newBalance);
  event FundsWithdrawn(uint256 amountWithdrawn, address payee);
  event UpkeepIntervalSet(uint256 oldUpkeepInterval, uint256 newUpkeepInterval);
  event MaxCheckSet(uint256 oldMaxCheck, uint256 newMaxCheck);
  event MaxPerformSet(uint256 oldMaxPerform, uint256 newMaxPerform);
  event MinWaitPeriodSet(uint256 s_minWaitPeriodSeconds, uint256 minWaitPeriodSeconds);
  event TopUpBlocked(address indexed topUpAddress);
  event TopUpFailed(address indexed recipient);
  event TopUpSucceeded(address indexed topUpAddress, uint256 amount);
  event TopUpUpdated(address indexed addr, uint256 oldTopUpAmount, uint256 newTopUpAmount);
  event WatchlistUpdated();

  error InvalidAddress(address target);
  error InvalidMaxCheck(uint16 maxCheck);
  error InvalixMaxPerform(uint16 maxPerform);
  error InvalidMinBalance(uint96 minBalance);
  error InvalidTopUpAmount(uint96 topUpAmount);
  error InvalidUpkeepInterval(uint8 upkeepInterval);
  error InvalidLinkTokenAddress(address lt);
  error InvalidWatchList();
  error InvalidChainSelector();
  error DuplicateAddress(address duplicate);
  error ReentrantCall();

  struct MonitoredAddress {
    uint96 minBalance;
    uint96 topUpAmount;
    uint56 lastTopUpTimestamp;
    bool isActive;
  }

  bytes32 private constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
  bytes32 private constant EXECUTOR_ROLE = keccak256("EXECUTOR_ROLE");
  uint96 private constant DEFAULT_TOP_UP_AMOUNT_JUELS = 3000000000000000000;
  uint96 private constant DEFAULT_MIN_BALANCE_JUELS = 1000000000000000000;
  IERC20 private immutable i_linkToken;

  uint256 private s_minWaitPeriodSeconds;
  uint16 private s_maxPerform;
  uint16 private s_maxCheck;
  uint8 private s_upkeepInterval;

  /// @notice s_watchList contains all the addresses watched by this monitor
  /// @dev It mainly provides the length() function
  EnumerableSet.AddressSet private s_watchList;

  /// @notice s_targets contains all the addresses watched by this monitor
  /// Each key points to a MonitoredAddress with all the needed metadata
  mapping(address targetAddress => MonitoredAddress targetProperties) private s_targets;

  /// @notice s_onRampAddresses represents a list of CCIP onRamp addresses watched on this contract
  /// There has to be only one onRamp per dstChainSelector.
  /// dstChainSelector is needed as we have to track the live onRamp, and delete the onRamp
  /// whenever a new one is deployed with the same dstChainSelector.
  EnumerableMap.UintToAddressMap private s_onRampAddresses;

  bool private reentrancyGuard;

  /// @param admin is the administrator address of this contract
  /// @param linkToken the LINK token address
  /// @param minWaitPeriodSeconds represents the amount of time that has to wait a contract to be funded
  /// @param maxPerform maximum amount of contracts to fund
  /// @param maxCheck maximum amount of contracts to check
  /// @param upkeepInterval randomizes the check for underfunded contracts
  constructor(
    address admin,
    IERC20 linkToken,
    uint256 minWaitPeriodSeconds,
    uint16 maxPerform,
    uint16 maxCheck,
    uint8 upkeepInterval
  ) {
    _setRoleAdmin(ADMIN_ROLE, ADMIN_ROLE);
    _setRoleAdmin(EXECUTOR_ROLE, ADMIN_ROLE);
    _grantRole(ADMIN_ROLE, admin);
    i_linkToken = linkToken;
    setMinWaitPeriodSeconds(minWaitPeriodSeconds);
    setMaxPerform(maxPerform);
    setMaxCheck(maxCheck);
    setUpkeepInterval(upkeepInterval);
    reentrancyGuard = false;
  }

  /// @notice Sets the list of subscriptions to watch and their funding parameters
  /// @param addresses the list of target addresses to watch (could be direct target or IAggregatorProxy)
  /// @param minBalances the list of corresponding minBalance for the target address
  /// @param topUpAmounts the list of corresponding minTopUp for the target address
  function setWatchList(
    address[] calldata addresses,
    uint96[] calldata minBalances,
    uint96[] calldata topUpAmounts,
    uint64[] calldata dstChainSelectors
  ) external onlyAdminOrExecutor {
    if (
      addresses.length != minBalances.length ||
      addresses.length != topUpAmounts.length ||
      addresses.length != dstChainSelectors.length
    ) {
      revert InvalidWatchList();
    }
    for (uint256 idx = s_watchList.length(); idx > 0; idx--) {
      address member = s_watchList.at(idx - 1);
      s_watchList.remove(member);
      delete s_targets[member];
    }
    // s_onRampAddresses is not the same length as s_watchList, so it has
    // to be clean in a separate loop
    for (uint256 idx = s_onRampAddresses.length(); idx > 0; idx--) {
      (uint256 key, ) = s_onRampAddresses.at(idx - 1);
      s_onRampAddresses.remove(key);
    }
    for (uint256 idx = 0; idx < addresses.length; idx++) {
      address targetAddress = addresses[idx];
      if (s_targets[targetAddress].isActive) revert DuplicateAddress(targetAddress);
      if (targetAddress == address(0)) revert InvalidWatchList();
      if (minBalances[idx] == 0) revert InvalidWatchList();
      if (topUpAmounts[idx] == 0) revert InvalidWatchList();
      s_targets[targetAddress] = MonitoredAddress({
        isActive: true,
        minBalance: minBalances[idx],
        topUpAmount: topUpAmounts[idx],
        lastTopUpTimestamp: 0
      });
      if (dstChainSelectors[idx] > 0) {
        s_onRampAddresses.set(dstChainSelectors[idx], targetAddress);
      }
      s_watchList.add(targetAddress);
    }
    emit WatchlistUpdated();
  }

  /// @notice Adds a new address to the watchlist
  /// @param targetAddress the address to be added to the watchlist
  /// @param dstChainSelector carries a non-zero value in case the targetAddress is an onRamp, otherwise it carries a 0
  /// @dev this function has to be compatible with the event onRampSet(address, dstChainSelector) emitted by
  /// the CCIP router. Important detail to know is this event is also emitted when an onRamp is decommissioned,
  /// in which case it will carry the proper dstChainSelector along with the 0x0 address
  function addToWatchListOrDecommission(address targetAddress, uint64 dstChainSelector) public onlyAdminOrExecutor {
    if (s_targets[targetAddress].isActive) revert DuplicateAddress(targetAddress);
    if (targetAddress == address(0) && dstChainSelector == 0) revert InvalidAddress(targetAddress);
    bool onRampExists = s_onRampAddresses.contains(dstChainSelector);
    // if targetAddress is an existing onRamp, there's a need of cleaning the previous onRamp associated to this dstChainSelector
    // there's no need to remove any other address that's not an onRamp
    if (dstChainSelector > 0 && onRampExists) {
      address oldAddress = s_onRampAddresses.get(dstChainSelector);
      removeFromWatchList(oldAddress);
    }
    // only add the new address if it's not 0x0
    if (targetAddress != address(0)) {
      s_targets[targetAddress] = MonitoredAddress({
        isActive: true,
        minBalance: DEFAULT_MIN_BALANCE_JUELS,
        topUpAmount: DEFAULT_TOP_UP_AMOUNT_JUELS,
        lastTopUpTimestamp: 0
      });
      s_watchList.add(targetAddress);
      // add the contract to onRampAddresses if it carries a valid dstChainSelector
      if (dstChainSelector > 0) {
        s_onRampAddresses.set(dstChainSelector, targetAddress);
      }
      // else if is redundant as this is the only corner case left, maintaining it for legibility
    } else if (targetAddress == address(0) && dstChainSelector > 0) {
      s_onRampAddresses.remove(dstChainSelector);
    }
  }

  /// @notice Delete an address from the watchlist and sets the target to inactive
  /// @param targetAddress the address to be deleted
  function removeFromWatchList(address targetAddress) public onlyAdminOrExecutor returns (bool) {
    if (s_watchList.remove(targetAddress)) {
      delete s_targets[targetAddress];
      return true;
    }
    return false;
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
    uint256 numTargets = s_watchList.length();
    uint256 idx = uint256(blockhash(block.number - (block.number % s_upkeepInterval) - 1)) % numTargets;
    uint256 numToCheck = numTargets < maxCheck ? numTargets : maxCheck;
    uint256 numFound = 0;
    uint256 minWaitPeriod = s_minWaitPeriodSeconds;
    address[] memory targetsToFund = new address[](maxPerform);
    MonitoredAddress memory contractToFund;
    address targetAddress;
    for (
      uint256 numChecked = 0;
      numChecked < numToCheck;
      (idx, numChecked) = ((idx + 1) % numTargets, numChecked + 1)
    ) {
      targetAddress = s_watchList.at(idx);
      contractToFund = s_targets[targetAddress];
      (bool fundingNeeded, ) = _needsFunding(
        targetAddress,
        contractToFund.lastTopUpTimestamp + minWaitPeriod,
        contractToFund.minBalance,
        contractToFund.isActive
      );
      if (fundingNeeded) {
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

  /// @notice tries to fund an array of target addresses, checking if they're underfunded in the process
  /// @param targetAddresses is an array of contract addresses to be funded in case they're underfunded
  function topUp(address[] memory targetAddresses) public whenNotPaused nonReentrant {
    MonitoredAddress memory contractToFund;
    uint256 minWaitPeriod = s_minWaitPeriodSeconds;
    uint256 localBalance = i_linkToken.balanceOf(address(this));
    for (uint256 idx = 0; idx < targetAddresses.length; idx++) {
      address targetAddress = targetAddresses[idx];
      contractToFund = s_targets[targetAddress];

      (bool fundingNeeded, address target) = _needsFunding(
        targetAddress,
        contractToFund.lastTopUpTimestamp + minWaitPeriod,
        contractToFund.minBalance,
        contractToFund.isActive
      );
      if (localBalance >= contractToFund.topUpAmount && fundingNeeded) {
        bool success = i_linkToken.transfer(target, contractToFund.topUpAmount);
        if (success) {
          localBalance -= contractToFund.topUpAmount;
          s_targets[targetAddress].lastTopUpTimestamp = uint56(block.timestamp);
          emit TopUpSucceeded(target, contractToFund.topUpAmount);
        } else {
          emit TopUpFailed(targetAddress);
        }
      } else {
        emit TopUpBlocked(targetAddress);
      }
    }
  }

  /// @notice checks the target (could be direct target or IAggregatorProxy), and determines
  /// if it is eligible for funding
  /// @param targetAddress the target to check
  /// @param minBalance minimum balance required for the target
  /// @param minWaitPeriodPassed the minimum wait period (target lastTopUpTimestamp + minWaitPeriod)
  /// @return bool whether the target needs funding or not
  /// @return address the address to fund. for DF, this is the aggregator address behind the proxy address.
  ///         for other products, it's the original target address
  function _needsFunding(
    address targetAddress,
    uint256 minWaitPeriodPassed,
    uint256 minBalance,
    bool contractIsActive
  ) private view returns (bool, address) {
    // Explicitly check if the targetAddress is the zero address
    // or if it's not a contract. In both cases return with false,
    // to prevent target.linkAvailableForPayment from running,
    // which would revert the operation.
    if (targetAddress == address(0) || targetAddress.code.length == 0) {
      return (false, address(0));
    }
    ILinkAvailable target;
    IAggregatorProxy proxy = IAggregatorProxy(targetAddress);
    try proxy.aggregator() returns (address aggregatorAddress) {
      // proxy.aggregator() can return a 0 address if the address is not an aggregator
      if (aggregatorAddress == address(0)) return (false, address(0));
      target = ILinkAvailable(aggregatorAddress);
    } catch {
      target = ILinkAvailable(targetAddress);
    }
    try target.linkAvailableForPayment() returns (int256 balance) {
      if (balance < int256(minBalance) && minWaitPeriodPassed <= block.timestamp && contractIsActive) {
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
    if (needsFunding.length == 0) {
      return (false, "");
    }
    uint96 total_batch_balance;
    for (uint256 idx = 0; idx < needsFunding.length; idx++) {
      address targetAddress = needsFunding[idx];
      total_batch_balance = total_batch_balance + s_targets[targetAddress].topUpAmount;
    }
    if (i_linkToken.balanceOf(address(this)) >= total_batch_balance) {
      return (true, abi.encode(needsFunding));
    }
    return (false, "");
  }

  /// @notice Called by the keeper to send funds to underfunded addresses.
  /// @param performData the abi encoded list of addresses to fund
  function performUpkeep(bytes calldata performData) external override {
    address[] memory needsFunding = abi.decode(performData, (address[]));
    topUp(needsFunding);
  }

  /// @notice Withdraws the contract balance in the LINK token.
  /// @param amount the amount of the LINK to withdraw
  /// @param payee the address to pay
  function withdraw(uint256 amount, address payable payee) external onlyAdminOrExecutor {
    if (payee == address(0)) revert InvalidAddress(payee);
    i_linkToken.transfer(payee, amount);
    emit FundsWithdrawn(amount, payee);
  }

  /// @notice Sets the minimum balance for the given target address
  function setMinBalance(address target, uint96 minBalance) external onlyRole(ADMIN_ROLE) {
    if (target == address(0)) revert InvalidAddress(target);
    if (minBalance == 0) revert InvalidMinBalance(minBalance);
    if (!s_targets[target].isActive) revert InvalidWatchList();
    uint256 oldBalance = s_targets[target].minBalance;
    s_targets[target].minBalance = minBalance;
    emit BalanceUpdated(target, oldBalance, minBalance);
  }

  /// @notice Sets the minimum balance for the given target address
  function setTopUpAmount(address target, uint96 topUpAmount) external onlyRole(ADMIN_ROLE) {
    if (target == address(0)) revert InvalidAddress(target);
    if (topUpAmount == 0) revert InvalidTopUpAmount(topUpAmount);
    if (!s_targets[target].isActive) revert InvalidWatchList();
    uint256 oldTopUpAmount = s_targets[target].topUpAmount;
    s_targets[target].topUpAmount = topUpAmount;
    emit BalanceUpdated(target, oldTopUpAmount, topUpAmount);
  }

  /// @notice Update s_maxPerform
  function setMaxPerform(uint16 maxPerform) public onlyRole(ADMIN_ROLE) {
    emit MaxPerformSet(s_maxPerform, maxPerform);
    s_maxPerform = maxPerform;
  }

  /// @notice Update s_maxCheck
  function setMaxCheck(uint16 maxCheck) public onlyRole(ADMIN_ROLE) {
    emit MaxCheckSet(s_maxCheck, maxCheck);
    s_maxCheck = maxCheck;
  }

  /// @notice Sets the minimum wait period (in seconds) for addresses between funding
  function setMinWaitPeriodSeconds(uint256 minWaitPeriodSeconds) public onlyRole(ADMIN_ROLE) {
    emit MinWaitPeriodSet(s_minWaitPeriodSeconds, minWaitPeriodSeconds);
    s_minWaitPeriodSeconds = minWaitPeriodSeconds;
  }

  /// @notice Update s_upkeepInterval
  function setUpkeepInterval(uint8 upkeepInterval) public onlyRole(ADMIN_ROLE) {
    if (upkeepInterval > 255) revert InvalidUpkeepInterval(upkeepInterval);
    emit UpkeepIntervalSet(s_upkeepInterval, upkeepInterval);
    s_upkeepInterval = upkeepInterval;
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

  /// @notice Gets upkeepInterval
  function getUpkeepInterval() external view returns (uint8) {
    return s_upkeepInterval;
  }

  /// @notice Gets the list of subscription ids being watched
  function getWatchList() external view returns (address[] memory) {
    return s_watchList.values();
  }

  /// @notice Gets the onRamp address with the specified dstChainSelector
  function getOnRampAddressAtChainSelector(uint64 dstChainSelector) external view returns (address) {
    if (dstChainSelector == 0) revert InvalidChainSelector();
    return s_onRampAddresses.get(dstChainSelector);
  }

  /// @notice Gets configuration information for an address on the watchlist
  function getAccountInfo(
    address targetAddress
  ) external view returns (bool isActive, uint96 minBalance, uint96 topUpAmount, uint56 lastTopUpTimestamp) {
    MonitoredAddress memory target = s_targets[targetAddress];
    return (target.isActive, target.minBalance, target.topUpAmount, target.lastTopUpTimestamp);
  }

  /// @dev Modifier to make a function callable only by executor role or the
  /// admin role.
  modifier onlyAdminOrExecutor() {
    address sender = _msgSender();
    if (!hasRole(ADMIN_ROLE, sender)) {
      _checkRole(EXECUTOR_ROLE, sender);
    }
    _;
  }

  modifier nonReentrant() {
    if (reentrancyGuard) revert ReentrantCall();
    reentrancyGuard = true;
    _;
    reentrancyGuard = false;
  }

  /// @notice Pause the contract, which prevents executing performUpkeep
  function pause() external onlyRole(ADMIN_ROLE) {
    _pause();
  }

  /// @notice Unpause the contract
  function unpause() external onlyRole(ADMIN_ROLE) {
    _unpause();
  }
}
