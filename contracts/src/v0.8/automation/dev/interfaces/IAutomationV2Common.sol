// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.4;

import {OnchainConfigLegacy} from "../AutomationConvenience.sol";

interface IAutomationV2Common {
  event UpkeepPerformed(
    uint256 indexed id,
    bool indexed success,
    uint96 totalPayment,
    uint256 gasUsed,
    uint256 gasOverhead,
    bytes trigger
  );
  event StaleUpkeepReport(uint256 indexed id, bytes trigger);
  event InsufficientFundsUpkeepReport(uint256 indexed id, bytes trigger);
  event ReorgedUpkeepReport(uint256 indexed id, bytes trigger);

  event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom);
  event UpkeepRegistered(uint256 indexed id, uint32 performGas, address admin);
  event UpkeepTriggerConfigSet(uint256 indexed id, bytes triggerConfig);
  event UpkeepPaused(uint256 indexed id);
  event UpkeepUnpaused(uint256 indexed id);
  event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination);
  event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight);

  event DedupKeyAdded(bytes32 indexed dedupKey);

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

/**
 * @notice the trigger structure conditional trigger type
 */
  struct ConditionalTrigger {
    uint32 blockNum;
    bytes32 blockHash;
  }

  /**
   * @notice state of the registry
   * @dev only used in params and return values
   * @dev this will likely be deprecated in a future version of the registry in favor of individual getters
   * @member nonce used for ID generation
   * @member ownerLinkBalance withdrawable balance of LINK by contract owner
   * @member expectedLinkBalance the expected balance of LINK of the registry
   * @member totalPremium the total premium collected on registry so far
   * @member numUpkeeps total number of upkeeps on the registry
   * @member configCount ordinal number of current config, out of all configs applied to this contract so far
   * @member latestConfigBlockNumber last block at which this config was set
   * @member latestConfigDigest domain-separation tag for current config
   * @member latestEpoch for which a report was transmitted
   * @member paused freeze on execution scoped to the entire registry
   */
  struct State {
    uint32 nonce;
    uint96 ownerLinkBalance;
    uint256 expectedLinkBalance;
    uint96 totalPremium;
    uint256 numUpkeeps;
    uint32 configCount;
    uint32 latestConfigBlockNumber;
    bytes32 latestConfigDigest;
    uint32 latestEpoch;
    bool paused;
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
      State memory state,
      OnchainConfigLegacy memory config,
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
  function setConfig(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfigBytes,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external;
  function owner() external view returns (address);
  function getTriggerType(uint256 upkeepId) external pure returns (uint8);
}
