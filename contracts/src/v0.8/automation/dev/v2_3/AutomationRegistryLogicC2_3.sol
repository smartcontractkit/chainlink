// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {AutomationRegistryBase2_3} from "./AutomationRegistryBase2_3.sol";
import {EnumerableSet} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import {Address} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";
import {IAutomationForwarder} from "../../interfaces/IAutomationForwarder.sol";
import {IChainModule} from "../../interfaces/IChainModule.sol";
import {IERC20Metadata as IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {IAutomationV21PlusCommon} from "../../interfaces/IAutomationV21PlusCommon.sol";

contract AutomationRegistryLogicC2_3 is AutomationRegistryBase2_3 {
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
    PayoutMode payoutMode,
    address wrappedNativeTokenAddress
  )
    AutomationRegistryBase2_3(
      link,
      linkUSDFeed,
      nativeUSDFeed,
      fastGasFeed,
      automationForwarderLogic,
      allowedReadOnlyAddress,
      payoutMode,
      wrappedNativeTokenAddress
    )
  {}

  // ================================================================
  // |                         NODE ACTIONS                         |
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
    s_reserveAmounts[IERC20(address(i_link))] = s_reserveAmounts[IERC20(address(i_link))] - balance;
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
    _onlyPrivilegeManagerAllowed();
    s_upkeepPrivilegeConfig[upkeepId] = newPrivilegeConfig;
    emit UpkeepPrivilegeConfigSet(upkeepId, newPrivilegeConfig);
  }

  /**
   * @notice this is used by the owner to set the initial payees for newly added transmitters. The owner is not allowed to change payees for existing transmitters.
   * @dev the IGNORE_ADDRESS is a "helper" that makes it easier to construct a list of payees when you only care about setting the payee for a small number of transmitters.
   */
  function setPayees(address[] calldata payees) external onlyOwner {
    if (s_transmittersList.length != payees.length) revert ParameterLengthError();
    for (uint256 i = 0; i < s_transmittersList.length; i++) {
      address transmitter = s_transmittersList[i];
      address oldPayee = s_transmitterPayees[transmitter];
      address newPayee = payees[i];

      if (
        (newPayee == ZERO_ADDRESS) || (oldPayee != ZERO_ADDRESS && oldPayee != newPayee && newPayee != IGNORE_ADDRESS)
      ) {
        revert InvalidPayee();
      }

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
    _onlyPrivilegeManagerAllowed();
    s_adminPrivilegeConfig[admin] = newPrivilegeConfig;
    emit AdminPrivilegeConfigSet(admin, newPrivilegeConfig);
  }

  /**
   * @notice settles NOPs' LINK payment offchain
   */
  function settleNOPsOffchain() external {
    _onlyFinanceAdminAllowed();
    if (s_payoutMode == PayoutMode.ON_CHAIN) revert MustSettleOnchain();

    uint96 totalPremium = s_hotVars.totalPremium;
    uint256 activeTransmittersLength = s_transmittersList.length;
    uint256 deactivatedTransmittersLength = s_deactivatedTransmitters.length();
    uint256 length = activeTransmittersLength + deactivatedTransmittersLength;
    uint256[] memory payments = new uint256[](length);
    address[] memory payees = new address[](length);

    for (uint256 i = 0; i < activeTransmittersLength; i++) {
      address transmitterAddr = s_transmittersList[i];
      uint96 balance = _updateTransmitterBalanceFromPool(
        transmitterAddr,
        totalPremium,
        uint96(activeTransmittersLength)
      );

      payments[i] = balance;
      payees[i] = s_transmitterPayees[transmitterAddr];
      s_transmitters[transmitterAddr].balance = 0;
    }

    for (uint256 i = 0; i < deactivatedTransmittersLength; i++) {
      address deactivatedAddr = s_deactivatedTransmitters.at(i);
      Transmitter memory transmitter = s_transmitters[deactivatedAddr];

      payees[i + activeTransmittersLength] = s_transmitterPayees[deactivatedAddr];
      payments[i + activeTransmittersLength] = transmitter.balance;
      s_transmitters[deactivatedAddr].balance = 0;
    }

    // reserve amount of LINK is reset to 0 since no user deposits of LINK are expected in offchain mode
    s_reserveAmounts[IERC20(address(i_link))] = 0;

    for (uint256 idx = s_deactivatedTransmitters.length(); idx > 0; idx--) {
      s_deactivatedTransmitters.remove(s_deactivatedTransmitters.at(idx - 1));
    }

    emit NOPsSettledOffchain(payees, payments);
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

  function getWrappedNativeTokenAddress() external view returns (address) {
    return address(i_wrappedNativeToken);
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

  function getPayoutMode() external view returns (PayoutMode) {
    return s_payoutMode;
  }

  function upkeepVersion() public pure returns (uint8) {
    return UPKEEP_VERSION_BASE;
  }

  /**
   * @notice gets the number of upkeeps on the registry
   */
  function getNumUpkeeps() external view returns (uint256) {
    return s_upkeepIDs.length();
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
   * @notice read the current on-chain config of the registry
   * @dev this function will change between versions, it should never be used where
   * backwards compatibility matters!
   */
  function getConfig() external view returns (OnchainConfig memory) {
    return
      OnchainConfig({
        checkGasLimit: s_storage.checkGasLimit,
        stalenessSeconds: s_hotVars.stalenessSeconds,
        gasCeilingMultiplier: s_hotVars.gasCeilingMultiplier,
        maxPerformGas: s_storage.maxPerformGas,
        maxCheckDataSize: s_storage.maxCheckDataSize,
        maxPerformDataSize: s_storage.maxPerformDataSize,
        maxRevertDataSize: s_storage.maxRevertDataSize,
        fallbackGasPrice: s_fallbackGasPrice,
        fallbackLinkPrice: s_fallbackLinkPrice,
        fallbackNativePrice: s_fallbackNativePrice,
        transcoder: s_storage.transcoder,
        registrars: s_registrars.values(),
        upkeepPrivilegeManager: s_storage.upkeepPrivilegeManager,
        chainModule: s_hotVars.chainModule,
        reorgProtectionEnabled: s_hotVars.reorgProtectionEnabled,
        financeAdmin: s_storage.financeAdmin
      });
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
    return getMaxPaymentForGas(id, _getTriggerType(id), upkeep.performGas, upkeep.billingToken);
  }

  /**
   * @notice calculates the maximum payment for a given gas limit
   * @param gasLimit the gas to calculate payment for
   */
  function getMaxPaymentForGas(
    uint256 id,
    Trigger triggerType,
    uint32 gasLimit,
    IERC20 billingToken
  ) public view returns (uint96 maxPayment) {
    HotVars memory hotVars = s_hotVars;
    (uint256 fastGasWei, uint256 linkUSD, uint256 nativeUSD) = _getFeedData(hotVars);
    return _getMaxPayment(id, hotVars, triggerType, gasLimit, fastGasWei, linkUSD, nativeUSD, billingToken);
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
   * @notice returns the admin's privilege config
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
   * @notice returns if the dedupKey exists or not
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
   * @notice returns the amount of a particular token that is reserved as
   * user deposits / NOP payments
   */
  function getReserveAmount(IERC20 billingToken) external view returns (uint256) {
    return s_reserveAmounts[billingToken];
  }

  /**
   * @notice returns the amount of a particular token that is withdraw-able by finance admin
   */
  function getAvailableERC20ForPayment(IERC20 billingToken) external view returns (uint256) {
    return billingToken.balanceOf(address(this)) - s_reserveAmounts[IERC20(address(billingToken))];
  }

  /**
   * @notice returns the size of the LINK liquidity pool
   */
  function linkAvailableForPayment() public view returns (int256) {
    return _linkAvailableForPayment();
  }

  /**
   * @notice returns the BillingOverrides config for a given upkeep
   */
  function getBillingOverrides(uint256 upkeepID) external view returns (BillingOverrides memory) {
    return s_billingOverrides[upkeepID];
  }

  /**
   * @notice returns the BillingConfig for a given billing token, this includes decimals and price feed etc
   */
  function getBillingConfig(IERC20 billingToken) external view returns (BillingConfig memory) {
    return s_billingConfigs[billingToken];
  }

  /**
   * @notice returns all active transmitters with their associated payees
   */
  function getTransmittersWithPayees() external view returns (TransmitterPayeeInfo[] memory) {
    uint256 transmitterCount = s_transmittersList.length;
    TransmitterPayeeInfo[] memory transmitters = new TransmitterPayeeInfo[](transmitterCount);

    for (uint256 i = 0; i < transmitterCount; i++) {
      address transmitterAddress = s_transmittersList[i];
      address payeeAddress = s_transmitterPayees[transmitterAddress];

      transmitters[i] = TransmitterPayeeInfo(transmitterAddress, payeeAddress);
    }

    return transmitters;
  }
}
