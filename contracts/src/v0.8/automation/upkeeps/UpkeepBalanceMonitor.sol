// SPDX-License-Identifier: MIT

pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {IAutomationRegistryConsumer} from "../interfaces/IAutomationRegistryConsumer.sol";
import {LinkTokenInterface} from "../../shared/interfaces/LinkTokenInterface.sol";
import {Pausable} from "@openzeppelin/contracts/security/Pausable.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

/// @title The UpkeepBalanceMonitor contract
/// @notice A keeper-compatible contract that monitors and funds Chainlink Automation upkeeps.
contract UpkeepBalanceMonitor is ConfirmedOwner, Pausable {
  using EnumerableSet for EnumerableSet.AddressSet;

  event ConfigSet(Config config);
  event ForwarderSet(address forwarderAddress);
  event FundsWithdrawn(uint256 amountWithdrawn, address payee);
  event TopUpFailed(uint256 indexed upkeepId);
  event TopUpSucceeded(uint256 indexed upkeepId, uint96 amount);
  event WatchListSet(address registryAddress);

  error InvalidConfig();
  error InvalidTopUpData();
  error OnlyForwarderOrOwner();

  /// @member maxBatchSize is the maximum number of upkeeps to fund in a single transaction
  /// @member minPercentage is the percentage of the upkeep's minBalance at which top-up occurs
  /// @member targetPercentage is the percentage of the upkeep's minBalance to top-up to
  /// @member maxTopUpAmount is the maximum amount of LINK to top-up an upkeep with
  struct Config {
    uint8 maxBatchSize;
    uint24 minPercentage;
    uint24 targetPercentage;
    uint96 maxTopUpAmount;
  }

  // ================================================================
  // |                           STORAGE                            |
  // ================================================================

  LinkTokenInterface private immutable LINK_TOKEN;

  mapping(address => uint256[]) s_registryWatchLists;
  EnumerableSet.AddressSet s_registries;
  Config private s_config;
  address private s_forwarderAddress;

  // ================================================================
  // |                         CONSTRUCTOR                          |
  // ================================================================

  /// @param linkToken the Link token address
  /// @param config the initial config for the contract
  constructor(LinkTokenInterface linkToken, Config memory config) ConfirmedOwner(msg.sender) {
    require(address(linkToken) != address(0));
    LINK_TOKEN = linkToken;
    setConfig(config);
  }

  // ================================================================
  // |                      CORE FUNCTIONALITY                      |
  // ================================================================

  /// @notice Gets a list of upkeeps that are underfunded
  /// @return needsFunding list of underfunded upkeepIDs
  /// @return registryAddresses list of registries that the upkeepIDs belong to
  /// @return topUpAmounts amount to top up each upkeep
  function getUnderfundedUpkeeps() public view returns (uint256[] memory, address[] memory, uint96[] memory) {
    Config memory config = s_config;
    uint256[] memory needsFunding = new uint256[](config.maxBatchSize);
    address[] memory registryAddresses = new address[](config.maxBatchSize);
    uint96[] memory topUpAmounts = new uint96[](config.maxBatchSize);
    uint256 availableFunds = LINK_TOKEN.balanceOf(address(this));
    uint256 count;
    for (uint256 i = 0; i < s_registries.length(); i++) {
      IAutomationRegistryConsumer registry = IAutomationRegistryConsumer(s_registries.at(i));
      for (uint256 j = 0; j < s_registryWatchLists[address(registry)].length; j++) {
        uint256 upkeepID = s_registryWatchLists[address(registry)][j];
        uint96 upkeepBalance = registry.getBalance(upkeepID);
        uint96 minBalance = registry.getMinBalance(upkeepID);
        uint96 topUpThreshold = (minBalance * config.minPercentage) / 100;
        uint96 topUpAmount = ((minBalance * config.targetPercentage) / 100) - upkeepBalance;
        if (topUpAmount > config.maxTopUpAmount) {
          topUpAmount = config.maxTopUpAmount;
        }
        if (upkeepBalance <= topUpThreshold && availableFunds >= topUpAmount) {
          needsFunding[count] = upkeepID;
          topUpAmounts[count] = topUpAmount;
          registryAddresses[count] = address(registry);
          count++;
          availableFunds -= topUpAmount;
        }
        if (count == config.maxBatchSize) {
          break;
        }
      }
      if (count == config.maxBatchSize) {
        break;
      }
    }
    if (count < config.maxBatchSize) {
      assembly {
        mstore(needsFunding, count)
        mstore(registryAddresses, count)
        mstore(topUpAmounts, count)
      }
    }
    return (needsFunding, registryAddresses, topUpAmounts);
  }

  /// @notice Called by the keeper/owner to send funds to underfunded upkeeps
  /// @param upkeepIDs the list of upkeep ids to fund
  /// @param registryAddresses the list of registries that the upkeepIDs belong to
  /// @param topUpAmounts the list of amounts to fund each upkeep with
  /// @dev We explicitly choose not to verify that input upkeepIDs are included in the watchlist. We also
  /// explicity permit any amount to be sent via topUpAmounts; it does not have to meet the criteria
  /// specified in getUnderfundedUpkeeps(). Here, we are relying on the security of automation's OCR to
  /// secure the output of getUnderfundedUpkeeps() as the input to topUp(), and we are treating the owner
  /// as a privileged user that can perform arbitrary top-ups to any upkeepID.
  function topUp(
    uint256[] memory upkeepIDs,
    address[] memory registryAddresses,
    uint96[] memory topUpAmounts
  ) public whenNotPaused {
    if (msg.sender != address(s_forwarderAddress) && msg.sender != owner()) revert OnlyForwarderOrOwner();
    if (upkeepIDs.length != registryAddresses.length || upkeepIDs.length != topUpAmounts.length)
      revert InvalidTopUpData();
    for (uint256 i = 0; i < upkeepIDs.length; i++) {
      try LINK_TOKEN.transferAndCall(registryAddresses[i], topUpAmounts[i], abi.encode(upkeepIDs[i])) returns (
        bool success
      ) {
        if (success) {
          emit TopUpSucceeded(upkeepIDs[i], topUpAmounts[i]);
          continue;
        }
      } catch {}
      emit TopUpFailed(upkeepIDs[i]);
    }
  }

  // ================================================================
  // |                    AUTOMATION COMPATIBLE                     |
  // ================================================================

  /// @notice Gets list of upkeeps ids that are underfunded and returns a keeper-compatible payload.
  /// @return upkeepNeeded signals if upkeep is needed, performData is an abi encoded list of subscription ids that need funds
  function checkUpkeep(bytes calldata) external view returns (bool upkeepNeeded, bytes memory performData) {
    (
      uint256[] memory needsFunding,
      address[] memory registryAddresses,
      uint96[] memory topUpAmounts
    ) = getUnderfundedUpkeeps();
    upkeepNeeded = needsFunding.length > 0;
    if (upkeepNeeded) {
      performData = abi.encode(needsFunding, registryAddresses, topUpAmounts);
    }
    return (upkeepNeeded, performData);
  }

  /// @notice Called by the keeper to send funds to underfunded addresses.
  /// @param performData the abi encoded list of addresses to fund
  function performUpkeep(bytes calldata performData) external {
    (uint256[] memory upkeepIDs, address[] memory registryAddresses, uint96[] memory topUpAmounts) = abi.decode(
      performData,
      (uint256[], address[], uint96[])
    );
    topUp(upkeepIDs, registryAddresses, topUpAmounts);
  }

  // ================================================================
  // |                            ADMIN                             |
  // ================================================================

  /// @notice Withdraws the contract balance in LINK.
  /// @param amount the amount of LINK (in juels) to withdraw
  /// @param payee the address to pay
  function withdraw(uint256 amount, address payee) external onlyOwner {
    require(payee != address(0));
    LINK_TOKEN.transfer(payee, amount);
    emit FundsWithdrawn(amount, payee);
  }

  /// @notice Pause the contract, which prevents executing performUpkeep.
  function pause() external onlyOwner {
    _pause();
  }

  /// @notice Unpause the contract.
  function unpause() external onlyOwner {
    _unpause();
  }

  // ================================================================
  // |                           SETTERS                            |
  // ================================================================

  /// @notice Sets the list of upkeeps to watch
  /// @param registryAddress the registry that this watchlist applies to
  /// @param watchlist the list of UpkeepIDs to watch
  function setWatchList(address registryAddress, uint256[] calldata watchlist) external onlyOwner {
    if (watchlist.length == 0) {
      s_registries.remove(registryAddress);
      delete s_registryWatchLists[registryAddress];
    } else {
      s_registries.add(registryAddress);
      s_registryWatchLists[registryAddress] = watchlist;
    }
    emit WatchListSet(registryAddress);
  }

  /// @notice Sets the contract config
  /// @param config the new config
  function setConfig(Config memory config) public onlyOwner {
    if (
      config.maxBatchSize == 0 ||
      config.minPercentage < 100 ||
      config.targetPercentage <= config.minPercentage ||
      config.maxTopUpAmount == 0
    ) {
      revert InvalidConfig();
    }
    s_config = config;
    emit ConfigSet(config);
  }

  /// @notice Sets the upkeep's forwarder contract
  /// @param forwarderAddress the new forwarder
  /// @dev this should only need to be called once, after registering the contract with the registry
  function setForwarder(address forwarderAddress) external onlyOwner {
    s_forwarderAddress = forwarderAddress;
    emit ForwarderSet(forwarderAddress);
  }

  // ================================================================
  // |                           GETTERS                            |
  // ================================================================

  /// @notice Gets the list of upkeeps ids being monitored
  function getWatchList() external view returns (address[] memory, uint256[][] memory) {
    address[] memory registryAddresses = s_registries.values();
    uint256[][] memory upkeepIDs = new uint256[][](registryAddresses.length);
    for (uint256 i = 0; i < registryAddresses.length; i++) {
      upkeepIDs[i] = s_registryWatchLists[registryAddresses[i]];
    }
    return (registryAddresses, upkeepIDs);
  }

  /// @notice Gets the contract config
  function getConfig() external view returns (Config memory) {
    return s_config;
  }

  /// @notice Gets the upkeep's forwarder contract
  function getForwarder() external view returns (address) {
    return s_forwarderAddress;
  }
}
