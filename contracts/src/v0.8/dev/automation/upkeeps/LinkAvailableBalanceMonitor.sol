// SPDX-License-Identifier: MIT

pragma solidity ^0.8.4;

import "../../../ConfirmedOwner.sol";
import "../../../interfaces/automation/KeeperCompatibleInterface.sol";
import "../../../vendor/openzeppelin-solidity/v4.7.0/contracts/security/Pausable.sol";
import "../../../vendor/openzeppelin-solidity/v4.7.0/contracts/token/ERC20/IERC20.sol";
import "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";

interface IAggregatorProxy {
  function aggregator() external view returns (address);
}

interface ILinkAvailable {
  function linkAvailableForPayment() external view returns (int256 availableBalance);
}

/**
 * @title The LinkAvailableBalanceMonitor contract.
 * @notice A keeper-compatible contract that monitors target contracts for balance from a custom
 * function linkAvailableForPayment() and funds them with LINK if it falls below a defined
 * threshold. Also supports aggregator proxy contracts monitoring which require fetching the actual
 * target contract through a predefined interface.
 * @dev with 30 addresses as the MAX_PERFORM, the measured max gas usage of performUpkeep is around 2M
 * therefore, we recommend an upkeep gas limit of 3M (this has a 33% margin of safety). Although, nothing
 * prevents us from using 5M gas and increasing MAX_PERFORM, 30 seems like a reasonable batch size that
 * is probably plenty for most needs.
 * @dev with 140 addresses as the MAX_CHECK, the measured max gas usage of checkUpkeep is around 3.5M,
 * which is 30% below the 5M limit.
 * Note that testing conditions DO NOT match live chain gas usage, hence the margins. Change
 * at your own risk!!!
 * @dev some areas for improvement / acknowledgement of limitations:
 * * validate that all addresses conform to interface when adding them to the watchlist
 * * this is a "trusless" upkeep, meaning it does not trust the caller of performUpkeep;
     we could save a fair amount of gas and re-write this upkeep for use with Automation v2.0+,
     which has significantly different trust assumptions
 */
contract LinkAvailableBalanceMonitor is ConfirmedOwner, Pausable, KeeperCompatibleInterface {
  using EnumerableSet for EnumerableSet.AddressSet;

  event FundsWithdrawn(uint256 amountWithdrawn, address payee);
  event TopUpSucceeded(address indexed topUpAddress);
  event TopUpBlocked(address indexed topUpAddress);
  event WatchlistUpdated();

  error InvalidWatchList();
  error DuplicateAddress(address duplicate);

  uint256 public constant MAX_PERFORM = 30; // max number to addresses to top up in a single batch
  uint256 public constant MAX_CHECK = 140; // max number of upkeeps to check (need to fit in 5M gas limit)
  IERC20 public immutable LINK_TOKEN;

  EnumerableSet.AddressSet private s_watchList;
  uint256 private s_minBalance;
  uint256 private s_topUpAmount;

  /**
   * @param linkTokenAddress the LINK token address
   * @param minBalance the minimum balance an aggregator can have before initiating a top up
   * @param topUpAmount the amount of LINK to top up an aggregator with at once
   */
  constructor(address linkTokenAddress, uint256 minBalance, uint256 topUpAmount) ConfirmedOwner(msg.sender) {
    require(linkTokenAddress != address(0));
    require(minBalance > 0);
    require(topUpAmount > 0);
    LINK_TOKEN = IERC20(linkTokenAddress);
    s_minBalance = minBalance;
    s_topUpAmount = topUpAmount;
  }

  /**
   * @notice Sets the list of subscriptions to watch and their funding parameters
   * @param addresses the list of target addresses to watch (could be direct target or IAggregatorProxy)
   */
  function setWatchList(address[] calldata addresses) external onlyOwner {
    // first, remove all existing addresses from list
    for (uint256 idx = s_watchList.length(); idx > 0; idx--) {
      require(s_watchList.remove(s_watchList.at(idx - 1)));
    }
    // then set new addresses
    for (uint256 idx = 0; idx < addresses.length; idx++) {
      if (s_watchList.contains(addresses[idx])) {
        revert DuplicateAddress(addresses[idx]);
      }
      if (addresses[idx] == address(0)) {
        revert InvalidWatchList();
      }
      s_watchList.add(addresses[idx]);
    }
    emit WatchlistUpdated();
  }

  /**
   * @notice Adds addresses to the watchlist without overwriting existing members
   * @param addresses the list of target addresses to watch (could be direct target or IAggregatorProxy)
   */
  function addToWatchList(address[] calldata addresses) external onlyOwner {
    for (uint256 idx = 0; idx < addresses.length; idx++) {
      if (s_watchList.contains(addresses[idx])) {
        revert DuplicateAddress(addresses[idx]);
      }
      if (addresses[idx] == address(0)) {
        revert InvalidWatchList();
      }
      s_watchList.add(addresses[idx]);
    }
    emit WatchlistUpdated();
  }

  /**
   * @notice Gets a list of proxies that are underfunded, up to the MAX_PERFORM size
   * @dev the function starts at a random index in the list to avoid biasing the first
   * addresses in the list over latter ones.
   * @dev the function will check at most MAX_CHECK proxies in a single call
   * @dev the function returns a list with a max length of MAX_PERFORM
   * @return list of target addresses which are underfunded
   */
  function sampleUnderfundedAddresses() public view returns (address[] memory) {
    uint256 numTargets = s_watchList.length();
    uint256 numChecked = 0;
    uint256 numToCheck = MAX_CHECK;
    uint256 idx = uint256(blockhash(block.number - 1)) % numTargets; // start at random index, to distribute load
    numToCheck = numTargets < MAX_CHECK ? numTargets : MAX_CHECK;
    uint256 numFound = 0;
    address[] memory targetsToFund = new address[](MAX_PERFORM);
    for (; numChecked < numToCheck; (idx, numChecked) = ((idx + 1) % numTargets, numChecked + 1)) {
      address target = s_watchList.at(idx);
      (bool needsFunding, ) = _needsFunding(target);
      if (needsFunding) {
        targetsToFund[numFound] = target;
        numFound++;
        if (numFound == MAX_PERFORM) {
          break; // max number of addresses in batch reached
        }
      }
    }
    if (numFound != MAX_PERFORM) {
      assembly {
        mstore(targetsToFund, numFound) // resize array to number of valid targets
      }
    }
    return targetsToFund;
  }

  /**
   * @notice Send funds to the targets provided.
   * @param targetAddresses the list of targets to fund
   */
  function topUp(address[] memory targetAddresses) public whenNotPaused {
    uint256 topUpAmount = s_topUpAmount;
    uint256 stopIdx = targetAddresses.length;
    uint256 numCanFund = LINK_TOKEN.balanceOf(address(this)) / topUpAmount;
    stopIdx = numCanFund < stopIdx ? numCanFund : stopIdx;
    for (uint256 idx = 0; idx < stopIdx; idx++) {
      (bool needsFunding, address target) = _needsFunding(targetAddresses[idx]);
      if (!s_watchList.contains(targetAddresses[idx]) || !needsFunding) {
        emit TopUpBlocked(targetAddresses[idx]);
        continue;
      }
      LINK_TOKEN.transfer(target, topUpAmount);
      emit TopUpSucceeded(targetAddresses[idx]);
    }
  }

  /**
   * @notice Gets list of subscription ids that are underfunded and returns a keeper-compatible payload.
   * @return upkeepNeeded signals if upkeep is needed
   * @return performData is an abi encoded list of subscription ids that need funds
   */
  function checkUpkeep(
    bytes calldata
  ) external view override whenNotPaused returns (bool upkeepNeeded, bytes memory performData) {
    address[] memory needsFunding = sampleUnderfundedAddresses();
    uint256 numCanFund = LINK_TOKEN.balanceOf(address(this)) / s_topUpAmount;
    if (numCanFund < needsFunding.length) {
      assembly {
        mstore(needsFunding, numCanFund) // resize
      }
    }
    upkeepNeeded = needsFunding.length > 0;
    performData = abi.encode(needsFunding);
    return (upkeepNeeded, performData);
  }

  /**
   * @notice Called by the keeper to send funds to underfunded addresses.
   * @param performData the abi encoded list of addresses to fund
   */
  function performUpkeep(bytes calldata performData) external override {
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
    LINK_TOKEN.transfer(payee, amount);
    emit FundsWithdrawn(amount, payee);
  }

  /**
   * @notice Sets the top up amount
   */
  function setTopUpAmount(uint256 topUpAmount) external onlyOwner returns (uint256) {
    require(topUpAmount > 0);
    return s_topUpAmount = topUpAmount;
  }

  /**
   * @notice Sets the minimum balance
   */
  function setMinBalance(uint256 minBalance) external onlyOwner returns (uint256) {
    require(minBalance > 0);
    return s_minBalance = minBalance;
  }

  /**
   * @notice Pause the contract, which prevents executing performUpkeep
   */
  function pause() external onlyOwner {
    _pause();
  }

  /**
   * @notice Unpause the contract
   */
  function unpause() external onlyOwner {
    _unpause();
  }

  /**
   * @notice Gets the list of subscription ids being watched
   */
  function getWatchList() external view returns (address[] memory) {
    return s_watchList.values();
  }

  /**
   * @notice Gets the configured top up amount
   */
  function getTopUpAmount() external view returns (uint256) {
    return s_topUpAmount;
  }

  /**
   * @notice Gets the configured minimum balance
   */
  function getMinBalance() external view returns (uint256) {
    return s_minBalance;
  }

  /**
   * @notice checks the target (could be direct target or IAggregatorProxy), and determines
   * if it is elligible for funding
   * @param targetAddress the target to check
   * @return bool whether the target needs funding or not
   * @return address the address of the contract needing funding
   */
  function _needsFunding(address targetAddress) private view returns (bool, address) {
    ILinkAvailable target;
    IAggregatorProxy proxy = IAggregatorProxy(targetAddress);
    try proxy.aggregator() returns (address aggregatorAddress) {
      target = ILinkAvailable(aggregatorAddress);
    } catch {
      target = ILinkAvailable(targetAddress);
    }
    try target.linkAvailableForPayment() returns (int256 balance) {
      if (balance < 0 || uint256(balance) > s_minBalance) {
        return (false, address(0));
      }
    } catch {
      return (false, address(0));
    }
    return (true, address(target));
  }
}
