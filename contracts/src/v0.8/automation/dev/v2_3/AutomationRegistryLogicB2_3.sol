// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {AutomationRegistryBase2_3} from "./AutomationRegistryBase2_3.sol";
import {EnumerableSet} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import {Address} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";
import {UpkeepFormat} from "../../interfaces/UpkeepTranscoderInterface.sol";
import {IAutomationForwarder} from "../../interfaces/IAutomationForwarder.sol";
import {IChainModule} from "../../interfaces/IChainModule.sol";
import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {IAutomationV21PlusCommon} from "../../interfaces/IAutomationV21PlusCommon.sol";

contract AutomationRegistryLogicB2_3 is AutomationRegistryBase2_3 {
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;
  using EnumerableSet for EnumerableSet.AddressSet;

  /**
   * @dev see AutomationRegistry master contract for constructor description
   */
  constructor(
    address link,
    address linkUSDFeed,
    address nativeUSDFeed,
    address fastGasFeed,
    address automationForwarderLogic,
    address allowedReadOnlyAddress,
    PayoutMode payoutMode
  )
    AutomationRegistryBase2_3(
      link,
      linkUSDFeed,
      nativeUSDFeed,
      fastGasFeed,
      automationForwarderLogic,
      allowedReadOnlyAddress,
      payoutMode
    )
  {}

  // ================================================================
  // |                      UPKEEP MANAGEMENT                       |
  // ================================================================

  /**
   * @notice transfers the address of an admin for an upkeep
   */
  function transferUpkeepAdmin(uint256 id, address proposed) external {
    _requireAdminAndNotCancelled(id);
    if (proposed == msg.sender) revert ValueNotChanged();

    if (s_proposedAdmin[id] != proposed) {
      s_proposedAdmin[id] = proposed;
      emit UpkeepAdminTransferRequested(id, msg.sender, proposed);
    }
  }

  /**
   * @notice accepts the transfer of an upkeep admin
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
   * @notice pauses an upkeep - an upkeep will be neither checked nor performed while paused
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
   * @notice unpauses an upkeep
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
   * @notice updates the checkData for an upkeep
   */
  function setUpkeepCheckData(uint256 id, bytes calldata newCheckData) external {
    _requireAdminAndNotCancelled(id);
    if (newCheckData.length > s_storage.maxCheckDataSize) revert CheckDataExceedsLimit();
    s_checkData[id] = newCheckData;
    emit UpkeepCheckDataSet(id, newCheckData);
  }

  /**
   * @notice updates the gas limit for an upkeep
   */
  function setUpkeepGasLimit(uint256 id, uint32 gasLimit) external {
    if (gasLimit < PERFORM_GAS_MIN || gasLimit > s_storage.maxPerformGas) revert GasLimitOutsideRange();
    _requireAdminAndNotCancelled(id);
    s_upkeep[id].performGas = gasLimit;

    emit UpkeepGasLimitSet(id, gasLimit);
  }

  /**
   * @notice updates the offchain config for an upkeep
   */
  function setUpkeepOffchainConfig(uint256 id, bytes calldata config) external {
    _requireAdminAndNotCancelled(id);
    s_upkeepOffchainConfig[id] = config;
    emit UpkeepOffchainConfigSet(id, config);
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

  /**
   * @notice withdraws an upkeep's funds from an upkeep
   * @dev note that an upkeep must be cancelled first!!
   */
  function withdrawFunds(uint256 id, address to) external nonReentrant {
    if (to == ZERO_ADDRESS) revert InvalidRecipient();
    Upkeep memory upkeep = s_upkeep[id];
    if (s_upkeepAdmin[id] != msg.sender) revert OnlyCallableByAdmin();
    if (upkeep.maxValidBlocknumber > s_hotVars.chainModule.blockNumber()) revert UpkeepNotCanceled();
    uint96 amountToWithdraw = s_upkeep[id].balance;
    s_reserveAmounts[address(upkeep.billingToken)] = s_reserveAmounts[address(upkeep.billingToken)] - amountToWithdraw;
    s_upkeep[id].balance = 0;
    bool success = upkeep.billingToken.transfer(to, amountToWithdraw);
    if (!success) revert TransferFailed();
    emit FundsWithdrawn(id, amountToWithdraw, to);
  }

  /**
   * @notice LINK available to withdraw by the finance team
   */
  function linkAvailableForPayment() public view returns (uint256) {
    return i_link.balanceOf(address(this)) - s_reserveAmounts[address(i_link)];
  }

  function withdrawLinkFees(address to, uint256 amount) external {
    _onlyFinanceAdminAllowed();
    if (to == ZERO_ADDRESS) revert InvalidRecipient();

    uint256 available = linkAvailableForPayment();
    if (amount > available) revert InsufficientBalance(available, amount);

    bool transferStatus = i_link.transfer(to, amount);
    if (!transferStatus) {
      revert TransferFailed();
    }
    emit FeesWithdrawn(to, address(i_link), amount);
  }

  function withdrawERC20Fees(address assetAddress, address to, uint256 amount) external {
    _onlyFinanceAdminAllowed();
    if (to == ZERO_ADDRESS) revert InvalidRecipient();

    bool transferStatus = IERC20(assetAddress).transfer(to, amount);
    if (!transferStatus) {
      revert TransferFailed();
    }

    emit FeesWithdrawn(to, assetAddress, amount);
  }

  // ================================================================
  // |                       NODE MANAGEMENT                        |
  // ================================================================

  /**
   * @notice transfers the address of payee for a transmitter
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
   * @notice accepts the transfer of the payee
   */
  function acceptPayeeship(address transmitter) external {
    if (s_proposedPayee[transmitter] != msg.sender) revert OnlyCallableByProposedPayee();
    address past = s_transmitterPayees[transmitter];
    s_transmitterPayees[transmitter] = msg.sender;
    s_proposedPayee[transmitter] = ZERO_ADDRESS;

    emit PayeeshipTransferred(transmitter, past, msg.sender);
  }

  /**
   * @notice withdraws LINK received as payment for work performed
   */
  function withdrawPayment(address from, address to) external {
    if (to == ZERO_ADDRESS) revert InvalidRecipient();
    if (s_payoutMode == PayoutMode.OFF_CHAIN) revert MustSettleOffchain();
    if (s_transmitterPayees[from] != msg.sender) revert OnlyCallableByPayee();
    uint96 balance = _updateTransmitterBalanceFromPool(from, s_hotVars.totalPremium, uint96(s_transmittersList.length));
    s_transmitters[from].balance = 0;
    s_reserveAmounts[address(i_link)] = s_reserveAmounts[address(i_link)] - balance;
    i_link.transfer(to, balance);
    emit PaymentWithdrawn(from, balance, to, msg.sender);
  }

  // ================================================================
  // |                   OWNER / MANAGER ACTIONS                    |
  // ================================================================

  /**
   * @notice sets the privilege config for an upkeep
   */
  function setUpkeepPrivilegeConfig(uint256 upkeepId, bytes calldata newPrivilegeConfig) external {
    if (msg.sender != s_storage.upkeepPrivilegeManager) {
      revert OnlyCallableByUpkeepPrivilegeManager();
    }
    s_upkeepPrivilegeConfig[upkeepId] = newPrivilegeConfig;
    emit UpkeepPrivilegeConfigSet(upkeepId, newPrivilegeConfig);
  }

  /**
   * @notice sets the payees for the transmitters
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
   * @notice sets the migration permission for a peer registry
   * @dev this must be done before upkeeps can be migrated to/from another registry
   */
  function setPeerRegistryMigrationPermission(address peer, MigrationPermission permission) external onlyOwner {
    s_peerRegistryMigrationPermission[peer] = permission;
  }

  /**
   * @notice pauses the entire registry
   */
  function pause() external onlyOwner {
    s_hotVars.paused = true;
    emit Paused(msg.sender);
  }

  /**
   * @notice unpauses the entire registry
   */
  function unpause() external onlyOwner {
    s_hotVars.paused = false;
    emit Unpaused(msg.sender);
  }

  /**
   * @notice sets a generic bytes field used to indicate the privilege that this admin address had
   * @param admin the address to set privilege for
   * @param newPrivilegeConfig the privileges that this admin has
   */
  function setAdminPrivilegeConfig(address admin, bytes calldata newPrivilegeConfig) external {
    if (msg.sender != s_storage.upkeepPrivilegeManager) {
      revert OnlyCallableByUpkeepPrivilegeManager();
    }
    s_adminPrivilegeConfig[admin] = newPrivilegeConfig;
    emit AdminPrivilegeConfigSet(admin, newPrivilegeConfig);
  }

  /**
   * @notice settles NOPs' LINK payment offchain
   */
  function settleNOPsOffchain() external {
    _onlyFinanceAdminAllowed();
    if (s_payoutMode == PayoutMode.ON_CHAIN) revert MustSettleOnchain();

    uint256 length = s_transmittersList.length;
    uint256[] memory balances = new uint256[](length);
    for (uint256 i = 0; i < length; i++) {
      address transmitterAddr = s_transmittersList[i];
      uint96 balance = _updateTransmitterBalanceFromPool(transmitterAddr, s_hotVars.totalPremium, uint96(length));
      balances[i] = balance;
      s_transmitters[transmitterAddr].balance = 0;
    }

    emit NOPsSettledOffchain(s_transmittersList, balances);
  }

  /**
   * @notice disables offchain payment for NOPs
   */
  function disableOffchainPayments() external onlyOwner {
    s_payoutMode = PayoutMode.ON_CHAIN;
  }

  // ================================================================
  // |                           GETTERS                            |
  // ================================================================

  function getConditionalGasOverhead() external pure returns (uint256) {
    return REGISTRY_CONDITIONAL_OVERHEAD;
  }

  function getLogGasOverhead() external pure returns (uint256) {
    return REGISTRY_LOG_OVERHEAD;
  }

  function getPerPerformByteGasOverhead() external pure returns (uint256) {
    return REGISTRY_PER_PERFORM_BYTE_GAS_OVERHEAD;
  }

  function getPerSignerGasOverhead() external pure returns (uint256) {
    return REGISTRY_PER_SIGNER_GAS_OVERHEAD;
  }

  function getTransmitCalldataFixedBytesOverhead() external pure returns (uint256) {
    return TRANSMIT_CALLDATA_FIXED_BYTES_OVERHEAD;
  }

  function getTransmitCalldataPerSignerBytesOverhead() external pure returns (uint256) {
    return TRANSMIT_CALLDATA_PER_SIGNER_BYTES_OVERHEAD;
  }

  function getCancellationDelay() external pure returns (uint256) {
    return CANCELLATION_DELAY;
  }

  function getLinkAddress() external view returns (address) {
    return address(i_link);
  }

  function getLinkUSDFeedAddress() external view returns (address) {
    return address(i_linkUSDFeed);
  }

  function getNativeUSDFeedAddress() external view returns (address) {
    return address(i_nativeUSDFeed);
  }

  function getFastGasFeedAddress() external view returns (address) {
    return address(i_fastGasFeed);
  }

  function getAutomationForwarderLogic() external view returns (address) {
    return i_automationForwarderLogic;
  }

  function getAllowedReadOnlyAddress() external view returns (address) {
    return i_allowedReadOnlyAddress;
  }

  function getPayoutMode() external view returns (PayoutMode) {
    return s_payoutMode;
  }

  function getBillingToken(uint256 upkeepID) external view returns (IERC20) {
    return s_upkeep[upkeepID].billingToken;
  }

  function getBillingTokens() external view returns (IERC20[] memory) {
    return s_billingTokens;
  }

  function supportsBillingToken(IERC20 token) external view returns (bool) {
    return address(s_billingConfigs[token].priceFeed) != address(0);
  }

  function getBillingTokenConfig(IERC20 token) external view returns (BillingConfig memory) {
    return s_billingConfigs[token];
  }

  function upkeepTranscoderVersion() public pure returns (UpkeepFormat) {
    return UPKEEP_TRANSCODER_VERSION_BASE;
  }

  function upkeepVersion() public pure returns (uint8) {
    return UPKEEP_VERSION_BASE;
  }

  /**
   * @notice read all of the details about an upkeep
   * @dev this function may be deprecated in a future version of automation in favor of individual
   * getters for each field
   */
  function getUpkeep(uint256 id) external view returns (IAutomationV21PlusCommon.UpkeepInfoLegacy memory upkeepInfo) {
    Upkeep memory reg = s_upkeep[id];
    address target = address(reg.forwarder) == address(0) ? address(0) : reg.forwarder.getTarget();
    upkeepInfo = IAutomationV21PlusCommon.UpkeepInfoLegacy({
      target: target,
      performGas: reg.performGas,
      checkData: s_checkData[id],
      balance: reg.balance,
      admin: s_upkeepAdmin[id],
      maxValidBlocknumber: reg.maxValidBlocknumber,
      lastPerformedBlockNumber: reg.lastPerformedBlockNumber,
      amountSpent: uint96(reg.amountSpent), // force casting to uint96 for backwards compatibility. Not an issue if it overflows.
      paused: reg.paused,
      offchainConfig: s_upkeepOffchainConfig[id]
    });
    return upkeepInfo;
  }

  /**
   * @notice retrieve active upkeep IDs. Active upkeep is defined as an upkeep which is not paused and not canceled.
   * @param startIndex starting index in list
   * @param maxCount max count to retrieve (0 = unlimited)
   * @dev the order of IDs in the list is **not guaranteed**, therefore, if making successive calls, one
   * should consider keeping the blockheight constant to ensure a holistic picture of the contract state
   */
  function getActiveUpkeepIDs(uint256 startIndex, uint256 maxCount) external view returns (uint256[] memory) {
    uint256 numUpkeeps = s_upkeepIDs.length();
    if (startIndex >= numUpkeeps) revert IndexOutOfRange();
    uint256 endIndex = startIndex + maxCount;
    endIndex = endIndex > numUpkeeps || maxCount == 0 ? numUpkeeps : endIndex;
    uint256[] memory ids = new uint256[](endIndex - startIndex);
    for (uint256 idx = 0; idx < ids.length; idx++) {
      ids[idx] = s_upkeepIDs.at(idx + startIndex);
    }
    return ids;
  }

  /**
   * @notice returns the upkeep's trigger type
   */
  function getTriggerType(uint256 upkeepId) external pure returns (Trigger) {
    return _getTriggerType(upkeepId);
  }

  /**
   * @notice returns the trigger config for an upkeeep
   */
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

    uint96 pooledShare = 0;
    if (transmitter.active) {
      uint96 totalDifference = s_hotVars.totalPremium - transmitter.lastCollected;
      pooledShare = totalDifference / uint96(s_transmittersList.length);
    }

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
   * @dev this function is deprecated
   */
  function getState()
    external
    view
    returns (
      IAutomationV21PlusCommon.StateLegacy memory state,
      IAutomationV21PlusCommon.OnchainConfigLegacy memory config,
      address[] memory signers,
      address[] memory transmitters,
      uint8 f
    )
  {
    state = IAutomationV21PlusCommon.StateLegacy({
      nonce: s_storage.nonce,
      ownerLinkBalance: 0, // deprecated
      expectedLinkBalance: 0, // deprecated
      totalPremium: s_hotVars.totalPremium,
      numUpkeeps: s_upkeepIDs.length(),
      configCount: s_storage.configCount,
      latestConfigBlockNumber: s_storage.latestConfigBlockNumber,
      latestConfigDigest: s_latestConfigDigest,
      latestEpoch: s_hotVars.latestEpoch,
      paused: s_hotVars.paused
    });

    config = IAutomationV21PlusCommon.OnchainConfigLegacy({
      paymentPremiumPPB: 0, // deprecated
      flatFeeMicroLink: 0, // deprecated
      checkGasLimit: s_storage.checkGasLimit,
      stalenessSeconds: s_hotVars.stalenessSeconds,
      gasCeilingMultiplier: s_hotVars.gasCeilingMultiplier,
      minUpkeepSpend: 0, // deprecated
      maxPerformGas: s_storage.maxPerformGas,
      maxCheckDataSize: s_storage.maxCheckDataSize,
      maxPerformDataSize: s_storage.maxPerformDataSize,
      maxRevertDataSize: s_storage.maxRevertDataSize,
      fallbackGasPrice: s_fallbackGasPrice,
      fallbackLinkPrice: s_fallbackLinkPrice,
      transcoder: s_storage.transcoder,
      registrars: s_registrars.values(),
      upkeepPrivilegeManager: s_storage.upkeepPrivilegeManager
    });

    return (state, config, s_signersList, s_transmittersList, s_hotVars.f);
  }

  /**
   * @notice read the Storage data
   * @dev this function signature will change with each version of automation
   * this should not be treated as a stable function
   */
  function getStorage() external view returns (Storage memory) {
    return s_storage;
  }

  /**
   * @notice read the HotVars data
   * @dev this function signature will change with each version of automation
   * this should not be treated as a stable function
   */
  function getHotVars() external view returns (HotVars memory) {
    return s_hotVars;
  }

  /**
   * @notice get the chain module
   */
  function getChainModule() external view returns (IChainModule chainModule) {
    return s_hotVars.chainModule;
  }

  /**
   * @notice if this registry has reorg protection enabled
   */
  function getReorgProtectionEnabled() external view returns (bool reorgProtectionEnabled) {
    return s_hotVars.reorgProtectionEnabled;
  }

  /**
   * @notice calculates the minimum balance required for an upkeep to remain eligible
   * @param id the upkeep id to calculate minimum balance for
   */
  function getBalance(uint256 id) external view returns (uint96 balance) {
    return s_upkeep[id].balance;
  }

  /**
   * @notice calculates the minimum balance required for an upkeep to remain eligible
   * @param id the upkeep id to calculate minimum balance for
   */
  function getMinBalance(uint256 id) external view returns (uint96) {
    return getMinBalanceForUpkeep(id);
  }

  /**
   * @notice calculates the minimum balance required for an upkeep to remain eligible
   * @param id the upkeep id to calculate minimum balance for
   * @dev this will be deprecated in a future version in favor of getMinBalance
   */
  function getMinBalanceForUpkeep(uint256 id) public view returns (uint96 minBalance) {
    Upkeep memory upkeep = s_upkeep[id];
    return getMaxPaymentForGas(_getTriggerType(id), upkeep.performGas, upkeep.billingToken);
  }

  /**
   * @notice calculates the maximum payment for a given gas limit
   * @param gasLimit the gas to calculate payment for
   */
  function getMaxPaymentForGas(
    Trigger triggerType,
    uint32 gasLimit,
    IERC20 billingToken
  ) public view returns (uint96 maxPayment) {
    HotVars memory hotVars = s_hotVars;
    (uint256 fastGasWei, uint256 linkUSD, uint256 nativeUSD) = _getFeedData(hotVars);
    return _getMaxPayment(hotVars, triggerType, gasLimit, fastGasWei, linkUSD, nativeUSD, billingToken);
  }

  /**
   * @notice retrieves the migration permission for a peer registry
   */
  function getPeerRegistryMigrationPermission(address peer) external view returns (MigrationPermission) {
    return s_peerRegistryMigrationPermission[peer];
  }

  /**
   * @notice returns the upkeep privilege config
   */
  function getUpkeepPrivilegeConfig(uint256 upkeepId) external view returns (bytes memory) {
    return s_upkeepPrivilegeConfig[upkeepId];
  }

  /**
   * @notice returns the upkeep privilege config
   */
  function getAdminPrivilegeConfig(address admin) external view returns (bytes memory) {
    return s_adminPrivilegeConfig[admin];
  }

  /**
   * @notice returns the upkeep's forwarder contract
   */
  function getForwarder(uint256 upkeepID) external view returns (IAutomationForwarder) {
    return s_upkeep[upkeepID].forwarder;
  }

  /**
   * @notice returns the upkeep's forwarder contract
   */
  function hasDedupKey(bytes32 dedupKey) external view returns (bool) {
    return s_dedupKeys[dedupKey];
  }

  /**
   * @notice returns the fallback native price
   */
  function getFallbackNativePrice() external view returns (uint256) {
    return s_fallbackNativePrice;
  }

  /**
   * @notice returns the fallback native price
   */
  function getReserveAmount(address billingToken) external view returns (uint256) {
    return s_reserveAmounts[billingToken];
  }
}
