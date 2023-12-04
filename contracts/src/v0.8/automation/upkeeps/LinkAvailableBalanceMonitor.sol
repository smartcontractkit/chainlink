// SPDX-License-Identifier: MIT

pragma solidity 0.8.19;

import {AccessControl} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/access/AccessControl.sol";
import {AutomationCompatibleInterface} from "../interfaces/AutomationCompatibleInterface.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

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
contract LinkAvailableBalanceMonitor is AccessControl, AutomationCompatibleInterface {
  event BalanceUpdated(address indexed addr, uint256 oldBalance, uint256 newBalance);
  event FundsWithdrawn(uint256 amountWithdrawn, address payee);
  event UpkeepIntervalSet(uint256 oldUpkeepInterval, uint256 newUpkeepInterval);
  event MaxCheckSet(uint256 oldMaxCheck, uint256 newMaxCheck);
  event MaxPerformSet(uint256 oldMaxPerform, uint256 newMaxPerform);
  event MinWaitPeriodSet(uint256 s_minWaitPeriodSeconds, uint256 minWaitPeriodSeconds);
  event TopUpBlocked(address indexed topUpAddress);
  event TopUpFailed(address indexed recipient);
  event TopUpSucceeded(address indexed topUpAddress);
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
  error DuplicateAddress(address duplicate);

  struct MonitoredAddress {
    uint96 minBalance;
    uint96 topUpAmount;
    uint56 lastTopUpTimestamp;
    bool isActive;
  }

  bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
  bytes32 public constant EXECUTOR_ROLE = keccak256("EXECUTOR_ROLE");
  uint96 private constant DEFAULT_TOP_UP_AMOUNT = 9;
  uint96 private constant DEFAULT_MIN_BALANCE = 1;
  IERC20 private immutable LINK_TOKEN;

  uint256 private s_minWaitPeriodSeconds;
  uint16 private s_maxPerform;
  uint16 private s_maxCheck;
  uint8 private s_upkeepInterval;
  address[] private s_watchList;
  mapping(address targetAddress => MonitoredAddress targetProperties) internal s_targets;
  mapping(uint64 dstChainSelector => address onRamp) internal s_onRampAddresses;

  /// @param linkTokenAddress the LINK token address
  constructor(
    address admin,
    address linkTokenAddress,
    uint256 minWaitPeriodSeconds,
    uint16 maxPerform,
    uint16 maxCheck,
    uint8 upkeepInterval
  ) {
    if (linkTokenAddress == address(0)) revert InvalidLinkTokenAddress(linkTokenAddress);
    _setRoleAdmin(ADMIN_ROLE, ADMIN_ROLE);
    _setRoleAdmin(EXECUTOR_ROLE, ADMIN_ROLE);
    _setupRole(ADMIN_ROLE, admin);
    LINK_TOKEN = IERC20(linkTokenAddress);
    setMinWaitPeriodSeconds(minWaitPeriodSeconds);
    setMaxPerform(maxPerform);
    setMaxCheck(maxCheck);
    setUpkeepInterval(upkeepInterval);
  }

  /// @notice Grants an address an executor role
  /// @param executor address to grant executor role to
  function granExecutorRole(address executor) public onlyRole(ADMIN_ROLE) {
    if (executor == address(0)) revert InvalidAddress(executor);
    _setupRole(EXECUTOR_ROLE, executor);
  }

  /// @notice Revokes the executor role from an address
  /// @param executor address to revoke executor role from
  function revokeExecutorRole(address executor) public onlyRole(ADMIN_ROLE) {
    if (executor == address(0)) revert InvalidAddress(executor);
    _revokeRole(EXECUTOR_ROLE, executor);
  }

  /// @notice Sets the list of subscriptions to watch and their funding parameters
  /// @param addresses the list of target addresses to watch (could be direct target or IAggregatorProxy)
  /// @param minBalances the list of corresponding minBalance for the target address
  /// @param topUpAmounts the list of corresponding minTopUp for the target address
  function setWatchList(
    address[] calldata addresses,
    uint96[] calldata minBalances,
    uint96[] calldata topUpAmounts
  ) external onlyRoleOrAdminRole(EXECUTOR_ROLE) {
    if (addresses.length != minBalances.length || addresses.length != topUpAmounts.length) {
      revert InvalidWatchList();
    }
    for (uint256 idx = 0; idx < s_watchList.length; idx++) {
      delete s_targets[s_watchList[idx]];
    }
    for (uint256 idx = 0; idx < addresses.length; idx++) {
      address targetAddress = addresses[idx];
      if (s_targets[targetAddress].isActive) revert DuplicateAddress(addresses[idx]);
      if (addresses[idx] == address(0)) revert InvalidWatchList();
      if (topUpAmounts[idx] == 0) revert InvalidWatchList();
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

  /// @notice Adds a new address to the watchlist
  /// @param targetAddress the address to be added to the watchlist
  /// @param dstChainSelector carries a non-zero value in case the targetAddress is an onRamp, otherwise it carries a 0
  /// @dev this function has to be compatible with the event onRampSet(address, dstChainSelector) emitted by
  /// the CCIP router. Important detail to know is this event is also emitted when an onRamp is decomissioned,
  /// in which case it will carry the proper dstChainSelector along with the 0x0 address
  function addToWatchList(address targetAddress, uint64 dstChainSelector) public onlyRoleOrAdminRole(EXECUTOR_ROLE) {
    if (s_targets[targetAddress].isActive) revert DuplicateAddress(targetAddress);
    address oldAddress = s_onRampAddresses[dstChainSelector];
    // if targetAddress is an existing onRamp, there's a need of cleaning the previous onRamp associated to this dstChainSelector
    // there's no need to remove any other address that's not an onRamp
    if (dstChainSelector > 0 && bytes(abi.encodePacked(oldAddress)).length > 0) {
      removeFromWatchList(oldAddress);
    }
    // only add the new address if it's not 0x0
    if (targetAddress != address(0)) {
      s_onRampAddresses[dstChainSelector] = targetAddress;
      s_targets[targetAddress] = MonitoredAddress({
        isActive: true,
        minBalance: DEFAULT_MIN_BALANCE,
        topUpAmount: DEFAULT_TOP_UP_AMOUNT,
        lastTopUpTimestamp: 0
      });
      s_watchList.push(targetAddress);
    } else {
      // if the address is 0x0, it means the onRamp has ben decomissioned and has to be cleaned
      delete s_onRampAddresses[dstChainSelector];
    }
  }

  /// @notice Delete an address from the watchlist and sets the target to inactive
  /// @param targetAddress the address to be deleted
  function removeFromWatchList(address targetAddress) public onlyRoleOrAdminRole(EXECUTOR_ROLE) returns (bool) {
    s_targets[targetAddress].isActive = false;
    for (uint i; i < s_watchList.length; i++) {
      if (s_watchList[i] == targetAddress) {
        s_watchList[i] = s_watchList[s_watchList.length - 1];
        s_watchList.pop();
        return true;
      }
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
    uint256 numTargets = s_watchList.length;
    uint256 idx = uint256(blockhash(block.number - (block.number % s_upkeepInterval) - 1)) % numTargets;
    uint256 numToCheck = numTargets < maxCheck ? numTargets : maxCheck;
    uint256 numFound = 0;
    address[] memory targetsToFund = new address[](maxPerform);
    MonitoredAddress memory target;
    for (
      uint256 numChecked = 0;
      numChecked < numToCheck;
      (idx, numChecked) = ((idx + 1) % numTargets, numChecked + 1)
    ) {
      address targetAddress = s_watchList[idx];
      target = s_targets[targetAddress];
      if (_needsFunding(targetAddress, target.minBalance)) {
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

  function topUp(address[] memory targetAddresses) public {
    MonitoredAddress memory target;
    uint256 localBalance = LINK_TOKEN.balanceOf(address(this));
    for (uint256 idx = 0; idx < targetAddresses.length; idx++) {
      address targetAddress = targetAddresses[idx];
      target = s_targets[targetAddress];
      if (localBalance >= target.topUpAmount && _needsFunding(targetAddress, target.minBalance)) {
        bool success = LINK_TOKEN.transfer(targetAddress, target.topUpAmount);
        if (success) {
          localBalance -= target.topUpAmount;
          target.lastTopUpTimestamp = uint56(block.timestamp);
          emit TopUpSucceeded(targetAddress);
        } else {
          emit TopUpFailed(targetAddress);
        }
      } else {
        emit TopUpBlocked(targetAddress);
      }
    }
  }

  /// @notice checks the target (could be direct target or IAggregatorProxy), and determines
  /// if it is elligible for funding
  /// @param targetAddress the target to check
  /// @param minBalance minimum balance required for the target
  /// @return bool whether the target needs funding or not
  function _needsFunding(address targetAddress, uint256 minBalance) private view returns (bool) {
    // Explicitly check if the targetAddress is the zero address
    // or if it's not a contract. In both cases return with false,
    // to prevent target.linkAvailableForPayment from running,
    // which would revert the operation.
    if (targetAddress == address(0) || targetAddress.code.length == 0) {
      return false;
    }
    MonitoredAddress memory addressToCheck = s_targets[targetAddress];
    ILinkAvailable target;
    IAggregatorProxy proxy = IAggregatorProxy(targetAddress);
    try proxy.aggregator() returns (address aggregatorAddress) {
      if (aggregatorAddress == address(0)) return false;
      target = ILinkAvailable(aggregatorAddress);
    } catch {
      target = ILinkAvailable(targetAddress);
    }
    try target.linkAvailableForPayment() returns (int256 balance) {
      if (
        balance < int256(minBalance) &&
        addressToCheck.lastTopUpTimestamp + s_minWaitPeriodSeconds <= block.timestamp &&
        addressToCheck.isActive
      ) {
        return true;
      }
    } catch {}
    return false;
  }

  /// @notice Gets list of subscription ids that are underfunded and returns a keeper-compatible payload.
  /// @return upkeepNeeded signals if upkeep is needed
  /// @return performData is an abi encoded list of subscription ids that need funds
  function checkUpkeep(bytes calldata) external view override returns (bool upkeepNeeded, bytes memory performData) {
    address[] memory needsFunding = sampleUnderfundedAddresses();
    upkeepNeeded = needsFunding.length > 0;
    performData = abi.encode(needsFunding);
    return (upkeepNeeded, performData);
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
  function withdraw(uint256 amount, address payable payee) external onlyRoleOrAdminRole(EXECUTOR_ROLE) {
    if (payee == address(0)) revert InvalidAddress(payee);
    LINK_TOKEN.transfer(payee, amount);
    emit FundsWithdrawn(amount, payee);
  }

  /// @notice Sets the minimum balance for the given target address
  function setMinBalance(address target, uint96 minBalance) external onlyRoleOrAdminRole(EXECUTOR_ROLE) {
    if (target == address(0)) revert InvalidAddress(target);
    if (minBalance == 0) revert InvalidMinBalance(minBalance);
    if (!s_targets[target].isActive) revert InvalidWatchList();
    uint256 oldBalance = s_targets[target].minBalance;
    s_targets[target].minBalance = minBalance;
    emit BalanceUpdated(target, oldBalance, minBalance);
  }

  /// @notice Sets the minimum balance for the given target address
  function setTopUpAmount(address target, uint96 topUpAmount) external onlyRoleOrAdminRole(EXECUTOR_ROLE) {
    if (target == address(0)) revert InvalidAddress(target);
    if (topUpAmount == 0) revert InvalidTopUpAmount(topUpAmount);
    if (!s_targets[target].isActive) revert InvalidWatchList();
    uint256 oldTopUpAmount = s_targets[target].topUpAmount;
    s_targets[target].topUpAmount = topUpAmount;
    emit BalanceUpdated(target, oldTopUpAmount, topUpAmount);
  }

  /// @notice Update s_maxPerform
  function setMaxPerform(uint16 maxPerform) public onlyRoleOrAdminRole(EXECUTOR_ROLE) {
    s_maxPerform = maxPerform;
    emit MaxPerformSet(s_maxPerform, maxPerform);
  }

  /// @notice Update s_maxCheck
  function setMaxCheck(uint16 maxCheck) public onlyRoleOrAdminRole(EXECUTOR_ROLE) {
    s_maxCheck = maxCheck;
    emit MaxCheckSet(s_maxCheck, maxCheck);
  }

  /// @notice Sets the minimum wait period (in seconds) for addresses between funding
  function setMinWaitPeriodSeconds(uint256 minWaitPeriodSeconds) public onlyRoleOrAdminRole(EXECUTOR_ROLE) {
    s_minWaitPeriodSeconds = minWaitPeriodSeconds;
    emit MinWaitPeriodSet(s_minWaitPeriodSeconds, minWaitPeriodSeconds);
  }

  /// @notice Update s_upkeepInterval
  function setUpkeepInterval(uint8 upkeepInterval) public onlyRoleOrAdminRole(EXECUTOR_ROLE) {
    if (upkeepInterval > 255) revert InvalidUpkeepInterval(upkeepInterval);
    s_upkeepInterval = upkeepInterval;
    emit UpkeepIntervalSet(s_upkeepInterval, upkeepInterval);
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
    return s_watchList;
  }

  /// @notice Gets configuration information for an address on the watchlist
  function getAccountInfo(
    address targetAddress
  ) external view returns (bool isActive, uint256 minBalance, uint256 topUpAmount) {
    MonitoredAddress memory target = s_targets[targetAddress];
    return (target.isActive, target.minBalance, target.topUpAmount);
  }

  /// @dev Modifier to make a function callable only by a certain role or the
  /// admin role.
  modifier onlyRoleOrAdminRole(bytes32 role) {
    address sender = _msgSender();
    if (!hasRole(ADMIN_ROLE, sender)) {
      _checkRole(role, sender);
    }
    _;
  }
}
