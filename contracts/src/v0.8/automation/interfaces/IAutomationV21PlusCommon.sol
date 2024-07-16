// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.4;

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
  struct UpkeepInfoLegacy {
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
  // solhint-disable-next-line gas-struct-packing
  struct StateLegacy {
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

  /**
   * @notice OnchainConfigLegacy of the registry
   * @dev only used in params and return values
   * @member paymentPremiumPPB payment premium rate oracles receive on top of
   * being reimbursed for gas, measured in parts per billion
   * @member flatFeeMicroLink flat fee paid to oracles for performing upkeeps,
   * priced in MicroLink; can be used in conjunction with or independently of
   * paymentPremiumPPB
   * @member checkGasLimit gas limit when checking for upkeep
   * @member stalenessSeconds number of seconds that is allowed for feed data to
   * be stale before switching to the fallback pricing
   * @member gasCeilingMultiplier multiplier to apply to the fast gas feed price
   * when calculating the payment ceiling for keepers
   * @member minUpkeepSpend minimum LINK that an upkeep must spend before cancelling
   * @member maxPerformGas max performGas allowed for an upkeep on this registry
   * @member maxCheckDataSize max length of checkData bytes
   * @member maxPerformDataSize max length of performData bytes
   * @member maxRevertDataSize max length of revertData bytes
   * @member fallbackGasPrice gas price used if the gas price feed is stale
   * @member fallbackLinkPrice LINK price used if the LINK price feed is stale
   * @member transcoder address of the transcoder contract
   * @member registrars addresses of the registrar contracts
   * @member upkeepPrivilegeManager address which can set privilege for upkeeps
   */
  // solhint-disable-next-line gas-struct-packing
  struct OnchainConfigLegacy {
    uint32 paymentPremiumPPB;
    uint32 flatFeeMicroLink; // min 0.000001 LINK, max 4294 LINK
    uint32 checkGasLimit;
    uint24 stalenessSeconds;
    uint16 gasCeilingMultiplier;
    uint96 minUpkeepSpend;
    uint32 maxPerformGas;
    uint32 maxCheckDataSize;
    uint32 maxPerformDataSize;
    uint32 maxRevertDataSize;
    uint256 fallbackGasPrice;
    uint256 fallbackLinkPrice;
    address transcoder;
    address[] registrars;
    address upkeepPrivilegeManager;
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
  function getUpkeep(uint256 id) external view returns (UpkeepInfoLegacy memory upkeepInfo);
  function getMinBalance(uint256 id) external view returns (uint96);
  function getState()
    external
    view
    returns (
      StateLegacy memory state,
      OnchainConfigLegacy memory config,
      address[] memory signers,
      address[] memory transmitters,
      uint8 f
    );
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
