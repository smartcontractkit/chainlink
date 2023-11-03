// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

import "../../shared/access/ConfirmedOwner.sol";
import {IKeeperRegistryMaster} from "../interfaces/v2_1/IKeeperRegistryMaster.sol";
import {LinkTokenInterface} from "../../shared/interfaces/LinkTokenInterface.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

/// @title The UpkeepBalanceMonitor contract.
/// @notice A keeper-compatible contract that monitors and funds Chainlink Automation upkeeps.
contract UpkeepBalanceMonitor is ConfirmedOwner, Pausable {
  uint256 private constant MIN_GAS_FOR_TRANSFER = 55_000;
  LinkTokenInterface public immutable LINK_TOKEN;

  event FundsAdded(uint256 amountAdded, uint256 newBalance, address sender);
  event FundsWithdrawn(uint256 amountWithdrawn, address payee);
  event KeeperRegistryAddressUpdated(IKeeperRegistryMaster oldAddress, IKeeperRegistryMaster newAddress);
  event LinkTokenAddressUpdated(address oldAddress, address newAddress);
  event ConfigSet(uint96 minPercentage, uint96 targetPercentage);
  event OutOfGas(uint256 lastId);
  event TopUpFailed(uint256 indexed upkeepId);
  event TopUpSucceeded(uint256 indexed upkeepId);

  error DuplicateSubcriptionId(uint256 duplicate);
  error InvalidKeeperRegistryVersion();
  error InvalidWatchList();
  error MinPercentageTooLow();
  error OnlyKeeperRegistry();

  struct Config {
    uint96 minPercentage;
    uint96 targetPercentage;
  }

  IKeeperRegistryMaster private s_registry;
  uint256[] private s_watchList; // the watchlist on which subscriptions are stored
  Config private s_config;

  /// @param linkTokenAddress the Link token address
  /// @param keeperRegistryAddress the address of the keeper registry contract
  /// @param minPercentage the percentage of the min balance at which to trigger top ups
  /// @param targetPercentage the percentage of the min balance to target during top ups
  constructor(
    address linkTokenAddress,
    IKeeperRegistryMaster keeperRegistryAddress,
    uint96 minPercentage,
    uint96 targetPercentage
  ) ConfirmedOwner(msg.sender) {
    require(linkTokenAddress != address(0));
    if (keccak256(bytes(keeperRegistryAddress.typeAndVersion())) != keccak256(bytes("KeeperRegistry 2.1.0"))) {
      revert InvalidKeeperRegistryVersion();
    }
    LINK_TOKEN = LinkTokenInterface(linkTokenAddress);
    setKeeperRegistryAddress(keeperRegistryAddress);
    setConfig(minPercentage, targetPercentage);
    LinkTokenInterface(linkTokenAddress).approve(address(keeperRegistryAddress), type(uint256).max);
  }

  /// @notice Sets the list of upkeeps to watch and their funding parameters.
  /// @param watchlist the list of subscription ids to watch
  function setWatchList(uint256[] calldata watchlist) external onlyOwner {
    s_watchList = watchlist;
  }

  /// @notice Gets a list of upkeeps that are underfunded.
  /// @return list of upkeeps that are underfunded
  function getUnderfundedUpkeeps() public view returns (uint256[] memory) {
    uint256 numUpkeeps = s_watchList.length;
    uint256[] memory needsFunding = new uint256[](numUpkeeps);
    Config memory config = s_config;
    uint256 availableFunds = LINK_TOKEN.balanceOf(address(this));
    uint256 count;
    uint256 upkeepID;
    for (uint256 i = 0; i < numUpkeeps; i++) {
      upkeepID = s_watchList[i];
      uint96 upkeepBalance = s_registry.getBalance(upkeepID);
      uint96 minBalance = s_registry.getMinBalance(upkeepID);
      uint96 topUpThreshold = (minBalance * config.minPercentage) / 100; // TODO - uint96?
      uint96 topUpAmount = (minBalance * config.targetPercentage) / 100;
      if (upkeepBalance <= topUpThreshold && availableFunds >= topUpAmount) {
        needsFunding[count] = upkeepID;
        count++;
        availableFunds -= topUpAmount;
      }
    }
    if (count < numUpkeeps) {
      assembly {
        mstore(needsFunding, count)
      }
    }
    return needsFunding;
  }

  /// @notice Send funds to the upkeeps provided.
  /// @param needsFunding the list of upkeeps to fund
  function topUp(uint256[] memory needsFunding) public whenNotPaused {
    uint256 contractBalance = LINK_TOKEN.balanceOf(address(this));
    // Target memory target;
    // for (uint256 i = 0; i < needsFunding.length; i++) {
    //   target = s_targets[needsFunding[i]];
    //   uint96 upkeepBalance = s_registry.getBalance(needsFunding[i]);
    //   uint96 minUpkeepBalance = s_registry.getMinBalanceForUpkeep(needsFunding[i]);
    //   uint96 minBalanceWithBuffer = addBuffer(minUpkeepBalance, buffer);
    //   if (
    //     target.isActive &&
    //     target.lastTopUpTimestamp + minWaitPeriodSeconds <= block.timestamp &&
    //     (upkeepBalance < target.minBalanceJuels ||
    //       //upkeepBalance < minUpkeepBalance) &&
    //       upkeepBalance < minBalanceWithBuffer) &&
    //     contractBalance >= target.topUpAmountJuels
    //   ) {
    //     s_registry.addFunds(needsFunding[i], target.topUpAmountJuels);
    //     s_targets[needsFunding[i]].lastTopUpTimestamp = uint56(block.timestamp);
    //     contractBalance -= target.topUpAmountJuels;
    //     emit TopUpSucceeded(needsFunding[i]);
    //   }
    //   if (gasleft() < MIN_GAS_FOR_TRANSFER) {
    //     emit OutOfGas(i);
    //     return;
    //   }
    // }
  }

  /// @notice Gets list of upkeeps ids that are underfunded and returns a keeper-compatible payload.
  /// @return upkeepNeeded signals if upkeep is needed, performData is an abi encoded list of subscription ids that need funds
  function checkUpkeep(
    bytes calldata
  ) external view whenNotPaused returns (bool upkeepNeeded, bytes memory performData) {
    uint256[] memory needsFunding = getUnderfundedUpkeeps();
    upkeepNeeded = needsFunding.length > 0;
    performData = abi.encode(needsFunding);
    return (upkeepNeeded, performData);
  }

  /// @notice Called by the keeper to send funds to underfunded addresses.
  /// @param performData the abi encoded list of addresses to fund
  function performUpkeep(bytes calldata performData) external whenNotPaused {
    // if (msg.sender != address(s_registry)) revert OnlyKeeperRegistry();
    // TODO - forwarder contract
    uint256[] memory needsFunding = abi.decode(performData, (uint256[]));
    topUp(needsFunding);
  }

  /// @notice Withdraws the contract balance in LINK.
  /// @param amount the amount of LINK (in juels) to withdraw
  /// @param payee the address to pay
  function withdraw(uint256 amount, address payee) external onlyOwner {
    require(payee != address(0));
    LINK_TOKEN.transfer(payee, amount);
    emit FundsWithdrawn(amount, payee);
  }

  /// @notice Sets the keeper registry address.
  function setKeeperRegistryAddress(IKeeperRegistryMaster keeperRegistryAddress) public onlyOwner {
    require(address(keeperRegistryAddress) != address(0));
    s_registry = keeperRegistryAddress;
    emit KeeperRegistryAddressUpdated(s_registry, keeperRegistryAddress);
  }

  /// @notice Sets the contract config
  function setConfig(uint96 minPercentage, uint96 targetPercentage) public onlyOwner {
    if (minPercentage < 100) revert MinPercentageTooLow();
    s_config = Config({minPercentage: minPercentage, targetPercentage: targetPercentage});
    emit ConfigSet(minPercentage, targetPercentage);
  }

  /// @notice Gets the keeper registry address
  function getKeeperRegistryAddress() external view returns (address) {
    return address(s_registry);
  }

  /// @notice Gets the list of upkeeps ids being watched.
  function getWatchList() external view returns (uint256[] memory) {
    return s_watchList;
  }

  /// @notice Gets the contract config
  function getConfig() external view returns (Config memory) {
    return s_config;
  }

  /// @notice Pause the contract, which prevents executing performUpkeep.
  function pause() external onlyOwner {
    _pause();
  }

  /// @notice Unpause the contract.
  function unpause() external onlyOwner {
    _unpause();
  }
}
