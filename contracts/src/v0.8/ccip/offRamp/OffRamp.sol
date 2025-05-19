// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IAny2EVMMessageReceiver} from "../interfaces/IAny2EVMMessageReceiver.sol";
import {IFeeQuoter} from "../interfaces/IFeeQuoter.sol";
import {IMessageInterceptor} from "../interfaces/IMessageInterceptor.sol";
import {INonceManager} from "../interfaces/INonceManager.sol";
import {IPoolV1} from "../interfaces/IPool.sol";
import {IRMNRemote} from "../interfaces/IRMNRemote.sol";
import {IRouter} from "../interfaces/IRouter.sol";
import {ITokenAdminRegistry} from "../interfaces/ITokenAdminRegistry.sol";

import {CallWithExactGas} from "../../shared/call/CallWithExactGas.sol";
import {Client} from "../libraries/Client.sol";
import {ERC165CheckerReverting} from "../libraries/ERC165CheckerReverting.sol";
import {Internal} from "../libraries/Internal.sol";
import {MerkleMultiProof} from "../libraries/MerkleMultiProof.sol";
import {Pool} from "../libraries/Pool.sol";
import {MultiOCR3Base} from "../ocr/MultiOCR3Base.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/token/ERC20/IERC20.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v5.0.2/contracts/utils/structs/EnumerableSet.sol";

/// @notice OffRamp enables OCR networks to execute multiple messages in an OffRamp in a single transaction.
/// @dev The OnRamp and OffRamp form a cross chain upgradeable unit. Any change to one of them results an
/// onchain upgrade of both contracts.
/// @dev MultiOCR3Base is used to store multiple OCR configs for the OffRamp. The execution plugin type has to be
/// configured without signature verification, and the commit plugin type with verification.
contract OffRamp is ITypeAndVersion, MultiOCR3Base {
  using ERC165CheckerReverting for address;
  using EnumerableSet for EnumerableSet.UintSet;

  error ZeroChainSelectorNotAllowed();
  error ExecutionError(bytes32 messageId, bytes err);
  error SourceChainNotEnabled(uint64 sourceChainSelector);
  error TokenDataMismatch(uint64 sourceChainSelector, uint64 sequenceNumber);
  error UnexpectedTokenData();
  error ManualExecutionNotYetEnabled(uint64 sourceChainSelector);
  error ManualExecutionGasLimitMismatch();
  error InvalidManualExecutionGasLimit(uint64 sourceChainSelector, bytes32 messageId, uint256 newLimit);
  error InvalidManualExecutionTokenGasOverride(
    bytes32 messageId, uint256 tokenIndex, uint256 oldLimit, uint256 tokenGasOverride
  );
  error ManualExecutionGasAmountCountMismatch(bytes32 messageId, uint64 sequenceNumber);
  error RootNotCommitted(uint64 sourceChainSelector);
  error RootAlreadyCommitted(uint64 sourceChainSelector, bytes32 merkleRoot);
  error InvalidRoot();
  error CanOnlySelfCall();
  error ReceiverError(bytes err);
  error TokenHandlingError(address target, bytes err);
  error ReleaseOrMintBalanceMismatch(uint256 amountReleased, uint256 balancePre, uint256 balancePost);
  error EmptyReport(uint64 sourceChainSelector);
  error EmptyBatch();
  error CursedByRMN(uint64 sourceChainSelector);
  error NotACompatiblePool(address notPool);
  error InvalidDataLength(uint256 expected, uint256 got);
  error InvalidNewState(uint64 sourceChainSelector, uint64 sequenceNumber, Internal.MessageExecutionState newState);
  error StaleCommitReport();
  error InvalidInterval(uint64 sourceChainSelector, uint64 min, uint64 max);
  error ZeroAddressNotAllowed();
  error InvalidMessageDestChainSelector(uint64 messageDestChainSelector);
  error SourceChainSelectorMismatch(uint64 reportSourceChainSelector, uint64 messageSourceChainSelector);
  error SignatureVerificationRequiredInCommitPlugin();
  error SignatureVerificationNotAllowedInExecutionPlugin();
  error CommitOnRampMismatch(bytes reportOnRamp, bytes configOnRamp);
  error InvalidOnRampUpdate(uint64 sourceChainSelector);
  error RootBlessingMismatch(uint64 sourceChainSelector, bytes32 merkleRoot, bool isBlessed);

  /// @dev Atlas depends on various events, if changing, please notify Atlas.
  event StaticConfigSet(StaticConfig staticConfig);
  event DynamicConfigSet(DynamicConfig dynamicConfig);
  event ExecutionStateChanged(
    uint64 indexed sourceChainSelector,
    uint64 indexed sequenceNumber,
    bytes32 indexed messageId,
    bytes32 messageHash,
    Internal.MessageExecutionState state,
    bytes returnData,
    uint256 gasUsed
  );
  event SourceChainSelectorAdded(uint64 sourceChainSelector);
  event SourceChainConfigSet(uint64 indexed sourceChainSelector, SourceChainConfig sourceConfig);
  event SkippedAlreadyExecutedMessage(uint64 sourceChainSelector, uint64 sequenceNumber);
  event AlreadyAttempted(uint64 sourceChainSelector, uint64 sequenceNumber);
  event CommitReportAccepted(
    Internal.MerkleRoot[] blessedMerkleRoots,
    Internal.MerkleRoot[] unblessedMerkleRoots,
    Internal.PriceUpdates priceUpdates
  );
  event RootRemoved(bytes32 root);
  event SkippedReportExecution(uint64 sourceChainSelector);

  /// @dev Struct that contains the static configuration. The individual components are stored as immutable variables.
  // solhint-disable-next-line gas-struct-packing
  struct StaticConfig {
    uint64 chainSelector; // ───────╮ Destination chainSelector
    uint16 gasForCallExactCheck; // | Gas for call exact check
    IRMNRemote rmnRemote; // ───────╯ RMN Verification Contract
    address tokenAdminRegistry; // Token admin registry address
    address nonceManager; // Nonce manager address
  }

  /// @dev Per-chain source config (defining a lane from a Source Chain -> Dest OffRamp).
  struct SourceChainConfig {
    IRouter router; // ─────────────────╮ Local router to use for messages coming from this source chain.
    bool isEnabled; //                  │ Flag whether the source chain is enabled or not.
    uint64 minSeqNr; //                 │ The min sequence number expected for future messages.
    bool isRMNVerificationDisabled; // ─╯ Flag whether the RMN verification is disabled or not.
    bytes onRamp; // OnRamp address on the source chain.
  }

  /// @dev Same as SourceChainConfig but with source chain selector so that an array of these
  /// can be passed in the constructor and the applySourceChainConfigUpdates function.
  struct SourceChainConfigArgs {
    IRouter router; // ─────────────────╮  Local router to use for messages coming from this source chain.
    uint64 sourceChainSelector; //      │  Source chain selector of the config to update.
    bool isEnabled; //                  │  Flag whether the source chain is enabled or not.
    bool isRMNVerificationDisabled; // ─╯ Flag whether the RMN verification is disabled or not.
    bytes onRamp; // OnRamp address on the source chain.
  }

  /// @dev Dynamic offRamp config.
  /// @dev Since DynamicConfig is part of DynamicConfigSet event, if changing it, we should update the ABI on Atlas.
  struct DynamicConfig {
    address feeQuoter; // ──────────────────────────────╮ FeeQuoter address on the local chain.
    uint32 permissionLessExecutionThresholdSeconds; // ─╯ Waiting time before manual execution is enabled.
    address messageInterceptor; // Optional, validates incoming messages (zero address = no interceptor).
  }

  /// @dev Report that is committed by the observing DON at the committing phase.
  struct CommitReport {
    Internal.PriceUpdates priceUpdates; // List of gas and price updates to commit.
    Internal.MerkleRoot[] blessedMerkleRoots; // List of merkle roots from source chains for which RMN is enabled.
    Internal.MerkleRoot[] unblessedMerkleRoots; // List of merkle roots from source chains for which RMN is disabled.
    IRMNRemote.Signature[] rmnSignatures; // RMN signatures on the merkle roots.
  }

  /// @dev Both receiverExecutionGasLimit and tokenGasOverrides are optional. To indicate no override, set the value
  /// to 0. The length of tokenGasOverrides must match the length of tokenAmounts, even if it only contains zeros.
  struct GasLimitOverride {
    uint256 receiverExecutionGasLimit; // Overrides EVM2EVMMessage.gasLimit.
    uint32[] tokenGasOverrides; // Overrides EVM2EVMMessage.sourceTokenData.destGasAmount, length must be same as tokenAmounts.
  }

  // STATIC CONFIG
  string public constant override typeAndVersion = "OffRamp 1.6.0";
  /// @dev Hash of encoded address(0) used for empty address checks.
  bytes32 internal constant EMPTY_ENCODED_ADDRESS_HASH = keccak256(abi.encode(address(0)));
  /// @dev ChainSelector of this chain.
  uint64 internal immutable i_chainSelector;
  /// @dev The RMN verification contract.
  IRMNRemote internal immutable i_rmnRemote;
  /// @dev The address of the token admin registry.
  address internal immutable i_tokenAdminRegistry;
  /// @dev The address of the nonce manager.
  address internal immutable i_nonceManager;
  /// @dev The minimum amount of gas to perform the call with exact gas.
  /// We include this in the offRamp so that we can redeploy to adjust it should a hardfork change the gas costs of
  /// relevant opcodes in callWithExactGas.
  uint16 internal immutable i_gasForCallExactCheck;

  // DYNAMIC CONFIG
  DynamicConfig internal s_dynamicConfig;

  /// @notice Set of source chain selectors.
  EnumerableSet.UintSet internal s_sourceChainSelectors;

  /// @notice SourceChainConfig per source chain selector.
  mapping(uint64 sourceChainSelector => SourceChainConfig sourceChainConfig) private s_sourceChainConfigs;

  // STATE
  /// @dev A mapping of sequence numbers (per source chain) to execution state using a bitmap with each execution
  /// state only taking up 2 bits of the uint256, packing 128 states into a single slot.
  /// Message state is tracked to ensure message can only be executed successfully once.
  mapping(uint64 sourceChainSelector => mapping(uint64 seqNum => uint256 executionStateBitmap)) internal
    s_executionStates;

  /// @notice Commit timestamp of merkle roots per source chain.
  mapping(uint64 sourceChainSelector => mapping(bytes32 merkleRoot => uint256 timestamp)) internal s_roots;
  /// @dev The sequence number of the last price update.
  uint64 private s_latestPriceSequenceNumber;

  constructor(
    StaticConfig memory staticConfig,
    DynamicConfig memory dynamicConfig,
    SourceChainConfigArgs[] memory sourceChainConfigs
  ) MultiOCR3Base() {
    if (
      address(staticConfig.rmnRemote) == address(0) || staticConfig.tokenAdminRegistry == address(0)
        || staticConfig.nonceManager == address(0)
    ) {
      revert ZeroAddressNotAllowed();
    }

    if (staticConfig.chainSelector == 0) {
      revert ZeroChainSelectorNotAllowed();
    }

    i_chainSelector = staticConfig.chainSelector;
    i_rmnRemote = staticConfig.rmnRemote;
    i_tokenAdminRegistry = staticConfig.tokenAdminRegistry;
    i_nonceManager = staticConfig.nonceManager;
    i_gasForCallExactCheck = staticConfig.gasForCallExactCheck;
    emit StaticConfigSet(staticConfig);

    _setDynamicConfig(dynamicConfig);
    _applySourceChainConfigUpdates(sourceChainConfigs);
  }

  // ================================================================
  // │                           Execution                          │
  // ================================================================

  // The size of the execution state in bits.
  uint256 private constant MESSAGE_EXECUTION_STATE_BIT_WIDTH = 2;
  // The mask for the execution state bits.
  uint256 private constant MESSAGE_EXECUTION_STATE_MASK = (1 << MESSAGE_EXECUTION_STATE_BIT_WIDTH) - 1;

  /// @notice Returns the current execution state of a message based on its sequenceNumber.
  /// @param sourceChainSelector The source chain to get the execution state for.
  /// @param sequenceNumber The sequence number of the message to get the execution state for.
  /// @return executionState The current execution state of the message.
  /// @dev We use the literal number 128 because using a constant increased gas usage.
  function getExecutionState(
    uint64 sourceChainSelector,
    uint64 sequenceNumber
  ) public view returns (Internal.MessageExecutionState) {
    return Internal.MessageExecutionState(
      (
        _getSequenceNumberBitmap(sourceChainSelector, sequenceNumber)
          >> ((sequenceNumber % 128) * MESSAGE_EXECUTION_STATE_BIT_WIDTH)
      ) & MESSAGE_EXECUTION_STATE_MASK
    );
  }

  /// @notice Sets a new execution state for a given sequence number. It will overwrite any existing state.
  /// @param sourceChainSelector The source chain to set the execution state for.
  /// @param sequenceNumber The sequence number for which the state will be saved.
  /// @param newState The new value the state will be in after this function is called.
  /// @dev We use the literal number 128 because using a constant increased gas usage.
  function _setExecutionState(
    uint64 sourceChainSelector,
    uint64 sequenceNumber,
    Internal.MessageExecutionState newState
  ) internal {
    uint256 offset = (sequenceNumber % 128) * MESSAGE_EXECUTION_STATE_BIT_WIDTH;
    uint256 bitmap = _getSequenceNumberBitmap(sourceChainSelector, sequenceNumber);
    // To unset any potential existing state we zero the bits of the section the state occupies,
    // then we do an AND operation to blank out any existing state for the section.
    bitmap &= ~(MESSAGE_EXECUTION_STATE_MASK << offset);
    // Set the new state.
    bitmap |= uint256(newState) << offset;

    s_executionStates[sourceChainSelector][sequenceNumber / 128] = bitmap;
  }

  /// @param sourceChainSelector remote source chain selector to get sequence number bitmap for.
  /// @param sequenceNumber sequence number to get bitmap for.
  /// @return bitmap Bitmap of the given sequence number for the provided source chain selector. One bitmap represents
  /// 128 sequence numbers.
  function _getSequenceNumberBitmap(
    uint64 sourceChainSelector,
    uint64 sequenceNumber
  ) internal view returns (uint256 bitmap) {
    return s_executionStates[sourceChainSelector][sequenceNumber / 128];
  }

  /// @notice Manually executes a set of reports.
  /// @param reports Internal.ExecutionReportSingleChain[] - list of reports to execute.
  /// @param gasLimitOverrides New gasLimit for each message per report. The outer array represents each report, the
  //  inner array represents each message in the report.
  //  i.e. gasLimitOverrides[report1][report1Message1] -> access message1 from report1
  /// @dev We permit gas limit overrides so that users may manually execute messages which failed due to insufficient
  /// gas provided. The reports do not have to contain all the messages (they can be omitted). Multiple reports can be
  /// passed in simultaneously.
  function manuallyExecute(
    Internal.ExecutionReport[] memory reports,
    GasLimitOverride[][] memory gasLimitOverrides
  ) external {
    // We do this here because the other _execute path is already covered by MultiOCR3Base.
    _whenChainNotForked();

    uint256 numReports = reports.length;
    if (numReports != gasLimitOverrides.length) revert ManualExecutionGasLimitMismatch();

    for (uint256 reportIndex = 0; reportIndex < numReports; ++reportIndex) {
      Internal.ExecutionReport memory report = reports[reportIndex];

      uint256 numMsgs = report.messages.length;
      GasLimitOverride[] memory msgGasLimitOverrides = gasLimitOverrides[reportIndex];

      // Gas override values need to be provided, even when no override is desired. We expect an array of the correct
      // size with all `0` values if no override is desired.
      if (numMsgs != msgGasLimitOverrides.length) revert ManualExecutionGasLimitMismatch();

      for (uint256 msgIndex = 0; msgIndex < numMsgs; ++msgIndex) {
        uint256 newLimit = msgGasLimitOverrides[msgIndex].receiverExecutionGasLimit;
        Internal.Any2EVMRampMessage memory message = report.messages[msgIndex];
        if (newLimit != 0) {
          // Checks to ensure messages will not be executed with less gas than specified.
          if (newLimit < message.gasLimit) {
            revert InvalidManualExecutionGasLimit(report.sourceChainSelector, message.header.messageId, newLimit);
          }
        }
        if (message.tokenAmounts.length != msgGasLimitOverrides[msgIndex].tokenGasOverrides.length) {
          revert ManualExecutionGasAmountCountMismatch(message.header.messageId, message.header.sequenceNumber);
        }

        // The gas limit can not be lowered as that could cause the message to fail. If manual execution is done
        // from an UNTOUCHED state and we would allow lower gas limit, anyone could grief by executing the message with
        // lower gas limit than the DON would have used. This results in the message being marked FAILURE and the DON
        // would not attempt it with the correct gas limit.
        for (uint256 tokenIndex = 0; tokenIndex < message.tokenAmounts.length; ++tokenIndex) {
          uint256 tokenGasOverride = msgGasLimitOverrides[msgIndex].tokenGasOverrides[tokenIndex];
          if (tokenGasOverride != 0) {
            uint256 destGasAmount = message.tokenAmounts[tokenIndex].destGasAmount;
            if (tokenGasOverride < destGasAmount) {
              revert InvalidManualExecutionTokenGasOverride(
                message.header.messageId, tokenIndex, destGasAmount, tokenGasOverride
              );
            }
          }
        }
      }
    }

    _batchExecute(reports, gasLimitOverrides);
  }

  /// @notice Transmit function for execution reports. The function takes no signatures, and expects the exec plugin
  /// type to be configured with no signatures.
  /// @param report serialized execution report.
  function execute(bytes32[2] calldata reportContext, bytes calldata report) external {
    _batchExecute(abi.decode(report, (Internal.ExecutionReport[])), new GasLimitOverride[][](0));

    bytes32[] memory emptySigs = new bytes32[](0);
    _transmit(uint8(Internal.OCRPluginType.Execution), reportContext, report, emptySigs, emptySigs, bytes32(""));
  }

  /// @notice Batch executes a set of reports, each report matching one single source chain.
  /// @param reports Set of execution reports (one per chain) containing the messages and proofs.
  /// @param manualExecGasLimits An array of gas limits to use for manual execution The outer array represents each
  //  report, the inner array represents each message in the report.
  //  i.e. gasLimitOverrides[report1][report1Message1] -> access message1 from report1.
  /// @dev The manualExecGasLimits array should either be empty, or match the length of the reports array.
  /// @dev If called from manual execution, each inner array's length has to match the number of messages.
  function _batchExecute(
    Internal.ExecutionReport[] memory reports,
    GasLimitOverride[][] memory manualExecGasOverrides
  ) internal {
    if (reports.length == 0) revert EmptyBatch();

    bool areManualGasLimitsEmpty = manualExecGasOverrides.length == 0;
    // Cache array for gas savings in the loop's condition.
    GasLimitOverride[] memory emptyGasLimits = new GasLimitOverride[](0);

    for (uint256 i = 0; i < reports.length; ++i) {
      _executeSingleReport(reports[i], areManualGasLimitsEmpty ? emptyGasLimits : manualExecGasOverrides[i]);
    }
  }

  /// @notice Executes a report, executing each message in order.
  /// @param report The execution report containing the messages and proofs.
  /// @param manualExecGasExecOverrides An array of gas limits to use for manual execution.
  /// @dev If called from the DON, this array is always empty.
  /// @dev If called from manual execution, this array is always same length as messages.
  /// @dev This function can fully revert in some cases, reverting potentially valid other reports with it. The reasons
  /// for these reverts are so severe that we prefer to revert the entire batch instead of silently failing.
  function _executeSingleReport(
    Internal.ExecutionReport memory report,
    GasLimitOverride[] memory manualExecGasExecOverrides
  ) internal {
    uint64 sourceChainSelector = report.sourceChainSelector;
    bool manualExecution = manualExecGasExecOverrides.length != 0;
    if (i_rmnRemote.isCursed(bytes16(uint128(sourceChainSelector)))) {
      if (manualExecution) {
        // For manual execution we don't want to silently fail so we revert.
        revert CursedByRMN(sourceChainSelector);
      }
      // For DON execution we do not revert as a single lane curse can revert the entire batch.
      emit SkippedReportExecution(sourceChainSelector);
      return;
    }

    uint256 numMsgs = report.messages.length;
    if (numMsgs == 0) revert EmptyReport(report.sourceChainSelector);
    if (numMsgs != report.offchainTokenData.length) revert UnexpectedTokenData();

    bytes32[] memory hashedLeaves = new bytes32[](numMsgs);

    {
      // We do this hash here instead of in _verify to avoid two separate loops over the same data. Hashing all of the
      // message fields ensures that the message being executed is correct and not tampered with. Including the known
      // OnRamp ensures that the message originates from the correct on ramp version. We know the sourceChainSelector
      // and i_destChainSelector are correct because we revert below when they are not.
      bytes32 metaDataHash = keccak256(
        abi.encode(
          Internal.ANY_2_EVM_MESSAGE_HASH,
          sourceChainSelector,
          i_chainSelector,
          keccak256(_getEnabledSourceChainConfig(sourceChainSelector).onRamp)
        )
      );

      for (uint256 i = 0; i < numMsgs; ++i) {
        Internal.Any2EVMRampMessage memory message = report.messages[i];

        // Commits do not verify the destChainSelector in the message since only the root is committed, so we
        // have to check it explicitly. This check is also important as we have assumed the metaDataHash above uses
        // the i_chainSelector as the destChainSelector.
        if (message.header.destChainSelector != i_chainSelector) {
          revert InvalidMessageDestChainSelector(message.header.destChainSelector);
        }
        // If the message source chain selector does not match the report's source chain selector and the root has not
        // been committed for the report source chain selector this will be caught by the root verification.
        // This acts as an extra check to ensure the message source chain selector matches the report's source chain.
        if (message.header.sourceChainSelector != sourceChainSelector) {
          revert SourceChainSelectorMismatch(sourceChainSelector, message.header.sourceChainSelector);
        }

        hashedLeaves[i] = Internal._hash(message, metaDataHash);
      }
    }

    // SECURITY CRITICAL CHECK.
    uint256 timestampCommitted = _verify(sourceChainSelector, hashedLeaves, report.proofs, report.proofFlagBits);
    if (timestampCommitted == 0) revert RootNotCommitted(sourceChainSelector);

    // Execute messages.
    for (uint256 i = 0; i < numMsgs; ++i) {
      uint256 gasStart = gasleft();
      Internal.Any2EVMRampMessage memory message = report.messages[i];
      message = _beforeExecuteSingleMessage(message);

      Internal.MessageExecutionState originalState =
        getExecutionState(sourceChainSelector, message.header.sequenceNumber);
      // Two valid cases here, we either have never touched this message before, or we tried to execute and failed. This
      // check protects against reentry and re-execution because the other state is IN_PROGRESS which should not be
      // allowed to execute.
      if (
        !(
          originalState == Internal.MessageExecutionState.UNTOUCHED
            || originalState == Internal.MessageExecutionState.FAILURE
        )
      ) {
        // If the message has already been executed, we skip it. We want to not revert on race conditions between
        // executing parties. This will allow us to open up manual exec while also attempting with the DON, without
        // reverting an entire DON batch when a user manually executes while the tx is inflight.
        emit SkippedAlreadyExecutedMessage(sourceChainSelector, message.header.sequenceNumber);
        continue;
      }
      uint32[] memory tokenGasOverrides;
      if (manualExecution) {
        tokenGasOverrides = manualExecGasExecOverrides[i].tokenGasOverrides;
        bool isOldCommitReport =
          (block.timestamp - timestampCommitted) > s_dynamicConfig.permissionLessExecutionThresholdSeconds;
        // Manually execution is fine if we previously failed or if the commit report is just too old.
        // Acceptable state transitions: UNTOUCHED->SUCCESS, UNTOUCHED->FAILURE, FAILURE->SUCCESS.
        if (!(isOldCommitReport || originalState == Internal.MessageExecutionState.FAILURE)) {
          revert ManualExecutionNotYetEnabled(sourceChainSelector);
        }

        // Manual execution gas limit can override gas limit specified in the message. Value of 0 indicates no override.
        if (manualExecGasExecOverrides[i].receiverExecutionGasLimit != 0) {
          message.gasLimit = manualExecGasExecOverrides[i].receiverExecutionGasLimit;
        }
      } else {
        // DON can only execute a message once.
        // Acceptable state transitions: UNTOUCHED->SUCCESS, UNTOUCHED->FAILURE.
        if (originalState != Internal.MessageExecutionState.UNTOUCHED) {
          emit AlreadyAttempted(sourceChainSelector, message.header.sequenceNumber);
          continue;
        }
      }

      // Nonce changes per state transition (these only apply for ordered messages):
      // UNTOUCHED -> FAILURE  nonce bump.
      // UNTOUCHED -> SUCCESS  nonce bump.
      // FAILURE   -> SUCCESS  no nonce bump.
      // UNTOUCHED messages MUST be executed in order always.
      // If nonce == 0 then out of order execution is allowed.
      if (message.header.nonce != 0) {
        if (originalState == Internal.MessageExecutionState.UNTOUCHED) {
          // If a nonce is not incremented, that means it was skipped, and we can ignore the message.
          if (
            !INonceManager(i_nonceManager).incrementInboundNonce(
              sourceChainSelector, message.header.nonce, message.sender
            )
          ) continue;
        }
      }

      // We expect only valid messages will be committed but we check when executing as a defense in depth measure.
      bytes[] memory offchainTokenData = report.offchainTokenData[i];
      if (message.tokenAmounts.length != offchainTokenData.length) {
        revert TokenDataMismatch(sourceChainSelector, message.header.sequenceNumber);
      }

      _setExecutionState(sourceChainSelector, message.header.sequenceNumber, Internal.MessageExecutionState.IN_PROGRESS);
      (Internal.MessageExecutionState newState, bytes memory returnData) =
        _trialExecute(message, offchainTokenData, tokenGasOverrides);
      _setExecutionState(sourceChainSelector, message.header.sequenceNumber, newState);

      // Since it's hard to estimate whether manual execution will succeed, we revert the entire transaction if it
      // fails. This will show the user if their manual exec will fail before they submit it.
      if (manualExecution) {
        if (newState == Internal.MessageExecutionState.FAILURE) {
          if (originalState != Internal.MessageExecutionState.UNTOUCHED) {
            // If manual execution fails, we revert the entire transaction, unless the originalState is UNTOUCHED as we
            // would still be making progress by changing the state from UNTOUCHED to FAILURE.
            revert ExecutionError(message.header.messageId, returnData);
          }
        }
      }

      // The only valid prior states are UNTOUCHED and FAILURE (checked above).
      // The only valid post states are FAILURE and SUCCESS (checked below).
      if (newState != Internal.MessageExecutionState.SUCCESS) {
        if (newState != Internal.MessageExecutionState.FAILURE) {
          revert InvalidNewState(sourceChainSelector, message.header.sequenceNumber, newState);
        }
      }

      emit ExecutionStateChanged(
        sourceChainSelector,
        message.header.sequenceNumber,
        message.header.messageId,
        hashedLeaves[i],
        newState,
        returnData,
        // This emit covers not only the execution through the router, but also all of the overhead in executing the
        // message. This gives the most accurate representation of the gas used in the execution.
        gasStart - gasleft()
      );
    }
  }

  /// @notice Try executing a message.
  /// @param message Internal.Any2EVMRampMessage memory message.
  /// @param offchainTokenData Data provided by the DON for token transfers.
  /// @return executionState The new state of the message, being either SUCCESS or FAILURE.
  /// @return errData Revert data in bytes if CCIP receiver reverted during execution.
  function _trialExecute(
    Internal.Any2EVMRampMessage memory message,
    bytes[] memory offchainTokenData,
    uint32[] memory tokenGasOverrides
  ) internal returns (Internal.MessageExecutionState executionState, bytes memory) {
    try this.executeSingleMessage(message, offchainTokenData, tokenGasOverrides) {}
    catch (bytes memory err) {
      if (msg.sender == Internal.GAS_ESTIMATION_SENDER) {
        if (
          CallWithExactGas.NOT_ENOUGH_GAS_FOR_CALL_SIG == bytes4(err)
            || CallWithExactGas.NO_GAS_FOR_CALL_EXACT_CHECK_SIG == bytes4(err)
            || ERC165CheckerReverting.InsufficientGasForStaticCall.selector == bytes4(err)
        ) {
          revert InsufficientGasToCompleteTx(bytes4(err));
        }
      }
      // return the message execution state as FAILURE and the revert data.
      // Max length of revert data is Router.MAX_RET_BYTES, max length of err is 4 + Router.MAX_RET_BYTES.
      return (Internal.MessageExecutionState.FAILURE, err);
    }
    // If message execution succeeded, no CCIP receiver return data is expected, return with empty bytes.
    return (Internal.MessageExecutionState.SUCCESS, "");
  }

  /// @notice hook for applying custom logic to the input message before executeSingleMessage()
  /// @param message initial message
  /// @return transformedMessage modified message
  function _beforeExecuteSingleMessage(
    Internal.Any2EVMRampMessage memory message
  ) internal virtual returns (Internal.Any2EVMRampMessage memory transformedMessage) {
    return message;
  }

  /// @notice Executes a single message.
  /// @param message The message that will be executed.
  /// @param offchainTokenData Token transfer data to be passed to TokenPool.
  /// @dev We make this external and callable by the contract itself, in order to try/catch
  /// its execution and enforce atomicity among successful message processing and token transfer.
  /// @dev We use ERC-165 to check for the ccipReceive interface to permit sending tokens to contracts, for example
  /// smart contract wallets, without an associated message.
  function executeSingleMessage(
    Internal.Any2EVMRampMessage memory message,
    bytes[] calldata offchainTokenData,
    uint32[] calldata tokenGasOverrides
  ) external {
    if (msg.sender != address(this)) revert CanOnlySelfCall();

    Client.EVMTokenAmount[] memory destTokenAmounts = new Client.EVMTokenAmount[](0);
    if (message.tokenAmounts.length > 0) {
      destTokenAmounts = _releaseOrMintTokens(
        message.tokenAmounts,
        message.sender,
        message.receiver,
        message.header.sourceChainSelector,
        offchainTokenData,
        tokenGasOverrides
      );
    }

    Client.Any2EVMMessage memory any2EvmMessage = Client.Any2EVMMessage({
      messageId: message.header.messageId,
      sourceChainSelector: message.header.sourceChainSelector,
      sender: message.sender,
      data: message.data,
      destTokenAmounts: destTokenAmounts
    });

    // The main message interceptor is the aggregate rate limiter, but we also allow for a custom interceptor. This is
    // why we always have to call into the contract when it's enabled, even when there are no tokens in the message.
    address messageInterceptor = s_dynamicConfig.messageInterceptor;
    if (messageInterceptor != address(0)) {
      try IMessageInterceptor(messageInterceptor).onInboundMessage(any2EvmMessage) {}
      catch (bytes memory err) {
        revert IMessageInterceptor.MessageValidationError(err);
      }
    }

    // There are three cases in which we skip calling the receiver:
    // 1. If the message data is empty AND the gas limit is 0.
    //          This indicates a message that only transfers tokens. It is valid to only send tokens to a contract
    //          that supports the IAny2EVMMessageReceiver interface, but without this first check we would call the
    //          receiver without any gas, which would revert the transaction.
    // 2. If the receiver is not a contract.
    // 3. If the receiver is a contract but it does not support the IAny2EVMMessageReceiver interface.
    //
    // The ordering of these checks is important, as the first check is the cheapest to execute.
    //
    // To prevent message delivery bypass issues, a modified version of the ERC165Checker is used
    // which checks for sufficient gas before making the external call.
    if (
      (message.data.length == 0 && message.gasLimit == 0) || message.receiver.code.length == 0
        || !message.receiver._supportsInterfaceReverting(type(IAny2EVMMessageReceiver).interfaceId)
    ) return;

    (bool success, bytes memory returnData,) = s_sourceChainConfigs[message.header.sourceChainSelector]
      .router
      .routeMessage(any2EvmMessage, i_gasForCallExactCheck, message.gasLimit, message.receiver);
    // If CCIP receiver execution is not successful, revert the call including token transfers.
    if (!success) revert ReceiverError(returnData);
  }

  // ================================================================
  // │                      Tokens and pools                        │
  // ================================================================

  /// @notice Uses a pool to release or mint a token to a receiver address, with balance checks before and after the
  /// transfer. This is done to ensure the exact number of tokens the pool claims to release are actually transferred.
  /// @dev The local token address is validated through the TokenAdminRegistry. If, due to some misconfiguration, the
  /// token is unknown to the registry, the offRamp will revert. The tx, and the tokens, can be retrieved by registering
  /// the token on this chain, and re-trying the msg.
  /// @param sourceTokenAmount Amount and source data of the token to be released/minted.
  /// @param originalSender The message sender on the source chain.
  /// @param receiver The address that will receive the tokens.
  /// @param sourceChainSelector The remote source chain selector
  /// @param offchainTokenData Data fetched offchain by the DON.
  /// @return destTokenAmount local token address with amount.
  function _releaseOrMintSingleToken(
    Internal.Any2EVMTokenTransfer memory sourceTokenAmount,
    bytes memory originalSender,
    address receiver,
    uint64 sourceChainSelector,
    bytes memory offchainTokenData
  ) internal returns (Client.EVMTokenAmount memory destTokenAmount) {
    // We need to safely decode the token address from the sourceTokenData as it could be wrong, in which case it
    // doesn't have to be a valid EVM address.
    address localToken = sourceTokenAmount.destTokenAddress;
    // We check with the token admin registry if the token has a pool on this chain.
    address localPoolAddress = ITokenAdminRegistry(i_tokenAdminRegistry).getPool(localToken);
    // This will call the supportsInterface through the ERC165Checker, and not directly on the pool address.
    // This is done to prevent a pool from reverting the entire transaction if it doesn't support the interface.
    // The call gets a max or 30k gas per instance, of which there are three. This means offchain gas estimations should
    // account for 90k gas overhead due to the interface check.
    if (localPoolAddress == address(0) || !localPoolAddress._supportsInterfaceReverting(Pool.CCIP_POOL_V1)) {
      revert NotACompatiblePool(localPoolAddress);
    }

    // We retrieve the local token balance of the receiver before the pool call.
    (uint256 balancePre, uint256 gasLeft) = _getBalanceOfReceiver(receiver, localToken, sourceTokenAmount.destGasAmount);

    // We determined that the pool address is a valid EVM address, but that does not mean the code at this address is a
    // (compatible) pool contract. _callWithExactGasSafeReturnData will check if the location contains a contract. If it
    // doesn't it reverts with a known error. We call the pool with exact gas  to increase resistance against malicious
    // tokens or token pools. We protect against return data bombs by capping the return data size at MAX_RET_BYTES.
    (bool success, bytes memory returnData, uint256 gasUsedReleaseOrMint) = CallWithExactGas
      ._callWithExactGasSafeReturnData(
      abi.encodeCall(
        IPoolV1.releaseOrMint,
        Pool.ReleaseOrMintInV1({
          originalSender: originalSender,
          receiver: receiver,
          amount: sourceTokenAmount.amount,
          localToken: localToken,
          remoteChainSelector: sourceChainSelector,
          sourcePoolAddress: sourceTokenAmount.sourcePoolAddress,
          sourcePoolData: sourceTokenAmount.extraData,
          offchainTokenData: offchainTokenData
        })
      ),
      localPoolAddress,
      gasLeft,
      i_gasForCallExactCheck,
      Internal.MAX_RET_BYTES
    );

    // Wrap and rethrow the error so we can catch it lower in the stack.
    if (!success) revert TokenHandlingError(localPoolAddress, returnData);

    // If the call was successful, the returnData should be the amount released or minted denominated in the local
    // token's decimals.
    if (returnData.length != Pool.CCIP_POOL_V1_RET_BYTES) {
      revert InvalidDataLength(Pool.CCIP_POOL_V1_RET_BYTES, returnData.length);
    }
    uint256 localAmount = abi.decode(returnData, (uint256));

    // We don't need to do balance checks if the pool is the receiver, as they would always fail in the case
    // of a lockRelease pool.
    if (receiver != localPoolAddress) {
      (uint256 balancePost,) = _getBalanceOfReceiver(receiver, localToken, gasLeft - gasUsedReleaseOrMint);

      // First we check if the subtraction would result in an underflow to ensure we revert with a clear error.
      if (balancePost < balancePre || balancePost - balancePre != localAmount) {
        revert ReleaseOrMintBalanceMismatch(localAmount, balancePre, balancePost);
      }
    }

    return Client.EVMTokenAmount({token: localToken, amount: localAmount});
  }

  /// @notice Retrieves the balance of a receiver address for a given token.
  /// @param receiver The address to check the balance of.
  /// @param token The token address.
  /// @param gasLimit The gas limit to use for the call.
  /// @return balance The balance of the receiver.
  /// @return gasLeft The gas left after the call.
  function _getBalanceOfReceiver(
    address receiver,
    address token,
    uint256 gasLimit
  ) internal returns (uint256 balance, uint256 gasLeft) {
    (bool success, bytes memory returnData, uint256 gasUsed) = CallWithExactGas._callWithExactGasSafeReturnData(
      abi.encodeCall(IERC20.balanceOf, (receiver)), token, gasLimit, i_gasForCallExactCheck, Internal.MAX_RET_BYTES
    );
    if (!success) revert TokenHandlingError(token, returnData);

    // If the call was successful, the returnData should contain only the balance.
    if (returnData.length != Internal.MAX_BALANCE_OF_RET_BYTES) {
      revert InvalidDataLength(Internal.MAX_BALANCE_OF_RET_BYTES, returnData.length);
    }

    // Return the decoded balance, which cannot fail as we checked the length, and the gas that is left
    // after this call.
    return (abi.decode(returnData, (uint256)), gasLimit - gasUsed);
  }

  /// @notice Uses pools to release or mint a number of different tokens to a receiver address.
  /// @param sourceTokenAmounts List of token amounts with source data of the tokens to be released/minted.
  /// @param originalSender The message sender on the source chain.
  /// @param receiver The address that will receive the tokens.
  /// @param sourceChainSelector The remote source chain selector.
  /// @param offchainTokenData Array of token data fetched offchain by the DON.
  /// @param tokenGasOverrides Array of override gas limits to use for token transfers. If empty, the normal gas limit
  /// as defined on the source chain is used.
  /// @return destTokenAmounts local token addresses with amounts.
  function _releaseOrMintTokens(
    Internal.Any2EVMTokenTransfer[] memory sourceTokenAmounts,
    bytes memory originalSender,
    address receiver,
    uint64 sourceChainSelector,
    bytes[] calldata offchainTokenData,
    uint32[] calldata tokenGasOverrides
  ) internal returns (Client.EVMTokenAmount[] memory destTokenAmounts) {
    destTokenAmounts = new Client.EVMTokenAmount[](sourceTokenAmounts.length);
    bool isTokenGasOverridesEmpty = tokenGasOverrides.length == 0;
    for (uint256 i = 0; i < sourceTokenAmounts.length; ++i) {
      if (!isTokenGasOverridesEmpty) {
        if (tokenGasOverrides[i] != 0) {
          sourceTokenAmounts[i].destGasAmount = tokenGasOverrides[i];
        }
      }
      destTokenAmounts[i] = _releaseOrMintSingleToken(
        sourceTokenAmounts[i], originalSender, receiver, sourceChainSelector, offchainTokenData[i]
      );
    }

    return destTokenAmounts;
  }

  // ================================================================
  // │                           Commit                             │
  // ================================================================

  /// @notice Transmit function for commit reports. The function requires signatures,
  /// and expects the commit plugin type to be configured with signatures.
  /// @param report serialized commit report.
  /// @dev A commitReport can have two distinct parts (batched together to amortize the cost of checking sigs):
  /// 1. Price updates
  /// 2. A batch of merkle root and sequence number intervals (per-source)
  /// Both have their own, separate, staleness checks, with price updates using the epoch and round number of the latest
  /// price update. The merkle root checks for staleness are based on the seqNums.  They need to be separate because
  /// a price report for round t+2 might be included before a report containing a merkle root for round t+1. This merkle
  /// root report for round t+1 is still valid and should not be rejected. When a report with a stale root but valid
  /// price updates is submitted, we are OK to revert to preserve the invariant that we always revert on invalid
  /// sequence number ranges. If that happens, prices will be updated in later rounds.
  function commit(
    bytes32[2] calldata reportContext,
    bytes calldata report,
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs
  ) external {
    CommitReport memory commitReport = abi.decode(report, (CommitReport));
    DynamicConfig storage dynamicConfig = s_dynamicConfig;

    // Verify RMN signatures
    if (commitReport.blessedMerkleRoots.length > 0) {
      i_rmnRemote.verify(address(this), commitReport.blessedMerkleRoots, commitReport.rmnSignatures);
    }

    // Check if the report contains price updates.
    if (commitReport.priceUpdates.tokenPriceUpdates.length > 0 || commitReport.priceUpdates.gasPriceUpdates.length > 0)
    {
      uint64 ocrSequenceNumber = uint64(uint256(reportContext[1]));

      // Check for price staleness based on the epoch and round.
      if (s_latestPriceSequenceNumber < ocrSequenceNumber) {
        // If prices are not stale, update the latest epoch and round.
        s_latestPriceSequenceNumber = ocrSequenceNumber;
        // And update the prices in the fee quoter.
        IFeeQuoter(dynamicConfig.feeQuoter).updatePrices(commitReport.priceUpdates);
      } else {
        // If prices are stale and the report doesn't contain a root, this report does not have any valid information
        // and we revert. If it does contain a merkle root, continue to the root checking section.
        if (commitReport.blessedMerkleRoots.length + commitReport.unblessedMerkleRoots.length == 0) {
          revert StaleCommitReport();
        }
      }
    }

    for (uint256 i = 0; i < commitReport.blessedMerkleRoots.length; ++i) {
      _commitRoot(commitReport.blessedMerkleRoots[i], true);
    }

    for (uint256 i = 0; i < commitReport.unblessedMerkleRoots.length; ++i) {
      _commitRoot(commitReport.unblessedMerkleRoots[i], false);
    }

    emit CommitReportAccepted(
      commitReport.blessedMerkleRoots, commitReport.unblessedMerkleRoots, commitReport.priceUpdates
    );

    _transmit(uint8(Internal.OCRPluginType.Commit), reportContext, report, rs, ss, rawVs);
  }

  /// @notice Commits a single merkle root. The blessing status has to match the source chain config.
  /// @dev An unblessed root means that RMN verification is disabled for the source chain. It does not mean there is
  /// some future point where the root will be blessed.
  /// @param root The merkle root to commit.
  /// @param isBlessed The blessing status of the root.
  function _commitRoot(Internal.MerkleRoot memory root, bool isBlessed) internal {
    uint64 sourceChainSelector = root.sourceChainSelector;

    if (i_rmnRemote.isCursed(bytes16(uint128(sourceChainSelector)))) {
      revert CursedByRMN(sourceChainSelector);
    }

    SourceChainConfig storage sourceChainConfig = _getEnabledSourceChainConfig(sourceChainSelector);

    // If the root is blessed but RMN blessing is disabled for the source chain, or if the root is not blessed but RMN
    // blessing is enabled, we revert.
    if (isBlessed == sourceChainConfig.isRMNVerificationDisabled) {
      revert RootBlessingMismatch(sourceChainSelector, root.merkleRoot, isBlessed);
    }

    if (keccak256(root.onRampAddress) != keccak256(sourceChainConfig.onRamp)) {
      revert CommitOnRampMismatch(root.onRampAddress, sourceChainConfig.onRamp);
    }

    if (sourceChainConfig.minSeqNr != root.minSeqNr || root.minSeqNr > root.maxSeqNr) {
      revert InvalidInterval(sourceChainSelector, root.minSeqNr, root.maxSeqNr);
    }

    bytes32 merkleRoot = root.merkleRoot;
    if (merkleRoot == bytes32(0)) revert InvalidRoot();
    // If we reached this section, the report should contain a valid root.
    // We disallow duplicate roots as that would reset the timestamp and delay potential manual execution.
    if (s_roots[sourceChainSelector][merkleRoot] != 0) {
      revert RootAlreadyCommitted(sourceChainSelector, merkleRoot);
    }

    sourceChainConfig.minSeqNr = root.maxSeqNr + 1;
    s_roots[sourceChainSelector][merkleRoot] = block.timestamp;
  }

  /// @notice Returns the sequence number of the last price update.
  /// @return sequenceNumber The latest price update sequence number.
  function getLatestPriceSequenceNumber() external view returns (uint64) {
    return s_latestPriceSequenceNumber;
  }

  /// @notice Returns the timestamp of a potentially previously committed merkle root.
  /// If the root was never committed 0 will be returned.
  /// @param sourceChainSelector The source chain selector.
  /// @param root The merkle root to check the commit status for.
  /// @return timestamp The timestamp of the committed root or zero in the case that it was never committed.
  function getMerkleRoot(uint64 sourceChainSelector, bytes32 root) external view returns (uint256) {
    return s_roots[sourceChainSelector][root];
  }

  /// @notice Returns timestamp of when root was accepted or 0 if verification fails.
  /// @dev This method uses a merkle tree within a merkle tree, with the hashedLeaves,
  /// proofs and proofFlagBits being used to get the root of the inner tree.
  /// This root is then used as the singular leaf of the outer tree.
  /// @return timestamp The commit timestamp of the root.
  function _verify(
    uint64 sourceChainSelector,
    bytes32[] memory hashedLeaves,
    bytes32[] memory proofs,
    uint256 proofFlagBits
  ) internal view virtual returns (uint256 timestamp) {
    bytes32 root = MerkleMultiProof._merkleRoot(hashedLeaves, proofs, proofFlagBits);
    return s_roots[sourceChainSelector][root];
  }

  /// @inheritdoc MultiOCR3Base
  function _afterOCR3ConfigSet(
    uint8 ocrPluginType
  ) internal override {
    bool isSignatureVerificationEnabled = s_ocrConfigs[ocrPluginType].configInfo.isSignatureVerificationEnabled;

    if (ocrPluginType == uint8(Internal.OCRPluginType.Commit)) {
      // Signature verification must be enabled for commit plugin.
      if (!isSignatureVerificationEnabled) {
        revert SignatureVerificationRequiredInCommitPlugin();
      }
      // When the OCR config changes, we reset the sequence number  since it is scoped per config digest.
      // Note that s_minSeqNr/roots do not need to be reset as the roots persist across reconfigurations
      // and are de-duplicated separately.
      s_latestPriceSequenceNumber = 0;
    } else if (ocrPluginType == uint8(Internal.OCRPluginType.Execution)) {
      // Signature verification must be disabled for execution plugin.
      if (isSignatureVerificationEnabled) {
        revert SignatureVerificationNotAllowedInExecutionPlugin();
      }
    }
  }

  // ================================================================
  // │                           Config                             │
  // ================================================================

  /// @notice Returns the static config.
  /// @dev This function will always return the same struct as the contents is static and can never change.
  /// @return staticConfig The static config.
  function getStaticConfig() external view returns (StaticConfig memory) {
    return StaticConfig({
      chainSelector: i_chainSelector,
      gasForCallExactCheck: i_gasForCallExactCheck,
      rmnRemote: i_rmnRemote,
      tokenAdminRegistry: i_tokenAdminRegistry,
      nonceManager: i_nonceManager
    });
  }

  /// @notice Returns the current dynamic config.
  /// @return dynamicConfig The current dynamic config.
  function getDynamicConfig() external view returns (DynamicConfig memory) {
    return s_dynamicConfig;
  }

  /// @notice Returns the source chain config for the provided source chain selector.
  /// @param sourceChainSelector chain to retrieve configuration for.
  /// @return sourceChainConfig The config for the source chain.
  function getSourceChainConfig(
    uint64 sourceChainSelector
  ) external view returns (SourceChainConfig memory) {
    return s_sourceChainConfigs[sourceChainSelector];
  }

  /// @notice Returns all source chain configs.
  /// @return sourceChainConfigs The source chain configs corresponding to all the supported chain selectors.
  function getAllSourceChainConfigs() external view returns (uint64[] memory, SourceChainConfig[] memory) {
    SourceChainConfig[] memory sourceChainConfigs = new SourceChainConfig[](s_sourceChainSelectors.length());
    uint64[] memory sourceChainSelectors = new uint64[](s_sourceChainSelectors.length());
    for (uint256 i = 0; i < s_sourceChainSelectors.length(); ++i) {
      sourceChainSelectors[i] = uint64(s_sourceChainSelectors.at(i));
      sourceChainConfigs[i] = s_sourceChainConfigs[sourceChainSelectors[i]];
    }
    return (sourceChainSelectors, sourceChainConfigs);
  }

  /// @notice Updates source configs.
  /// @param sourceChainConfigUpdates Source chain configs.
  function applySourceChainConfigUpdates(
    SourceChainConfigArgs[] memory sourceChainConfigUpdates
  ) external onlyOwner {
    _applySourceChainConfigUpdates(sourceChainConfigUpdates);
  }

  /// @notice Updates source configs.
  /// @param sourceChainConfigUpdates Source chain configs.
  function _applySourceChainConfigUpdates(
    SourceChainConfigArgs[] memory sourceChainConfigUpdates
  ) internal {
    for (uint256 i = 0; i < sourceChainConfigUpdates.length; ++i) {
      SourceChainConfigArgs memory sourceConfigUpdate = sourceChainConfigUpdates[i];
      uint64 sourceChainSelector = sourceConfigUpdate.sourceChainSelector;

      if (sourceChainSelector == 0) {
        revert ZeroChainSelectorNotAllowed();
      }

      if (address(sourceConfigUpdate.router) == address(0)) {
        revert ZeroAddressNotAllowed();
      }

      SourceChainConfig storage currentConfig = s_sourceChainConfigs[sourceChainSelector];
      bytes memory newOnRamp = sourceConfigUpdate.onRamp;

      if (currentConfig.onRamp.length == 0) {
        currentConfig.minSeqNr = 1;
        emit SourceChainSelectorAdded(sourceChainSelector);
      } else {
        if (currentConfig.minSeqNr != 1 && keccak256(currentConfig.onRamp) != keccak256(newOnRamp)) {
          // OnRamp updates should only happens due to a misconfiguration.
          // If an OnRamp is misconfigured, no reports should have been committed and no messages should have been
          // executed. This is enforced by the onRamp address check in the commit function.
          revert InvalidOnRampUpdate(sourceChainSelector);
        }
      }

      // OnRamp can never be zero - if it is, then the source chain has been added for the first time.
      if (newOnRamp.length == 0 || keccak256(newOnRamp) == EMPTY_ENCODED_ADDRESS_HASH) {
        revert ZeroAddressNotAllowed();
      }

      currentConfig.onRamp = newOnRamp;
      currentConfig.isEnabled = sourceConfigUpdate.isEnabled;
      currentConfig.router = sourceConfigUpdate.router;
      currentConfig.isRMNVerificationDisabled = sourceConfigUpdate.isRMNVerificationDisabled;

      // We don't need to check the return value, as inserting the item twice has no effect.
      s_sourceChainSelectors.add(sourceChainSelector);

      emit SourceChainConfigSet(sourceChainSelector, currentConfig);
    }
  }

  /// @notice Sets the dynamic config.
  /// @param dynamicConfig The new dynamic config.
  function setDynamicConfig(
    DynamicConfig memory dynamicConfig
  ) external onlyOwner {
    _setDynamicConfig(dynamicConfig);
  }

  /// @notice Sets the dynamic config.
  /// @param dynamicConfig The dynamic config.
  function _setDynamicConfig(
    DynamicConfig memory dynamicConfig
  ) internal {
    if (dynamicConfig.feeQuoter == address(0)) {
      revert ZeroAddressNotAllowed();
    }

    s_dynamicConfig = dynamicConfig;

    emit DynamicConfigSet(dynamicConfig);
  }

  /// @notice Returns a source chain config with a check that the config is enabled.
  /// @param sourceChainSelector Source chain selector to check for cursing.
  /// @return sourceChainConfig The source chain config storage pointer.
  function _getEnabledSourceChainConfig(
    uint64 sourceChainSelector
  ) internal view returns (SourceChainConfig storage) {
    SourceChainConfig storage sourceChainConfig = s_sourceChainConfigs[sourceChainSelector];
    if (!sourceChainConfig.isEnabled) {
      revert SourceChainNotEnabled(sourceChainSelector);
    }

    return sourceChainConfig;
  }

  // ================================================================
  // │                            Access                            │
  // ================================================================

  /// @notice Reverts as this contract should not be able to receive CCIP messages.
  function ccipReceive(
    Client.Any2EVMMessage calldata
  ) external pure {
    // solhint-disable-next-line
    revert();
  }
}
