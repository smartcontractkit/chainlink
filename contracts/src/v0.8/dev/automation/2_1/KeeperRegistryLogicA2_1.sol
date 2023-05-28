// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./KeeperRegistryBase2_1.sol";
import "./KeeperRegistryLogicB2_1.sol";
import "./Chainable.sol";
import {AutomationForwarder} from "./AutomationForwarder.sol";
import "../../../interfaces/automation/UpkeepTranscoderInterfaceV2.sol";

// TODO - we can probably combine these interfaces
import "../../../interfaces/automation/MigratableKeeperRegistryInterface.sol";
import "../../../interfaces/automation/MigratableKeeperRegistryInterfaceV2.sol";

/**
 * @notice Logic contract, works in tandem with KeeperRegistry as a proxy
 */
contract KeeperRegistryLogicA2_1 is
  KeeperRegistryBase2_1,
  Chainable,
  MigratableKeeperRegistryInterface,
  MigratableKeeperRegistryInterfaceV2
{
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;

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
      logicB.getFastGasFeedAddress()
    )
    Chainable(address(logicB))
  {}

  UpkeepFormat public constant override upkeepTranscoderVersion = UPKEEP_TRANSCODER_VERSION_BASE;

  uint8 public constant override upkeepVersion = UPKEEP_VERSION_BASE;

  /**
   * @dev this function will be deprecated in a future version of chainlink automation
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
      uint256 fastGasWei,
      uint256 linkNative
    )
  {
    return checkUpkeep(id, s_checkData[id]);
  }

  function checkUpkeep(
    uint256 id,
    bytes memory checkData
  )
    public
    cannotExecute
    returns (
      bool upkeepNeeded,
      bytes memory performData,
      UpkeepFailureReason upkeepFailureReason,
      uint256 gasUsed,
      uint256 fastGasWei,
      uint256 linkNative
    )
  {
    Trigger triggerType = getTriggerType(id);
    HotVars memory hotVars = s_hotVars;
    Upkeep memory upkeep = s_upkeep[id];
    if (upkeep.maxValidBlocknumber != UINT32_MAX)
      return (false, bytes(""), UpkeepFailureReason.UPKEEP_CANCELLED, gasUsed, 0, 0);
    if (upkeep.paused) return (false, bytes(""), UpkeepFailureReason.UPKEEP_PAUSED, gasUsed, 0, 0);

    (fastGasWei, linkNative) = _getFeedData(hotVars);
    uint96 maxLinkPayment = _getMaxLinkPayment(
      hotVars,
      upkeep.executeGas,
      s_storage.maxPerformDataSize,
      fastGasWei,
      linkNative,
      false
    );
    if (upkeep.balance < maxLinkPayment) {
      return (false, bytes(""), UpkeepFailureReason.INSUFFICIENT_BALANCE, gasUsed, fastGasWei, linkNative);
    }

    upkeepNeeded = true;

    if (triggerType == Trigger.CONDITION || triggerType == Trigger.LOG) {
      gasUsed = gasleft();
      bytes memory callData = abi.encodeWithSelector(CHECK_SELECTOR, checkData);
      (upkeepNeeded, performData) = upkeep.target.call{gas: s_storage.checkGasLimit}(callData);
      gasUsed = gasUsed - gasleft();
      if (!upkeepNeeded) {
        upkeepFailureReason = UpkeepFailureReason.TARGET_CHECK_REVERTED;
      } else {
        (upkeepNeeded, performData) = abi.decode(performData, (bool, bytes));
        if (!upkeepNeeded)
          return (false, bytes(""), UpkeepFailureReason.UPKEEP_NOT_NEEDED, gasUsed, fastGasWei, linkNative);
        if (performData.length > s_storage.maxPerformDataSize)
          return (false, bytes(""), UpkeepFailureReason.PERFORM_DATA_EXCEEDS_LIMIT, gasUsed, fastGasWei, linkNative);
      }
    }
    // TODO - consider doing cron validation

    return (upkeepNeeded, performData, upkeepFailureReason, gasUsed, fastGasWei, linkNative);
  }

  function registerUpkeep(
    address target,
    uint32 gasLimit, // TODO - we may want to allow 0 for "unlimited"
    address admin,
    bytes calldata checkData, // TODO - this should be included in the trigger
    bytes calldata extraData
  ) external returns (uint256 id) {
    if (msg.sender != owner() && msg.sender != s_storage.registrar) revert OnlyCallableByOwnerOrRegistrar();
    (Trigger triggerType, bytes memory triggerConfig, bytes memory offchainConfig) = abi.decode(
      extraData,
      (Trigger, bytes, bytes)
    );
    validateTrigger(triggerType, triggerConfig);
    id = _createID(triggerType);
    AutomationForwarder forwarder = new AutomationForwarder(id, target);
    _createUpkeep(id, target, gasLimit, admin, 0, checkData, false, triggerConfig, offchainConfig, forwarder);
    s_storage.nonce++;
    emit UpkeepRegistered(id, gasLimit, admin);
    emit UpkeepTriggerConfigSet(id, triggerConfig);
    emit UpkeepOffchainConfigSet(id, offchainConfig);
    return (id);
  }

  function addFunds(uint256 id, uint96 amount) external {
    Upkeep memory upkeep = s_upkeep[id];
    if (upkeep.maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();

    s_upkeep[id].balance = upkeep.balance + amount;
    s_expectedLinkBalance = s_expectedLinkBalance + amount;
    i_link.transferFrom(msg.sender, address(this), amount);
    emit FundsAdded(id, msg.sender, amount);
  }

  /**
   * @notice creates an ID for the upkeep based on the upkeep's type
   * @dev the format of the ID looks like this:
   * ****00000000000X****************
   * 4 bytes of entropy
   * 11 bytes of zeros
   * 1 identifying byte for the trigger type
   * 16 bytes of entropy
   * @dev this maintains the same level of entropy as eth addresses, so IDs will still be unique
   * @dev we add the "identifying" part in the middle so that it is mostly hidden from users who usually only
   * see the first 4 and last 4 hex values ex 0x1234...ABCD
   */
  function _createID(Trigger triggerType) private view returns (uint256) {
    bytes1 empty;
    bytes memory idBytes = abi.encodePacked(
      keccak256(abi.encode(_blockHash(_blockNum() - 1), address(this), s_storage.nonce))
    );
    for (uint256 idx = 4; idx < 15; idx++) {
      idBytes[idx] = empty;
    }
    idBytes[15] = bytes1(uint8(triggerType));
    return uint256(bytes32(idBytes));
  }

  function validateTrigger(Trigger triggerType, bytes memory offchainConfig) public {
    if (msg.sender == address(this)) {
      // separate call stack, we can revert with any reason here
      if (triggerType == Trigger.CONDITION || triggerType == Trigger.READY) {
        require(offchainConfig.length == 0);
      } else if (triggerType == Trigger.LOG) {
        // will revert if data isn't in the correct format
        LogTriggerConfig memory trigger = abi.decode(offchainConfig, (LogTriggerConfig));
        require(trigger.contractAddress != ZERO_ADDRESS);
        require(trigger.topic0 != bytes32(0));
        require(uint8(trigger.filterSelector) < 8); // 8 corresponds to 1000 in binary, max is 111
      } else if (triggerType == Trigger.CRON) {
        // will revert if data isn't in the correct format
        CronTriggerConfig memory trigger = abi.decode(offchainConfig, (CronTriggerConfig));
        require(trigger.payload.length % 32 == 4);
        // TODO - gas analysis to see if it's feasible to validate cron string
      } else {
        revert();
      }
    } else {
      // called directly, only revert with proper message
      try KeeperRegistryLogicA2_1(address(this)).validateTrigger(triggerType, offchainConfig) {
        return;
      } catch {
        revert InvalidTrigger();
      }
    }
  }

  /**
   * @dev Called through KeeperRegistry main contract
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

  function setUpkeepTriggerConfig(uint256 id, bytes calldata triggerConfig) external {
    _requireAdminAndNotCancelled(id);
    Trigger triggerType = getTriggerType(id);
    validateTrigger(triggerType, triggerConfig);
    s_upkeepTriggerConfig[id] = triggerConfig;
    emit UpkeepTriggerConfigSet(id, triggerConfig);
  }

  function migrateUpkeeps(
    uint256[] calldata ids,
    address destination
  ) external override(MigratableKeeperRegistryInterface, MigratableKeeperRegistryInterfaceV2) {
    if (
      s_peerRegistryMigrationPermission[destination] != MigrationPermission.OUTGOING &&
      s_peerRegistryMigrationPermission[destination] != MigrationPermission.BIDIRECTIONAL
    ) revert MigrationNotPermitted();
    if (s_storage.transcoder == ZERO_ADDRESS) revert TranscoderNotSet();
    if (ids.length == 0) revert ArrayHasNoEntries();
    uint256 id;
    Upkeep memory upkeep;
    uint256 totalBalanceRemaining;
    bytes[] memory checkDatas = new bytes[](ids.length);
    address[] memory admins = new address[](ids.length);
    Upkeep[] memory upkeeps = new Upkeep[](ids.length);
    bytes[] memory triggerConfigs = new bytes[](ids.length);
    bytes[] memory offchainConfigs = new bytes[](ids.length);
    for (uint256 idx = 0; idx < ids.length; idx++) {
      id = ids[idx];
      upkeep = s_upkeep[id];
      _requireAdminAndNotCancelled(id);
      upkeeps[idx] = upkeep;
      checkDatas[idx] = s_checkData[id];
      admins[idx] = s_upkeepAdmin[id];
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
    bytes memory encodedUpkeeps = abi.encode(ids, upkeeps, checkDatas, admins, triggerConfigs, offchainConfigs);
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
   * @dev Called through KeeperRegistry main contract
   */
  function receiveUpkeeps(
    bytes calldata encodedUpkeeps
  ) external override(MigratableKeeperRegistryInterface, MigratableKeeperRegistryInterfaceV2) {
    if (
      s_peerRegistryMigrationPermission[msg.sender] != MigrationPermission.INCOMING &&
      s_peerRegistryMigrationPermission[msg.sender] != MigrationPermission.BIDIRECTIONAL
    ) revert MigrationNotPermitted();
    (
      uint256[] memory ids,
      Upkeep[] memory upkeeps,
      bytes[] memory checkDatas,
      address[] memory upkeepAdmins,
      bytes[] memory triggerConfigs,
      bytes[] memory offchainConfigs
    ) = abi.decode(encodedUpkeeps, (uint256[], Upkeep[], bytes[], address[], bytes[], bytes[]));
    // TODO - we should be creating the forwarder in the transcoder, not here
    for (uint256 idx = 0; idx < ids.length; idx++) {
      if (address(upkeeps[idx].forwarder) == ZERO_ADDRESS) {
        upkeeps[idx].forwarder = new AutomationForwarder(ids[idx], upkeeps[idx].target);
      }
      _createUpkeep(
        ids[idx],
        upkeeps[idx].target,
        upkeeps[idx].executeGas,
        upkeepAdmins[idx],
        upkeeps[idx].balance,
        checkDatas[idx],
        upkeeps[idx].paused,
        triggerConfigs[idx],
        offchainConfigs[idx],
        upkeeps[idx].forwarder
      );
      emit UpkeepReceived(ids[idx], upkeeps[idx].balance, msg.sender);
    }
  }
}
