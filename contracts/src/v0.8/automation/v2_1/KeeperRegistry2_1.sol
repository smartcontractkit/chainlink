// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.16;

import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import {Address} from "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";
import {KeeperRegistryBase2_1} from "./KeeperRegistryBase2_1.sol";
import {KeeperRegistryLogicB2_1} from "./KeeperRegistryLogicB2_1.sol";
import {Chainable} from "../Chainable.sol";
import {IERC677Receiver} from "../../shared/interfaces/IERC677Receiver.sol";
import {OCR2Abstract} from "../../shared/ocr2/OCR2Abstract.sol";
import {IAutomationV21PlusCommon} from "../interfaces/IAutomationV21PlusCommon.sol";

/**
 * @notice Registry for adding work for Chainlink Keepers to perform on client
 * contracts. Clients must support the Upkeep interface.
 */
contract KeeperRegistry2_1 is KeeperRegistryBase2_1, OCR2Abstract, Chainable, IERC677Receiver {
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;
  using EnumerableSet for EnumerableSet.AddressSet;

  /**
   * @notice versions:
   * - KeeperRegistry 2.1.0: introduces support for log, cron, and ready triggers
                           : removes the need for "wrapped perform data"
   * - KeeperRegistry 2.0.2: pass revert bytes as performData when target contract reverts
   *                       : fixes issue with arbitrum block number
   *                       : does an early return in case of stale report instead of revert
   * - KeeperRegistry 2.0.1: implements workaround for buggy migrate function in 1.X
   * - KeeperRegistry 2.0.0: implement OCR interface
   * - KeeperRegistry 1.3.0: split contract into Proxy and Logic
   *                       : account for Arbitrum and Optimism L1 gas fee
   *                       : allow users to configure upkeeps
   * - KeeperRegistry 1.2.0: allow funding within performUpkeep
   *                       : allow configurable registry maxPerformGas
   *                       : add function to let admin change upkeep gas limit
   *                       : add minUpkeepSpend requirement
   *                       : upgrade to solidity v0.8
   * - KeeperRegistry 1.1.0: added flatFeeMicroLink
   * - KeeperRegistry 1.0.0: initial release
   */
  string public constant override typeAndVersion = "KeeperRegistry 2.1.0";

  /**
   * @param logicA the address of the first logic contract, but cast as logicB in order to call logicB functions
   */
  constructor(
    KeeperRegistryLogicB2_1 logicA
  )
    KeeperRegistryBase2_1(
      logicA.getMode(),
      logicA.getLinkAddress(),
      logicA.getLinkNativeFeedAddress(),
      logicA.getFastGasFeedAddress(),
      logicA.getAutomationForwarderLogic()
    )
    Chainable(address(logicA))
  {}

  // ================================================================
  // |                           ACTIONS                            |
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
    HotVars memory hotVars = s_hotVars;

    if (hotVars.paused) revert RegistryPaused();
    if (!s_transmitters[msg.sender].active) revert OnlyActiveTransmitters();

    // Verify signatures
    if (s_latestConfigDigest != reportContext[0]) revert ConfigDigestMismatch();
    if (rs.length != hotVars.f + 1 || rs.length != ss.length) revert IncorrectNumberOfSignatures();
    _verifyReportSignature(reportContext, rawReport, rs, ss, rawVs);

    Report memory report = _decodeReport(rawReport);
    UpkeepTransmitInfo[] memory upkeepTransmitInfo = new UpkeepTransmitInfo[](report.upkeepIds.length);
    uint16 numUpkeepsPassedChecks;

    for (uint256 i = 0; i < report.upkeepIds.length; i++) {
      upkeepTransmitInfo[i].upkeep = s_upkeep[report.upkeepIds[i]];
      upkeepTransmitInfo[i].triggerType = _getTriggerType(report.upkeepIds[i]);
      upkeepTransmitInfo[i].maxLinkPayment = _getMaxLinkPayment(
        hotVars,
        upkeepTransmitInfo[i].triggerType,
        uint32(report.gasLimits[i]),
        uint32(report.performDatas[i].length),
        report.fastGasWei,
        report.linkNative,
        true
      );
      (upkeepTransmitInfo[i].earlyChecksPassed, upkeepTransmitInfo[i].dedupID) = _prePerformChecks(
        report.upkeepIds[i],
        report.triggers[i],
        upkeepTransmitInfo[i]
      );

      if (upkeepTransmitInfo[i].earlyChecksPassed) {
        numUpkeepsPassedChecks += 1;
      } else {
        continue;
      }

      // Actually perform the target upkeep
      (upkeepTransmitInfo[i].performSuccess, upkeepTransmitInfo[i].gasUsed) = _performUpkeep(
        upkeepTransmitInfo[i].upkeep.forwarder,
        report.gasLimits[i],
        report.performDatas[i]
      );

      // Deduct that gasUsed by upkeep from our running counter
      gasOverhead -= upkeepTransmitInfo[i].gasUsed;

      // Store last perform block number / deduping key for upkeep
      _updateTriggerMarker(report.upkeepIds[i], upkeepTransmitInfo[i]);
    }
    // No upkeeps to be performed in this report
    if (numUpkeepsPassedChecks == 0) {
      return;
    }

    // This is the overall gas overhead that will be split across performed upkeeps
    // Take upper bound of 16 gas per callData bytes, which is approximated to be reportLength
    // Rest of msg.data is accounted for in accounting overheads
    gasOverhead =
      (gasOverhead - gasleft() + 16 * rawReport.length) +
      ACCOUNTING_FIXED_GAS_OVERHEAD +
      (ACCOUNTING_PER_SIGNER_GAS_OVERHEAD * (hotVars.f + 1));
    gasOverhead = gasOverhead / numUpkeepsPassedChecks + ACCOUNTING_PER_UPKEEP_GAS_OVERHEAD;

    uint96 totalReimbursement;
    uint96 totalPremium;
    {
      uint96 reimbursement;
      uint96 premium;
      for (uint256 i = 0; i < report.upkeepIds.length; i++) {
        if (upkeepTransmitInfo[i].earlyChecksPassed) {
          upkeepTransmitInfo[i].gasOverhead = _getCappedGasOverhead(
            gasOverhead,
            upkeepTransmitInfo[i].triggerType,
            uint32(report.performDatas[i].length),
            hotVars.f
          );

          (reimbursement, premium) = _postPerformPayment(
            hotVars,
            report.upkeepIds[i],
            upkeepTransmitInfo[i],
            report.fastGasWei,
            report.linkNative,
            numUpkeepsPassedChecks
          );
          totalPremium += premium;
          totalReimbursement += reimbursement;

          emit UpkeepPerformed(
            report.upkeepIds[i],
            upkeepTransmitInfo[i].performSuccess,
            reimbursement + premium,
            upkeepTransmitInfo[i].gasUsed,
            upkeepTransmitInfo[i].gasOverhead,
            report.triggers[i]
          );
        }
      }
    }
    // record payments
    s_transmitters[msg.sender].balance += totalReimbursement;
    s_hotVars.totalPremium += totalPremium;

    uint40 epochAndRound = uint40(uint256(reportContext[1]));
    uint32 epoch = uint32(epochAndRound >> 8);
    if (epoch > hotVars.latestEpoch) {
      s_hotVars.latestEpoch = epoch;
    }
  }

  /**
   * @notice simulates the upkeep with the perform data returned from checkUpkeep
   * @param id identifier of the upkeep to execute the data with.
   * @param performData calldata parameter to be passed to the target upkeep.
   * @return success whether the call reverted or not
   * @return gasUsed the amount of gas the target contract consumed
   */
  function simulatePerformUpkeep(
    uint256 id,
    bytes calldata performData
  ) external cannotExecute returns (bool success, uint256 gasUsed) {
    if (s_hotVars.paused) revert RegistryPaused();
    Upkeep memory upkeep = s_upkeep[id];
    (success, gasUsed) = _performUpkeep(upkeep.forwarder, upkeep.performGas, performData);
    return (success, gasUsed);
  }

  /**
   * @notice uses LINK's transferAndCall to LINK and add funding to an upkeep
   * @dev safe to cast uint256 to uint96 as total LINK supply is under UINT96MAX
   * @param sender the account which transferred the funds
   * @param amount number of LINK transfer
   */
  function onTokenTransfer(address sender, uint256 amount, bytes calldata data) external override {
    if (msg.sender != address(i_link)) revert OnlyCallableByLINKToken();
    if (data.length != 32) revert InvalidDataLength();
    uint256 id = abi.decode(data, (uint256));
    if (s_upkeep[id].maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();
    s_upkeep[id].balance = s_upkeep[id].balance + uint96(amount);
    s_expectedLinkBalance = s_expectedLinkBalance + amount;
    emit FundsAdded(id, sender, uint96(amount));
  }

  // ================================================================
  // |                           SETTERS                            |
  // ================================================================

  /**
   * @inheritdoc OCR2Abstract
   * @dev prefer the type-safe version of setConfig (below) whenever possible
   */
  function setConfig(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfigBytes,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) external override {
    setConfigTypeSafe(
      signers,
      transmitters,
      f,
      abi.decode(onchainConfigBytes, (IAutomationV21PlusCommon.OnchainConfigLegacy)),
      offchainConfigVersion,
      offchainConfig
    );
  }

  function setConfigTypeSafe(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    IAutomationV21PlusCommon.OnchainConfigLegacy memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) public onlyOwner {
    if (signers.length > MAX_NUM_ORACLES) revert TooManyOracles();
    if (f == 0) revert IncorrectNumberOfFaultyOracles();
    if (signers.length != transmitters.length || signers.length <= 3 * f) revert IncorrectNumberOfSigners();

    // move all pooled payments out of the pool to each transmitter's balance
    uint96 totalPremium = s_hotVars.totalPremium;
    uint96 oldLength = uint96(s_transmittersList.length);
    for (uint256 i = 0; i < oldLength; i++) {
      _updateTransmitterBalanceFromPool(s_transmittersList[i], totalPremium, oldLength);
    }

    // remove any old signer/transmitter addresses
    address signerAddress;
    address transmitterAddress;
    for (uint256 i = 0; i < oldLength; i++) {
      signerAddress = s_signersList[i];
      transmitterAddress = s_transmittersList[i];
      delete s_signers[signerAddress];
      // Do not delete the whole transmitter struct as it has balance information stored
      s_transmitters[transmitterAddress].active = false;
    }
    delete s_signersList;
    delete s_transmittersList;

    // add new signer/transmitter addresses
    {
      Transmitter memory transmitter;
      address temp;
      for (uint256 i = 0; i < signers.length; i++) {
        if (s_signers[signers[i]].active) revert RepeatedSigner();
        if (signers[i] == ZERO_ADDRESS) revert InvalidSigner();
        s_signers[signers[i]] = Signer({active: true, index: uint8(i)});

        temp = transmitters[i];
        if (temp == ZERO_ADDRESS) revert InvalidTransmitter();
        transmitter = s_transmitters[temp];
        if (transmitter.active) revert RepeatedTransmitter();
        transmitter.active = true;
        transmitter.index = uint8(i);
        // new transmitters start afresh from current totalPremium
        // some spare change of premium from previous pool will be forfeited
        transmitter.lastCollected = totalPremium;
        s_transmitters[temp] = transmitter;
      }
    }
    s_signersList = signers;
    s_transmittersList = transmitters;

    s_hotVars = HotVars({
      f: f,
      paymentPremiumPPB: onchainConfig.paymentPremiumPPB,
      flatFeeMicroLink: onchainConfig.flatFeeMicroLink,
      stalenessSeconds: onchainConfig.stalenessSeconds,
      gasCeilingMultiplier: onchainConfig.gasCeilingMultiplier,
      paused: s_hotVars.paused,
      reentrancyGuard: s_hotVars.reentrancyGuard,
      totalPremium: totalPremium,
      latestEpoch: 0 // DON restarts epoch
    });

    s_storage = Storage({
      checkGasLimit: onchainConfig.checkGasLimit,
      minUpkeepSpend: onchainConfig.minUpkeepSpend,
      maxPerformGas: onchainConfig.maxPerformGas,
      transcoder: onchainConfig.transcoder,
      maxCheckDataSize: onchainConfig.maxCheckDataSize,
      maxPerformDataSize: onchainConfig.maxPerformDataSize,
      maxRevertDataSize: onchainConfig.maxRevertDataSize,
      upkeepPrivilegeManager: onchainConfig.upkeepPrivilegeManager,
      nonce: s_storage.nonce,
      configCount: s_storage.configCount,
      latestConfigBlockNumber: s_storage.latestConfigBlockNumber,
      ownerLinkBalance: s_storage.ownerLinkBalance
    });
    s_fallbackGasPrice = onchainConfig.fallbackGasPrice;
    s_fallbackLinkPrice = onchainConfig.fallbackLinkPrice;

    uint32 previousConfigBlockNumber = s_storage.latestConfigBlockNumber;
    s_storage.latestConfigBlockNumber = uint32(_blockNum());
    s_storage.configCount += 1;

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

    for (uint256 idx = 0; idx < s_registrars.length(); idx++) {
      s_registrars.remove(s_registrars.at(idx));
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

  // ================================================================
  // |                           GETTERS                            |
  // ================================================================

  /**
   * @inheritdoc OCR2Abstract
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
