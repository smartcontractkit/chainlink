// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./KeeperRegistryBase2_1.sol";
import {UpkeepInfo} from "./interfaces/AutomationRegistryInterface2_1.sol";

contract KeeperRegistryLogicB2_1 is KeeperRegistryBase2_1 {
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;

  /**
   * @dev see KeeperRegistry master contract for constructor description
   */
  constructor(
    Mode mode,
    address link,
    address linkNativeFeed,
    address fastGasFeed
  ) KeeperRegistryBase2_1(mode, link, linkNativeFeed, fastGasFeed) {}

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

  ///////////////////
  // OWNER ACTIONS //
  ///////////////////

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function withdrawOwnerFunds() external onlyOwner {
    uint96 amount = s_storage.ownerLinkBalance;

    s_expectedLinkBalance = s_expectedLinkBalance - amount;
    s_storage.ownerLinkBalance = 0;

    emit OwnerFundsWithdrawn(amount);
    i_link.transfer(msg.sender, amount);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function recoverFunds() external onlyOwner {
    uint256 total = i_link.balanceOf(address(this));
    i_link.transfer(msg.sender, total - s_expectedLinkBalance);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function setPayees(address[] calldata payees) external onlyOwner {
    if (s_transmittersList.length != payees.length) revert ParameterLengthError();
    for (uint256 i = 0; i < s_transmittersList.length; i++) {
      address transmitter = s_transmittersList[i];
      address oldPayee = s_transmitterPayees[transmitter];
      address newPayee = payees[i];
      if (
        (newPayee == ZERO_ADDRESS) || (oldPayee != ZERO_ADDRESS && oldPayee != newPayee && newPayee != IGNORE_ADDRESS)
      ) revert InvalidPayee();
      if (newPayee != IGNORE_ADDRESS) {
        s_transmitterPayees[transmitter] = newPayee;
      }
    }
    emit PayeesUpdated(s_transmittersList, payees);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function pause() external onlyOwner {
    s_hotVars.paused = true;

    emit Paused(msg.sender);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function unpause() external onlyOwner {
    s_hotVars.paused = false;

    emit Unpaused(msg.sender);
  }

  /**
   * @dev Called through KeeperRegistry main contract
   */
  function setPeerRegistryMigrationPermission(address peer, MigrationPermission permission) external onlyOwner {
    s_peerRegistryMigrationPermission[peer] = permission;
  }

  /////////////
  // GETTERS //
  /////////////

  /**
   * @notice read all of the details about an upkeep
   */
  function getUpkeep(uint256 id) external view returns (UpkeepInfo memory upkeepInfo) {
    Upkeep memory reg = s_upkeep[id];
    upkeepInfo = UpkeepInfo({
      target: reg.target,
      forwarder: address(reg.forwarder),
      executeGas: reg.executeGas,
      checkData: s_checkData[id],
      balance: reg.balance,
      admin: s_upkeepAdmin[id],
      maxValidBlocknumber: reg.maxValidBlocknumber,
      lastPerformBlockNumber: reg.lastPerformBlockNumber,
      amountSpent: reg.amountSpent,
      paused: reg.paused,
      offchainConfig: s_upkeepOffchainConfig[id]
    });
    return upkeepInfo;
  }

  /**
   * @notice retrieve active upkeep IDs. Active upkeep is defined as an upkeep which is not paused and not canceled.
   * @param startIndex starting index in list
   * @param endIndex end index. non-inclusive
   * @dev the order of IDs in the list is **not guaranteed**, therefore, if making successive calls, one
   * should consider keeping the blockheight constant to ensure a holistic picture of the contract state
   * @dev this function will not error if an endIndex is provided that is too large, instead, it will resolve gracefully
   * by suppling as many results as it can within the bounds of the upkeep set
   */
  function getActiveUpkeepIDs(uint256 startIndex, uint256 endIndex) external view returns (uint256[] memory) {
    uint256 maxIndex = s_upkeepIDs.length();
    endIndex = endIndex > maxIndex ? maxIndex : endIndex;
    if (startIndex > endIndex) revert IndexOutOfRange();
    uint256 count = endIndex - startIndex;
    uint256[] memory ids = new uint256[](count);
    for (uint256 idx = 0; idx < count; idx++) {
      ids[idx] = s_upkeepIDs.at(idx + startIndex);
    }
    return ids;
  }

  /**
   * @notice retrieve active upkeep IDs within the range, filtered by the provided type
   * @param startIndex starting index in list
   * @param endIndex end index. non-inclusive
   * @dev the order of IDs in the list is **not guaranteed**, therefore, if making successive calls, one
   * should consider keeping the blockheight constant to ensure a holistic picture of the contract state
   * @dev this function will not error if an endIndex is provided that is too large, instead, it will resolve gracefully
   * by suppling as many results as it can within the bounds of the upkeep set
   */
  function getActiveUpkeepIDsByType(
    uint256 startIndex,
    uint256 endIndex,
    Trigger trigger
  ) external view returns (uint256[] memory) {
    uint256 maxIndex = s_upkeepIDs.length();
    endIndex = endIndex > maxIndex ? maxIndex : endIndex;
    if (startIndex > endIndex) revert IndexOutOfRange();
    uint256 maxCount = endIndex - startIndex;
    uint256 numFound;
    uint256[] memory ids = new uint256[](maxCount);
    for (uint256 idx = 0; idx < maxCount; idx++) {
      uint256 id = s_upkeepIDs.at(idx + startIndex);
      if (getTriggerType(id) == trigger) {
        ids[numFound] = id;
        numFound++;
      }
    }
    if (numFound != maxCount) {
      assembly {
        mstore(ids, numFound)
      }
    }
    return ids;
  }

  function getUpkeepTriggerConfig(uint256 upkeepId) public view returns (bytes memory) {
    return s_upkeepTriggerConfig[upkeepId];
  }

  /**
   * @notice read the current info about any transmitter address
   */
  function getTransmitterInfo(
    address query
  ) external view returns (bool active, uint8 index, uint96 balance, uint96 lastCollected, address payee) {
    Transmitter memory transmitter = s_transmitters[query];
    uint96 totalDifference = s_hotVars.totalPremium - transmitter.lastCollected;
    uint96 pooledShare = totalDifference / uint96(s_transmittersList.length);

    return (
      transmitter.active,
      transmitter.index,
      (transmitter.balance + pooledShare),
      transmitter.lastCollected,
      s_transmitterPayees[query]
    );
  }

  /**
   * @notice read the current info about any signer address
   */
  function getSignerInfo(address query) external view returns (bool active, uint8 index) {
    Signer memory signer = s_signers[query];
    return (signer.active, signer.index);
  }

  /**
   * @notice read the current state of the registry
   */
  function getState()
    external
    view
    returns (
      State memory state,
      OnchainConfig memory config,
      address[] memory signers,
      address[] memory transmitters,
      uint8 f
    )
  {
    state = State({
      nonce: s_storage.nonce,
      ownerLinkBalance: s_storage.ownerLinkBalance,
      expectedLinkBalance: s_expectedLinkBalance,
      totalPremium: s_hotVars.totalPremium,
      numUpkeeps: s_upkeepIDs.length(),
      configCount: s_storage.configCount,
      latestConfigBlockNumber: s_storage.latestConfigBlockNumber,
      latestConfigDigest: s_latestConfigDigest,
      latestEpoch: s_hotVars.latestEpoch,
      paused: s_hotVars.paused
    });

    config = OnchainConfig({
      paymentPremiumPPB: s_hotVars.paymentPremiumPPB,
      flatFeeMicroLink: s_hotVars.flatFeeMicroLink,
      checkGasLimit: s_storage.checkGasLimit,
      stalenessSeconds: s_hotVars.stalenessSeconds,
      gasCeilingMultiplier: s_hotVars.gasCeilingMultiplier,
      minUpkeepSpend: s_storage.minUpkeepSpend,
      maxPerformGas: s_storage.maxPerformGas,
      maxCheckDataSize: s_storage.maxCheckDataSize,
      maxPerformDataSize: s_storage.maxPerformDataSize,
      fallbackGasPrice: s_fallbackGasPrice,
      fallbackLinkPrice: s_fallbackLinkPrice,
      transcoder: s_storage.transcoder,
      registrar: s_storage.registrar
    });

    return (state, config, s_signersList, s_transmittersList, s_hotVars.f);
  }

  /**
   * @notice calculates the minimum balance required for an upkeep to remain eligible
   * @param id the upkeep id to calculate minimum balance for
   */
  function getMinBalanceForUpkeep(uint256 id) external view returns (uint96 minBalance) {
    return getMaxPaymentForGas(s_upkeep[id].executeGas);
  }

  /**
   * @notice calculates the maximum payment for a given gas limit
   * @param gasLimit the gas to calculate payment for
   */
  function getMaxPaymentForGas(uint32 gasLimit) public view returns (uint96 maxPayment) {
    HotVars memory hotVars = s_hotVars;
    (uint256 fastGasWei, uint256 linkNative) = _getFeedData(hotVars);
    return _getMaxLinkPayment(hotVars, gasLimit, s_storage.maxPerformDataSize, fastGasWei, linkNative, false);
  }

  /**
   * @notice retrieves the migration permission for a peer registry
   */
  function getPeerRegistryMigrationPermission(address peer) external view returns (MigrationPermission) {
    return s_peerRegistryMigrationPermission[peer];
  }
}
