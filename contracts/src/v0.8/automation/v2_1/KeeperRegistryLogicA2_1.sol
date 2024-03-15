// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.16;

import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import {Address} from "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";
import {KeeperRegistryBase2_1} from "./KeeperRegistryBase2_1.sol";
import {KeeperRegistryLogicB2_1} from "./KeeperRegistryLogicB2_1.sol";
import {Chainable} from "../Chainable.sol";
import {AutomationForwarder} from "../AutomationForwarder.sol";
import {IAutomationForwarder} from "../interfaces/IAutomationForwarder.sol";
import {UpkeepTranscoderInterfaceV2} from "../interfaces/UpkeepTranscoderInterfaceV2.sol";
import {MigratableKeeperRegistryInterfaceV2} from "../interfaces/MigratableKeeperRegistryInterfaceV2.sol";

/**
 * @notice Logic contract, works in tandem with KeeperRegistry as a proxy
 */
contract KeeperRegistryLogicA2_1 is KeeperRegistryBase2_1, Chainable {
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;
  using EnumerableSet for EnumerableSet.AddressSet;

  /**
   * @param logicB the address of the second logic contract
   */
  constructor(
    KeeperRegistryLogicB2_1 logicB
  )
    KeeperRegistryBase2_1(
      logicB.getMode(),
      logicB.getLinkAddress(),
      logicB.getLinkNativeFeedAddress(),
      logicB.getFastGasFeedAddress(),
      logicB.getAutomationForwarderLogic()
    )
    Chainable(address(logicB))
  {}

  /**
   * @notice called by the automation DON to check if work is needed
   * @param id the upkeep ID to check for work needed
   * @param triggerData extra contextual data about the trigger (not used in all code paths)
   * @dev this one of the core functions called in the hot path
   * @dev there is a 2nd checkUpkeep function (below) that is being maintained for backwards compatibility
   * @dev there is an incongruency on what gets returned during failure modes
   * ex sometimes we include price data, sometimes we omit it depending on the failure
   */
  function checkUpkeep(
    uint256 id,
    bytes memory triggerData
  )
    public
    cannotExecute
    returns (
      bool upkeepNeeded,
      bytes memory performData,
      UpkeepFailureReason upkeepFailureReason,
      uint256 gasUsed,
      uint256 gasLimit,
      uint256 fastGasWei,
      uint256 linkNative
    )
  {
    Trigger triggerType = _getTriggerType(id);
    HotVars memory hotVars = s_hotVars;
    Upkeep memory upkeep = s_upkeep[id];

    if (hotVars.paused) return (false, bytes(""), UpkeepFailureReason.REGISTRY_PAUSED, 0, upkeep.performGas, 0, 0);
    if (upkeep.maxValidBlocknumber != UINT32_MAX)
      return (false, bytes(""), UpkeepFailureReason.UPKEEP_CANCELLED, 0, upkeep.performGas, 0, 0);
    if (upkeep.paused) return (false, bytes(""), UpkeepFailureReason.UPKEEP_PAUSED, 0, upkeep.performGas, 0, 0);

    (fastGasWei, linkNative) = _getFeedData(hotVars);
    uint96 maxLinkPayment = _getMaxLinkPayment(
      hotVars,
      triggerType,
      upkeep.performGas,
      s_storage.maxPerformDataSize,
      fastGasWei,
      linkNative,
      false
    );
    if (upkeep.balance < maxLinkPayment) {
      return (false, bytes(""), UpkeepFailureReason.INSUFFICIENT_BALANCE, 0, upkeep.performGas, 0, 0);
    }

    bytes memory callData = _checkPayload(id, triggerType, triggerData);

    gasUsed = gasleft();
    (bool success, bytes memory result) = upkeep.forwarder.getTarget().call{gas: s_storage.checkGasLimit}(callData);
    gasUsed = gasUsed - gasleft();

    if (!success) {
      // User's target check reverted. We capture the revert data here and pass it within performData
      if (result.length > s_storage.maxRevertDataSize) {
        return (
          false,
          bytes(""),
          UpkeepFailureReason.REVERT_DATA_EXCEEDS_LIMIT,
          gasUsed,
          upkeep.performGas,
          fastGasWei,
          linkNative
        );
      }
      return (
        upkeepNeeded,
        result,
        UpkeepFailureReason.TARGET_CHECK_REVERTED,
        gasUsed,
        upkeep.performGas,
        fastGasWei,
        linkNative
      );
    }

    (upkeepNeeded, performData) = abi.decode(result, (bool, bytes));
    if (!upkeepNeeded)
      return (
        false,
        bytes(""),
        UpkeepFailureReason.UPKEEP_NOT_NEEDED,
        gasUsed,
        upkeep.performGas,
        fastGasWei,
        linkNative
      );

    if (performData.length > s_storage.maxPerformDataSize)
      return (
        false,
        bytes(""),
        UpkeepFailureReason.PERFORM_DATA_EXCEEDS_LIMIT,
        gasUsed,
        upkeep.performGas,
        fastGasWei,
        linkNative
      );

    return (upkeepNeeded, performData, upkeepFailureReason, gasUsed, upkeep.performGas, fastGasWei, linkNative);
  }

  /**
   * @notice see other checkUpkeep function for description
   * @dev this function may be deprecated in a future version of chainlink automation
   */
  function checkUpkeep(
    uint256 id
  )
    external
    returns (
      bool upkeepNeeded,
      bytes memory performData,
      UpkeepFailureReason upkeepFailureReason,
      uint256 gasUsed,
      uint256 gasLimit,
      uint256 fastGasWei,
      uint256 linkNative
    )
  {
    return checkUpkeep(id, bytes(""));
  }

  /**
   * @dev checkCallback is used specifically for automation data streams lookups (see StreamsLookupCompatibleInterface.sol)
   * @param id the upkeepID to execute a callback for
   * @param values the values returned from the data streams lookup
   * @param extraData the user-provided extra context data
   */
  function checkCallback(
    uint256 id,
    bytes[] memory values,
    bytes calldata extraData
  )
    external
    cannotExecute
    returns (bool upkeepNeeded, bytes memory performData, UpkeepFailureReason upkeepFailureReason, uint256 gasUsed)
  {
    bytes memory payload = abi.encodeWithSelector(CHECK_CALLBACK_SELECTOR, values, extraData);
    return executeCallback(id, payload);
  }

  /**
   * @notice this is a generic callback executor that forwards a call to a user's contract with the configured
   * gas limit
   * @param id the upkeepID to execute a callback for
   * @param payload the data (including function selector) to call on the upkeep target contract
   */
  function executeCallback(
    uint256 id,
    bytes memory payload
  )
    public
    cannotExecute
    returns (bool upkeepNeeded, bytes memory performData, UpkeepFailureReason upkeepFailureReason, uint256 gasUsed)
  {
    Upkeep memory upkeep = s_upkeep[id];
    gasUsed = gasleft();
    (bool success, bytes memory result) = upkeep.forwarder.getTarget().call{gas: s_storage.checkGasLimit}(payload);
    gasUsed = gasUsed - gasleft();
    if (!success) {
      return (false, bytes(""), UpkeepFailureReason.CALLBACK_REVERTED, gasUsed);
    }
    (upkeepNeeded, performData) = abi.decode(result, (bool, bytes));
    if (!upkeepNeeded) {
      return (false, bytes(""), UpkeepFailureReason.UPKEEP_NOT_NEEDED, gasUsed);
    }
    if (performData.length > s_storage.maxPerformDataSize) {
      return (false, bytes(""), UpkeepFailureReason.PERFORM_DATA_EXCEEDS_LIMIT, gasUsed);
    }
    return (upkeepNeeded, performData, upkeepFailureReason, gasUsed);
  }

  /**
   * @notice adds a new upkeep
   * @param target address to perform upkeep on
   * @param gasLimit amount of gas to provide the target contract when
   * performing upkeep
   * @param admin address to cancel upkeep and withdraw remaining funds
   * @param triggerType the trigger for the upkeep
   * @param checkData data passed to the contract when checking for upkeep
   * @param triggerConfig the config for the trigger
   * @param offchainConfig arbitrary offchain config for the upkeep
   */
  function registerUpkeep(
    address target,
    uint32 gasLimit,
    address admin,
    Trigger triggerType,
    bytes calldata checkData,
    bytes memory triggerConfig,
    bytes memory offchainConfig
  ) public returns (uint256 id) {
    if (msg.sender != owner() && !s_registrars.contains(msg.sender)) revert OnlyCallableByOwnerOrRegistrar();
    if (!target.isContract()) revert NotAContract();
    id = _createID(triggerType);
    IAutomationForwarder forwarder = IAutomationForwarder(
      address(new AutomationForwarder(target, address(this), i_automationForwarderLogic))
    );
    _createUpkeep(
      id,
      Upkeep({
        performGas: gasLimit,
        balance: 0,
        maxValidBlocknumber: UINT32_MAX,
        lastPerformedBlockNumber: 0,
        amountSpent: 0,
        paused: false,
        forwarder: forwarder
      }),
      admin,
      checkData,
      triggerConfig,
      offchainConfig
    );
    s_storage.nonce++;
    emit UpkeepRegistered(id, gasLimit, admin);
    emit UpkeepCheckDataSet(id, checkData);
    emit UpkeepTriggerConfigSet(id, triggerConfig);
    emit UpkeepOffchainConfigSet(id, offchainConfig);
    return (id);
  }

  /**
   * @notice this function registers a conditional upkeep, using a backwards compatible function signature
   * @dev this function is backwards compatible with versions <=2.0, but may be removed in a future version
   */
  function registerUpkeep(
    address target,
    uint32 gasLimit,
    address admin,
    bytes calldata checkData,
    bytes calldata offchainConfig
  ) external returns (uint256 id) {
    return registerUpkeep(target, gasLimit, admin, Trigger.CONDITION, checkData, bytes(""), offchainConfig);
  }

  /**
   * @notice cancels an upkeep
   * @param id the upkeepID to cancel
   * @dev if a user cancels an upkeep, their funds are locked for CANCELLATION_DELAY blocks to
   * allow any pending performUpkeep txs time to get confirmed
   */
  function cancelUpkeep(uint256 id) external {
    Upkeep memory upkeep = s_upkeep[id];
    bool canceled = upkeep.maxValidBlocknumber != UINT32_MAX;
    bool isOwner = msg.sender == owner();

    if (canceled && !(isOwner && upkeep.maxValidBlocknumber > _blockNum())) revert CannotCancel();
    if (!isOwner && msg.sender != s_upkeepAdmin[id]) revert OnlyCallableByOwnerOrAdmin();

    uint256 height = _blockNum();
    if (!isOwner) {
      height = height + CANCELLATION_DELAY;
    }
    s_upkeep[id].maxValidBlocknumber = uint32(height);
    s_upkeepIDs.remove(id);

    // charge the cancellation fee if the minUpkeepSpend is not met
    uint96 minUpkeepSpend = s_storage.minUpkeepSpend;
    uint96 cancellationFee = 0;
    // cancellationFee is supposed to be min(max(minUpkeepSpend - amountSpent,0), amountLeft)
    if (upkeep.amountSpent < minUpkeepSpend) {
      cancellationFee = minUpkeepSpend - upkeep.amountSpent;
      if (cancellationFee > upkeep.balance) {
        cancellationFee = upkeep.balance;
      }
    }
    s_upkeep[id].balance = upkeep.balance - cancellationFee;
    s_storage.ownerLinkBalance = s_storage.ownerLinkBalance + cancellationFee;

    emit UpkeepCanceled(id, uint64(height));
  }

  /**
   * @notice adds fund to an upkeep
   * @param id the upkeepID
   * @param amount the amount of LINK to fund, in jules (jules = "wei" of LINK)
   */
  function addFunds(uint256 id, uint96 amount) external {
    Upkeep memory upkeep = s_upkeep[id];
    if (upkeep.maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();
    s_upkeep[id].balance = upkeep.balance + amount;
    s_expectedLinkBalance = s_expectedLinkBalance + amount;
    i_link.transferFrom(msg.sender, address(this), amount);
    emit FundsAdded(id, msg.sender, amount);
  }

  /**
   * @notice migrates upkeeps from one registry to another
   * @param ids the upkeepIDs to migrate
   * @param destination the destination registry address
   * @dev a transcoder must be set in order to enable migration
   * @dev migration permissions must be set on *both* sending and receiving registries
   * @dev only an upkeep admin can migrate their upkeeps
   */
  function migrateUpkeeps(uint256[] calldata ids, address destination) external {
    if (
      s_peerRegistryMigrationPermission[destination] != MigrationPermission.OUTGOING &&
      s_peerRegistryMigrationPermission[destination] != MigrationPermission.BIDIRECTIONAL
    ) revert MigrationNotPermitted();
    if (s_storage.transcoder == ZERO_ADDRESS) revert TranscoderNotSet();
    if (ids.length == 0) revert ArrayHasNoEntries();
    uint256 id;
    Upkeep memory upkeep;
    uint256 totalBalanceRemaining;
    address[] memory admins = new address[](ids.length);
    Upkeep[] memory upkeeps = new Upkeep[](ids.length);
    bytes[] memory checkDatas = new bytes[](ids.length);
    bytes[] memory triggerConfigs = new bytes[](ids.length);
    bytes[] memory offchainConfigs = new bytes[](ids.length);
    for (uint256 idx = 0; idx < ids.length; idx++) {
      id = ids[idx];
      upkeep = s_upkeep[id];
      _requireAdminAndNotCancelled(id);
      upkeep.forwarder.updateRegistry(destination);
      upkeeps[idx] = upkeep;
      admins[idx] = s_upkeepAdmin[id];
      checkDatas[idx] = s_checkData[id];
      triggerConfigs[idx] = s_upkeepTriggerConfig[id];
      offchainConfigs[idx] = s_upkeepOffchainConfig[id];
      totalBalanceRemaining = totalBalanceRemaining + upkeep.balance;
      delete s_upkeep[id];
      delete s_checkData[id];
      delete s_upkeepTriggerConfig[id];
      delete s_upkeepOffchainConfig[id];
      // nullify existing proposed admin change if an upkeep is being migrated
      delete s_proposedAdmin[id];
      s_upkeepIDs.remove(id);
      emit UpkeepMigrated(id, upkeep.balance, destination);
    }
    s_expectedLinkBalance = s_expectedLinkBalance - totalBalanceRemaining;
    bytes memory encodedUpkeeps = abi.encode(
      ids,
      upkeeps,
      new address[](ids.length),
      admins,
      checkDatas,
      triggerConfigs,
      offchainConfigs
    );
    MigratableKeeperRegistryInterfaceV2(destination).receiveUpkeeps(
      UpkeepTranscoderInterfaceV2(s_storage.transcoder).transcodeUpkeeps(
        UPKEEP_VERSION_BASE,
        MigratableKeeperRegistryInterfaceV2(destination).upkeepVersion(),
        encodedUpkeeps
      )
    );
    i_link.transfer(destination, totalBalanceRemaining);
  }

  /**
   * @notice received upkeeps migrated from another registry
   * @param encodedUpkeeps the raw upkeep data to import
   * @dev this function is never called directly, it is only called by another registry's migrate function
   */
  function receiveUpkeeps(bytes calldata encodedUpkeeps) external {
    if (
      s_peerRegistryMigrationPermission[msg.sender] != MigrationPermission.INCOMING &&
      s_peerRegistryMigrationPermission[msg.sender] != MigrationPermission.BIDIRECTIONAL
    ) revert MigrationNotPermitted();
    (
      uint256[] memory ids,
      Upkeep[] memory upkeeps,
      address[] memory targets,
      address[] memory upkeepAdmins,
      bytes[] memory checkDatas,
      bytes[] memory triggerConfigs,
      bytes[] memory offchainConfigs
    ) = abi.decode(encodedUpkeeps, (uint256[], Upkeep[], address[], address[], bytes[], bytes[], bytes[]));
    for (uint256 idx = 0; idx < ids.length; idx++) {
      if (address(upkeeps[idx].forwarder) == ZERO_ADDRESS) {
        upkeeps[idx].forwarder = IAutomationForwarder(
          address(new AutomationForwarder(targets[idx], address(this), i_automationForwarderLogic))
        );
      }
      _createUpkeep(
        ids[idx],
        upkeeps[idx],
        upkeepAdmins[idx],
        checkDatas[idx],
        triggerConfigs[idx],
        offchainConfigs[idx]
      );
      emit UpkeepReceived(ids[idx], upkeeps[idx].balance, msg.sender);
    }
  }

  /**
   * @notice sets the upkeep trigger config
   * @param id the upkeepID to change the trigger for
   * @param triggerConfig the new trigger config
   */
  function setUpkeepTriggerConfig(uint256 id, bytes calldata triggerConfig) external {
    _requireAdminAndNotCancelled(id);
    s_upkeepTriggerConfig[id] = triggerConfig;
    emit UpkeepTriggerConfigSet(id, triggerConfig);
  }
}
