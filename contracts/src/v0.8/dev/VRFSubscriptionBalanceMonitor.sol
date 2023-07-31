// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import "../shared/access/ConfirmedOwner.sol";
import "../automation/interfaces/KeeperCompatibleInterface.sol";
import "../interfaces/VRFCoordinatorV2Interface.sol";
import "../shared/interfaces/LinkTokenInterface.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

/**
 * @title The VRFSubscriptionBalanceMonitor contract.
 * @notice A keeper-compatible contract that monitors and funds VRF subscriptions.
 */
contract VRFSubscriptionBalanceMonitor is ConfirmedOwner, Pausable, KeeperCompatibleInterface {
  VRFCoordinatorV2Interface public COORDINATOR;
  LinkTokenInterface public LINKTOKEN;

  uint256 private constant MIN_GAS_FOR_TRANSFER = 55_000;

  event FundsAdded(uint256 amountAdded, uint256 newBalance, address sender);
  event FundsWithdrawn(uint256 amountWithdrawn, address payee);
  event TopUpSucceeded(uint64 indexed subscriptionId);
  event TopUpFailed(uint64 indexed subscriptionId);
  event KeeperRegistryAddressUpdated(address oldAddress, address newAddress);
  event VRFCoordinatorV2AddressUpdated(address oldAddress, address newAddress);
  event LinkTokenAddressUpdated(address oldAddress, address newAddress);
  event MinWaitPeriodUpdated(uint256 oldMinWaitPeriod, uint256 newMinWaitPeriod);
  event OutOfGas(uint256 lastId);

  error InvalidWatchList();
  error OnlyKeeperRegistry();
  error DuplicateSubcriptionId(uint64 duplicate);

  struct Target {
    bool isActive;
    uint96 minBalanceJuels;
    uint96 topUpAmountJuels;
    uint56 lastTopUpTimestamp;
  }

  address public s_keeperRegistryAddress; // the address of the keeper registry
  uint256 public s_minWaitPeriodSeconds; // minimum time to wait between top-ups
  uint64[] public s_watchList; // the watchlist on which subscriptions are stored
  mapping(uint64 => Target) internal s_targets;

  /**
   * @param linkTokenAddress the Link token address
   * @param coordinatorAddress the address of the vrf coordinator contract
   * @param keeperRegistryAddress the address of the keeper registry contract
   * @param minWaitPeriodSeconds the minimum wait period for addresses between funding
   */
  constructor(
    address linkTokenAddress,
    address coordinatorAddress,
    address keeperRegistryAddress,
    uint256 minWaitPeriodSeconds
  ) ConfirmedOwner(msg.sender) {
    setLinkTokenAddress(linkTokenAddress);
    setVRFCoordinatorV2Address(coordinatorAddress);
    setKeeperRegistryAddress(keeperRegistryAddress);
    setMinWaitPeriodSeconds(minWaitPeriodSeconds);
  }

  /**
   * @notice Sets the list of subscriptions to watch and their funding parameters.
   * @param subscriptionIds the list of subscription ids to watch
   * @param minBalancesJuels the minimum balances for each subscription
   * @param topUpAmountsJuels the amount to top up each subscription
   */
  function setWatchList(
    uint64[] calldata subscriptionIds,
    uint96[] calldata minBalancesJuels,
    uint96[] calldata topUpAmountsJuels
  ) external onlyOwner {
    if (subscriptionIds.length != minBalancesJuels.length || subscriptionIds.length != topUpAmountsJuels.length) {
      revert InvalidWatchList();
    }
    uint64[] memory oldWatchList = s_watchList;
    for (uint256 idx = 0; idx < oldWatchList.length; idx++) {
      s_targets[oldWatchList[idx]].isActive = false;
    }
    for (uint256 idx = 0; idx < subscriptionIds.length; idx++) {
      if (s_targets[subscriptionIds[idx]].isActive) {
        revert DuplicateSubcriptionId(subscriptionIds[idx]);
      }
      if (subscriptionIds[idx] == 0) {
        revert InvalidWatchList();
      }
      if (topUpAmountsJuels[idx] <= minBalancesJuels[idx]) {
        revert InvalidWatchList();
      }
      s_targets[subscriptionIds[idx]] = Target({
        isActive: true,
        minBalanceJuels: minBalancesJuels[idx],
        topUpAmountJuels: topUpAmountsJuels[idx],
        lastTopUpTimestamp: 0
      });
    }
    s_watchList = subscriptionIds;
  }

  /**
   * @notice Gets a list of subscriptions that are underfunded.
   * @return list of subscriptions that are underfunded
   */
  function getUnderfundedSubscriptions() public view returns (uint64[] memory) {
    uint64[] memory watchList = s_watchList;
    uint64[] memory needsFunding = new uint64[](watchList.length);
    uint256 count = 0;
    uint256 minWaitPeriod = s_minWaitPeriodSeconds;
    uint256 contractBalance = LINKTOKEN.balanceOf(address(this));
    Target memory target;
    for (uint256 idx = 0; idx < watchList.length; idx++) {
      target = s_targets[watchList[idx]];
      (uint96 subscriptionBalance, , , ) = COORDINATOR.getSubscription(watchList[idx]);
      if (
        target.lastTopUpTimestamp + minWaitPeriod <= block.timestamp &&
        contractBalance >= target.topUpAmountJuels &&
        subscriptionBalance < target.minBalanceJuels
      ) {
        needsFunding[count] = watchList[idx];
        count++;
        contractBalance -= target.topUpAmountJuels;
      }
    }
    if (count < watchList.length) {
      assembly {
        mstore(needsFunding, count)
      }
    }
    return needsFunding;
  }

  /**
   * @notice Send funds to the subscriptions provided.
   * @param needsFunding the list of subscriptions to fund
   */
  function topUp(uint64[] memory needsFunding) public whenNotPaused {
    uint256 minWaitPeriodSeconds = s_minWaitPeriodSeconds;
    uint256 contractBalance = LINKTOKEN.balanceOf(address(this));
    Target memory target;
    for (uint256 idx = 0; idx < needsFunding.length; idx++) {
      target = s_targets[needsFunding[idx]];
      (uint96 subscriptionBalance, , , ) = COORDINATOR.getSubscription(needsFunding[idx]);
      if (
        target.isActive &&
        target.lastTopUpTimestamp + minWaitPeriodSeconds <= block.timestamp &&
        subscriptionBalance < target.minBalanceJuels &&
        contractBalance >= target.topUpAmountJuels
      ) {
        bool success = LINKTOKEN.transferAndCall(
          address(COORDINATOR),
          target.topUpAmountJuels,
          abi.encode(needsFunding[idx])
        );
        if (success) {
          s_targets[needsFunding[idx]].lastTopUpTimestamp = uint56(block.timestamp);
          contractBalance -= target.topUpAmountJuels;
          emit TopUpSucceeded(needsFunding[idx]);
        } else {
          emit TopUpFailed(needsFunding[idx]);
        }
      }
      if (gasleft() < MIN_GAS_FOR_TRANSFER) {
        emit OutOfGas(idx);
        return;
      }
    }
  }

  /**
   * @notice Gets list of subscription ids that are underfunded and returns a keeper-compatible payload.
   * @return upkeepNeeded signals if upkeep is needed, performData is an abi encoded list of subscription ids that need funds
   */
  function checkUpkeep(
    bytes calldata
  ) external view override whenNotPaused returns (bool upkeepNeeded, bytes memory performData) {
    uint64[] memory needsFunding = getUnderfundedSubscriptions();
    upkeepNeeded = needsFunding.length > 0;
    performData = abi.encode(needsFunding);
    return (upkeepNeeded, performData);
  }

  /**
   * @notice Called by the keeper to send funds to underfunded addresses.
   * @param performData the abi encoded list of addresses to fund
   */
  function performUpkeep(bytes calldata performData) external override onlyKeeperRegistry whenNotPaused {
    uint64[] memory needsFunding = abi.decode(performData, (uint64[]));
    topUp(needsFunding);
  }

  /**
   * @notice Withdraws the contract balance in LINK.
   * @param amount the amount of LINK (in juels) to withdraw
   * @param payee the address to pay
   */
  function withdraw(uint256 amount, address payable payee) external onlyOwner {
    require(payee != address(0));
    emit FundsWithdrawn(amount, payee);
    LINKTOKEN.transfer(payee, amount);
  }

  /**
   * @notice Sets the LINK token address.
   */
  function setLinkTokenAddress(address linkTokenAddress) public onlyOwner {
    require(linkTokenAddress != address(0));
    emit LinkTokenAddressUpdated(address(LINKTOKEN), linkTokenAddress);
    LINKTOKEN = LinkTokenInterface(linkTokenAddress);
  }

  /**
   * @notice Sets the VRF coordinator address.
   */
  function setVRFCoordinatorV2Address(address coordinatorAddress) public onlyOwner {
    require(coordinatorAddress != address(0));
    emit VRFCoordinatorV2AddressUpdated(address(COORDINATOR), coordinatorAddress);
    COORDINATOR = VRFCoordinatorV2Interface(coordinatorAddress);
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
   * @notice Gets configuration information for a subscription on the watchlist.
   */
  function getSubscriptionInfo(
    uint64 subscriptionId
  ) external view returns (bool isActive, uint96 minBalanceJuels, uint96 topUpAmountJuels, uint56 lastTopUpTimestamp) {
    Target memory target = s_targets[subscriptionId];
    return (target.isActive, target.minBalanceJuels, target.topUpAmountJuels, target.lastTopUpTimestamp);
  }

  /**
   * @notice Gets the list of subscription ids being watched.
   */
  function getWatchList() external view returns (uint64[] memory) {
    return s_watchList;
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
