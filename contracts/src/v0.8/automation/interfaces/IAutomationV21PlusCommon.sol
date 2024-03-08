// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.4;

import {AutomationRegistryBase2_2} from "../v2_2/AutomationRegistryBase2_2.sol";

interface IAutomationV21PlusCommon {
  // registry events
  event AdminPrivilegeConfigSet(address indexed admin, bytes privilegeConfig);
  event CancelledUpkeepReport(uint256 indexed id, bytes trigger);
  event ConfigSet(
    uint32 previousConfigBlockNumber,
    bytes32 configDigest,
    uint64 configCount,
    address[] signers,
    address[] transmitters,
    uint8 f,
    bytes onchainConfig,
    uint64 offchainConfigVersion,
    bytes offchainConfig
  );
  event DedupKeyAdded(bytes32 indexed dedupKey);
  event InsufficientFundsUpkeepReport(uint256 indexed id, bytes trigger);
  event OwnerFundsWithdrawn(uint96 amount);
  event OwnershipTransferred(address indexed from, address indexed to);
  event OwnershipTransferRequested(address indexed from, address indexed to);
  event Paused(address account);
  event PayeeshipTransferred(address indexed transmitter, address indexed from, address indexed to);
  event PayeeshipTransferRequested(address indexed transmitter, address indexed from, address indexed to);
  event PayeesUpdated(address[] transmitters, address[] payees);
  event PaymentWithdrawn(address indexed transmitter, uint256 indexed amount, address indexed to, address payee);
  event ReorgedUpkeepReport(uint256 indexed id, bytes trigger);
  event StaleUpkeepReport(uint256 indexed id, bytes trigger);
  event Transmitted(bytes32 configDigest, uint32 epoch);
  event Unpaused(address account);

  // upkeep events
  event FundsAdded(uint256 indexed id, address indexed from, uint96 amount);
  event FundsWithdrawn(uint256 indexed id, uint256 amount, address to);
  event UpkeepAdminTransferred(uint256 indexed id, address indexed from, address indexed to);
  event UpkeepAdminTransferRequested(uint256 indexed id, address indexed from, address indexed to);
  event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight);
  event UpkeepCheckDataSet(uint256 indexed id, bytes newCheckData);
  event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit);
  event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination);
  event UpkeepOffchainConfigSet(uint256 indexed id, bytes offchainConfig);
  event UpkeepPaused(uint256 indexed id);
  event UpkeepPerformed(
    uint256 indexed id,
    bool indexed success,
    uint96 totalPayment,
    uint256 gasUsed,
    uint256 gasOverhead,
    bytes trigger
  );
  event UpkeepPrivilegeConfigSet(uint256 indexed id, bytes privilegeConfig);
  event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom);
  event UpkeepRegistered(uint256 indexed id, uint32 performGas, address admin);
  event UpkeepTriggerConfigSet(uint256 indexed id, bytes triggerConfig);
  event UpkeepUnpaused(uint256 indexed id);

  /**
   * @notice structure of trigger for log triggers
   */
  struct LogTriggerConfig {
    address contractAddress;
    uint8 filterSelector; // denotes which topics apply to filter ex 000, 101, 111...only last 3 bits apply
    bytes32 topic0;
    bytes32 topic1;
    bytes32 topic2;
    bytes32 topic3;
  }

  /// @dev Report transmitted by OCR to transmit function
  struct Report {
    uint256 fastGasWei;
    uint256 linkNative;
    uint256[] upkeepIds;
    uint256[] gasLimits;
    bytes[] triggers;
    bytes[] performDatas;
  }

  /**
   * @notice all information about an upkeep
   * @dev only used in return values
   * @dev this will likely be deprecated in a future version of the registry
   * @member target the contract which needs to be serviced
   * @member performGas the gas limit of upkeep execution
   * @member checkData the checkData bytes for this upkeep
   * @member balance the balance of this upkeep
   * @member admin for this upkeep
   * @member maxValidBlocknumber until which block this upkeep is valid
   * @member lastPerformedBlockNumber the last block number when this upkeep was performed
   * @member amountSpent the amount this upkeep has spent
   * @member paused if this upkeep has been paused
   * @member offchainConfig the off-chain config of this upkeep
   */
  struct UpkeepInfo {
    address target;
    uint32 performGas;
    bytes checkData;
    uint96 balance;
    address admin;
    uint64 maxValidBlocknumber;
    uint32 lastPerformedBlockNumber;
    uint96 amountSpent;
    bool paused;
    bytes offchainConfig;
  }

  /**
   * @notice the trigger structure conditional trigger type
   */
  struct ConditionalTrigger {
    uint32 blockNum;
    bytes32 blockHash;
  }

  /**
   * @notice the trigger structure of log upkeeps
   * @dev NOTE that blockNum / blockHash describe the block used for the callback,
   * not necessarily the block number that the log was emitted in!!!!
   */
  struct LogTrigger {
    bytes32 logBlockHash;
    bytes32 txHash;
    uint32 logIndex;
    uint32 blockNum;
    bytes32 blockHash;
  }

  function checkUpkeep(
    uint256 id,
    bytes memory triggerData
  )
    external
    view
    returns (
      bool upkeepNeeded,
      bytes memory performData,
      uint8 upkeepFailureReason,
      uint256 gasUsed,
      uint256 gasLimit,
      uint256 fastGasWei,
      uint256 linkNative
    );
  function checkUpkeep(
    uint256 id
  )
    external
    view
    returns (
      bool upkeepNeeded,
      bytes memory performData,
      uint8 upkeepFailureReason,
      uint256 gasUsed,
      uint256 gasLimit,
      uint256 fastGasWei,
      uint256 linkNative
    );
  function simulatePerformUpkeep(
    uint256 id,
    bytes memory performData
  ) external view returns (bool success, uint256 gasUsed);
  function executeCallback(
    uint256 id,
    bytes memory payload
  ) external returns (bool upkeepNeeded, bytes memory performData, uint8 upkeepFailureReason, uint256 gasUsed);
  function checkCallback(
    uint256 id,
    bytes[] memory values,
    bytes memory extraData
  ) external view returns (bool upkeepNeeded, bytes memory performData, uint8 upkeepFailureReason, uint256 gasUsed);
  function typeAndVersion() external view returns (string memory);
  function addFunds(uint256 id, uint96 amount) external;
  function cancelUpkeep(uint256 id) external;

  function getUpkeepPrivilegeConfig(uint256 upkeepId) external view returns (bytes memory);
  function hasDedupKey(bytes32 dedupKey) external view returns (bool);
  function getUpkeepTriggerConfig(uint256 upkeepId) external view returns (bytes memory);
  function getUpkeep(uint256 id) external view returns (UpkeepInfo memory upkeepInfo);
  function getMinBalance(uint256 id) external view returns (uint96);
  function getState()
    external
    view
    returns (
      AutomationRegistryBase2_2.State memory state,
      AutomationRegistryBase2_2.OnchainConfigLegacy memory config,
      address[] memory signers,
      address[] memory transmitters,
      uint8 f
    );
  function registerUpkeep(
    address target,
    uint32 gasLimit,
    address admin,
    uint8 triggerType,
    bytes memory checkData,
    bytes memory triggerConfig,
    bytes memory offchainConfig
  ) external returns (uint256 id);
  function setUpkeepGasLimit(uint256 id, uint32 gasLimit) external;
  function setUpkeepPrivilegeConfig(uint256 upkeepId, bytes memory newPrivilegeConfig) external;
  function pauseUpkeep(uint256 id) external;
  function unpauseUpkeep(uint256 id) external;
  function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) external view returns (uint256[] memory);
  function pause() external;
  function setUpkeepCheckData(uint256 id, bytes memory newCheckData) external;
  function setUpkeepTriggerConfig(uint256 id, bytes memory triggerConfig) external;
  function owner() external view returns (address);
  function getTriggerType(uint256 upkeepId) external pure returns (uint8);
}
