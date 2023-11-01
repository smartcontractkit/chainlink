// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import "../../../ConfirmedOwner.sol";
import {IKeeperRegistryMaster} from "../2_1/interfaces/IKeeperRegistryMaster.sol";
import {LinkTokenInterface} from "../../../interfaces/LinkTokenInterface.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

/**
 * @title The UpkeepBalanceMonitor contract.
 * @notice A keeper-compatible contract that monitors and funds Chainlink Automation upkeeps.
 */
contract UpkeepBalanceMonitor is ConfirmedOwner, Pausable {
  LinkTokenInterface public LINKTOKEN;
  IKeeperRegistryMaster public REGISTRY;

  uint256 private constant MIN_GAS_FOR_TRANSFER = 55_000;

  bytes4 fundSig = REGISTRY.addFunds.selector;

  event FundsAdded(uint256 amountAdded, uint256 newBalance, address sender);
  event FundsWithdrawn(uint256 amountWithdrawn, address payee);
  event TopUpSucceeded(uint256 indexed upkeepId);
  event TopUpFailed(uint256 indexed upkeepId);
  event KeeperRegistryAddressUpdated(address oldAddress, address newAddress);
  event LinkTokenAddressUpdated(address oldAddress, address newAddress);
  event MinWaitPeriodUpdated(uint256 oldMinWaitPeriod, uint256 newMinWaitPeriod);
  event OutOfGas(uint256 lastId);

  error InvalidWatchList();
  error OnlyKeeperRegistry();
  error DuplicateSubcriptionId(uint256 duplicate);

  struct Target {
    bool isActive;
    uint96 minBalanceJuels;
    uint96 topUpAmountJuels;
    uint56 lastTopUpTimestamp;
  }

  address public s_keeperRegistryAddress; // the address of the keeper registry
  uint256 public s_minWaitPeriodSeconds; // minimum time to wait between top-ups
  uint256[] public s_watchList; // the watchlist on which subscriptions are stored
  mapping(uint256 => Target) internal s_targets;

  /**
   * @param linkTokenAddress the Link token address
   * @param keeperRegistryAddress the address of the keeper registry contract
   * @param minWaitPeriodSeconds the minimum wait period for addresses between funding
   */
  constructor(
    address linkTokenAddress,
    address keeperRegistryAddress,
    uint256 minWaitPeriodSeconds
  ) ConfirmedOwner(msg.sender) {
    require(linkTokenAddress != address(0));
    LINKTOKEN = LinkTokenInterface(linkTokenAddress);
    setKeeperRegistryAddress(keeperRegistryAddress); // 0xE16Df59B887e3Caa439E0b29B42bA2e7976FD8b2
    setMinWaitPeriodSeconds(minWaitPeriodSeconds); //0
    LINKTOKEN.approve(keeperRegistryAddress, type(uint256).max);
  }

  /**
   * @notice Sets the list of upkeeps to watch and their funding parameters.
   * @param upkeepIDs the list of subscription ids to watch
   * @param minBalancesJuels the minimum balances for each upkeep
   * @param topUpAmountsJuels the amount to top up each upkeep
   */
  function setWatchList(
    uint256[] calldata upkeepIDs,
    uint96[] calldata minBalancesJuels,
    uint96[] calldata topUpAmountsJuels
  ) external onlyOwner {
    if (upkeepIDs.length != minBalancesJuels.length || upkeepIDs.length != topUpAmountsJuels.length) {
      revert InvalidWatchList();
    }
    uint256[] memory oldWatchList = s_watchList;
    for (uint256 idx = 0; idx < oldWatchList.length; idx++) {
      s_targets[oldWatchList[idx]].isActive = false;
    }
    for (uint256 idx = 0; idx < upkeepIDs.length; idx++) {
      if (s_targets[upkeepIDs[idx]].isActive) {
        revert DuplicateSubcriptionId(upkeepIDs[idx]);
      }
      if (upkeepIDs[idx] == 0) {
        revert InvalidWatchList();
      }
      if (topUpAmountsJuels[idx] <= minBalancesJuels[idx]) {
        revert InvalidWatchList();
      }
      s_targets[upkeepIDs[idx]] = Target({
        isActive: true,
        minBalanceJuels: minBalancesJuels[idx],
        topUpAmountJuels: topUpAmountsJuels[idx],
        lastTopUpTimestamp: 0
      });
    }
    s_watchList = upkeepIDs;
  }

  /**
   * @notice Gets a list of upkeeps that are underfunded.
   * @return list of upkeeps that are underfunded
   */
  function getUnderfundedUpkeeps() public view returns (uint256[] memory) {
    uint256[] memory watchList = s_watchList;
    uint256[] memory needsFunding = new uint256[](watchList.length);
    uint256 count = 0;
    uint256 minWaitPeriod = s_minWaitPeriodSeconds;
    uint256 contractBalance = LINKTOKEN.balanceOf(address(this));
    Target memory target;
    for (uint256 idx = 0; idx < watchList.length; idx++) {
      target = s_targets[watchList[idx]];
      //( , , , uint96 upkeepBalance, , , ,) = REGISTRY.getUpkeep(watchList[idx]); <- for 1.2
      UpkeepInfo memory upkeepInfo; //2.0
      upkeepInfo = REGISTRY.getUpkeep(watchList[idx]); //2.0
      uint96 upkeepBalance = upkeepInfo.balance; //2.0
      uint96 minUpkeepBalance = REGISTRY.getMinBalanceForUpkeep(watchList[idx]);
      uint96 minBalanceWithBuffer = getBalanceWithBuffer(minUpkeepBalance);
      if (
        target.lastTopUpTimestamp + minWaitPeriod <= block.timestamp &&
        contractBalance >= target.topUpAmountJuels &&
        (upkeepBalance < target.minBalanceJuels ||
          //upkeepBalance < minUpkeepBalance)
          upkeepBalance < minBalanceWithBuffer)
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
   * @notice Send funds to the upkeeps provided.
   * @param needsFunding the list of upkeeps to fund
   */
  function topUp(uint256[] memory needsFunding) public whenNotPaused {
    uint256 minWaitPeriodSeconds = s_minWaitPeriodSeconds;
    uint256 contractBalance = LINKTOKEN.balanceOf(address(this));
    Target memory target;
    for (uint256 idx = 0; idx < needsFunding.length; idx++) {
      target = s_targets[needsFunding[idx]];
      //( , , , uint96 upkeepBalance, , , ,) = REGISTRY.getUpkeep(needsFunding[idx]); <- for 1.2
      UpkeepInfo memory upkeepInfo; //2.0
      upkeepInfo = REGISTRY.getUpkeep(needsFunding[idx]); //2.0
      uint96 upkeepBalance = upkeepInfo.balance; //2.0
      uint96 minUpkeepBalance = REGISTRY.getMinBalanceForUpkeep(needsFunding[idx]);
      uint96 minBalanceWithBuffer = getBalanceWithBuffer(minUpkeepBalance);
      if (
        target.isActive &&
        target.lastTopUpTimestamp + minWaitPeriodSeconds <= block.timestamp &&
        (upkeepBalance < target.minBalanceJuels ||
          //upkeepBalance < minUpkeepBalance) &&
          upkeepBalance < minBalanceWithBuffer) &&
        contractBalance >= target.topUpAmountJuels
      ) {
        REGISTRY.addFunds(needsFunding[idx], target.topUpAmountJuels);
        s_targets[needsFunding[idx]].lastTopUpTimestamp = uint56(block.timestamp);
        contractBalance -= target.topUpAmountJuels;
        emit TopUpSucceeded(needsFunding[idx]);
      }
      if (gasleft() < MIN_GAS_FOR_TRANSFER) {
        emit OutOfGas(idx);
        return;
      }
    }
  }

  /**
   * @notice Gets list of upkeeps ids that are underfunded and returns a keeper-compatible payload.
   * @return upkeepNeeded signals if upkeep is needed, performData is an abi encoded list of subscription ids that need funds
   */
  function checkUpkeep(
    bytes calldata
  ) external view whenNotPaused returns (bool upkeepNeeded, bytes memory performData) {
    uint256[] memory needsFunding = getUnderfundedUpkeeps();
    upkeepNeeded = needsFunding.length > 0;
    performData = abi.encode(needsFunding);
    return (upkeepNeeded, performData);
  }

  /**
   * @notice Called by the keeper to send funds to underfunded addresses.
   * @param performData the abi encoded list of addresses to fund
   */
  function performUpkeep(bytes calldata performData) external onlyKeeperRegistry whenNotPaused {
    uint256[] memory needsFunding = abi.decode(performData, (uint256[]));
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
   * @notice Sets the keeper registry address.
   */
  function setKeeperRegistryAddress(address keeperRegistryAddress) public onlyOwner {
    require(keeperRegistryAddress != address(0));
    emit KeeperRegistryAddressUpdated(s_keeperRegistryAddress, keeperRegistryAddress);
    s_keeperRegistryAddress = keeperRegistryAddress;
    REGISTRY = IKeeperRegistryMaster(keeperRegistryAddress);
  }

  /**
   * @notice Sets the minimum wait period (in seconds) for upkeep ids between funding.
   */
  function setMinWaitPeriodSeconds(uint256 period) public onlyOwner {
    emit MinWaitPeriodUpdated(s_minWaitPeriodSeconds, period);
    s_minWaitPeriodSeconds = period;
  }

  /**
   * @notice Gets configuration information for a upkeep on the watchlist.
   */
  function getUpkeepInfo(
    uint256 upkeepId
  ) external view returns (bool isActive, uint96 minBalanceJuels, uint96 topUpAmountJuels, uint56 lastTopUpTimestamp) {
    Target memory target = s_targets[upkeepId];
    return (target.isActive, target.minBalanceJuels, target.topUpAmountJuels, target.lastTopUpTimestamp);
  }

  /**
   * @notice Gets the list of upkeeps ids being watched.
   */
  function getWatchList() external view returns (uint256[] memory) {
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

  /**
   * @notice Called to add buffer to minimum balance of upkeeps
   * @param num the current minimum balance
   */
  function getBalanceWithBuffer(uint96 num) internal pure returns (uint96) {
    uint96 buffer = 20;
    uint96 result = uint96((uint256(num) * (100 + buffer)) / 100); // convert to uint256 to prevent overflow
    return result;
  }

  modifier onlyKeeperRegistry() {
    if (msg.sender != s_keeperRegistryAddress) {
      revert OnlyKeeperRegistry();
    }
    _;
  }
}
