// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import {Address} from "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";
import {AutomationRegistryBase2_3} from "./AutomationRegistryBase2_3.sol";
import {AutomationRegistryLogicA2_3} from "./AutomationRegistryLogicA2_3.sol";
import {AutomationRegistryLogicC2_3} from "./AutomationRegistryLogicC2_3.sol";
import {Chainable} from "../Chainable.sol";
import {OCR2Abstract} from "../../shared/ocr2/OCR2Abstract.sol";
import {IERC20Metadata as IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/extensions/IERC20Metadata.sol";

/**
 * @notice Registry for adding work for Chainlink nodes to perform on client
 * contracts. Clients must support the AutomationCompatibleInterface interface.
 */
contract AutomationRegistry2_3 is AutomationRegistryBase2_3, OCR2Abstract, Chainable {
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;
  using EnumerableSet for EnumerableSet.AddressSet;

  /**
   * @notice versions:
   * AutomationRegistry 2.3.0: supports native and ERC20 billing
   *                           changes flat fee to USD-denominated
   *                           adds support for custom billing overrides
   * AutomationRegistry 2.2.0: moves chain-specific integration code into a separate module
   * KeeperRegistry 2.1.0:     introduces support for log triggers
   *                           removes the need for "wrapped perform data"
   * KeeperRegistry 2.0.2:     pass revert bytes as performData when target contract reverts
   *                           fixes issue with arbitrum block number
   *                           does an early return in case of stale report instead of revert
   * KeeperRegistry 2.0.1:     implements workaround for buggy migrate function in 1.X
   * KeeperRegistry 2.0.0:     implement OCR interface
   * KeeperRegistry 1.3.0:     split contract into Proxy and Logic
   *                           account for Arbitrum and Optimism L1 gas fee
   *                           allow users to configure upkeeps
   * KeeperRegistry 1.2.0:     allow funding within performUpkeep
   *                           allow configurable registry maxPerformGas
   *                           add function to let admin change upkeep gas limit
   *                           add minUpkeepSpend requirement
   *                           upgrade to solidity v0.8
   * KeeperRegistry 1.1.0:     added flatFeeMicroLink
   * KeeperRegistry 1.0.0:     initial release
   */
  string public constant override typeAndVersion = "AutomationRegistry 2.3.0";

  /**
   * @param logicA the address of the first logic contract
   * @dev we cast the contract to logicC in order to call logicC functions (via fallback)
   */
  constructor(
    AutomationRegistryLogicA2_3 logicA
  )
    AutomationRegistryBase2_3(
      AutomationRegistryLogicC2_3(address(logicA)).getLinkAddress(),
      AutomationRegistryLogicC2_3(address(logicA)).getLinkUSDFeedAddress(),
      AutomationRegistryLogicC2_3(address(logicA)).getNativeUSDFeedAddress(),
      AutomationRegistryLogicC2_3(address(logicA)).getFastGasFeedAddress(),
      AutomationRegistryLogicC2_3(address(logicA)).getAutomationForwarderLogic(),
      AutomationRegistryLogicC2_3(address(logicA)).getAllowedReadOnlyAddress(),
      AutomationRegistryLogicC2_3(address(logicA)).getPayoutMode(),
      AutomationRegistryLogicC2_3(address(logicA)).getWrappedNativeTokenAddress()
    )
    Chainable(address(logicA))
  {}

  /**
   * @notice holds the variables used in the transmit function, necessary to avoid stack too deep errors
   */
  struct TransmitVars {
    uint16 numUpkeepsPassedChecks;
    uint96 totalReimbursement;
    uint96 totalPremium;
    uint256 totalCalldataWeight;
  }

  // ================================================================
  // |                       HOT PATH ACTIONS                       |
  // ================================================================

  /**
   * @inheritdoc OCR2Abstract
   */
  function transmit(
    bytes32[3] calldata reportContext,
    bytes calldata rawReport,
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs
  ) external override {
    uint256 gasOverhead = gasleft();
    // use this msg.data length check to ensure no extra data is included in the call
    // 4 is first 4 bytes of the keccak-256 hash of the function signature. ss.length == rs.length so use one of them
    // 4 + (32 * 3) + (rawReport.length + 32 + 32) + (32 * rs.length + 32 + 32) + (32 * ss.length + 32 + 32) + 32
    uint256 requiredLength = 324 + rawReport.length + 64 * rs.length;
    if (msg.data.length != requiredLength) revert InvalidDataLength();
    HotVars memory hotVars = s_hotVars;

    if (hotVars.paused) revert RegistryPaused();
    if (!s_transmitters[msg.sender].active) revert OnlyActiveTransmitters();

    // Verify signatures
    if (s_latestConfigDigest != reportContext[0]) revert ConfigDigestMismatch();
    if (rs.length != hotVars.f + 1 || rs.length != ss.length) revert IncorrectNumberOfSignatures();
    _verifyReportSignature(reportContext, rawReport, rs, ss, rawVs);

    Report memory report = _decodeReport(rawReport);

    uint40 epochAndRound = uint40(uint256(reportContext[1]));
    uint32 epoch = uint32(epochAndRound >> 8);

    _handleReport(hotVars, report, gasOverhead);

    if (epoch > hotVars.latestEpoch) {
      s_hotVars.latestEpoch = epoch;
    }
  }

  /**
   * @notice handles the report by performing the upkeeps and updating the state
   * @param hotVars the hot variables of the registry
   * @param report the report to be handled (already verified and decoded)
   * @param gasOverhead the running tally of gas overhead to be split across the upkeeps
   * @dev had to split this function from transmit() to avoid stack too deep errors
   * @dev all other internal / private functions are generally defined in the Base contract
   * we leave this here because it is essentially a continuation of the transmit() function,
   */
  function _handleReport(HotVars memory hotVars, Report memory report, uint256 gasOverhead) private {
    UpkeepTransmitInfo[] memory upkeepTransmitInfo = new UpkeepTransmitInfo[](report.upkeepIds.length);
    TransmitVars memory transmitVars = TransmitVars({
      numUpkeepsPassedChecks: 0,
      totalCalldataWeight: 0,
      totalReimbursement: 0,
      totalPremium: 0
    });

    uint256 blocknumber = hotVars.chainModule.blockNumber();
    uint256 l1Fee = hotVars.chainModule.getCurrentL1Fee(msg.data.length);

    for (uint256 i = 0; i < report.upkeepIds.length; i++) {
      upkeepTransmitInfo[i].upkeep = s_upkeep[report.upkeepIds[i]];
      upkeepTransmitInfo[i].triggerType = _getTriggerType(report.upkeepIds[i]);

      (upkeepTransmitInfo[i].earlyChecksPassed, upkeepTransmitInfo[i].dedupID) = _prePerformChecks(
        report.upkeepIds[i],
        blocknumber,
        report.triggers[i],
        upkeepTransmitInfo[i],
        hotVars
      );

      if (upkeepTransmitInfo[i].earlyChecksPassed) {
        transmitVars.numUpkeepsPassedChecks += 1;
      } else {
        continue;
      }

      // Actually perform the target upkeep
      (upkeepTransmitInfo[i].performSuccess, upkeepTransmitInfo[i].gasUsed) = _performUpkeep(
        upkeepTransmitInfo[i].upkeep.forwarder,
        report.gasLimits[i],
        report.performDatas[i]
      );

      // To split L1 fee across the upkeeps, assign a weight to this upkeep based on the length
      // of the perform data and calldata overhead
      upkeepTransmitInfo[i].calldataWeight =
        report.performDatas[i].length +
        TRANSMIT_CALLDATA_FIXED_BYTES_OVERHEAD +
        (TRANSMIT_CALLDATA_PER_SIGNER_BYTES_OVERHEAD * (hotVars.f + 1));
      transmitVars.totalCalldataWeight += upkeepTransmitInfo[i].calldataWeight;

      // Deduct the gasUsed by upkeep from the overhead tally - upkeeps pay for their own gas individually
      gasOverhead -= upkeepTransmitInfo[i].gasUsed;

      // Store last perform block number / deduping key for upkeep
      _updateTriggerMarker(report.upkeepIds[i], blocknumber, upkeepTransmitInfo[i]);
    }
    // No upkeeps to be performed in this report
    if (transmitVars.numUpkeepsPassedChecks == 0) {
      return;
    }

    // This is the overall gas overhead that will be split across performed upkeeps
    // Take upper bound of 16 gas per callData bytes
    gasOverhead = (gasOverhead - gasleft()) + (16 * msg.data.length) + ACCOUNTING_FIXED_GAS_OVERHEAD;
    gasOverhead = gasOverhead / transmitVars.numUpkeepsPassedChecks + ACCOUNTING_PER_UPKEEP_GAS_OVERHEAD;

    {
      BillingTokenPaymentParams memory billingTokenParams;
      uint256 nativeUSD = _getNativeUSD(hotVars);
      for (uint256 i = 0; i < report.upkeepIds.length; i++) {
        if (upkeepTransmitInfo[i].earlyChecksPassed) {
          billingTokenParams = _getBillingTokenPaymentParams(hotVars, upkeepTransmitInfo[i].upkeep.billingToken);
          PaymentReceipt memory receipt = _handlePayment(
            hotVars,
            PaymentParams({
              gasLimit: upkeepTransmitInfo[i].gasUsed,
              gasOverhead: gasOverhead,
              l1CostWei: (l1Fee * upkeepTransmitInfo[i].calldataWeight) / transmitVars.totalCalldataWeight,
              fastGasWei: report.fastGasWei,
              linkUSD: report.linkUSD,
              nativeUSD: nativeUSD,
              billingToken: upkeepTransmitInfo[i].upkeep.billingToken,
              billingTokenParams: billingTokenParams,
              isTransaction: true
            }),
            report.upkeepIds[i],
            upkeepTransmitInfo[i].upkeep
          );
          transmitVars.totalPremium += receipt.premiumInJuels;
          transmitVars.totalReimbursement += receipt.gasReimbursementInJuels;

          emit UpkeepPerformed(
            report.upkeepIds[i],
            upkeepTransmitInfo[i].performSuccess,
            receipt.gasChargeInBillingToken + receipt.premiumInBillingToken,
            upkeepTransmitInfo[i].gasUsed,
            gasOverhead,
            report.triggers[i]
          );
        }
      }
    }
    // record payments to NOPs, all in LINK
    s_transmitters[msg.sender].balance += transmitVars.totalReimbursement;
    s_hotVars.totalPremium += transmitVars.totalPremium;
    s_reserveAmounts[IERC20(address(i_link))] += transmitVars.totalReimbursement + transmitVars.totalPremium;
  }

  // ================================================================
  // |                         OCR2ABSTRACT                         |
  // ================================================================

  /**
   * @inheritdoc OCR2Abstract
   * @dev prefer the type-safe version of setConfig (below) whenever possible. The OnchainConfig could differ between registry versions
   * @dev this function takes up precious space on the root contract, but must be implemented to conform to the OCR2Abstract interface
   */
  function setConfig(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfigBytes,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external override {
    (OnchainConfig memory config, IERC20[] memory billingTokens, BillingConfig[] memory billingConfigs) = abi.decode(
      onchainConfigBytes,
      (OnchainConfig, IERC20[], BillingConfig[])
    );

    setConfigTypeSafe(
      signers,
      transmitters,
      f,
      config,
      offchainConfigVersion,
      offchainConfig,
      billingTokens,
      billingConfigs
    );
  }

  /**
   * @notice sets the configuration for the registry
   * @param signers the list of permitted signers
   * @param transmitters the list of permitted transmitters
   * @param f the maximum tolerance for faulty nodes
   * @param onchainConfig configuration values that are used on-chain
   * @param offchainConfigVersion the version of the offchainConfig
   * @param offchainConfig configuration values that are used off-chain
   * @param billingTokens the list of valid billing tokens
   * @param billingConfigs the configurations for each billing token
   */
  function setConfigTypeSafe(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    OnchainConfig memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig,
    IERC20[] memory billingTokens,
    BillingConfig[] memory billingConfigs
  ) public onlyOwner {
    if (signers.length > MAX_NUM_ORACLES) revert TooManyOracles();
    if (f == 0) revert IncorrectNumberOfFaultyOracles();
    if (signers.length != transmitters.length || signers.length <= 3 * f) revert IncorrectNumberOfSigners();
    if (billingTokens.length != billingConfigs.length) revert ParameterLengthError();
    // set billing config for tokens
    _setBillingConfig(billingTokens, billingConfigs);

    _updateTransmitters(signers, transmitters);

    s_hotVars = HotVars({
      f: f,
      stalenessSeconds: onchainConfig.stalenessSeconds,
      gasCeilingMultiplier: onchainConfig.gasCeilingMultiplier,
      paused: s_hotVars.paused,
      reentrancyGuard: s_hotVars.reentrancyGuard,
      totalPremium: s_hotVars.totalPremium,
      latestEpoch: 0, // DON restarts epoch
      reorgProtectionEnabled: onchainConfig.reorgProtectionEnabled,
      chainModule: onchainConfig.chainModule
    });

    uint32 previousConfigBlockNumber = s_storage.latestConfigBlockNumber;
    uint32 newLatestConfigBlockNumber = uint32(onchainConfig.chainModule.blockNumber());
    uint32 newConfigCount = s_storage.configCount + 1;

    s_storage = Storage({
      checkGasLimit: onchainConfig.checkGasLimit,
      maxPerformGas: onchainConfig.maxPerformGas,
      transcoder: onchainConfig.transcoder,
      maxCheckDataSize: onchainConfig.maxCheckDataSize,
      maxPerformDataSize: onchainConfig.maxPerformDataSize,
      maxRevertDataSize: onchainConfig.maxRevertDataSize,
      upkeepPrivilegeManager: onchainConfig.upkeepPrivilegeManager,
      financeAdmin: onchainConfig.financeAdmin,
      nonce: s_storage.nonce,
      configCount: newConfigCount,
      latestConfigBlockNumber: newLatestConfigBlockNumber
    });
    s_fallbackGasPrice = onchainConfig.fallbackGasPrice;
    s_fallbackLinkPrice = onchainConfig.fallbackLinkPrice;
    s_fallbackNativePrice = onchainConfig.fallbackNativePrice;

    bytes memory onchainConfigBytes = abi.encode(onchainConfig);

    s_latestConfigDigest = _configDigestFromConfigData(
      block.chainid,
      address(this),
      s_storage.configCount,
      signers,
      transmitters,
      f,
      onchainConfigBytes,
      offchainConfigVersion,
      offchainConfig
    );

    for (uint256 idx = s_registrars.length(); idx > 0; idx--) {
      s_registrars.remove(s_registrars.at(idx - 1));
    }

    for (uint256 idx = 0; idx < onchainConfig.registrars.length; idx++) {
      s_registrars.add(onchainConfig.registrars[idx]);
    }

    emit ConfigSet(
      previousConfigBlockNumber,
      s_latestConfigDigest,
      s_storage.configCount,
      signers,
      transmitters,
      f,
      onchainConfigBytes,
      offchainConfigVersion,
      offchainConfig
    );
  }

  /**
   * @inheritdoc OCR2Abstract
   * @dev this function takes up precious space on the root contract, but must be implemented to conform to the OCR2Abstract interface
   */
  function latestConfigDetails()
    external
    view
    override
    returns (uint32 configCount, uint32 blockNumber, bytes32 configDigest)
  {
    return (s_storage.configCount, s_storage.latestConfigBlockNumber, s_latestConfigDigest);
  }

  /**
   * @inheritdoc OCR2Abstract
   * @dev this function takes up precious space on the root contract, but must be implemented to conform to the OCR2Abstract interface
   */
  function latestConfigDigestAndEpoch()
    external
    view
    override
    returns (bool scanLogs, bytes32 configDigest, uint32 epoch)
  {
    return (false, s_latestConfigDigest, s_hotVars.latestEpoch);
  }
}
