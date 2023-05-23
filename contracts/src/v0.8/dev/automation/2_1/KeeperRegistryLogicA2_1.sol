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
    if (upkeep.balance < maxLinkPayment)
      return (false, bytes(""), UpkeepFailureReason.INSUFFICIENT_BALANCE, gasUsed, fastGasWei, linkNative);

    gasUsed = gasleft();
    bytes memory callData = abi.encodeWithSelector(CHECK_SELECTOR, checkData);
    (bool success, bytes memory result) = upkeep.target.call{gas: s_storage.checkGasLimit}(callData);
    gasUsed = gasUsed - gasleft();

    if (!success) {
      upkeepFailureReason = UpkeepFailureReason.TARGET_CHECK_REVERTED;
    } else {
      (upkeepNeeded, result) = abi.decode(result, (bool, bytes));
      if (!upkeepNeeded)
        return (false, bytes(""), UpkeepFailureReason.UPKEEP_NOT_NEEDED, gasUsed, fastGasWei, linkNative);
      if (result.length > s_storage.maxPerformDataSize)
        return (false, bytes(""), UpkeepFailureReason.PERFORM_DATA_EXCEEDS_LIMIT, gasUsed, fastGasWei, linkNative);
    }

    return (success, result, upkeepFailureReason, gasUsed, fastGasWei, linkNative);
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

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function withdrawFunds(uint256 id, address to) external nonReentrant {
    if (to == ZERO_ADDRESS) revert InvalidRecipient();
    Upkeep memory upkeep = s_upkeep[id];
    if (s_upkeepAdmin[id] != msg.sender) revert OnlyCallableByAdmin();
    if (upkeep.maxValidBlocknumber > _blockNum()) revert UpkeepNotCanceled();

    uint96 amountToWithdraw = s_upkeep[id].balance;
    s_expectedLinkBalance = s_expectedLinkBalance - amountToWithdraw;
    s_upkeep[id].balance = 0;
    i_link.transfer(to, amountToWithdraw);
    emit FundsWithdrawn(id, amountToWithdraw, to);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function setUpkeepGasLimit(uint256 id, uint32 gasLimit) external {
    if (gasLimit < PERFORM_GAS_MIN || gasLimit > s_storage.maxPerformGas) revert GasLimitOutsideRange();
    _requireAdminAndNotCancelled(id);
    s_upkeep[id].executeGas = gasLimit;

    emit UpkeepGasLimitSet(id, gasLimit);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function setUpkeepOffchainConfig(uint256 id, bytes calldata config) external {
    _requireAdminAndNotCancelled(id);

    s_upkeepOffchainConfig[id] = config;

    emit UpkeepOffchainConfigSet(id, config);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function withdrawPayment(address from, address to) external {
    if (to == ZERO_ADDRESS) revert InvalidRecipient();
    if (s_transmitterPayees[from] != msg.sender) revert OnlyCallableByPayee();

    uint96 balance = _updateTransmitterBalanceFromPool(from, s_hotVars.totalPremium, uint96(s_transmittersList.length));
    s_transmitters[from].balance = 0;
    s_expectedLinkBalance = s_expectedLinkBalance - balance;

    i_link.transfer(to, balance);

    emit PaymentWithdrawn(from, balance, to, msg.sender);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function transferPayeeship(address transmitter, address proposed) external {
    if (s_transmitterPayees[transmitter] != msg.sender) revert OnlyCallableByPayee();
    if (proposed == msg.sender) revert ValueNotChanged();

    if (s_proposedPayee[transmitter] != proposed) {
      s_proposedPayee[transmitter] = proposed;
      emit PayeeshipTransferRequested(transmitter, msg.sender, proposed);
    }
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function acceptPayeeship(address transmitter) external {
    if (s_proposedPayee[transmitter] != msg.sender) revert OnlyCallableByProposedPayee();
    address past = s_transmitterPayees[transmitter];
    s_transmitterPayees[transmitter] = msg.sender;
    s_proposedPayee[transmitter] = ZERO_ADDRESS;

    emit PayeeshipTransferred(transmitter, past, msg.sender);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function transferUpkeepAdmin(uint256 id, address proposed) external {
    _requireAdminAndNotCancelled(id);
    if (proposed == msg.sender) revert ValueNotChanged();
    if (proposed == ZERO_ADDRESS) revert InvalidRecipient();

    if (s_proposedAdmin[id] != proposed) {
      s_proposedAdmin[id] = proposed;
      emit UpkeepAdminTransferRequested(id, msg.sender, proposed);
    }
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function acceptUpkeepAdmin(uint256 id) external {
    Upkeep memory upkeep = s_upkeep[id];
    if (upkeep.maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();
    if (s_proposedAdmin[id] != msg.sender) revert OnlyCallableByProposedAdmin();
    address past = s_upkeepAdmin[id];
    s_upkeepAdmin[id] = msg.sender;
    s_proposedAdmin[id] = ZERO_ADDRESS;

    emit UpkeepAdminTransferred(id, past, msg.sender);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function pauseUpkeep(uint256 id) external {
    _requireAdminAndNotCancelled(id);
    Upkeep memory upkeep = s_upkeep[id];
    if (upkeep.paused) revert OnlyUnpausedUpkeep();
    s_upkeep[id].paused = true;
    s_upkeepIDs.remove(id);
    emit UpkeepPaused(id);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function unpauseUpkeep(uint256 id) external {
    _requireAdminAndNotCancelled(id);
    Upkeep memory upkeep = s_upkeep[id];
    if (!upkeep.paused) revert OnlyPausedUpkeep();
    s_upkeep[id].paused = false;
    s_upkeepIDs.add(id);
    emit UpkeepUnpaused(id);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function updateCheckData(uint256 id, bytes calldata newCheckData) external {
    _requireAdminAndNotCancelled(id);
    if (newCheckData.length > s_storage.maxCheckDataSize) revert CheckDataExceedsLimit();
    s_checkData[id] = newCheckData;
    emit UpkeepCheckDataUpdated(id, newCheckData);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
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
    bytes[] memory offchainConfigs = new bytes[](ids.length);
    for (uint256 idx = 0; idx < ids.length; idx++) {
      id = ids[idx];
      upkeep = s_upkeep[id];
      _requireAdminAndNotCancelled(id);
      upkeeps[idx] = upkeep;
      checkDatas[idx] = s_checkData[id];
      admins[idx] = s_upkeepAdmin[id];
      offchainConfigs[idx] = s_upkeepOffchainConfig[id];
      totalBalanceRemaining = totalBalanceRemaining + upkeep.balance;
      delete s_upkeep[id];
      delete s_checkData[id];
      // nullify existing proposed admin change if an upkeep is being migrated
      delete s_proposedAdmin[id];
      s_upkeepIDs.remove(id);
      emit UpkeepMigrated(id, upkeep.balance, destination);
    }
    s_expectedLinkBalance = s_expectedLinkBalance - totalBalanceRemaining;
    bytes memory encodedUpkeeps = abi.encode(ids, upkeeps, checkDatas, admins, offchainConfigs);
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
      bytes[] memory offchainConfigs
    ) = abi.decode(encodedUpkeeps, (uint256[], Upkeep[], bytes[], address[], bytes[]));
    // TODO - we should be creating the forwarder in the transcoder, not here
    for (uint256 idx = 0; idx < ids.length; idx++) {
      if (address(upkeeps[idx].forwarder) == address(0)) {
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
        offchainConfigs[idx],
        upkeeps[idx].forwarder
      );
      emit UpkeepReceived(ids[idx], upkeeps[idx].balance, msg.sender);
    }
  }

  /**
   * @dev ensures the upkeep is not cancelled and the caller is the upkeep admin
   */
  function _requireAdminAndNotCancelled(uint256 upkeepId) internal view {
    if (msg.sender != s_upkeepAdmin[upkeepId]) revert OnlyCallableByAdmin();
    if (s_upkeep[upkeepId].maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();
  }
}
